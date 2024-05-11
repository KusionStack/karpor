<div align="center">
<p></p><p></p>
<p>
    <img width="60%" src="https://kusionstack-io-dev.vercel.app/karpor/assets/logo/logo-full.svg">
</p>

<h2 style="font-size: 1.5em;">
  Intuitive Discovery 🔍 ,<br>
  Limitless Insight 📊.<span style="color: gray; font-weight: normal;"> With AI.<sup style="color: gray; font-weight: normal; font-size: 0.5em"> coming soon</sup></span>✨
</h2>

<p align="center">
  <a href="https://karpor-demo.kusionstack.io" target="_blank"><b>🎮 Live Demo</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/installation" target="_blank"><b>⚙️ Install</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/quick-start" target="_blank"><b>⚡️ Quick Start</b></a> •
  <a href="https://kusionstack.io/karpor" target="_blank"><b>📚 Docs</b></a>
</p>


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

Karpor is a Modern **Kubernetes Explorer** focusing on **🔍 Search**, **💡 Insight** and **🤖 Intelligence**. It has features such as `non-invasive`, `read-only`, `secure`, and `multi-cloud` and `multi-cluster` support, and can serve as a **Kubernetes Data Plane** to reduce the cost of discovering and understanding kubernetes resources.

https://github.com/KusionStack/karpor/assets/9360247/c5050dfa-23f3-49ac-ba4a-1026ab043e6c

## Why Karpor?

<h3 align="center">🤝 User Friendly</h3>

<table>
  <tr>
    <td style="vertical-align: middle;">
      <b>⚡️ Lightweight and Easy to Setup</b><br />One `helm` is done.<br /><br/>
      <b>📦 Out of the Box</b><br />Built-in security and compliance policies, resource sync strategy, resource transform rule, and topology relationship definitions.<br /><br />
      <b>🔍 Quickly search and locate resource(s)</b><br />Quickly search and locate resource(s) of interest across a large number of clusters in a user-friendly way.
    </td>
    <td style="width: 60%; vertical-align: middle;">
      <img src="https://kusionstack-io-dev.vercel.app/karpor/assets/overview/user-friendly.png" alt="User Friendly" />
    </td>
  </tr>
</table>

<h3 align="center">✨ Intelligent</h3>

<table>
  <tr>
    <td style="width: 60%; vertical-align: middle;">
      <img src="https://kusionstack-io-dev.vercel.app/karpor/assets/overview/intelligent.png" alt="Intelligent" />
    </td>
    <td style="vertical-align: middle;">
      <b>🔒 Compliance Protection</b><br />Automatically identify potential risks and receive AI suggestions for remediation.<br /><br/>
      <b>📊 Logical and topological views</b><br />Logical and topological views of relevant resources within their operational context.
    </td>
  </tr>
</table>

<h3 align="center">⚡️ Low Burden</h3>

<table>
  <tr>
    <td style="vertical-align: middle;">
      <b>🔒 Read-Only, Non-Invasive</b><br />Read-Only data plane, Non-Invasive to user cluster.<br /><br />
      <b>⚙️ Kubernetes Native API Compatible</b><br />Seamless integration of existing kubernetes tool chain.<br /><br />
      <b>☁️ Multi-Cluster and Multi-Cloud/Hybrid-Cloud</b><br />Natively supports Multi-Cluster and Multi-Cloud/Hybrid-Cloud.
    </td>
    <td style="width: 60%; vertical-align: middle;">
      <img src="https://kusionstack-io-dev.vercel.app/karpor/assets/overview/low-burden.png" alt="Low Burden" />
    </td>
  </tr>
</table>

</br>

## ⚙️ Installation

### Install with Helm

[Helm](https://github.com/helm/helm) is a tool for managing packages of pre-configured kubernetes resources.

```bash
$ helm repo add kusionstack https://kusionstack.github.io/charts
$ helm repo update
$ helm install karpor kusionstack/karpor
```

For more information about installation, please check the [Installation Guide](https://kusionstack.io/karpor/getting-started/installation) on official website.

## 📖 Documentation

Detailed documentation is available at [Karpor Website](https://kusionstack.io/karpor).

## 🤝 How to contribute

Karpor is still in the initial stage, and there are many capabilities that need to be made up, so we welcome everyone to participate in construction with us. Visit the [Contribution Guide](CONTRIBUTING.md) to understand how to participate in the contribution Karpor project. If you have any questions, please [Submit the Issue](https://github.com/KusionStack/karpor/issues).

## 🎖︎ Contributors

Thanks all! Come and join us! 🍻

<a href="https://github.com/KusionStack/karpor/graphs/contributors">
<img src="https://contrib.rocks/image?repo=KusionStack/kusion" />
</a>

<!-- ## 🏛️ License -->

<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Felliotxx%2Fkusion.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Felliotxx%2Fkusion?ref=badge_shield&issueType=license) -->
