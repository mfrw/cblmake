Summary:        Build CBL-Mariner derrivates
Name:           cblmake
Version:        0.0.1
Release:        1%{?dist}
License:        MIT
Vendor:         Microsoft Corporation
Distribution:   Mariner
URL:            https://github.com/mfrw/cblmake
Source:         https://github.com/mfrw/cblmake/archive/release/%{name}-%{version}.tar.gz

# Vendor All the sources required for the build to work.
# Source1:      %{name}-vendor-%{version}.tar.gz

BuildRequires:  golang

%description
A tool to build build on top of the CBL-Mariner toolkit to enable rpm/iso/vhd builds.

%prep
%autosetup

%build
go build -v -tags netgo ./cmd/cblamake

%install
# Just ship the the static binary
# TODO

%files
%license
%{_bindir}/cblmake

%changelog
* Thu May 12 2022 Muhammad Falak <mwani@microsoft.com> - 0.0.1-1
- Inital CBL-Mariner import
