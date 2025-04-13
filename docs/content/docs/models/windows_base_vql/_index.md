# Windows Base VQL Sigma Model

This model is designed for triage of dead disk, or file based live
analysis using VQL rules. The rules that use this model will be
evaluated once on all events.

After all relevant rules are evaluated, the collection is complete.

Rules that utilize this model may include a `vql` section which may
contain a VQL lambda to dictates how the event is generated. This
allows the rule itself to generate all relevant fields.

For example:

```yaml
vql: |
x=>dict(
  Timestamp=timestamp(epoch=now()),
  EventData=dict(
    Files=SearchFiles(Glob='C:/Users/*/AppData/Roaming/rclone/rclone.conf')
  ))
```

The following utility functions are defined:

* `SearchFiles(Glob)`: Allows searching for files with a
  glob. Returns the file size as well as the first 100 bytes.

* `SearchRegistryKeys(Glob)`: Allows searching for registry keys -
  returns a dict with key/value pairs from the registry.


# Log Sources

Following is a list of recognized log sources.


## `vql/windows/*`

This log source emits a single event. All rules using the log
source will receive this event, where they can run arbitrary VQL
queries to build the event themselves.

This is most useful for rules that want to generate their own
event data.




Sample use in a sigma rule:
```yaml
log_sources:
  category: vql
  product: windows
```


## `filesystem/windows/glob`

This log source searches for all files on the drive - it takes a
long time but allows rules to check for presence of a particular
filename.




Sample use in a sigma rule:
```yaml
log_sources:
  category: filesystem
  product: windows
  service: glob
```


## `*/windows/schtasks`

Enumerates All Scheduled tasks



Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: schtasks
```


## `*/windows/services`





Sample use in a sigma rule:
```yaml
log_sources:
  product: windows
  service: services
```


## `persistence/windows/services`





Sample use in a sigma rule:
```yaml
log_sources:
  category: persistence
  product: windows
  service: services
```


## `process_creation/vql/execution`





Sample use in a sigma rule:
```yaml
log_sources:
  category: process_creation
  product: vql
  service: execution
```


## `webserver/windows/*`





Sample use in a sigma rule:
```yaml
log_sources:
  category: webserver
  product: windows
```


## `process_creation/windows/pslist`





Sample use in a sigma rule:
```yaml
log_sources:
  category: process_creation
  product: windows
  service: pslist
```


## `image_load/vql/pslist`





Sample use in a sigma rule:
```yaml
log_sources:
  category: image_load
  product: vql
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


