Summary: Name-indexed data processing tool
Name: mlr
Version: 3.3.2
Release: 1%{?dist}
License: BSD2
Group: Applications/Text
Source: https://github.com/johnkerl/miller/releases/download/v%{version}/%{name}-%{version}.tar.gz
URL: http://johnkerl.org/miller/doc
Distribution: Fedora Project
BuildRequires: flex >= 2.5.35

%description
Miller (mlr) allows name-indexed data such as CSV and JSON files to be
processed with functions equivalent to sed, awk, cut, join, sort etc. It can
convert between formats, preserves headers when sorting or reversing, and
streams data where possible so its memory requirements stay small. It works
well with pipes and can feed "tail -f".

%prep
%setup -q

%build
%configure
make

%check
make check

%install
make install
make clean

%clean
make clean

%files
%defattr(755, root, root, -)
%{_bindir}/mlr
%defattr(644, root, root, -)
%{_mandir}/man1/mlr.1
%attr(644, root, root) 
%license LICENSE.txt

%doc README.md

%changelog
* Sun Feb 07 2016 John Kerl <kerl.john.r@gmail.com>
- Initial RedHat/Fedoraa submission of Miller
