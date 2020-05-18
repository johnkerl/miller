Summary: Name-indexed data processing tool
Name: miller
Version: 5.8.0
Release: 1%{?dist}
License: BSD
Source: https://github.com/johnkerl/miller/releases/download/%{version}/mlr-%{version}.tar.gz
URL: http://johnkerl.org/miller/doc
BuildRequires: gcc
BuildRequires: flex >= 2.5.35

%description
Miller (mlr) allows name-indexed data such as CSV and JSON files to be
processed with functions equivalent to sed, awk, cut, join, sort etc. It can
convert between formats, preserves headers when sorting or reversing, and
streams data where possible so its memory requirements stay small. It works
well with pipes and can feed "tail -f".

%prep
%setup -q -n mlr-%{version}

%build
%configure
%make_build

%check
make check

%install
%make_install

%files
%license LICENSE.txt
%doc README.md
%{_bindir}/mlr
%{_mandir}/man1/mlr.1*

%changelog
* Sun May 17 2020 John Kerl <kerl.john.r@gmail.com> - 5.8.0-1
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
