# Windows Base Sigma Model

This model is designed for triage of dead disk, or file based live
analysis. The rules that use this model will be evaluated once on
all events.

After all relevant rules are evaluated, the collection is complete.

# Log Sources

Following is a list of recognized log sources.


## `*/windows/application`

This Log Source generates events from the Application Channel, usually stored in the file `C:\Windows\System32\WinEvt\Logs\Application.evtx`

The channel stores a wide variety of system events from multiple
services.




Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: application
```


## `*/windows/applocker`

This Log Source generates combined events from the Windows `AppLocker service`. Events are usually stored in the files:
  * `C:\Windows\System32\WinEvt\Logs\Microsoft-Windows-AppLocker%4MSI and Script.evtx`
  * `C:\Windows\System32\WinEvt\Logs\Microsoft-Windows-AppLocker%4EXE and DLL.evtx`
  * `C:\Windows\System32\WinEvt\Logs\Microsoft-Windows-AppLocker%4Packaged app-Deployment.evtx`
  * `C:\Windows\System32\WinEvt\Logs\Microsoft-Windows-AppLocker%4Packaged app-Execution.evtx`




Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: applocker
```


## `*/windows/appmodel-runtime`

This Log Source generates combined events from the Windows `AppModel Runtime`. Events are usually stored in the files:
  * `C:\Windows\System32\WinEvt\Logs\Microsoft-Windows-AppModel-Runtime%4Admin.evtx`




Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: appmodel-runtime
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


## `*/windows/bits-client`

This Log Source generates combined events from the Windows Bits Client service. Events are usually stored in the files:
  * `C:\Windows\System32\WinEvt\Logs\Microsoft-Windows-Bits-Client%4Operational.evtx`

The BITS service is used to download files and it is often misused by threat actors to download malicious payloads.



#### Sample Events


##### EventID 3 - New Job Creation
<pre class="json-renderer">
{"Timestamp":"2025-01-13T13:48:20.745705604Z","System":{"Provider":{"Name":"Microsoft-Windows-Bits-Client","Guid":"EF1CC15B-46C1-414E-BB95-E76B077BD51E"},"EventID":{"Value":3},"Version":3,"Level":4,"Task":0,"Opcode":0,"Keywords":4611686018427387904,"TimeCreated":{"SystemTime":1736776100.7457056},"EventRecordID":1320,"Correlation":{},"Execution":{"ProcessID":8936,"ThreadID":9100},"Channel":"Microsoft-Windows-Bits-Client/Operational","Computer":"WIN-SJE0CKQO83P","Security":{"UserID":"S-1-5-18"}},"EventData":{"jobTitle":"Chrome Component Updater","jobId":"B73C90F1-5FA7-4445-8E49-6C40870E4502","jobOwner":"WIN-SJE0CKQO83P\\Administrator","processPath":"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe","processId":3616,"ClientProcessStartKey":1407374883553491},"Message":"The BITS service created a new job.\nTransfer job: Chrome Component Updater\nJob ID: B73C90F1-5FA7-4445-8E49-6C40870E4502\nOwner: WIN-SJE0CKQO83P\\Administrator\nProcess Path: C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe\nProcess ID: 3616\r\n"}

</pre>




Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: bits-client
```


## `*/windows/capi2`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: capi2
```


## `*/windows/certificateservicesclient-lifecycle-system`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: certificateservicesclient-lifecycle-system
```


## `*/windows/codeintegrity-operational`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: codeintegrity-operational
```


## `*/windows/diagnosis-scripted`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: diagnosis-scripted
```


## `*/windows/dns-client`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: dns-client
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


## `*/windows/driver-framework`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: driver-framework
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


## `*/windows/lsa-server`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: lsa-server
```


## `*/windows/microsoft-servicebus-client`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: microsoft-servicebus-client
```


## `*/windows/msexchange-management`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: msexchange-management
```


## `*/windows/ntlm`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: ntlm
```


## `*/windows/openssh`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: openssh
```


## `*/windows/powershell`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: powershell
```


## `*/windows/powershell-classic`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: powershell-classic
```


## `*/windows/schtasks`

Enumerates All Scheduled tasks



Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: schtasks
```


## `*/windows/security`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: security
```


## `*/windows/security-mitigations`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: security-mitigations
```


## `*/windows/services`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: services
```


## `*/windows/shell-core`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: shell-core
```


## `*/windows/smbclient-security`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: smbclient-security
```


## `*/windows/sysmon`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: sysmon
```


## `*/windows/system`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: system
```


## `*/windows/taskscheduler`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: taskscheduler
```


## `*/windows/terminalservices-localsessionmanager`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: terminalservices-localsessionmanager
```


## `*/windows/vhdmp`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: vhdmp
```


## `*/windows/windefend`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: windefend
```


## `*/windows/wmi`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: wmi
```


## `antivirus/windows/windefend`





Sample use in a sigma rule:
```yaml
log_sources:
  category: antivirus
  product: windows
  service: windefend
```


## `image_load/windows/pslist`





Sample use in a sigma rule:
```yaml
log_sources:
  category: image_load
  product: windows
  service: pslist
```


## `network_connection/windows/netstat`





Sample use in a sigma rule:
```yaml
log_sources:
  category: network_connection
  product: windows
  service: netstat
```


## `persistence/windows/services`





Sample use in a sigma rule:
```yaml
log_sources:
  category: persistence
  product: windows
  service: services
```


## `process_creation/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: process_creation
  product: windows
  service: *
```


## `process_creation/windows/execution`





Sample use in a sigma rule:
```yaml
log_sources:
  category: process_creation
  product: windows
  service: execution
```


## `process_creation/windows/pslist`





Sample use in a sigma rule:
```yaml
log_sources:
  category: process_creation
  product: windows
  service: pslist
```


## `ps_classic_provider_start/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_classic_provider_start
  product: windows
  service: *
```


## `ps_classic_start/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_classic_start
  product: windows
  service: *
```


## `ps_module/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_module
  product: windows
  service: *
```


## `ps_script/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: ps_script
  product: windows
  service: *
```


## `registry_add/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: registry_add
  product: windows
  service: *
```


## `registry_event/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: registry_event
  product: windows
  service: *
```


## `registry_set/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: registry_set
  product: windows
  service: *
```


## `vql/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: vql
  product: windows
  service: *
```


## `webserver/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: webserver
  product: windows
  service: *
```
