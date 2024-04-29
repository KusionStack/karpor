---
title: Release Process And Cadence
---

## Release Planning

We will establish and continuously follow up on the release plan through [GitHub Milestones](https://github.com/KusionStack/karpor/milestones). Each release milestone will include two types of tasks:

- Tasks Maintainers commit to complete. Maintainers will decide on the features they are committed to implementing before the next release based on their available time and effort. Usually, tasks are finalized after offline discussions and then added to the milestone. These tasks will be assigned to the Maintainer who plans to implement or test them.
- Additional items contributed by community contributors, typically non-urgent features or optimizations. Maintainers do not commit to completing these issues within the release cycle but will commit to reviewing submissions from the community. 

The milestones will clearly describe the most important features and their expected completion dates. This will clearly inform end-users about the timing and contents of the next release. 

In addition to the next milestone, we will also maintain drafts of future release milestones.

## Release Standards

- All **official releases** should be tagged on the `main` branch, with optional pre-release version suffixes such as: `alpha`, `beta`, `rc`, for example, a regular official release version might be `v1.2.3`, `v1.2.3-alpha.0`. For instance, if we want to perform some validations before releasing the official version `v1.2.3`, we could first release a pre-release version like `v1.2.3-alpha.0`, followed by `v1.2.3` after the validation is complete.
- Maintainers commit to completing certain features and enhancements, tracking progress through [GitHub Milestones](https://github.com/KusionStack/karpor/milestones).
- We will do our best to avoid release delays; thus, if we cannot complete a feature on time, it will be moved to the next release. 
- A new version will be released every **1 month**.

## Release Standard Procedure

Maintainers are responsible for driving the release process and following standard operating procedures to ensure the quality of the release.

1. Tag the git commit designated for release and push it upstream; the tag needs to comply with [Semantic Versioning](#semantic-versioning).
2. Ensure that the triggered Github Actions pipeline is executed successfully. Once successful, it will automatically generate a new Github Release, which includes the Changelog calculated from commit messages, as well as artifacts such as images and tar.gz files.
3. Write clear release notes based on the **Github Release**, including:
   - User-friendly release highlights.
   - Deprecated and incompatible changes.
   - Brief instructions on how to install and upgrade.

## Gate Testing
Before creating the release branch, we will have a **1-week** code freeze period. During this period, we will refrain from merging any feature PRs and will only fix bugs.

Maintainers will test and fix these last-minute issues before each release.

## Semantic Versioning

`Karpor` adopts [Semantic Versioning](https://semver.org/) for its version numbers.

The version format: `MAJOR.MINOR.PATCH`, for example, `v1.2.3`. The version number **incrementing rules** are as follows:
- MAJOR version when you make incompatible API changes.
- MINOR version when you add functionality in a backwards-compatible manner.
- PATCH version when you make backwards-compatible bug fixes. 

**Pre-release version numbers and build metadata** can be added to the `MAJOR.MINOR.PATCH` as an extension, like `v1.2.3-alpha.0`, `v1.2.3-beta.1`, `v1.2.3-rc.2`, where `-alpha.0`, `-beta.1`, `-rc.2` are pre-release versions.
