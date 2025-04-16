# Linux Base eBPF Model

This model is designed to follow eBPF events on Linux.

Events are provided by various eBPF functions in the `watch_ebpf()`
plugin. These events can provide security critical information which
can be detected using Sigma Rules.

# Log Sources

Following is a list of recognized log sources.


## `network_connection/linux/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: network_connection
  product: linux
```


## `file_event/linux/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: file_event
  product: linux
```


## `process_creation/linux/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: process_creation
  product: linux
```


## `ebpf/linux/*`

Reports events from the ebpf subsystem.

NOTE: Events are enriched using the process tracker. You
probably want to also collect the `Linux.Events.TrackProcesses`
monitoring artifact.



#### Sample Events


##### security_file_open: A file is opened
<pre class="json-renderer">
{"Timestamp":"2025-02-12T17:11:39+10:00","EventType":"security_file_open","System":{"Timestamp":"2025-02-12T17:11:39.322765542+10:00","EventID":732,"EventName":"security_file_open","ThreadStartTime":"2025-02-12T17:11:00.93207523+10:00","ProcessorID":3,"ProcessID":1495852,"ThreadID":1495852,"ParentProcessID":2920,"HostProcessID":1495852,"HostThreadID":1495852,"HostParentProcessID":2920,"UserID":1000,"MountNS":4026531841,"ProcessName":"python","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"pathname":"/etc/passwd","flags":32768,"dev":265289728,"inode":525345,"ctime":1717126653049616364,"syscall_pathname":"/etc/passwd","ProcInfo":{"Name":"python","CommandLine":"python","CreateTime":"2025-02-12T17:11:00+10:00","Exe":"/usr/bin/python3.10","Cwd":"/home/mic/projects/velociraptor/gui/velociraptor","Username":"mic"}}}

</pre>


##### bpf_attach: A program is loading a new eBPF program into the kernel.
<pre class="json-renderer">
{"Timestamp":"2025-02-12T17:20:36+10:00","EventType":"bpf_attach","System":{"Timestamp":"2025-02-12T17:20:36.494004508+10:00","EventID":770,"EventName":"bpf_attach","ThreadStartTime":"2025-02-12T17:20:18.731905267+10:00","ProcessorID":0,"ProcessID":1502227,"ThreadID":1502237,"ParentProcessID":1497772,"HostProcessID":1502227,"HostThreadID":1502237,"HostParentProcessID":1497772,"UserID":0,"MountNS":4026531841,"ProcessName":"test","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"prog_type":17,"prog_name":"tracepoint__raw","prog_id":638,"prog_helpers":[0,0,0,0],"symbol_name":"sys_enter","symbol_addr":0,"attach_type":0,"ProcInfo":{"Name":"test","CommandLine":"./test dump fstat","CreateTime":"2025-02-12T17:20:17+10:00","Exe":"/home/mic/projects/tracee_velociraptor/test","Cwd":"/home/mic/projects/tracee_velociraptor","Username":"root"}}}

</pre>


##### kill: Kill another process
<pre class="json-renderer">
{"Timestamp":"2025-02-12T17:34:20+10:00","EventType":"kill","System":{"Timestamp":"2025-02-12T17:34:20.195398538+10:00","EventID":62,"EventName":"kill","ThreadStartTime":"2025-02-12T17:25:32.909343797+10:00","ProcessorID":2,"ProcessID":1505752,"ThreadID":1505752,"ParentProcessID":1497772,"HostProcessID":1505752,"HostThreadID":1505752,"HostParentProcessID":1497772,"UserID":0,"MountNS":4026531841,"ProcessName":"python","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"pid":1511451,"sig":15,"ProcInfo":{"Name":"python","CommandLine":"python","CreateTime":"2025-02-12T17:25:32+10:00","Exe":"/usr/bin/python3.10","Cwd":"/home/mic/projects/tracee_velociraptor","Username":"root"}}}

</pre>


##### module_load: A module is loaded into the kernel
<pre class="json-renderer">
{"Timestamp":"2025-02-13T01:19:47+10:00","EventType":"module_load","System":{"Timestamp":"2025-02-13T01:19:47.949336796+10:00","EventID":783,"EventName":"module_load","ThreadStartTime":"2025-02-13T01:19:47.936944684+10:00","ProcessorID":11,"ProcessID":1637035,"ThreadID":1637035,"ParentProcessID":1497772,"HostProcessID":1637035,"HostThreadID":1637035,"HostParentProcessID":1497772,"UserID":0,"MountNS":4026531841,"ProcessName":"modprobe","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"name":"raid0","version":"","src_version":"F80BCFC5A37DB69C992A616","pathname":"/usr/lib/modules/5.15.0-131-generic/kernel/drivers/md/raid0.ko","dev":265289728,"inode":2927026,"ctime":1738101387402437264,"ProcInfo":{"Name":"modprobe","Username":0,"Exe":"/usr/sbin/modprobe","CommandLine":"modprobe raid0","CreateTime":"2025-02-13T01:19:47.939353257+10:00"}}}

</pre>


##### mount: A filesystem is mounted
<pre class="json-renderer">
{"Timestamp":"2025-02-12T21:15:32+10:00","EventType":"mount","System":{"Timestamp":"2025-02-12T18:21:23.131524017+10:00","EventID":165,"EventName":"mount","ThreadStartTime":"2025-02-12T18:21:23.121883153+10:00","ProcessorID":11,"ProcessID":1544183,"ThreadID":1544183,"ParentProcessID":1544169,"HostProcessID":1544183,"HostThreadID":1544183,"HostParentProcessID":1544169,"UserID":0,"MountNS":4026531841,"ProcessName":"mount.cifs","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"source":"//192.168.1.153/Shared","target":".","filesystemtype":"cifs","mountflags":0,"data":94258787270320,"ProcInfo":{"Name":"mount.cifs","CommandLine":"/sbin/mount.cifs //192.168.1.153/Shared /home/mic/windows -o rw,credentials=/home/mic/.credentials,uid=1000,gid=1000","CreateTime":"2025-02-12T21:15:25+10:00","Exe":"/usr/sbin/mount.cifs","Cwd":"/home/mic/windows","Username":"root"}}}

</pre>


##### sched_process_exec: A process starts
<pre class="json-renderer">
{"Timestamp":"2025-02-12T21:20:34+10:00","EventType":"sched_process_exec","System":{"Timestamp":"2025-02-12T18:26:31.694485948+10:00","EventID":715,"EventName":"sched_process_exec","ThreadStartTime":"2025-02-12T18:26:31.687074227+10:00","ProcessorID":0,"ProcessID":1547749,"ThreadID":1547749,"ParentProcessID":1506235,"HostProcessID":1547749,"HostThreadID":1547749,"HostParentProcessID":1506235,"UserID":1000,"MountNS":4026531841,"ProcessName":"ssh","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"cmdpath":"/usr/bin/ssh","pathname":"/usr/bin/ssh","dev":265289728,"inode":2627189,"ctime":1719890760862661949,"inode_mode":33261,"interpreter_pathname":"","interpreter_dev":0,"interpreter_inode":0,"interpreter_ctime":0,"argv":["ssh","localhost"],"interp":"/usr/bin/ssh","stdin_type":8192,"stdin_path":"/dev/pts/13","invoked_from_kernel":0,"prev_comm":"bash","env":[],"ProcInfo":{"CreateTime":null}}}

</pre>


##### security_socket_connect: A process is making an outbound connection
<pre class="json-renderer">
{"Timestamp":"2025-02-13T00:41:08+10:00","EventType":"security_socket_connect","System":{"Timestamp":"2025-02-12T19:51:17.61437999+10:00","EventID":736,"EventName":"security_socket_connect","ThreadStartTime":"2025-02-12T19:51:17.600313258+10:00","ProcessorID":15,"ProcessID":1606982,"ThreadID":1606982,"ParentProcessID":2569,"HostProcessID":1606982,"HostThreadID":1606982,"HostParentProcessID":2569,"UserID":1000,"MountNS":4026531841,"ProcessName":"ssh","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"sockfd":3,"type":1,"remote_addr":{"sa_family":"AF_INET","sin_addr":"192.168.1.1","sin_port":"22"},"ProcInfo":{"Name":"ssh","CommandLine":"ssh hostname","CreateTime":"2025-02-13T00:41:06+10:00","Exe":"/usr/bin/ssh","Cwd":"/home/mic/","Username":"mic"}}}

</pre>


##### security_socket_bind: A process is binding to a socket
<pre class="json-renderer">
{"Timestamp":"2025-02-13T00:34:48+10:00","EventType":"security_socket_bind","System":{"Timestamp":"2025-02-12T19:44:58.527067089+10:00","EventID":738,"EventName":"security_socket_bind","ThreadStartTime":"2025-02-12T19:44:58.521605416+10:00","ProcessorID":8,"ProcessID":1602622,"ThreadID":1602622,"ParentProcessID":2569,"HostProcessID":1602622,"HostThreadID":1602622,"HostParentProcessID":2569,"UserID":1000,"MountNS":4026531841,"ProcessName":"nc","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"sockfd":3,"local_addr":{"sa_family":"AF_INET","sin_addr":"0.0.0.0","sin_port":"6666"},"ProcInfo":{"Name":"nc","CommandLine":"nc -l -p 6666","CreateTime":"2025-02-13T00:34:47+10:00","Exe":"/usr/bin/nc.openbsd","Cwd":"/home/mic/projects/velociraptor-sigma-rules","Username":"mic"}}}

</pre>


##### security_inode_unlink: A file is deleted
<pre class="json-renderer">
{"Timestamp":"2025-02-13T01:05:47+10:00","EventType":"security_inode_unlink","System":{"Timestamp":"2025-02-13T01:05:47.24461195+10:00","EventID":733,"EventName":"security_inode_unlink","ThreadStartTime":"2025-02-13T01:05:47.239887186+10:00","ProcessorID":14,"ProcessID":1627429,"ThreadID":1627429,"ParentProcessID":1497772,"HostProcessID":1627429,"HostThreadID":1627429,"HostParentProcessID":1497772,"UserID":0,"MountNS":4026531841,"ProcessName":"rm","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"pathname":"/tmp/owned.sh","inode":2142763,"dev":265289728,"ctime":1739372746203242108,"ProcInfo":{"Name":"rm","Username":0,"Exe":"/usr/bin/rm","CommandLine":"rm /tmp/owned.sh","CreateTime":"2025-02-13T01:05:47.241443048+10:00"}}}

</pre>


##### chown: A file is changing ownership
<pre class="json-renderer">
{"Timestamp":"2025-02-13T01:10:41+10:00","EventType":"fchownat","System":{"Timestamp":"2025-02-13T01:10:41.543170814+10:00","EventID":260,"EventName":"fchownat","ThreadStartTime":"2025-02-13T01:10:41.536500565+10:00","ProcessorID":13,"ProcessID":1630757,"ThreadID":1630757,"ParentProcessID":1497772,"HostProcessID":1630757,"HostThreadID":1630757,"HostParentProcessID":1497772,"UserID":0,"MountNS":4026531841,"ProcessName":"chown","HostName":"devbox","CgroupID":6344,"MainHostname":"devbox"},"EventData":{"dirfd":-100,"pathname":"/tmp/owned.sh","owner":1000,"group":1000,"flags":0,"ProcInfo":{"Name":"chown","Username":0,"Exe":"/usr/bin/chown","CommandLine":"chown mic:mic /tmp/owned.sh","CreateTime":"2025-02-13T01:10:41.538574768+10:00"}}}

</pre>




Sample use in a sigma rule:
```yaml
logsource:
  category: ebpf
  product: linux
```


