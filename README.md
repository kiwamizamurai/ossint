# Ossint

This tool is for performing OSINT on GitHub OSS.

## Installation

```bash
go get github.com/kiwamizamurai/ossint/cmd/ossint
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

## Releases

Prebuilt binaries for Linux, macOS, and Windows are available in the [Releases](https://github.com/kiwamizamurai/ossint/releases) section.