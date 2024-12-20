---
title: Code Contribution Guide
---

In this code contribution guide, you will learn about the following:

- [How to run Karpor locally](#running-karpor-locally)
- [How to create a pull request](#creating-a-pull-request)
- [Code review guidelines](#code-review)
- [Formatting guidelines for pull requests](#formatting-guidelines)
- [Updating Documentation and Website](#updating-documentation-and-website)

## Running Karpor Locally

This guide will help you get started with Karpor development.

### Prerequisites

* Golang version 1.22+

<details>
  <summary>Installing Golang</summary>

1. Install go1.22 from the [official website](https://go.dev/dl/). Extract the binary files and place them at a location, assuming it is located under the home directory `~/go/`, here is an example command, you should choose the correct binary file for your system.
  ```
  wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
  tar xzf go1.20.2.linux-amd64.tar.gz
  ```

If you would like to maintain multiple versions of golang in your local development environment, you can download the package and extract it to a location, like `~/go/go1.22.1`, and then alter the path in the command below accordingly.

1. Set environment variables for Golang

  ```
  export PATH=~/go/bin/:$PATH
  export GOROOT=~/go/
  export GOPATH=~/gopath/
  ```

If the `gopath` folder does not exist, create it with `mkdir ~/gopath`. These commands will add the go binary folder to the `PATH` environment variable (making it the primary choice for go) and set the `GOROOT` environment to this go folder. Please add these lines to your `~/.bashrc` or `~/.zshrc` file, so you won't need to set these environment variables every time you open a new terminal.

1. (Optional) Some regions, such as China, may have slow connection to the default go registry; you can configure GOPROXY to speed up the download process.
  ```
  go env -w GOPROXY=https://goproxy.cn,direct
  ```

</details>

* Kubernetes version v1.20+ configured with `~/.kube/config`.

* golangci-lint version v1.52.2+, it will be installed automatically if you run `make lint`, if the installation fails, you can install it manually.

<details>
  <summary>Manually installing golangci-lint</summary>

You can install it manually following the [guide](https://golangci-lint.run/welcome/install), or use the command:

```
cd ~/go/ && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2
```

</details>

### Building

- Clone this project

```shell
git clone git@github.com:KusionStack/karpor.git
```

- Build locally

Executing `make build-all` will build the executables for all platforms; if you only want to build for a specific platform, execute `make build-${PlatformName}`, e.g., `make build-darwin`. To see all available commands, execute `make help`.

### Testing

It's essential to write tests to maintain code quality, you can run all unit tests by executing the following command in the project root directory:
```shell
make test
```

If you need to generate extra coverage report files, execute:
```shell
make cover
```

Then you can view the content of the coverage report in a browser by running:
```shell
make cover-html
```

## Creating a Pull Request

We are thrilled that you are considering contributing to the Karpor project!

This document will guide you through the process of creating a [pull request](https://help.github.com/en/articles/about-pull-requests/).


### Before you begin

We know you are excited to create your first pull request. Before we get started, make sure your code follows the relevant [code conventions](../conventions/code-conventions.md).

### Your First Pull Request

Before submitting your PR, run the following commands to ensure they all succeed:
```
make test
make lint
```

If this is your first time contributing to an open-source project on GitHub, please make sure to read the instructions on [creating a pull request](https://help.github.com/en/articles/creating-a-pull-request).

To increase the chances of your pull request being accepted, please ensure your pull request follows these guidelines:

- The title and description match the implementation.
- The commits in the pull request follow the [formatting guidelines](#Formatting-guidelines).
- The pull request closes a related issue.
- The pull request includes necessary tests to verify the expected behavior.
- If your pull request has conflicts, please rebase your branch onto the main branch.

If the pull request fixes a bug:

- The pull request description must contain `Closes #<issue number>` or `Fixes #<issue number>`.
- To prevent regressions, the pull request should include tests that replicate the bug being fixed.

## Code Review
Once you have created a pull request, the next step is to have others review your changes. Review is a learning opportunity for both reviewers and the author of the pull request.

If you believe a specific person should review your pull request, you can tag them in the description or a comment.
Tag a user by typing an `@` symbol followed by their username.

## Formatting Guidelines

A well-crafted pull request can minimize the time to get your changes accepted. These guidelines will help you write well-formulated commit messages and descriptions for your pull requests.

### Commit Message Format

Karpor adheres to [conventional-commits](https://www.conventionalcommits.org/en/v1.0.0/).

Commit messages should be organized following this structure:

```
<type>[optional scope]: <subject>

[optional body]
```

#### Example

Commit message with scope:

```
feat(lang): add polish language
```

Commit message without body:

```
docs: correct spelling of CHANGELOG
```

Commit message with multiple body paragraphs:：

```
fix: correct minor typos in code

see the issue for details

on typos fixed.

reviewed-by: Z
refs #133
```

#### `<type>`(Required)
The required type helps better capture the area of the commit, based on the [Angular guidelines](https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md#-commit-message-guidelines).

We use lowercase for `<type>` to avoid spending time on case-sensitive issues. `<type>` can be one of the following:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **build**: Changes that affect the build system or external dependencies
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation

#### `<scope>`(Optional)

Scope is optional and can be provided to the type of commit to provide additional contextual information, enclosed in parentheses. It can be anything specifying the place of the commit change. Github issue links are also valid scopes e.g., fix(ui), feat(api), fix(#233), etc.

When the change affects multiple scopes, `*` can be used.

#### `<subject>`(Required)

The subject must come immediately after the type/scope prefix, followed by a colon and space. It is a concise summary of the code changes, for example, "fix: array parsing issue when multiple spaces were contained in string", rather than "fix: bug".

#### `<body>`(Required)

A longer commit body can be provided after the brief subject, giving additional context information about the code change. The body must begin one line after the description.

### Pull Request Title

When accepting pull requests, the Karpor team merges all commits into one.

The pull request title becomes the subject line of the merged commit message.

We still encourage contributors to write informative commit messages, as they will be part of the Git commit body.

We use the pull request titles when generating change logs for releases. Hence, we strive to make the titles as informative as possible.

Make sure your pull request title uses the same format as the commit message subject line. If the format is not followed, we will add a `title-needs-formatting` label on the pull request.

### Passing All CI Checks
Before merging, all testing CIs should pass:
- Coverage should not drop. Currently, the pull request coverage should be at least 70%.
- Karpor uses a **CLA** for the contributor agreement. It requires you to sign for every commit before merging the pull request.

## Updating Documentation and Website

If your pull request has been merged, and it is a new feature or enhancement, you need to update the documentation and send a pull request to the [kusionstack.io](https://github.com/KusionStack/kusionstack.io) repository.

Learn how to write documentation through the following guide:

- [kusionstack.io Developer Guide](https://github.com/KusionStack/kusionstack.io/blob/main/README.md)

Awesome, you've completed the lifecycle of code contribution!
