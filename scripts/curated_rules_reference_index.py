import urllib.request
import json
import html
import yaml
import re
import os
import zipfile

# Where we generate the search index.
commits_url = "https://api.github.com/repos/Velocidex/velociraptor-sigma-rules/commits"
output_data_path = "docs/static/curated_rules/data.json"
archive_path = "docs/static/curated_rules/curated_rules.zip"
artifact_root_directory = "docs/content/curated_rules/rules"
artifact_page_directory = "docs/content/curated_rules/rules/pages"

org = "Velocidex"
project = "velociraptor-sigma-rules"

# Each yaml file will be converted to a markdown if needed.
template = """---
title: %s
hidden: true
tags: %s
editURL: https://github.com/%s/%s/edit/master/%s
---

%s

<pre><code class="language-yaml">
%s
</code></pre>

"""

date_regex = re.compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}")
hash_regex = re.compile("#([0-9_a-z]+)", re.I | re.M | re.S)

# ToDo: Check for descriptions larger than 1 paragraph. If so, warn
def cleanDescription(description):
  description = description.replace("\r\n", "\n")
  top_paragraph = description.split("\n\n")[0]
  return top_paragraph

def getTags(tags):
  result = []
  for m in hash_regex.finditer(tags):
    result.append(m.group(1))

  return result

# Create a zip file with all the artifacts in it.
def make_archive(archive_path):
  with zipfile.ZipFile(archive_path, mode='w',
                       compression=zipfile.ZIP_DEFLATED) as archive:
    for root, dirs, files in os.walk(artifact_root_directory):
      for name in files:
        if not name.endswith(".yaml"):
          continue

        filename = os.path.join(root, name)
        with open(filename) as fd:
          data = fd.read()

        # Force archive member to have a fixed timestamp - therefore
        # the zip hash will depend only on the content, allowing git
        # to deduplicate it properly across commits.
        archive.writestr(zipfile.ZipInfo(filename=filename), data)

def build_markdown():
  index = []

  for root, dirs, files in os.walk(artifact_root_directory):
    files.sort()

    for name in files:
      if not name.endswith(".yaml"):
        continue

      yaml_filename = os.path.join(root, name)
      try:
        with open(yaml_filename) as stream:
          content = stream.read()
          data = yaml.safe_load(content)

          base_name = name
          filename_name = os.path.join(artifact_page_directory, base_name.lower())

          description = data.get("description", "")

          record = {
            "title": data["title"],
            "author": data.get("author"),
            "description": cleanDescription(description),
            "link": os.path.join("/curated_ruleset/rules/pages/",
                                base_name.lower()).replace("\\", "/"),
            "tags": data.get("tags", []),
          }

          index.append(record)

          md_filename = filename_name + ".md"
          with open(md_filename, "w") as fd:
            fd.write(template % (
              data["title"],
              json.dumps(record["tags"]),
              org, project,
              yaml_filename,
              data["description"],
              html.escape(content, quote=False)))
      except Exception as e:
        print("Error processing %s: %s" % (yaml_filename, e))
  index = sorted(index, key=lambda x: x["title"],
                 reverse=True)

  with open(output_data_path, "w") as fd:
    fd.write(json.dumps(index, indent=4))
    print("Writing data.json in %s" % output_data_path)

if __name__ == "__main__":
  build_markdown()
  make_archive(archive_path)

if os.getenv('CI'):
   # Remove this file so the site may be pushed correctly.
   os.remove(artifact_root_directory + "/.gitignore")
   os.remove(os.path.dirname(output_data_path) + "/.gitignore")
