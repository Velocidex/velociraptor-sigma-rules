---
title: "Velociraptor rule set"
date: 2023-10-16T15:05:44+10:00
---

# Velociraptor Ruleset

These are hand picked rules to help with the initial triaging of a machine. These rules get the most out of Velociraptor, by querying a large variety of
forensic artifacts, rather than just scanning event logs. These rules do not include standard informational rules, instead only including high confidence
rules to reduce noise for incident responders. 

{{< ruleset "/index/velociraptor_index.json" >}}
