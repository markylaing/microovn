package ovn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/canonical/microcluster/state"
	"github.com/lxc/lxd/shared/logger"

	"github.com/canonical/microovn/microovn/database"
	"github.com/canonical/microovn/microovn/ovn/paths"
)

var ovnEnvTpl = template.Must(template.New("ovnEnvTpl").Parse(`# # Generated by MicroOVN, DO NOT EDIT.
OVN_INITIAL_NB="{{ .nbInitial }}"
OVN_INITIAL_SB="{{ .sbInitial }}"
OVN_NB_CONNECT="{{ .nbConnect }}"
OVN_SB_CONNECT="{{ .sbConnect }}"
OVN_LOCAL_IP="{{ .localAddr }}"
`))

// networkProtocol returns appropriate network protocol that should be used
// by OVN services.
func networkProtocol(s *state.State) string {
	_, _, err := getCA(s)
	if err != nil {
		return "tcp"
	} else {
		return "ssl"
	}

}

// localServiceActive function accepts service names (like "central" or "switch") and returns true/false based
// on whether the selected service is running on this node.
func localServiceActive(s *state.State, serviceName string) (bool, error) {
	serviceActive := false
	err := s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		// Get list of all active local services.
		name := s.Name()
		services, err := database.GetServices(ctx, tx, database.ServiceFilter{Member: &name})
		if err != nil {
			return err
		}

		// Check if the specified service is among active local services.
		for _, srv := range services {
			if srv.Service == serviceName {
				serviceActive = true
			}
		}

		return nil
	})

	return serviceActive, err
}

func connectString(s *state.State, port int) (string, error) {
	var err error
	var servers []database.Service

	err = s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		serviceName := "central"
		servers, err = database.GetServices(ctx, tx, database.ServiceFilter{Service: &serviceName})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	addresses := make([]string, 0, len(servers))
	remotes := s.Remotes().RemotesByName()
	protocol := networkProtocol(s)
	for _, server := range servers {
		remote, ok := remotes[server.Member]
		if !ok {
			continue
		}

		addresses = append(
			addresses,
			fmt.Sprintf("%s:%s",
				protocol,
				netip.AddrPortFrom(remote.Address.Addr(), uint16(port)).String(),
			),
		)
	}

	return strings.Join(addresses, ","), nil
}

func generateEnvironment(s *state.State) error {
	// Get the servers.
	nbConnect, err := connectString(s, 6641)
	if err != nil {
		return err
	}

	sbConnect, err := connectString(s, 6642)
	if err != nil {
		return err
	}

	// Get the initial (first server).
	var nbInitial string
	var sbInitial string
	err = s.Database.Transaction(s.Context, func(ctx context.Context, tx *sql.Tx) error {
		serviceName := "central"
		servers, err := database.GetServices(ctx, tx, database.ServiceFilter{Service: &serviceName})
		if err != nil {
			return err
		}

		server := servers[0]

		remotes := s.Remotes().RemotesByName()
		remote, ok := remotes[server.Member]
		if !ok {
			return fmt.Errorf("Remote couldn't be found for %q", server.Member)
		}

		addrString := remote.Address.Addr().String()
		if remote.Address.Addr().Is6() {
			addrString = "[" + addrString + "]"
		}

		nbInitial = addrString
		sbInitial = addrString

		return nil
	})
	if err != nil {
		return err
	}

	// Generate ovn.env.
	fd, err := os.OpenFile(paths.OvnEnvFile(), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't open ovn.env: %w", err)
	}
	defer fd.Close()

	localAddr := s.Address().Hostname()
	if ip, err := netip.ParseAddr(localAddr); err == nil && ip.Is6() {
		localAddr = "[" + localAddr + "]"
	}

	err = ovnEnvTpl.Execute(fd, map[string]any{
		"localAddr": localAddr,
		"nbInitial": nbInitial,
		"sbInitial": sbInitial,
		"nbConnect": nbConnect,
		"sbConnect": sbConnect,
	})
	if err != nil {
		return fmt.Errorf("Couldn't render ovn.env: %w", err)
	}

	return nil
}

func createPaths() error {
	// Create our various paths.
	for _, path := range paths.RequiredDirs() {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("Unable to create %q: %w", path, err)
		}
	}

	return nil
}

// cleanupPaths backs up directories defined by paths.BackupDirs and then removes directories
// created by createPaths function. This effectively removes any data created during MicroOVN runtime.
func cleanupPaths() error {
	var errs []error

	// Create timestamped backup dir
	backupDir := fmt.Sprintf("backup_%d", time.Now().Unix())
	backupPath := filepath.Join(paths.Root(), backupDir)
	err := os.Mkdir(backupPath, 0750)
	if err != nil {
		errs = append(
			errs,
			fmt.Errorf(
				"failed to create backup directory '%s'. Refusing to continue with data removal: %s",
				backupPath,
				err,
			),
		)
		return errors.Join(errs...)
	}

	// Backup selected directories
	for _, dir := range paths.BackupDirs() {
		_, fileName := filepath.Split(dir)
		destination := filepath.Join(backupPath, fileName)
		err = os.Rename(dir, destination)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Return if any backups failed
	if len(errs) > 0 {
		errs = append(
			errs,
			fmt.Errorf("failures occured during backup. Refusing to continue with data removal"),
		)
		return errors.Join(errs...)
	}
	logger.Infof("MicroOVN data backed up to %s", backupPath)

	// Remove rest of the directories
	for _, dir := range paths.RequiredDirs() {
		err = os.RemoveAll(dir)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to remove directory '%s': %w", dir, err))
		}
	}

	return errors.Join(errs...)
}
