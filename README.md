# go-git-release (prototype)

go-git-release is an opinionated release tool designed to work with the Monorepos we use at Terminal. It leverages `git-cliff`
to create release notes off of specific project directories within a monorepo. This is currently a prototype project and is not
ready for use in real environments.

## Prerequisites

* Git
* Git Cliff

## Install 

### Locally

```
go install ./cmd/go-git-release
```

## Usage

```
NAME:
   go-git-release - Simple opinionated release tooling for monorepos.

USAGE:
   go-git-release [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value      Specifies the name of the path to create release notes. (default: ".") [$RELEASE_PATH]
   --notes value    Specifies the name of the file to export release notes. (default: "RELEASE_NOTES.md") [$RELEASE_NOTES_PATH]
   --project value  Specifies the name of the project for release notes and release commits. (e.g. project-name) [$RELEASE_PROJECT]
   --tag value      Specifies the name of the tag for release notes and release commits. (e.g. v2024.12.01) [$RELEASE_TAG]
   --dry-run        Prints the commands that would be executed without running them. (default: false) [$RELEASE_DRY_RUN]
   --help, -h       show help
```

## Example

### Dry Run

1. Navigate to the repository root.

```
cd /path/to/repository/root
```

2. Use dry-run to review the expected changes.

```
go-git-release --dir /path/to/project/root --project project-name --tag v2024.12.01 --dry-run
```

3. If satisfied, run the actual command removing `--dry-run`

```
go-git-release --dir /path/to/project/root --project project-name --tag v2024.12.01
```

4. Push your branch and tag

```
git push release/project-name-v2024.12.01
git push origin project-name-v2024.12.01
```