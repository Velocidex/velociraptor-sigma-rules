---
title: "Velociraptor Hayabusa Live Detection"
date: 2023-10-15T00:14:44+10:00
weight: 10
---

# Velociraptor Hayabusa Live Detection

These rules are automatically imported from the [Hayabusa
project](https://github.com/Yamato-Security/hayabusa) and
[SigmaHQ](https://github.com/SigmaHQ/sigma).

Rules are automatically compiled using the [configuation
file](https://github.com/Velocidex/velociraptor-sigma-rules/blob/master/config/velociraptor_windows_event_monitoring.yaml)
into an artifact pack.

To download the latest version of the artifact pack [click
here](https://sigma.velocidex.com/Velociraptor-Hayabusa-Monitoring.zip)

These rules are designed to work in live mode. This means the rules
will match events on the endpoint directly as the event logs are
written to the system event logs. Velociraptor will therefore forward
only matching detections rather than all events.

This artifact relies on the `watch_evtx()` plugin. This plugin will
follow the windows event logs (similar to `tail -f` on Linux) and
parse events periodically.

There are some differences between this approach as compared to real
event log sources (such as `ETW`):

1. The log file following approach amortizes CPU load over time, if
   the check period is not too frequent CPU load can be lower than
   ETW.
2. The `watch_evtx()` approach does not rely on Windows ETW and
   therefore does not impact limited ETW session resources.
3. ETW based monitoring is more real time, which will send detections
   to the server sooner.


### Search the Hayabusa ruleset

{{< ruleset "/index/hayabusa_index.json" >}}
