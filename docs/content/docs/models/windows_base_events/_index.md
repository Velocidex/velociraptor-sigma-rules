# Windows Base Sigma Model

This model is designed for triage of dead disk, or file based live
analysis. The rules that use this model will be evaluated once on
all events.

After all relevant rules are evaluated, the collection is complete.

# Log Sources

Following is a list of recognized log sources.


## `*/windows/diagnosis-scripted`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: diagnosis-scripted
```


## `*/windows/driver-framework`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: driver-framework
```


## `*/windows/msexchange-management`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: msexchange-management
```


## `*/windows/powershell-classic`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: powershell-classic
```


## `*/windows/dns-client`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: dns-client
```


## `*/windows/lsa-server`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: lsa-server
```


## `*/windows/sysmon`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: sysmon
```


## `ps_classic_provider_start/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_classic_provider_start
  product: windows
```


## `registry_add/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: registry_add
  product: windows
```


## `*/windows/capi2`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: capi2
```


## `*/windows/firewall-as`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: firewall-as
```


## `*/windows/ldap_debug`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: ldap_debug
```


## `*/windows/powershell`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: powershell
```


## `ps_script/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_script
  product: windows
```


## `*/windows/application`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: application
```


## `*/windows/applocker`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: applocker
```


## `*/windows/dns-server`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: dns-server
```


## `*/windows/dns-server-analytic`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: dns-server-analytic
```


## `*/windows/microsoft-servicebus-client`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: microsoft-servicebus-client
```


## `*/windows/openssh`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: openssh
```


## `*/windows/security`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: security
```


## `*/windows/taskscheduler`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: taskscheduler
```


## `*/windows/wmi`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: wmi
```


## `registry_set/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: registry_set
  product: windows
```


## `*/windows/certificateservicesclient-lifecycle-system`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: certificateservicesclient-lifecycle-system
```


## `process_creation/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: process_creation
  product: windows
```


## `ps_classic_start/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_classic_start
  product: windows
```


## `registry_event/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: registry_event
  product: windows
```


## `antivirus/windows/windefend`





Sample use in a sigma rule:
```yaml
log_sources:
  category: antivirus
  product: windows
  service: windefend
```


## `*/windows/appmodel-runtime`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: appmodel-runtime
```


## `*/windows/bits-client`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: bits-client
```


## `*/windows/codeintegrity-operational`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: codeintegrity-operational
```


## `*/windows/security-mitigations`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: security-mitigations
```


## `*/windows/system`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: system
```


## `*/windows/windefend`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: windefend
```


## `*/windows/appxdeployment-server`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: appxdeployment-server
```


## `*/windows/appxpackaging-om`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: appxpackaging-om
```


## `*/windows/shell-core`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: shell-core
```


## `*/windows/vhdmp`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: vhdmp
```


## `ps_module/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_module
  product: windows
```


## `*/windows/ntlm`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: ntlm
```


## `*/windows/smbclient-security`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: smbclient-security
```


## `*/windows/terminalservices-localsessionmanager`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: terminalservices-localsessionmanager
```


