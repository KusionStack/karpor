# Code Contribution Guide

You will learn the following things in the code contribution guide:

- [How to Run Karbour Locally](#run-karbour-locally)
- [How to Create a pull request](#create-a-pull-request)
- [Code Review Guide](#code-review)
- [Formatting guidelines of pull request](#formatting-guidelines)

## Run Karbour Locally

This guide helps you get started developing Karbour.

### Prerequisites

* Golang version 1.19+

<details>
  <summary>Install Golang</summary>

1. Install go1.19 from [official site](https://go.dev/dl/). Unpack the binary and place it somewhere, assume it's in the home path `~/go/`, below is an example command, you should choose the right binary according to your system.
  ```
  wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
  tar xzf go1.20.2.linux-amd64.tar.gz
  ```

If you want to keep multiple golang version in your local develop environment, you can download the package and unfold it into some place, like `~/go/go1.19.1`, then the following commands should also change according to the path.

1. Set environment variables for Golang

  ```
  export PATH=~/go/bin/:$PATH
  export GOROOT=~/go/
  export GOPATH=~/gopath/
  ```

  Create a gopath folder if not exist `mkdir ~/gopath`. These commands will add the go binary folder to the `PATH` environment (let it to be the primary choice for go), and set the `GOROOT` environment to this go folder. Please add these lines to your `~/.bashrc` or `~/.zshrc` file, so that you don't need to set these environment variables every time you open a new terminal.

1. (Optional) Some area like China may be too slow to connect to the default go registry, you can configure GOPROXY to speed up the download process. 
  ```
  go env -w GOPROXY=https://goproxy.cn,direct
  ```


</details>


* golangci-lint 1.49.0+, it will install automatically if you run `make`, you can install it manually if the installation broken.

<details>
  <summary>Install golangci-lint manually</summary>

You can install it manually follow [the guide](https://golangci-lint.run/usage/install/#local-installation) or the following command:

```
cd ~/go/ && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2
```

</details>

### Build

- Clone this project

```shell script
git clone git@github.com:KusionStack/karbour.git
```

For local development, we probably need to build both of them.

- Build Karbour CLI

```shell script
make build-all
```

After the karbour cli built successfully, `make build-all` command will create `karbour` binary to `_build/${OS}` under the project.

- Configure `karbour` binary to System PATH

```shell script
export PATH=$PATH:/your/path/to/project/karbour/${OS}
```

Then you can use `karbour` command directly.

### Testing

It's necessary to write tests for good code quality, please refer to [the principle of test](./principle-of-test) before you start.

#### Unit test

```shell script
make cover
```

## Create a pull request

We're excited that you're considering making a contribution to the Karbour project!
This document guides you through the process of creating a [pull request](https://help.github.com/en/articles/about-pull-requests/).

### Before you begin

We know you're excited to create your first pull request. Before we get started, make sure your code follows the relevant [code conventions](./code-conventions).

### Your first pull request

Before you submit a PR, run this command to ensure it is ready:
```
make reviewable
```

If this is your first time contributing to an open-source project on GitHub, make sure you read about [Creating a pull request](https://help.github.com/en/articles/creating-a-pull-request).

To increase the chance of having your pull request accepted, make sure your pull request follows these guidelines:

- Title and description matches the implementation.
- Commits within the pull request follow the [Formatting guidelines](#Formatting-guidelines).
- The pull request closes one related issue.
- The pull request contains necessary tests that verify the intended behavior.
- If your pull request has conflicts, rebase your branch onto the main branch.

If the pull request fixes a bug:

- The pull request description must include `Closes #<issue number>` or `Fixes #<issue number>`.
- To avoid regressions, the pull request should include tests that replicate the fixed bug.
- Generally, we will maintain the last 2 releases for bugfix. You should add `backport release-x.x` label or comment `/backport release-x.y` for the releases contained the bug, github bot will automatically backport this PR to the specified release branch after PR merged. If there're any conflicts, you should cherry-pick it manually.

## Code review

Once you've created a pull request, the next step is to have someone review your change.
A review is a learning opportunity for both the reviewer and the author of the pull request.

If you think a specific person needs to review your pull request, then you can tag them in the description or in a comment.
Tag a user by typing the `@` symbol followed by their GitHub username.

We recommend that you read [How to do a code review](https://google.github.io/eng-practices/review/reviewer/) to learn more about code reviews.

## Formatting guidelines

A well-written pull request minimizes the time to get your change accepted.
These guidelines help you write good commit messages and descriptions for your pull requests.

### Commit message format

Karbour follows the [conventional-commits](https://www.conventionalcommits.org/en/v1.0.0/) to improve better history information.

The commit message should be structured as follows:

```
<type>[optional scope]: <subject>

[optional body]
```

#### Examples:

Commit message with scope:

```
feat(lang): add polish language
```

Commit message with no body:

```
docs: correct spelling of CHANGELOG
```

Commit message with multi-paragraph body:

```
fix: correct minor typos in code

see the issue for details

on typos fixed.

Reviewed-by: Z
Refs #133
```

#### `<type>` (required)

Type is required to better capture the area of the commit, based on the [Angular convention](https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md#-commit-message-guidelines).

`<type>` can be one of the following:

* **feat**: A new feature
* **fix**: A bug fix
* **docs**: Documentation only changes
* **build**: Changes that affect the build system or external dependencies 
* **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **perf**: A code change that improves performance
* **test**: Adding missing or correcting existing tests
* **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation
* **ci**: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs, GithubAction)

#### `<scope>` (optional)

Scope is optional, it may be provided to a commit’s type, to provide additional contextual information and is contained within parenthesis, it is could be anything specifying place of the commit change. Github issue link is
also a valid scope. For example: fix(cli), feat(api), fix(#233), etc.

You can use `*` when the change affects more than a single scope.

#### `<subject>` (required)

The subject MUST immediately follow the colon and space after the type/scope prefix. The description is a short summary of the code changes, e.g., "fix: array parsing issue when multiple spaces were contained in string", instead of "fix: bug".

#### `<body>` (optional)

A longer commit body may be provided after the short subject, providing additional contextual information about the code changes. The body MUST begin one blank line after the description.

### Pull request titles

The Karbour team _squashes_ all commits into one when we accept a pull request.
The title of the pull request becomes the subject line of the squashed commit message.
We still encourage contributors to write informative commit messages, as they become a part of the Git commit body.

We use the pull request title when we generate change logs for releases. As such, we strive to make the title as informative as possible.

Make sure that the title for your pull request uses the same format as the subject line in the commit message. If the format is not followed, we will add a label `title-needs-formatting` on the pull request.

### Pass all the CI checks

Before merge, All test CI should pass green.
- [Unit test](https://github.com/KusionStack/karbour/blob/main/.github/workflows/check.yaml#L12)
- [Golang Lint](https://github.com/KusionStack/karbour/blob/main/.github/workflows/check.yaml#L37)
- [Commit Lint](https://github.com/KusionStack/karbour/blob/main/.github/workflows/check.yaml#L55)
- [CLA](https://github.com/KusionStack/karbour/blob/main/.github/workflows/cla.yaml)

## Update the docs & website

If your pull request merged and this is a new feature or enhancement, it's necessary to update the docs and send a pull request to [karbour.com](https://github.com/KusionStack/karbour.com) repo.

Learn how to write the docs by the following guide:

* [karbour.com Developer Guide](https://github.com/KusionStack/karbour.com/blob/main/README.md)

Great, you have complete the lifecycle of code contribution, try to [join the community as a member](https://github.com/KusionStack/community/blob/main/ROLES.md) if you're interested.
