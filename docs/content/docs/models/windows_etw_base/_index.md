# Windows Base ETW Model

This model is designed to follow ETW event sources.

ETW events are provided by various ETW Providers in the Windows
Kernel. These events can provide security critical information which
can be detected using Sigma Rules.

# Log Sources

Following is a list of recognized log sources.


## `etw/windows/kernel`

Events from the `NT Kernel Logger` provider

The `NT Kernel Logger` ETW source is a special purpose ETW
provider that reports details about network/registry and file.

This provider enriches events with process information from the
process tracker.

This provider is special: Enabling this provider implicitly
triggers many other ETW providers such as File, Process,
Registry and Network monitoring. Velociraptor's ETW subsystem
recognizes the `Kernel Logger` provider automatically and
performs additional processing:

- Resolves full files paths from kernel space (uses device
  notation) to regular filesystem paths (e.g. `C:\Windows`).

- Collects rundown events to determine the initial system
  state. This allows Velociraptor to resolve file and registry
  paths from events that refer to kernel object addresses.

For these reasons it is preferable to use this provider over
the `Microsoft-Windows-Kernel-File` or
`Microsoft-Windows-Kernel-Registry` providers.



#### Sample Events


##### WriteFile
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:13:08Z","EventType":"WriteFile","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"WriteFile","EventID":{"Value":0}},"EventData":{"Offset":"450608","IrpPtr":"0xFFFF8203632DFB48","FileObject":"0xFFFF82036AFDE270","FileKey":"0xFFFFB189591D4700","TTID":"11240", "IoSize":"318","IoFlags":"395776","FileName":"C:\\datastore\\clients\\C.34365d02e4e1aa77\\monitoring_logs\\Windows.ETW.KernelFile\\2025-01-30.json","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### ReadFile
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:13:08Z","EventType":"ReadFile","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"ReadFile","EventID":{"Value":0}},"EventData":{"Offset":"0","IrpPtr":"0xFFFF82036378AB48","FileObject":"0xFFFF82036AFDE270","FileKey":"0xFFFFB189591D4700","TTID":"11240","IoSize":"2","IoFlags":"0","FileName":"C:\\datastore\\clients\\C.34365d02e4e1aa77\\monitoring_logs\\Windows.ETW.KernelFile\\2025-01-30.json","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### CloseFile
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:13:08Z","EventType":"CloseFile","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"CloseFile","EventID":{"Value":0}},"EventData":{"IrpPtr":"0xFFFF82036378AB48","FileObject":"0xFFFF82035B7263B0","FileKey":"0xFFFFB189591D4700","TTID":"11240","FileName":"C:\\datastore\\clients\\C.34365d02e4e1aa77\\monitoring_logs\\Windows.ETW.KernelFile\\2025-01-30.json","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### ReleaseFile
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:21Z","EventType":"ReleaseFile","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"ReleaseFile","EventID":{"Value":0}},"EventData":{"IrpPtr":"0xFFFF8203597030F8","FileObject":"0xFFFF820378327250","FileKey":"0xFFFFB1893BC871B0","TTID":"4928","FileName":"C:\\datastore\\1\\VelociraptorClient_info.log.202501270000","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### CreateFile
<pre class="json-renderer">
{
    "Timestamp": "2025-02-10T14:53:28Z",
    "EventType": "CreateFile",
    "System": {
        "Channel": "NT-Kernel-Logger",
        "Computer": "Hostname",
        "EventType": "CreateFile",
        "EventID": {
            "Value": 0
        }
    },
    "EventData": {
        "IrpPtr": "0xFFFFBB0CCAF9D0F8",
        "FileObject": "0xFFFFBB0CCED8D220",
        "TTID": "6772",
        "CreateOptions": "21119008",
        "FileAttributes": "0",
        "ShareAccess": "7",
        "OpenPath": "C:\\Windows\\System32\\psapi.dll",
        "ProcInfo": {
            "Pid": 3896,
            "Ppid": 1160,
            "Name": "MsMpEng.exe",
            "Threads": 86,
            "Username": "NT AUTHORITY\\SYSTEM",
            "OwnerSid": "S-1-5-18",
            "CommandLine": "\"C:\\ProgramData\\Microsoft\\Windows Defender\\Platform\\4.18.24090.11-0\\MsMpEng.exe\"",
            "Exe": "C:\\ProgramData\\Microsoft\\Windows Defender\\Platform\\4.18.24090.11-0\\MsMpEng.exe",
            "TokenIsElevated": true,
            "CreateTime": "2025-02-09T12:57:23.9023563Z",
            "User": 343.9949132,
            "System": 2.9596658,
            "IoCounters": {
                "ReadOperationCount": 267063,
                "WriteOperationCount": 78945,
                "OtherOperationCount": 3306593,
                "ReadTransferCount": 7729991705,
                "WriteTransferCount": 793027416,
                "OtherTransferCount": 843561371
            },
            "Memory": {
                "PageFaultCount": 108663785,
                "PeakWorkingSetSize": 998694912,
                "WorkingSetSize": 209084416,
                "QuotaPeakPagedPoolUsage": 1543600,
                "QuotaPagedPoolUsage": 698424,
                "QuotaPeakNonPagedPoolUsage": 535768,
                "QuotaNonPagedPoolUsage": 254496,
                "PagefileUsage": 331485184,
                "PeakPagefileUsage": 1064894464
            },
            "PebBaseAddress": 513193275392,
            "IsWow64": false
        }
    }
}

</pre>


##### RegQueryValue
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:21Z","EventType":"RegQueryValue","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"RegQueryValue","EventID":{"Value":0}},"EventData":{"InitialTime":"4597740950432","Status":"0","Index":"1","KeyHandle":"0xFFFFB189314B9200","KeyName":"StandardName","RegistryPath":"\\REGISTRY\\MACHINE\\SYSTEM\\ControlSet001\\Control\\TimeZoneInformation\\StandardName","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### RegOpenKey
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:28Z","EventType":"RegOpenKey","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"RegOpenKey","EventID":{"Value":0}},"EventData":{"InitialTime":"4597800728147","Status":"0","Index":"0","KeyHandle":"0xFFFFB1892A2E3050","KeyName":"SOFTWARE\\Microsoft\\Ole\\Extensions","RegistryPath":"\\REGISTRY\\MACHINE\\SOFTWARE\\Microsoft\\Ole\\Extensions","ProcInfo":{"Name":"chrome.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe","CommandLine":"\"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe\" "}}}

</pre>


##### RegCloseKey
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:28Z","EventType":"RegCloseKey","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"RegCloseKey","EventID":{"Value":0}},"EventData":{"InitialTime":"4597800729177","Status":"0","Index":"0","KeyHandle":"0xFFFFB18933C408B0","KeyName":"","RegistryPath":"\\REGISTRY\\MACHINE\\SOFTWARE\\Microsoft\\Ole\\Extensions","ProcInfo":{"Name":"chrome.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe","CommandLine":"\"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe\" "}}}

</pre>


##### RegCreateKey
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:36Z","EventType":"RegCreateKey","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"RegCreateKey","EventID":{"Value":0}},"EventData":{"InitialTime":"4597827001762","Status":"0","Index":"0","KeyHandle":"0xFFFFB189314C06D0","KeyName":"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections","RegistryPath":"\\REGISTRY\\USER\\S-1-5-21-241402409-3571345782-2557608070-500\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings\\Connections","ProcInfo":{"Name":"chrome.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe","CommandLine":"\"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe\" "}}}

</pre>


##### SendTCPv4
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:22Z","EventType":"SendTCPv4","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"SendTCPv4","EventID":{"Value":0}},"EventData":{"PID":"1596","size":"1652","daddr":"192.168.1.5","saddr":"192.168.1.237","dport":"37890","sport":"3389","startime":"45977585","endtime":"45977585","seqnum":"0","connid":null,"ProcInfo":null}}

</pre>


##### RecvTCPv4
<pre class="json-renderer">
{"Timestamp":"2025-01-30T14:06:22Z","EventType":"RecvTCPv4","System":{"Channel":"NT Kernel Logger","Computer":"WIN-SJE0CKQO83P","EventType":"RecvTCPv4","EventID":{"Value":0}},"EventData":{"PID":"1596","size":"74","daddr":"192.168.1.5","saddr":"192.168.1.237","dport":"37890","sport":"3389","seqnum":"0","connid":null,"ProcInfo":null}}

</pre>


##### CreateProcess
<pre class="json-renderer">
{
    "Timestamp": "2025-01-30T13:59:46Z",
    "EventType": "CreateProcess",
    "System": {
        "Channel": "NT Kernel Logger",
        "Computer": "WIN-SJE0CKQO83P",
        "EventType": "CreateProcess",
        "EventID": {
            "Value": 0
        }
    },
    "EventData": {
        "UniqueProcessKey": "0xFFFF82035F210380",
        "ProcessId": "0x1C5C",
        "ParentId": "0x524",
        "SessionId": "0",
        "ExitStatus": "259",
        "DirectoryTableBase": "0x1050BD000",
        "Flags": "0",
        "UserSID": "\\\\NT AUTHORITY\\SYSTEM",
        "ImageFileName": "WmiPrvSE.exe",
        "CommandLine": "C:\\Windows\\system32\\wbem\\wmiprvse.exe -Embedding",
        "PackageFullName": "",
        "ApplicationId": "",
        "ProcInfo": {
            "Pid": 1316,
            "Ppid": 1172,
            "Name": "svchost.exe",
            "Threads": 14,
            "Username": "NT AUTHORITY\\SYSTEM",
            "OwnerSid": "S-1-5-18",
            "CommandLine": "C:\\Windows\\system32\\svchost.exe -k DcomLaunch -p",
            "Exe": "C:\\Windows\\System32\\svchost.exe",
            "TokenIsElevated": true,
            "CreateTime": "2025-01-03T15:14:38.8669202Z",
            "User": 5.859375,
            "System": 14.65625,
            "IoCounters": {
                "ReadOperationCount": 12,
                "WriteOperationCount": 0,
                "OtherOperationCount": 161602,
                "ReadTransferCount": 49152,
                "WriteTransferCount": 0,
                "OtherTransferCount": 4860956
            },
            "Memory": {
                "PageFaultCount": 72766,
                "PeakWorkingSetSize": 27328512,
                "WorkingSetSize": 25452544,
                "QuotaPeakPagedPoolUsage": 743528,
                "QuotaPagedPoolUsage": 743288,
                "QuotaPeakNonPagedPoolUsage": 37128,
                "QuotaNonPagedPoolUsage": 21304,
                "PagefileUsage": 8261632,
                "PeakPagefileUsage": 8978432
            },
            "PebBaseAddress": 93408378880,
            "IsWow64": false
        }
    }
}

</pre>




Sample use in a sigma rule:
```yaml
log_sources:
  category: etw
  product: windows
  service: kernel
```


## `etw/windows/file`

Log source based on the `Microsoft-Windows-Kernel-File` provider.

See `etw/windows/kernel` for a better ETW provider.



#### Sample Events


##### FileOpen Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T07:33:12Z","System":{"Channel":"Microsoft-Windows-Kernel-File","Computer":"WIN-SJE0CKQO83P","EventType":"FileOpen","EventID":{"Value":12}},"EventData":{"Irp":"0xFFFF8203746EAB08","FileObject":"0xFFFF8203783249B0","IssuingThreadId":"8924","CreateOptions":"0x1200000","CreateAttributes":"0x0","ShareAccess":"0x7","FileName":"\\Device\\HarddiskVolume3\\datastore\\1\\VelociraptorClient_info.log.202501270000","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### CreateNewFile Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T07:33:07Z","System":{"Channel":"Microsoft-Windows-Kernel-File","Computer":"WIN-SJE0CKQO83P","EventType":"CreateNewFile","EventID":{"Value":30}},"EventData":{"Irp":"0xFFFF82035D7F9C88","FileObject":"0xFFFF82035D145A80","IssuingThreadId":"8096","CreateOptions":"0x5000060","CreateAttributes":"0x80","ShareAccess":"0x3","FileName":"\\Device\\HarddiskVolume3\\datastore\\hunts\\H.CUDEISQSB5K60\\notebook\\N.H.CUDEISQSB5K60\\NC.CUDIKP2J653IQ-CUDIMC0CRIMHS.json.db","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>


##### NameCreate Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T07:33:07Z","System":{"Channel":"Microsoft-Windows-Kernel-File","Computer":"WIN-SJE0CKQO83P","EventType":"NameCreate","EventID":{"Value":10}},"EventData":{"FileKey":"0xFFFFB18937D871B0","FileName":"\\Device\\HarddiskVolume3\\datastore\\hunts\\H.CUDEISQSB5K60\\notebook\\N.H.CUDEISQSB5K60\\NC.CUDIKP2J653IQ-CUDIMC0CRIMHS.json.db","ProcInfo":{"Name":"velociraptor.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\velociraptor.exe","CommandLine":"c:\\velociraptor.exe  gui --datastore c:\\datastore\\ --nobrowser --debug -v"}}}

</pre>




Sample use in a sigma rule:
```yaml
log_sources:
  category: etw
  product: windows
  service: file
```


## `etw/windows/registry`

Log source based on the `Microsoft-Windows-Kernel-Registry` provider.

See `etw/windows/kernel` for a better ETW provider.



#### Sample Events


##### CreateKey Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T08:10:20Z","System":{"Channel":"Microsoft-Windows-Kernel-Registry","Computer":"WIN-SJE0CKQO83P","EventType":"CreateKey","EventID":{"Value":1}},"EventData":{"BaseObject":"0xFFFFB1893D043AA0","KeyObject":"0x0","Status":"0x104","Disposition":"0","BaseName":"","RelativeName":"\\REGISTRY\\USER\\S-1-5-21-241402409-3571345782-2557608070-500_Classes\\Local Settings\\Software\\Microsoft\\Windows\\CurrentVersion\\AppModel\\SystemAppData\\Microsoft.Windows.StartMenuExperienceHost_cw5n1h2txyewy","ProcInfo":null,"ParentProcInfo":{"Name":"svchost.exe","Username":"NT AUTHORITY\\SYSTEM","Exe":"C:\\Windows\\System32\\svchost.exe","CommandLine":"C:\\Windows\\system32\\svchost.exe -k DcomLaunch -p"}}}

</pre>


##### DeleteKey Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T08:12:50Z","System":{"Channel":"Microsoft-Windows-Kernel-Registry","Computer":"WIN-SJE0CKQO83P","EventType":"DeleteKey","EventID":{"Value":3}},"EventData":{"KeyObject":"0xFFFFB189385B9BB0","Status":"0x0","KeyName":"","ProcInfo":null,"ParentProcInfo":{"Name":"regedit.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\Windows\\regedit.exe","CommandLine":"\"C:\\Windows\\regedit.exe\" "}}}

</pre>


##### DeleteValueKey Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T08:12:33Z","System":{"Channel":"Microsoft-Windows-Kernel-Registry","Computer":"WIN-SJE0CKQO83P","EventType":"DeleteValueKey","EventID":{"Value":6}},"EventData":{"KeyObject":"0xFFFFB1893D0460E0","Status":"0xC0000034","KeyName":"","ValueName":"CachedFeatureString","ProcInfo":null,"ParentProcInfo":{"Name":"SearchApp.exe","Username":"WIN-SJE0CKQO83P\\Administrator","Exe":"C:\\Windows\\SystemApps\\Microsoft.Windows.Search_cw5n1h2txyewy\\SearchApp.exe","CommandLine":"\"C:\\Windows\\SystemApps\\Microsoft.Windows.Search_cw5n1h2txyewy\\SearchApp.exe\" -ServerName:CortanaUI.AppX8z9r6jm96hw4bsbneegw0kyxx296wr9t.mca"}}}

</pre>


##### OpenKey Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T08:08:51Z","System":{"Channel":"Microsoft-Windows-Kernel-Registry","Computer":"WIN-SJE0CKQO83P","EventType":"OpenKey","EventID":{"Value":2}},"EventData":{"BaseObject":"0xFFFFB1892A262AD0","KeyObject":"0xFFFFB1893C344100","Status":"0x0","Disposition":"0","BaseName":"","RelativeName":"\\Registry\\Machine\\Hardware\\DeviceMap\\VIDEO","ProcInfo":null,"ParentProcInfo":{"Name":"vm3dservice.exe","Username":"NT AUTHORITY\\SYSTEM","Exe":"C:\\Windows\\System32\\vm3dservice.exe","CommandLine":"vm3dservice.exe -n"}}}

</pre>


##### SetValueKey Event
<pre class="json-renderer">
{"Timestamp":"2025-01-30T08:10:20Z","System":{"Channel":"Microsoft-Windows-Kernel-Registry","Computer":"WIN-SJE0CKQO83P","EventType":"SetValueKey","EventID":{"Value":5}},"EventData":{"KeyObject":"0xFFFFB1893D043AA0","Status":"0x0","Type":"11","DataSize":"8","KeyName":"","ValueName":"PCT","CapturedDataSize":"0","CapturedData":"","PreviousDataType":"0","PreviousDataSize":"0","PreviousDataCapturedSize":"0","PreviousData":null,"ProcInfo":null,"ParentProcInfo":{"Name":"svchost.exe","Username":"NT AUTHORITY\\SYSTEM","Exe":"C:\\Windows\\System32\\svchost.exe","CommandLine":"C:\\Windows\\system32\\svchost.exe -k DcomLaunch -p"}}}

</pre>




Sample use in a sigma rule:
```yaml
log_sources:
  category: etw
  product: windows
  service: registry
```


## `etw/windows/process`

Log source based on the `Microsoft-Windows-Kernel-Registry` provider

See `etw/windows/kernel` for a better ETW provider.




Sample use in a sigma rule:
```yaml
log_sources:
  category: etw
  product: windows
  service: process
```


