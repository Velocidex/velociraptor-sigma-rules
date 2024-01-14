import urllib.request
import argparse
import pathlib
import json
import html
import yaml
import re
import os
import zipfile

parser = argparse.ArgumentParser(description='Build an index of sigma rules.')
parser.add_argument(
  'path', type=str, nargs="+",
  help='Path to directory containing sigma rules.')

parser.add_argument(
  'output_path',  type=str,
  help='Path to write the output index file on')

parser.add_argument(
  '--root',  type=str, default='.', dest='root',
  help='Path to the root directory of the repo')

parser.add_argument(
  '--github_link',  default='https://github.com/Yamato-Security/hayabusa-rules/blob/main/',
  help="Link to Github for root of the repo")

date_regex = re.compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}")

def cleanDescription(description):
  description = description or ""
  description = description.replace("\r\n", "\n")
  top_paragraph = description.split("\n\n")[0]
  return top_paragraph

def build_index(rule_root_directories, repo_root='.',
                github_link='https://github.com/Velocidex/velociraptor-sigma-rules/blob/master/'):
  index = []

  for rule_root_directory in rule_root_directories:
    for root, dirs, files in os.walk(rule_root_directory):
      files.sort()

      for name in sorted(files):
        if not name.endswith(".yaml") and not name.endswith(".yml"):
          continue

        yaml_filename = os.path.join(root, name)
        print(yaml_filename)

        try:
          with open(yaml_filename) as stream:
            content = stream.read()
            data = yaml.safe_load(content)

            base_name = name
            description = data.get("description", "")
            posix_path = pathlib.PurePosixPath(yaml_filename)
            record = {
              "title": data["title"],
              "author": data.get("author"),
              "description": cleanDescription(description),
              "link": github_link+str(posix_path.relative_to(repo_root)),
              "tags": data.get("tags") or [],
            }

            index.append(record)
        except Exception as e:
          print("Error processing %s: %s" % (yaml_filename, e))

  index = sorted(index, key=lambda x: x["title"], reverse=True)

  return index


if __name__ == "__main__":
  args = parser.parse_args()
  index = build_index(args.path, repo_root=args.root, github_link=args.github_link)

  with open(args.output_path, "w") as fd:
    fd.write(json.dumps(index, sort_keys=True))
    print("Writing index data in %s" % args.output_path)
