#!/usr/bin/python3

"""Refresh the lolrmm rules.

RMM tools allow attackers to remotely access a system. There are a
large number of different such tools, some are legitimate.

The LolRMM project keeps a database of indicators we can use to detect
such tools. We build several sets of rules around this information
depending on the detection method.

1. Detecting executable names in prefetch output.
2. Detecting executable names in the filesystem. (This is very slow).
3. Detecting executable name in currently running processes.

"""

import json
import urllib.request

process_template = '''
title: RMM %s (Pslist)
description: |
%s

logsource:
  category: process_creation
  product: windows
  service: pslist

detection:
  condition: selection
  selection:
%s

details: Detected RMM process %%EventData.Image%%
enrichment: |
  x=>dict(
   CallChain=GetCallChain(Pid=x.EventData.Pid),
   ImageHash=x.EventData.Image && hash(path=x.EventData.Image),
   ParentHash=x.EventData.ParentImage && hash(path=x.EventData.ParentImage))

'''


file_template = '''
title: RMM %s (filesystem)
description: |
%s

logsource:
  category: filesystem
  product: windows
  service: glob

detection:
  condition: selection
  selection:
%s

details: Detected RMM executable %%EventData.OSPath%%
enrichment: |
  x=>dict(
   Hash=hash(path=x.EventData.OSPath))

'''


def get_lolrmm(url):
    with urllib.request.urlopen(url) as response:
        return response.read()


def indent(text):
    result = []
    for line in text.splitlines():
        result.append("   " + line)
    return "\n".join(result)

raw_data = get_lolrmm('https://lolrmm.io/api/rmm_tools.json')
data = json.loads(raw_data)

rules_idx = {}

rules = []

def build_pslist_selection(paths):
    endswith = []
    globs = []
    for path in paths:
        if "*" in path:
            globs.append("     - " + path.replace("*", ".*").replace("\\", "\\\\"))
        else:
            endswith.append("     - \\" + path)


    result = []
    if endswith:
        result.append("   - Image|endswith:")
        result += endswith

    if globs:
        result.append("   - Image|re:")
        result += globs

    return "\n".join(result)


def build_filesystem_selection(paths):
    endswith = []
    globs = []
    for path in paths:
        if "*" in path:
            globs.append("     - " + path.replace("*", ".*").replace("\\", "\\\\"))
        else:
            endswith.append("     - " + path)


    result = []
    if endswith:
        result.append("   - Name:")
        result += endswith

    if globs:
        result.append("   - OSPath|re:")
        result += globs

    return "\n".join(result)


for entry in data:
    lines = []
    if not entry["Details"] or not entry["Details"]["InstallationPaths"]:
        continue

    name = entry["Name"]
    if name + "P" in rules_idx:
        continue

    rules_idx[name + "P"] = 1

    rules.append(process_template % (
        name, indent(entry["Description"]),
        build_pslist_selection(entry["Details"]["InstallationPaths"])))


for entry in data:
    lines = []
    if not entry["Details"] or not entry["Details"]["InstallationPaths"]:
        continue

    name = entry["Name"]
    if name + "Files" in rules_idx:
        continue

    rules_idx[name + "Files"] = 1

    rules.append(file_template % (
        name, indent(entry["Description"]),
        build_filesystem_selection(entry["Details"]["InstallationPaths"])))

print("\n---\n".join(rules))
