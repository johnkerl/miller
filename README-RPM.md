## Instructions to make Miller source/binary RPMs for the RPM-experienced

Grab `miller.spec` and go to town.

## Instructions to make Miller source/binary RPMs for the RPM-inexperienced


### One-time setup
Change `3.3.2` to desired version. Release-package example:

https://github.com/johnkerl/miller/releases/download/v3.3.2/mlr-3.3.2.tar.gz


```
cd
mkdir ~/rpmbuild
mkdir ~/rpmbuild/SPECS
mkdir ~/rpmbuild/SOURCES
cp /your/path/to/miller/clone/miller.spec ~/rpmbuild/SPECS
cp /your/path/to/mlr-3.3.2.tar.gz ~/rpmbuild/SOURCES

cd ~/rpmbuild/SPECS
```

You may need to do
```
sudo yum install rpmbuild
```

### Linting
If you have changed the `miller.spec` file:
```
rpmlint miller.spec
```
You may need to do
```
sudo yum install rpmlint
```

### Build source-RPM only
```
rpmbuild -bs miller.spec
Wrote: /your/home/dir/rpmbuild/SRPMS/mlr-3.3.2-1.el6.src.rpm
```

```
rpm -qpl ../SRPMS/mlr-3.3.2-1.el6.src.rpm
mlr-3.3.2.tar.gz
miller.spec
```

```
rpm -qpi ../SRPMS/mlr-3.3.2-1.el6.src.rpm
Name        : mlr                          Relocations: (not relocatable)
Version     : 3.3.2                             Vendor: (none)
Release     : 1.el6                         Build Date: Sun 07 Feb 2016 09:43:39 PM EST
Install Date: (not installed)               Build Host: host.name.goes.here
Group       : Applications/Text             Source RPM: (none)
Size        : 774430                           License: BSD2
Signature   : (none)
URL         : http://johnkerl.org/miller/doc
Summary     : Name-indexed data processing tool
Description :
Miller (mlr) allows name-indexed data such as CSV and JSON files to be
processed with functions equivalent to sed, awk, cut, join, sort etc. It can
convert between formats, preserves headers when sorting or reversing, and
streams data where possible so its memory requirements stay small. It works
well with pipes and can feed "tail -f".
```

### Build source and binary RPMs

```
rpmbuild -ba miller.spec
```

```
rpm -qpl ../RPMS/x86_64//mlr-3.3.2-1.el6.x86_64.rpm
/usr/bin/mlr
/usr/share/man/man1/mlr.1.gz
```

```
sudo rpm -ivh ../RPMS/x86_64/mlr-3.3.2-1.el6.x86_64.rpm 
Preparing...                ########################################### [100%]
   1:mlr                    ########################################### [100%]
```

```
/usr/bin/mlr --version
Miller 3.3.2

man -M /usr/share/man mlr
```
and check the version in the DESCRIPTION section.

### Some handy references

* https://github.com/bonzini/grep/blob/master/grep.spec
* http://www.rpm.org/max-rpm/s1-rpm-build-creating-spec-file.html
* http://www.rpm.org/max-rpm/s1-rpm-inside-files-list-directives.html
* http://www.tldp.org/HOWTO/RPM-HOWTO/build.html
* http://www.tldp.org/LDP/solrhe/Securing-Optimizing-Linux-RH-Edition-v1.3/chap3sec20.html
* https://fedoraproject.org/wiki/How_to_create_a_GNU_Hello_RPM_package
