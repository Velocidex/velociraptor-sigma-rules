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
import zipfile


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
    return re.sub(r'base64decode\(string=".+"',
          r'base64decode(string="' + encoded + '"',
          artifact_data)

if __name__ == "__main__":
    argument_parser = argparse.ArgumentParser()
    argument_parser.add_argument("artifact", help="The artifact file to modify")
    argument_parser.add_argument("rule_dir", help="Path to the top level rules directory")
    argument_parser.add_argument(
        "--zip", type=str,
        help="Write the artifact in a zipfile")

    argument_parser.add_argument(
        "--name", type=str,
        help="Set a new namefor the new artifact")

    args = argument_parser.parse_args()
    artifact = build_artifact(args.artifact, args.rule_dir)

    if args.name:
        artifact = re.sub("^(name: ).+", "\\1" + args.name, artifact)

    if args.zip:
        artifact_name = "artifact.yaml"
        if args.name:
            artifact_name = args.name + ".yaml"

        with zipfile.ZipFile(args.zip, 'w', compression=zipfile.ZIP_DEFLATED) as z:
            z.writestr(artifact_name, artifact)

        print("Added %s to zip file %s" % (artifact_name, args.zip))

    else:
        print(artifact)

`
