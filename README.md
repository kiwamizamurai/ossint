# Ossint

This tool is for performing OSINT on GitHub OSS.

## Installation

```bash
brew tap kiwamizamurai/tap
brew install kiwamizamurai/tap/ossint
```

## Usage

Set up authentication using one of the following methods:

1. Set the `GITHUB_TOKEN` environment variable:

```bash
export GITHUB_TOKEN=YOUR_GITHUB_PAT
```

2. Install and authenticate with the GitHub CLI (`gh`):

```bash
gh auth token
```

If the `GITHUB_TOKEN` environment variable is not set, the tool will attempt to use the GitHub CLI to obtain an authentication token.

3. run the command

```bash
ossint --username kiwamizamurai
                               PR                               │ Stars │                            Title                            │ Additions │ Deletions │ Changed Files │ File Extensions 
────────────────────────────────────────────────────────────────┼───────┼─────────────────────────────────────────────────────────────┼───────────┼───────────┼───────────────┼─────────────────
 https://github.com/singerdmx/flutter-quill/pull/1739           │ 2548  │ feat: move cursor after inserting video/image               │ 4         │ 2         │ 1             │ .dart           
 https://github.com/singerdmx/flutter-quill/pull/1421           │ 2548  │ update document: build arguments                            │ 2         │ 0         │ 2             │ .md             
 https://github.com/Hiflylabs/awesome-dbt/pull/36               │ 1137  │ add sample for dbt-athena                                   │ 1         │ 0         │ 1             │ .md             
 https://github.com/fal-ai/dbt-fal/pull/851                     │ 851   │ docs: fix example file path                                 │ 1         │ 1         │ 1             │ .md             
 https://github.com/GoogleCloudPlatform/magic-modules/pull/7932 │ 803   │ Add field s3 path to google_storage_transfer_job            │ 9         │ 0         │ 2             │ .go, .markdown  
 https://github.com/echo724/notion2md/pull/49                   │ 677   │ fix: constructor in README.md                               │ 2         │ 2         │ 1             │ .md             
 https://github.com/unpluggedcoder/awesome-rust-tools/pull/9    │ 355   │ remove duplicate item                                       │ 0         │ 1         │ 1             │ .md             
 https://github.com/dbt-labs/dbt-bigquery/pull/1139             │ 216   │ fix: alter table description without full_refresh           │ 9         │ 0         │ 2             │ .yaml, .sql     
 https://github.com/dbt-labs/docs.getdbt.com/pull/5368          │ 118   │ [review and merge in aug 2024] remove fal related resources │ 11        │ 153       │ 7             │ .md, .js, .json 
```

## Releases

Prebuilt binaries for Linux, macOS, and Windows are available in the [Releases](https://github.com/kiwamizamurai/ossint/releases) section.