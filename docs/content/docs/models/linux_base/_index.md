# Linux Base Sigma Model

This model is designed for triage of dead disk, or file based live
analysis. The rules that use this model will be evaluated once on
all events.

After all relevant rules are evaluated, the collection is complete.

# Log Sources

Following is a list of recognized log sources.


## `*/linux/*`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
```


## `*/linux/sshd`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
  service: sshd
```


## `*/linux/cron`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
  service: cron
```


## `*/linux/auth`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
  service: auth
```


## `*/linux/syslog`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
  service: syslog
```


## `*/linux/sudo`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
  service: sudo
```


## `*/linux/auditd`





Sample use in a sigma rule:
```yaml
logsource:
  product: linux
  service: auditd
```


## `network_connection/linux/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: network_connection
  product: linux
```


## `process_creation/linux/*`





Sample use in a sigma rule:
```yaml
logsource:
  category: process_creation
  product: linux
```


