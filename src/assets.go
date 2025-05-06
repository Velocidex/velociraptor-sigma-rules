package main

const CompierPython = `#!/usr/bin/python3

"""Recompile a Velociraptor Sigma artifact with reconstituted rules.

This script will update the artifact with a set of rules from the
local filesystem. It is mostly used to modify the artifact in order to
remove a set of noisy rules from the curated set.

## Example:

python3 compile.py Windows.Hayabusa.Rules.yaml ./rules/ > Modified.yaml

"""

import argparse
import re
import os
from io import BytesIO
import gzip
import base64

def build_artifact(artifact_path, rule_path):
    with open(artifact_path) as fd:
        artifact_data = fd.read()

    # Read all the rule files and join them into a big multi doc yaml
    # file.
    compressed_rules = BytesIO()
    raw_rules = BytesIO()
    with gzip.GzipFile(fileobj=compressed_rules, mode="w", compresslevel=9) as out_fd:
        for root, dirs, files in sorted(os.walk(rule_path, topdown=False)):
            for name in sorted(files):
                if not name.endswith(".yml") and not name.endswith(".yaml"):
                    continue

                full_path = os.path.join(root, name)
                with open(full_path, "rb") as fd:
                    data = fd.read()
                    out_fd.write(b"\n---\n")
                    out_fd.write(data)
                    raw_rules.write(b"\n---\n")
                    raw_rules.write(data)

    encoded = base64.standard_b64encode(compressed_rules.getvalue()).decode("utf8")
    print(re.sub(r'base64decode\(string=".+"',
          r'base64decode(string="' + encoded + '"',
          artifact_data))

if __name__ == "__main__":
    argument_parser = argparse.ArgumentParser()
    argument_parser.add_argument("artifact", help="The artifact file to modify")
    argument_parser.add_argument("rule_dir", help="Path to the top level rules directory")

    args = argument_parser.parse_args()
    build_artifact(args.artifact, args.rule_dir)
`
