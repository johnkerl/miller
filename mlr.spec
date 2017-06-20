Summary: Name-indexed data processing tool
Name: mlr
Version: 5.2.1
Release: 1%{?dist}
License: BSD2
Group: Applications/Text
Source: https://github.com/johnkerl/miller/releases/download/v%{version}/%{name}-%{version}.tar.gz
URL: http://johnkerl.org/miller/doc
Buildroot: %{_tmppath}/%{name}-%{version}-root
BuildRequires: flex >= 2.5.35

%description
Miller (mlr) allows name-indexed data such as CSV and JSON files to be
processed with functions equivalent to sed, awk, cut, join, sort etc. It can
convert between formats, preserves headers when sorting or reversing, and
streams data where possible so its memory requirements stay small. It works
well with pipes and can feed "tail -f".

%prep
%setup -q -n %{name}-%{version}

%build
%configure
make

%check
make check

%install
rm -rf ${RPM_BUILD_ROOT}
make install DESTDIR="$RPM_BUILD_ROOT"

%clean
rm -rf ${RPM_BUILD_ROOT}

%files
%defattr(755, root, root, -)
%{_bindir}/mlr
%defattr(644, root, root, -)
%{_mandir}/man1/mlr.1.gz
%defattr(-,root,root)

%changelog
* Mon Jun 19 2017 John Kerl <kerl.john.r@gmail.com>
- 5.2.1 release
* Sun Jun 11 2017 John Kerl <kerl.john.r@gmail.com>
- 5.2.0 release
* Thu Apr 13 2017 John Kerl <kerl.john.r@gmail.com>
- 5.1.0 release
* Sat Mar 11 2017 John Kerl <kerl.john.r@gmail.com>
- 5.0.1 release
* Mon Feb 27 2017 John Kerl <kerl.john.r@gmail.com>
- 5.0.0 release
* Sun Aug 21 2016 John Kerl <kerl.john.r@gmail.com>
- 4.5.0 release
* Mon Apr 04 2016 John Kerl <kerl.john.r@gmail.com>
- 3.5.0 release
* Sun Feb 14 2016 John Kerl <kerl.john.r@gmail.com>
- 3.4.0 release
* Sun Feb 07 2016 John Kerl <kerl.john.r@gmail.com>
- Initial spec-file submission for Miller
