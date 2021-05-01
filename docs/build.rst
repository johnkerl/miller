..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Building from source
================================================================

Please also see :doc:`install` for information about pre-built executables.

Miller license
----------------------------------------------------------------

Two-clause BSD license https://github.com/johnkerl/miller/blob/master/LICENSE.txt.

From release tarball using autoconfig
----------------------------------------------------------------

Miller allows you the option of using GNU ``autoconfigure`` to build portably.

Grateful acknowledgement: Miller's GNU autoconfig work was done by the generous and expert efforts of `Thomas Klausner <https://github.com/0-wiz-0/>`_.

* Obtain ``mlr-i.j.k.tar.gz`` from https://github.com/johnkerl/miller/tags, replacing ``i.j.k`` with the desired release, e.g. ``2.2.1``.
* ``tar zxvf mlr-i.j.k.tar.gz``
* ``cd mlr-i.j.k``
* Install the following packages using your system's package manager (``apt-get``, ``yum install``, etc.): **flex**
* Various configuration options of your choice, e.g.

  * ``./configure``
  * ``./configure --prefix=/usr/local``
  * ``./configure --prefix=$HOME/pkgs``
  * ``./configure CC=clang``
  * ``./configure --disable-shared`` (to make a statically linked executable)
  * ``./configure 'CFLAGS=-Wall -std=gnu99 -O3'``
  * etc.

* ``make`` creates the ``c/mlr`` executable
* ``make check``
* ``make install`` copies the ``c/mlr`` executable to your prefix's ``bin`` subdirectory.

From git clone using autoconfig
----------------------------------------------------------------

* ``git clone https://github.com/johnkerl/miller``
* ``cd miller``
* Install the following packages using your system's package manager (``apt-get``, ``yum install``, etc.): **automake autoconf libtool flex**
* Run ``autoreconf -fiv``. (This is necessary when building from head as discussed in https://github.com/johnkerl/miller/issues/131.)
* Then continue from "Install the following ... " as above.

Without using autoconfig
----------------------------------------------------------------

GNU autoconfig is familiar to many users, and indeed plenty of folks won't bother to use an open-source software package which doesn't have autoconfig support. And this is for good reason: GNU autoconfig allows us to build software on a wide diversity of platforms. For this reason I'm happy that Miller supports autoconfig.

But, many others (myself included!) find autoconfig confusing: if it works without errors, great, but if not, the ``./configure && make`` output can be exceedingly difficult to decipher. And this also can be a turn-off for using open-source software: if you can't figure out the build errors, you may just keep walking. For this reason I'm happy that Miller allows you to build without autoconfig. (Of course, if you have any build errors, feel free to contact me at mailto:kerl.john.r+miller@gmail.com -- or, better, open an issue with "New Issue" at https://github.com/johnkerl/miller/issues.)

Steps:

* Obtain a release tarball or git clone.
* ``cd`` into the ``c`` subdirectory.
* Edit the ``INSTALLDIR`` in ``Makefile.no-autoconfig``.
* To change the C compiler, edit the ``CC=`` lines in ``Makefile.no-autoconfig`` and ``dsls/Makefile.no-autoconfig``.
* ``make -f Makefile.no-autoconfig`` creates the ``mlr`` executable and runs unit/regression tests (i.e. the equivalent of both ``make`` and ``make check`` using autoconfig).
* ``make install`` copies the ``mlr`` executable to your install directory.

The ``Makefile.no-autoconfig`` is simple: little more than ``gcc *.c``.  Customzing is less automatic than autoconfig, but more transparent. I expect this makefile to work with few modifications on a large fraction of modern Linux/BSD-like systems: I'm aware of successful use with ``gcc`` and ``clang``, on Ubuntu 12.04 LTS, SELinux, Darwin (MacOS Yosemite), and FreeBSD.

Windows
----------------------------------------------------------------

*Disclaimer: I'm now relying exclusively on* `Appveyor <https://ci.appveyor.com/project/johnkerl/miller>`_ *for Windows builds; I haven't built from source using MSYS in quite a while.*

Miller has been built on Windows using MSYS2: http://www.msys2.org/.  You can install MSYS2 and build Miller from its source code within MSYS2, and then you can use the binary from outside MSYS2.  You can also use a precompiled binary (see above).

You will first need to install MSYS2: http://www.msys2.org/.  Then, start an MSYS2 shell, e.g. (supposing you installed MSYS2 to ``C:\msys2\``) run ``C:\msys2\mingw64.exe``.  Within the MSYS2 shell, you can run the following to install dependent packages:

::

    pacman -Syu
    pacman -Su
    pacman -S base-devel
    pacman -S msys2-devel
    pacman -S mingw-w64-x86_64-toolchain
    pacman -S mingw-w64-x86_64-pcre
    pacman -S msys2-runtime

The list of dependent packages may be also found in **appveyor.yml** in the Miller base directory.

Then, simply run **msys2-build.sh** which is a thin wrapper around ``./configure && make`` which accommodates certain Windows/MSYS2 idiosyncracies.

There is a unit-test false-negative issue involving the semantics of the ``mkstemp`` library routine but a ``make -k`` in the ``c`` subdirectory has been producing a ``mlr.exe`` for me.

Within MSYS2 you can run ``mlr``: simply copy it from the ``c`` subdirectory to your desired location somewhere within your MSYS2 ``$PATH``.  To run ``mlr`` outside of MSYS2, just as with precompiled binaries as described above, you'll need ``msys-2.0.dll``.  One way to do this is to augment your path:

::

    C:\> set PATH=%PATH%;\msys64\mingw64\bin

Another way to do it is to copy the Miller executable and the DLL to the same directory:

::

    C:\> mkdir \mbin
    C:\> copy \msys64\mingw64\bin\msys-2.0.dll \mbin
    C:\> copy \msys64\wherever\you\installed\miller\c\mlr.exe \mbin
    C:\> set PATH=%PATH%;\mbin


In case of problems
----------------------------------------------------------------

If you have any build errors, feel free to contact me at mailto:kerl.john.r+miller@gmail.com -- or, better, open an issue with "New Issue" at https://github.com/johnkerl/miller/issues.

Dependencies
----------------------------------------------------------------

Required external dependencies
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are necessary to produce the ``mlr`` executable.

* ``gcc``, ``clang``, etc. (or presumably other compilers; please open an issue or send me a pull request if you have information for me about other 21st-century compilers)
* The standard C library
* ``flex``
* ``automake``, ``autoconf``, and ``libtool``, if you build with autoconfig

Optional external dependencies
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

This documentation pageset is built using Sphinx. Please see `./README.md` for details.

Internal dependencies
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

These are included within the `Miller source tree <https://github.com/johnkerl/miller>`_ and do not need to be separately installed (and in fact any separate installation will not be picked up in the Miller build):

* `Mersenne Twister <http://en.wikipedia.org/wiki/Mersenne_Twister>`_ for pseudorandom-number generation: `C implementation by Nishimura and Matsumoto <https://github.com/johnkerl/miller/blob/master/c/lib/mtrand.c>`_ with license terms respected.
* `MinUnit <http://www.jera.com/techinfo/jtns/jtn002.html>`_ for unit-testing, with as-is-no-warranty license http://www.jera.com/techinfo/jtns/jtn002.html#License, https://github.com/johnkerl/miller/blob/master/c/lib/minunit.h.
* The `Lemon parser-generator <http://www.hwaci.com/sw/lemon/>`_, the author of which explicitly disclaims copyright.
* The `udp JSON parser <https://github.com/udp/json-parser>`_, with BSD2 license.
* The `sheredom UTF-8 library <https://github.com/sheredom/utf8.h>`_, which is free and unencumbered software released into the public domain.
* The NetBSD ``strptime`` (needed for the Windows/MSYS2 port since MSYS2 lacks this), with BSD license.

Creating a new release: for developers
----------------------------------------------------------------

At present I'm the primary developer so this is just my checklist for making new releases.

In this example I am using version 3.4.0; of course that will change for subsequent revisions.

* Update version found in ``mlr --version`` and ``man mlr``:

  * Edit ``configure.ac``, ``c/mlrvers.h``, ``miller.spec``, and ``docs/conf.py`` from ``3.3.2-dev`` to ``3.4.0``.
  * Do a fresh ``autoreconf -fiv`` and commit the output. (Preferably on a Linux host, rather than MacOS, to reduce needless diffs in autogen build files.)
  * ``make -C c -f Makefile.no-autoconfig installhome && make -C man -f Makefile.no-autoconfig installhome && make -C docs -f Makefile.no-autoconfig html``
  * The ordering is important: the first build creates ``mlr``; the second runs ``mlr`` to create ``manpage.txt``; the third includes ``manpage.txt`` into one of its outputs.
  * Commit and push.

* Create the release tarball and SRPM:

  * On buildbox: ``./configure && make distcheck``
  * On buildbox: make SRPM as in https://github.com/johnkerl/miller/blob/master/README-RPM.md
  * On all buildboxes: ``cd c`` and ``make -f Makefile.no-autoconfig mlr.static``. Then copy ``mlr.static`` to ``../mlr.{arch}``. (This may require as prerequisite ``sudo yum install glibc-static`` or the like.)
  * For static binaries, please do ``ldd mlr.static`` and make sure it says ``not a dynamic executable``.
  * Then ``mv mlr.static ../mlr.linux_x86_64``
  * Pull back release tarball ``mlr-3.4.0.tar.gz`` and SRPM ``miller-3.4.0-1.el6.src.rpm`` from buildbox, and ``mlr.{arch}`` binaries from whatever buildboxes.
  * Download ``mlr.exe`` and ``msys-2.0.dll`` from https://ci.appveyor.com/project/johnkerl/miller/build/artifacts.

* Create the Github release tag:

  * Don't forget the ``v`` in ``v3.4.0``
  * Write the release notes
  * Attach the release tarball, SRPM, and binaries. Double-check assets were successfully uploaded.
  * Publish the release

* Check the release-specific docs:

  * Look at https://miller.readthedocs.io for new-version docs, after a few minutes' propagation time.

* Notify:

  * Submit ``brew`` pull request; notify any other distros which don't appear to have autoupdated since the previous release (notes below)
  * Similarly for ``macports``: https://github.com/macports/macports-ports/blob/master/textproc/miller/Portfile.
  * Social-media updates.

::

    git remote add upstream https://github.com/Homebrew/homebrew-core # one-time setup only
    git fetch upstream
    git rebase upstream/master
    git checkout -b miller-3.4.0
    shasum -a 256 /path/to/mlr-3.4.0.tar.gz
    edit Formula/miller.rb
    # Test the URL from the line like
    #   url "https://github.com/johnkerl/miller/releases/download/v3.4.0/mlr-3.4.0.tar.gz"
    # in a browser for typos
    # A '@BrewTestBot Test this please' comment within the homebrew-core pull request will restart the homebrew travis build
    git add Formula/miller.rb
    git commit -m 'miller 3.4.0'
    git push -u origin miller-3.4.0
    (submit the pull request)

* Afterwork:

  * Edit ``configure.ac`` and ``c/mlrvers.h`` to change version from ``3.4.0`` to ``3.4.0-dev``.
  * ``make -C c -f Makefile.no-autoconfig installhome && make -C doc -f Makefile.no-autoconfig all installhome``
  * Commit and push.


Misc. development notes
----------------------------------------------------------------

I use terminal width 120 and tabwidth 4.
