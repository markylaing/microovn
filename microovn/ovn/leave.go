package ovn

import (
	"github.com/canonical/microcluster/state"
	"github.com/lxc/lxd/shared/logger"

	"github.com/canonical/microovn/microovn/ovn/paths"
)

// Leave function gracefully departs from the OVN cluster before the member is removed from MicroOVN
// cluster. It ensures that:
//   - OVN chassis is stopped and removed from SB database
//   - OVN NB cluster is cleanly departed
//   - OVN SB cluster is cleanly departed
//
// Note (mkalcok): At this point, database table `services` no longer contains entries
// for departing cluster member, so we'll try to exit/leave/stop all possible services
// ignoring any errors from services that are not actually running.
func Leave(s *state.State) error {
	var err error
	chassisName := s.Name()

	// Gracefully exit OVN controller causing chassis to be automatically removed.
	logger.Infof("Stopping OVN Controller and removing Chassis '%s' from OVN SB database.", chassisName)
	_, err = ControllerCtl(s, "exit")
	if err != nil {
		logger.Warnf("Failed to gracefully stop OVN Controller: %s", err)
	}

	err = snapStop("chassis", true)
	if err != nil {
		logger.Warnf("Failed to stop Chassis service: %s", err)
	}

	err = snapStop("switch", true)
	if err != nil {
		logger.Warnf("Failed to stop Switch service: %s", err)
	}

	// Leave SB and NB clusters
	logger.Info("Leaving OVN Northbound cluster")
	_, err = AppCtl(s, paths.OvnNBControlSock(), "cluster/leave", "OVN_Northbound")
	if err != nil {
		logger.Warnf("Failed to leave OVN Northbound cluster: %s", err)
	}

	logger.Info("Leaving OVN Southbound cluster")
	_, err = AppCtl(s, paths.OvnSBControlSock(), "cluster/leave", "OVN_Southbound")
	if err != nil {
		logger.Warnf("Failed to leave OVN Southbound cluster: %s", err)
	}

	// Wait for NB and SB cluster members to complete departure process
	nbDatabase, err := newOvsdbSpec(OvsdbTypeNBLocal)
	if err == nil {
		err = waitForDBState(s, nbDatabase, OvsdbRemoved, defaultDBConnectWait)
		if err != nil {
			logger.Warnf("Failed to wait for NB cluster departure: %s", err)
		}
	} else {
		logger.Warnf("Failed to get NB database specification: %s", err)
	}

	sbDatabase, err := newOvsdbSpec(OvsdbTypeSBLocal)
	if err == nil {
		err = waitForDBState(s, sbDatabase, OvsdbRemoved, defaultDBConnectWait)
		if err != nil {
			logger.Warnf("Failed to wait for SB cluster departure: %s", err)
		}
	} else {
		logger.Warnf("Failed to get SB database specification: %s", err)
	}

	err = snapStop("central", true)
	if err != nil {
		logger.Warnf("Failed to stop Central service: %s", err)
	}

	logger.Info("Cleaning up runtime and data directories.")
	err = cleanupPaths()
	if err != nil {
		logger.Warn(err.Error())
	}

	return nil
}
