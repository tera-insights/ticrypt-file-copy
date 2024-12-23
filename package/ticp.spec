%global version 1.0.0

Name:       ticp
Version:    %{version}
Release:    el
Summary:    Simple CLI/Library to copy files at high speed
Vendor:     Tera Insights, LLC
Packager:   Vishisht Khilariwal <vishy@terainsights.com>

Group:      Applications/System
License:    tiCrypt VM Controller License
URL:        https://github.com/tera-insights/ticrypt-file-copy
Source0:    https://github.com/tera-insights/ticrypt-file-copy

BuildRoot:  /tmp/ticp

%description
High speed file copy
Daemon mode (Accepts requests over a tcp websocket)[In testing]
Benchmark mode
Progress bar
CLI interface
Recovery Mode

%install

# Install executables in /bin/ticp
install -d %{buildroot}/bin
install -p -m 755 bin/ticp %{buildroot}/bin/ticp

# Install configs in etc
install -d %{buildroot}/%{_sysconfdir}/ticrypt
install -p -m 600 doc/image.example.toml %{buildroot}/%{_sysconfdir}/ticrypt/controller.toml

# Install SystemD Unit file
install -d %{buildroot}%{_unitdir}
install -p -m 644 package/ticp.service %{buildroot}%{_unitdir}

%post
%systemd_post ticp.service

%preun
%systemd_preun ticp.service

%postun
%systemd_postun ticp.service

%doc doc/licenses

%dir    %attr(755, root, root)      %{_sysconfdir}/ticrypt
%dir    %attr(700, root, root)      %{_libexecdir}/ticrypt

%attr(755, root, root)  %{_libexecdir}/ticrypt/stub
%attr(644, root, root)  %{_unitdir}/ticp.service

%changelog
* Thu Dec 19 2024 Vishisht Khilariwal <vishy@terainsights.com> 1.0.0
- First release