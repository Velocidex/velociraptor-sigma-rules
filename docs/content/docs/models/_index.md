---
title: "Sigma Models"
date: 2023-10-15T00:14:44+10:00
weight: 100
IconClass: fa-solid fa-desktop
---

# Sigma Models

The Sigma standard specifies how to write detection rules in terms of
abstract `Log Sources`. The standard itself does not specify what log
sources are available and what kind of events these actually emit.

Therefore it is not enough to consider just a Sigma rule in
isolation - Sigma rules are always written inside a `Sigma Model`.

The `Model` is a specification of what `Log Sources` are available on
the system and what kind of events (i.e. what fields are present in
each field).

A Sigma model is the combination of `Log Sources` (which provide
events into the Sigma Rule) and `Field Mappings` which allow those
fields to be referenced in the `Sigma Rule`.

## What are models used for?

Depending on the use case different `Log Sources` are defined - for
example in a file focused forensic context (e.g. Forensic file Triage
or dead disk image), Log sources extract event data from files such as
event log files, or other forensic artifacts. Applying Sigma rules on
a forensic triage can surface interesting events quickly to guide
investigations.

On the other hand when applying the Sigma Rules to live events, `Log
Sources` might include live data such as eBPF events or ETW sourced
events not usually present in a dead disk image. Sigma rules in such a
live setting can be used to detect and monitor anomalous conditions in
real time.

In the above two examples, the Sigma rules must be written in context
of the model in use.
