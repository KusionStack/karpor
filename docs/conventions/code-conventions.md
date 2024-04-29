---
title: Code Conventions
---

In this section, you will find the code conventions for all kinds of code Karpor project related. It's not necessary to learn all of them at once, but make sure you have read corresponding parts before you start to code.

- [Go Code Conventions](#go-code-conventions)
- [Bash or Script Conventions](#bash-or-script-conventions)
- [Directory and File Conventions](#directory-and-file-conventions)
- [Logging Conventions](#logging-conventions)

## Go Code Conventions

  - [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

  - [Effective Go](https://golang.org/doc/effective_go.html)

  - Know and avoid [Go landmines](https://gist.github.com/lavalamp/4bd23295a9f32706a48f)

  - Comment your code.
    - [Go's commenting conventions](https://go.dev/blog/godoc)
    - If reviewers ask questions about why the code is the way it is, that's a
      sign that comments might be helpful.

  - Command-line flags should use dashes, not underscores

  - API
    - According to RFC3986, URLs are "case sensitive". Karpor uses `kebab-case` for API URLs.
      - e.g.: `POST /rest-api/v1/resource-group-rule`

  - Naming
    - Please consider package name when selecting an interface name, and avoid
      redundancy.

      - e.g.: `storage.Interface` is better than `storage.StorageInterface`.

    - Do not use uppercase characters, underscores, or dashes in package
      names.
    - Please consider parent directory name when choosing a package name.

      - so pkg/manager/cluster/foo.go should say `package cluster`
        not `package clustermanager`.
      - Unless there's a good reason, the `package foo` line should match
        the name of the directory in which the .go file exists.
      - Importers can use a different name if they need to disambiguate.

    - Locks should be called `lock` and should never be embedded (always `lock
      sync.Mutex`). When multiple locks are present, give each lock a distinct name
      following Go conventions - `stateLock`, `mapLock` etc.

## Bash or Script Conventions

  - https://google.github.io/styleguide/shell.xml

  - Ensure that build, release, test, and cluster-management scripts run on
    macOS

## Directory and File Conventions

- Avoid package sprawl. Find an appropriate subdirectory for new packages.
  - Libraries with no more appropriate home belong in new package
    subdirectories of pkg/util

- Avoid general utility packages. Packages called "util" are suspect. Instead,
  derive a name that describes your desired function. For example, the utility
  functions dealing with waiting for operations are in the "wait" package and
  include functionality like Poll. So the full name is wait.Poll

- All filenames should be lowercase

- Go source files and directories use underscores, not dashes
  - Package directories should generally avoid using separators as much as
    possible (when packages are multiple words, they usually should be in nested
    subdirectories).

- Document directories and filenames should use dashes rather than underscores

- Contrived examples that illustrate system features belong in
  `/docs/user-guide` or `/docs/admin`, depending on whether it is a feature primarily
  intended for users that deploy applications or cluster administrators,
  respectively. Actual application examples belong in /examples.
  - Examples should also illustrate [best practices for configuration and using the system](https://kubernetes.io/docs/concepts/configuration/overview/)

- Third-party code

  - Go code for normal third-party dependencies is managed using
    [go modules](https://github.com/golang/go/wiki/Modules)

  - Other third-party code belongs in `/third_party`
    - forked third party Go code goes in `/third_party/forked`
    - forked _golang stdlib_ code goes in `/third_party/forked/golang`

  - Third-party code must include licenses

  - This includes modified third-party code and excerpts, as well

## Linting and Formatting

To ensure consistency across the Go codebase, we require all code to pass a number of linter checks.

To run all linters, use the `lint` Makefile target:

```shell
make lint
```

The command will clean code along with some lint checks. Please remember to check in all changes after that.
