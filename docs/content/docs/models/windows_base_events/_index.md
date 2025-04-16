# Windows Base Sigma Model

This model is designed for triage of dead disk, or file based live
analysis. The rules that use this model will be evaluated once on
all events.

After all relevant rules are evaluated, the collection is complete.

# Log Sources

Following is a list of recognized log sources.


## `*/windows/application`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: application
```


## `*/windows/applocker`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: applocker
```


## `*/windows/appmodel-runtime`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: appmodel-runtime
```


## `*/windows/appxdeployment-server`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: appxdeployment-server
```


## `*/windows/appxpackaging-om`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: appxpackaging-om
```


## `*/windows/bits-client`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: bits-client
```


## `*/windows/capi2`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: capi2
```


## `*/windows/certificateservicesclient-lifecycle-system`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: certificateservicesclient-lifecycle-system
```


## `*/windows/codeintegrity-operational`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: codeintegrity-operational
```


## `*/windows/diagnosis-scripted`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: diagnosis-scripted
```


## `*/windows/dns-client`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: dns-client
```


## `*/windows/dns-server`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: dns-server
```


## `*/windows/dns-server-analytic`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: dns-server-analytic
```


## `*/windows/driver-framework`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: driver-framework
```


## `*/windows/firewall-as`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: firewall-as
```


## `*/windows/ldap_debug`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: ldap_debug
```


## `*/windows/lsa-server`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: lsa-server
```


## `*/windows/microsoft-servicebus-client`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: microsoft-servicebus-client
```


## `*/windows/msexchange-management`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: msexchange-management
```


## `*/windows/ntlm`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: ntlm
```


## `*/windows/openssh`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: openssh
```


## `*/windows/powershell`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: powershell
```


## `*/windows/powershell-classic`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: powershell-classic
```


## `*/windows/security`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: security
```


## `*/windows/security-mitigations`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: security-mitigations
```


## `*/windows/shell-core`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: shell-core
```


## `*/windows/smbclient-security`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: smbclient-security
```


## `*/windows/sysmon`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: sysmon
```


## `*/windows/system`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: system
```


## `*/windows/taskscheduler`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: taskscheduler
```


## `*/windows/terminalservices-localsessionmanager`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: terminalservices-localsessionmanager
```


## `*/windows/vhdmp`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: vhdmp
```


## `*/windows/windefend`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: windefend
```


## `*/windows/wmi`





Sample use in a sigma rule:
```yaml
logsource:
  product: windows
  service: wmi
```


## `process_creation/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: process_creation
  product: windows
```


## `ps_classic_provider_start/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: ps_classic_provider_start
  product: windows
```


## `ps_classic_start/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: ps_classic_start
  product: windows
```


## `ps_module/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: ps_module
  product: windows
```


## `ps_script/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: ps_script
  product: windows
```


## `registry_add/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: registry_add
  product: windows
```


## `registry_event/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: registry_event
  product: windows
```


## `registry_set/windows/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: registry_set
  product: windows
```


## `antivirus/windows/windefend`





Sample use in a sigma rule:
```yaml
logsource:
  category: antivirus
  product: windows
  service: windefend
```


