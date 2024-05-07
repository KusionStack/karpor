<div align="center">
<p></p><p></p>
<p>
    <img width="50%" src="assets/img/logo-full.svg">
</p>
<h2>Cross-Cluster Discovery 🔍 ,</br>Limitless Insight 📊.<span style="color: gray; font-weight: normal;"> With AI.<sup style="color: gray; font-weight: normal; font-size: 10px"> coming soon</sup></span>✨</h2>

[👉 Live Demo](https://karpor-demo.kusionstack.io) | [简体中文](https://github.com/KusionStack/karpor/blob/main/README-zh.md) | [English](https://github.com/KusionStack/karpor/blob/main/README.md)

[![Kusion](https://github.com/KusionStack/kusion/actions/workflows/release.yaml/badge.svg)](https://github.com/KusionStack/kusion/actions/workflows/release.yaml)
[![GitHub release](https://img.shields.io/github/release/KusionStack/kusion.svg)](https://github.com/KusionStack/kusion/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/KusionStack/kusion)](https://goreportcard.com/report/github.com/KusionStack/kusion)
[![Coverage Status](https://coveralls.io/repos/github/KusionStack/kusion/badge.svg)](https://coveralls.io/github/KusionStack/kusion)
[![Go Reference](https://pkg.go.dev/badge/github.com/KusionStack/kusion.svg)](https://pkg.go.dev/github.com/KusionStack/kusion)
[![license](https://img.shields.io/github/license/KusionStack/kusion.svg)](https://github.com/KusionStack/kusion/blob/main/LICENSE)

<!-- TODO: Uncomment when the repository is publicly. -->

<!-- [![Karpor](https://github.com/KarporStack/karpor/actions/workflows/release.yaml/badge.svg)](https://github.com/KarporStack/karpor/actions/workflows/release.yaml) -->

<!-- [![GitHub release](https://img.shields.io/github/release/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/releases) -->

<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/KusionStack/karpor)](https://goreportcard.com/report/github.com/KusionStack/karpor) -->

<!-- [![Coverage Status](https://coveralls.io/repos/github/KusionStack/karpor/badge.svg)](https://coveralls.io/github/KusionStack/karpor) -->

<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/KusionStack/karpor.svg)](https://pkg.go.dev/github.com/KusionStack/karpor) -->

<!-- [![license](https://img.shields.io/github/license/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/blob/main/LICENSE) -->

</div>

## What is Karpor?

Karpor is a **Kubernetes Explorer** focusing on **🔍 Search**, **💡 Insight** and **🤖 Intelligence**. It has features such as non-invasive, read-only, secure, and multi cloud and multi cluster support, and can serve as a **Kubernetes Data Plane** to reduce the cost of discovering and understanding kubernetes resources.

## Why Karpor?

- ⚡️ **Lightweight and Easy to Setup**. One `helm` is done.
- 📦 **Out of the Box**. Built-in security and compliance policies, resource sync strategy, resource transform rule, and topology relationship definitions.
- 💰 **Self-Hosted, Cost-Effective**. Bring your own server, scale when you need.
- 🔒 **Read-Only** data plane, **Non-Invasive** to user cluster.
- ⚙️ **Kubernetes Native API Compatible**. Seamless integration of existing kubernetes tool chain.
- ☁️ Natural support for **Multi-Cluster and Multi-Cloud**.

## ⚙️ Installation

### Install with Helm

[Helm](https://github.com/helm/helm) is a tool for managing packages of pre-configured kubernetes resources.

```bash
$ helm repo add kusionstack https://kusionstack.github.io/charts
$ helm repo update kusionstack
$ helm install karpor kusionstack/karpor
```

For more information about installation, please check the [Installation Guide](https://kusionstack.io/karpor/getting-started/installation) on official website.

## 📖 Documentation

Detailed documentation is available at [Karpor Website](https://kusionstack.io/karpor).

## 🤝 Contribution Guide

Karpor is still in the initial stage, and there are many capabilities that need to be made up, so we welcome everyone to participate in construction with us. Visit the [Contribution Guide](CONTRIBUTING.md) to understand how to participate in the contribution Karpor project. If you have any questions, please [Submit the Issue](https://github.com/KusionStack/karpor/issues).

## 🎖︎ Contributors

<a href="https://github.com/KusionStack/karpor/graphs/contributors">
<img src="https://contrib.rocks/image?repo=KusionStack/kusion" />
</a>

## 🌐 Contact Us

- Twitter: [KusionStack](https://twitter.com/KusionStack)
- Slack: [KusionStack](https://join.slack.com/t/karpor/shared_invite/zt-19lqcc3a9-_kTNwagaT5qwBE~my5Lnxg)
- DingTalk (Chinese): 42753001
- Wechat Group (Chinese)

  <img src="./assets/img/wechat.png" width="200" height="200"/>

<!-- ## 🏛️ License -->

<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Felliotxx%2Fkusion.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Felliotxx%2Fkusion?ref=badge_shield&issueType=license) -->
