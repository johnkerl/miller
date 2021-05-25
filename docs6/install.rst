..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Installation
================================================================

Prebuilt executables via package managers
----------------------------------------------------------------

`Homebrew <https://brew.sh/>`_ installation support for OSX is available via

::

    brew update && brew install miller

...and also via `MacPorts <https://www.macports.org/>`_:

::

    sudo port selfupdate && sudo port install miller

You may already have the ``mlr`` executable available in your platform's package manager on NetBSD, Debian Linux, Ubuntu Xenial and upward, Arch Linux, or perhaps other distributions. For example, on various Linux distributions you might do one of the following:

::

    sudo apt-get install miller

::

    sudo apt install miller

::

    sudo yum install miller

On Windows, Miller is available via `Chocolatey <https://chocolatey.org/>`_:

::

    choco install miller

Prebuilt executables via GitHub per release
----------------------------------------------------------------

Please see https://github.com/johnkerl/miller/releases where there are builds for OSX Yosemite, Linux x86-64 (dynamically linked), and Windows.

Miller is autobuilt for **Linux**, **MacOS**, and **Windows** using **GitHub Actions** on every commit (https://github.com/johnkerl/miller/actions).

Building from source
----------------------------------------------------------------

Please see :doc:`build`.
