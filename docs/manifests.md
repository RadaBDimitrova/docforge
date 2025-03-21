# Documentation Manifests

Manifests are yaml configuration files that describe the structure of a documentation bundle and the rules how to construct it. This document specifies the options for designing a manifest. All manifests start with `structure:`. All related code can be found in [the manifest package](../pkg/manifest)

## Structural elements

### File element

Manifest: files.yaml
```yaml
structure:
# file defined by github URL
- file: https://github.com/gardener/docforge/blob/master/docs/manifests.md
# rename README.md to overview.md. Note that if file extension is not specified it is assumed to be .md
- file: overview
  source: https://github.com/gardener/docforge/blob/master/docs/README.md
# defining a file that is the concatenation of multiple files
- file: combined
  multiSource:
  - https://github.com/gardener/docforge/blob/master/docs/cmd-ref/docforge.md
  - https://github.com/gardener/docforge/blob/master/docs/cmd-ref/docforge_version.md
# define a section file with no content and only frontmatter properties
- file: _index.md
  frontmatter:
    title: Section A
    description: Description of Section A

```
Result:
```
docforge-docs
├── manifests.md
├── overview.md
├── combined.md
└── _index_.md
```


### Directory element

Manifest: dir.yaml
```yaml
structure:
- dir: section
  structure:
  # the content of the directory is listed under structure the same way as the manifest
  - file: https://github.com/gardener/docforge/blob/master/docs/cmd-ref/docforge.md
- file: overview
  source: https://github.com/gardener/docforge/blob/master/docs/README.md
```
Result:
```
docforge-docs
├── section
|   └── docforge.md
└── overview.md
```

### FileTree element
Manifest: fileTree.yaml
```yaml
structure:
- dir: section
  structure:
  # loads all files from a given github directory (tree)
  - fileTree: https://github.com/gardener/docforge/tree/master/docs
    # relative paths of files to exclude
    excludeFiles:
    - cmd-ref/docforge.md
    - cmd-ref/docforge_version.md
    - cmd-ref/docforge_completion.md
    - cmd-ref/docforge_gen-cmd-docs.md
```
Result:
```
docforge-docs
└── section
    |── consistency.md
    |── manifest-ref.md
    |── manifests.md
    └── user-index.md
```

### Manifest element
Manifest: manifestElement.yaml
```yaml
structure:
- manifest: dir.yaml
# strucures get merged
- manifest: fileTree.yaml
```
Result:
```
docforge-docs
├── section
|   |── docforge.md
|   |── consistency.md
|   |── manifest-ref.md
|   |── manifests.md
|   └── user-index.md
└── overview.md
```

## Relative manifest links

If path starts with a `/` its considered from the repo root. Else its considered from the manifest position.

Manifest: docs/links.yaml

```yaml
structure:
# resolves to https://github.com/gardener/docforge/blob/master/README.md
- file: ../README.md
# resolves to https://github.com/gardener/docforge/tree/master/docs
- fileTree: /docs
```
## Frontmatter

Every node in the structural tree can define frontmatter. Dirs propagate their frontmatter to their children where children override frontmatter values if there is a collision

Manifest: frontmatter.yaml
```yaml
structure:
- dir: section
  # this will be propagated to apiserver.md 
  frontmatter:
    weight: 1
    tag: dev
  structure:
  - file: https://github.com/gardener/docforge/blob/master/docs/manifests.md
    frontmatter:
      title: Guides
      description: Walkthroughs of common activities
      # this will override the frontmatter property from section
      weight: 5
- file: overview
  source: https://github.com/gardener/docforge/blob/master/docs/README.md
```
Result: section/apiserver.md
``` yaml
---
title: Guides
description: Walkthroughs of common activities
weight: 5
tag: dev
---
```