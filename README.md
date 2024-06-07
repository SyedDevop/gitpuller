# GetPuller

A CLI Tool for Simplified GitHub Content Management.

## Installation:

To install `gitpuller`, run the following command:

```bash
go install github.com/SyedDevop/gitpuller
```

## Features:

- [ ] View files
- [ ] Search Repos
  - [ ] select branch

## Configuration:

  * (Required) Set the `GIT_TOKEN` environment variable to a env for access to the GitHub API.
  * (Optional) Create a config file at `$HOME/.config/gitpuller.yml` and set a keys, like:

```
email: example@email.com
userName:  SyedDevop
token: ******************** // replace this with tour git token
```
