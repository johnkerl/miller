Summary: Name-indexed data processing tool
Name: miller
Version: 6.5.0
Release: 1%{?dist}
License: BSD
Source: https://github.com/johnkerl/miller/releases/download/%{version}/miller-%{version}.tar.gz
URL: https://miller.readthedocs.io
# gcc for cgo transitive dependency
BuildRequires: golang
BuildRequires: gcc
BuildRequires: systemd-rpm-macros

%description
Miller (mlr) allows name-indexed data such as CSV and JSON files to be
processed with functions equivalent to sed, awk, cut, join, sort etc. It can
convert between formats, preserves headers when sorting or reversing, and
streams data where possible so its memory requirements stay small. It works
well with pipes and can feed "tail -f".

%prep
%autosetup

%build
make build

%check
make check

%install
make install

%files
%license LICENSE.txt
%doc README.md
%{_bindir}/mlr
%{_mandir}/man1/mlr.1*

%changelog
* Sun Nov 27 2022 John Kerl <kerl.john.r@gmail.com> - 6.5.0-1
- 6.5.0 release

* Sat Aug 20 2022 John Kerl <kerl.john.r@gmail.com> - 6.4.0-1
- 6.4.0 release

* Thu Jul 7 2022 John Kerl <kerl.john.r@gmail.com> - 6.3.0-1
- 6.3.0 release

* Fri Mar 18 2022 John Kerl <kerl.john.r@gmail.com> - 6.2.0-1
- 6.2.0 release

* Mon Mar 7 2022 John Kerl <kerl.john.r@gmail.com> - 6.1.0-1
- 6.1.0 release

* Sun Jan 9 2022 John Kerl <kerl.john.r@gmail.com> - 6.0.0-1
- 6.0.0 release

* Tue Mar 23 2021 John Kerl <kerl.john.r@gmail.com> - 5.10.2-1
- 5.10.2 release

* Sun Mar 21 2021 John Kerl <kerl.john.r@gmail.com> - 5.10.1-1
- 5.10.1 release

* Sun Nov 29 2020 John Kerl <kerl.john.r@gmail.com> - 5.10.0-1
- 5.10.0 release

* Wed Sep 02 2020 John Kerl <kerl.john.r@gmail.com> - 5.9.1-1
- 5.9.1 release

* Wed Aug 19 2020 John Kerl <kerl.john.r@gmail.com> - 5.9.0-1
- 5.9.0 release

* Mon Aug 03 2020 John Kerl <kerl.john.r@gmail.com> - 5.8.0-1
- 5.8.0 release

* Mon Mar 16 2020 John Kerl <kerl.john.r@gmail.com> - 5.7.0-1
- 5.7.0 release

* Sat Sep 21 2019 John Kerl <kerl.john.r@gmail.com> - 5.6.2-1
- 5.6.2 release

* Mon Sep 16 2019 John Kerl <kerl.john.r@gmail.com> - 5.6.1-1
- 5.6.1 release

* Thu Sep 12 2019 John Kerl <kerl.john.r@gmail.com> - 5.6.0-1
- 5.6.0 release

* Sat Aug 31 2019 John Kerl <kerl.john.r@gmail.com> - 5.5.0-1
- 5.5.0 release

* Tue May 28 2019 Stephen Kitt <steve@sk2.org> - 5.4.0-1
- Fix up for Fedora

* Sun Oct 14 2018 John Kerl <kerl.john.r@gmail.com> - 5.4.0-1
- 5.4.0 release

* Sat Jan 06 2018 John Kerl <kerl.john.r@gmail.com> - 5.3.0-1
- 5.3.0 release

* Thu Jul 20 2017 John Kerl <kerl.john.r@gmail.com> - 5.2.2-1
- 5.2.2 release

* Mon Jun 19 2017 John Kerl <kerl.john.r@gmail.com> - 5.2.1-1
- 5.2.1 release

* Sun Jun 11 2017 John Kerl <kerl.john.r@gmail.com> - 5.2.0-1
- 5.2.0 release

* Thu Apr 13 2017 John Kerl <kerl.john.r@gmail.com> - 5.1.0-1
- 5.1.0 release

* Sat Mar 11 2017 John Kerl <kerl.john.r@gmail.com> - 5.0.1-1
- 5.0.1 release

* Mon Feb 27 2017 John Kerl <kerl.john.r@gmail.com> - 5.0.0-1
- 5.0.0 release

* Sun Aug 21 2016 John Kerl <kerl.john.r@gmail.com> - 4.5.0-1
- 4.5.0 release

* Mon Apr 04 2016 John Kerl <kerl.john.r@gmail.com> - 3.5.0-1
- 3.5.0 release

* Sun Feb 14 2016 John Kerl <kerl.john.r@gmail.com> - 3.4.0-1
- 3.4.0 release

* Sun Feb 07 2016 John Kerl <kerl.john.r@gmail.com> - 3.3.2-1
- Initial spec-file submission for Miller
