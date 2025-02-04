:orphan:

.. _cheat-sheet:

reStructuredText cheat sheet
============================

This file contains the syntax for commonly used reST markup.
See the `reStructuredText style guide`_ for detailed information and conventions.

Also see the `Sphinx reStructuredText Primer`_ for more details on reST, and the `Canonical Documentation Style Guide`_ for general style conventions.

H2 heading
----------

H3 heading
~~~~~~~~~~

H4 heading
^^^^^^^^^^

H5 heading
..........

Inline formatting
-----------------

- :guilabel:`UI element`
- ``code``
- :file:`file path`
- :command:`command`
- :kbd:`Key`
- *Italic*
- **Bold**

Code blocks
-----------

Start a code block::

     code:
       - example: true

.. code::

     # Demonstrate a code block
     code:
       - example: true

.. code:: yaml

     # Demonstrate a code block
     code:
       - example: true

.. _a_section_target:

Links
-----

- `Canonical website <https://canonical.com/>`_
- `Canonical website`_ (defined in ``reuse/links.txt``)
- `Link defined at the bottom of the page`_
- https:\ //canonical.com/
- :ref:`a_section_target`
- :ref:`Link text <a_section_target>`
- :doc:`doc-cheat-sheet-myst`
- :doc:`Link text <doc-cheat-sheet-myst>`


Navigation
----------

Use the following syntax::

  .. toctree::
     :hidden:

     sub-page1
     sub-page2


Lists
-----

1. Step 1

   - Item 1

     * Sub-item
   - Item 2

     i. Sub-step 1
     #. Sub-step 2
#. Step 2

   a. Sub-step 1

      - Item
   #. Sub-step 2

Term 1:
  Definition
Term 2:
  Definition

Tables
------

+----------------------+------------+
| Header 1             | Header 2   |
+======================+============+
| Cell 1               | Cell 2     |
|                      |            |
| Second paragraph     |            |
+----------------------+------------+
| Cell 3               | Cell 4     |
+----------------------+------------+

.. list-table::
   :header-rows: 1

   * - Header 1
     - Header 2
   * - Cell 1

       Second paragraph
     - Cell 2
   * - Cell 3
     - Cell 4

.. rst-class:: align-center

   +----------------------+------------+
   | Header 1             | Header 2   |
   +======================+============+
   | Cell 1               | Cell 2     |
   |                      |            |
   | Second paragraph     |            |
   +----------------------+------------+
   | Cell 3               | Cell 4     |
   +----------------------+------------+

.. list-table::
   :header-rows: 1
   :align: center

   * - Header 1
     - Header 2
   * - Cell 1

       Second paragraph
     - Cell 2
   * - Cell 3
     - Cell 4

Notes
-----

.. note::
   A note.

.. tip::
   A tip.

.. important::
   Important information

.. caution::
   This might damage your hardware!

Images
------

.. image:: https://assets.ubuntu.com/v1/b3b72cb2-canonical-logo-166.png

.. figure:: https://assets.ubuntu.com/v1/b3b72cb2-canonical-logo-166.png
   :width: 100px
   :alt: Alt text

   Figure caption

Reuse
-----

.. |reuse_key| replace:: This is **included** text.

|reuse_key|

.. include:: index.rst
   :start-after: include_start
   :end-before: include_end

Tabs
----

.. tabs::

   .. group-tab:: Tab 1

      Content Tab 1

   .. group-tab:: Tab 2

      Content Tab 2


Glossary
--------

.. glossary::

   example term
     Definition of the example term.

:term:`example term`

More useful markup
------------------

- .. versionadded:: X.Y
- | Line 1
  | Line 2
  | Line 3
- .. This is a comment
- :abbr:`API (Application Programming Interface)`

----

Custom extensions
-----------------

Related links at the top of the page::

  :relatedlinks: https://github.com/canonical/lxd-sphinx-extensions, [RTFM](https://www.google.com)
  :discourse: 12345

Terms that should not be checked by the spelling checker: :spellexception:`PurposelyWrong`

A terminal view with input and output:

.. terminal::
   :input: command
   :user: root
   :host: vampyr

   the output

A link to a YouTube video:

.. youtube:: https://www.youtube.com/watch?v=iMLiK1fX4I0
          :title: Demo



.. LINKS
.. _Link defined at the bottom of the page: https://canonical.com/
