## Instructions to make Miller source/binary RPMs for the RPM-experienced

Grab `miller.spec` and go to town.

## Instructions to make Miller source/binary RPMs for the RPM-inexperienced

### One-time setup
Change `6.2.0` to desired version. Release-package example:

https://github.com/johnkerl/miller/releases/download/v6.2.0/miller-6.2.0.tar.gz

```
cd
mkdir ~/rpmbuild
mkdir ~/rpmbuild/SPECS
mkdir ~/rpmbuild/SOURCES
cp /your/path/to/miller/clone/miller.spec ~/rpmbuild/SPECS
cp /your/path/to/miller-6.2.0.tar.gz ~/rpmbuild/SOURCES

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
Wrote: /your/home/dir/rpmbuild/SRPMS/miller-6.2.0-1.el6.src.rpm
```

```
rpm -qpl ../SRPMS/miller-6.2.0-1.el6.src.rpm
miller-6.2.0.tar.gz
miller.spec
```

```
rpm -qpi ../SRPMS/miller-6.2.0-1.el6.src.rpm
Name        : mlr                          Relocations: (not relocatable)
Version     : 6.2.0                             Vendor: (none)
...
```

### Build source and binary RPMs

```
rpmbuild -ba miller.spec
```

```
rpm -qpl ../RPMS/x86_64//miller-6.2.0-1.el6.x86_64.rpm
/usr/bin/mlr
/usr/share/man/man1/mlr.1.gz
```

```
sudo rpm -ivh ../RPMS/x86_64/miller-6.2.0-1.el6.x86_64.rpm 
Preparing...                ########################################### [100%]
   1:mlr                    ########################################### [100%]
```

```
/usr/bin/mlr --version
Miller 6.2.0

man -M /usr/share/man mlr
```
and check the version in the DESCRIPTION section.
