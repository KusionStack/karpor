<div align="center">
<p></p><p></p>
<p>
    <img width="60%" src="https://kusionstack.io/karpor/assets/logo/logo-full.svg">
</p>

<h1 style="font-size: 1.5em;">
    Intelligence for Kubernetes
</h1>

<p align="center">
  <a href="https://karpor-demo.kusionstack.io" target="_blank"><b>🎮 Live Demo</b></a> •
  <a href="https://kusionstack.io/karpor/" target="_blank"><b>🌐 Website</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/quick-start" target="_blank"><b>⚡️ Quick Start</b></a> •
  <a href="https://kusionstack.io/karpor" target="_blank"><b>📚 Docs</b></a> •
  <a href="https://github.com/KusionStack/karpor/discussions" target="_blank"><b>💬 Discussions</b></a><br>
  [English] 
  <a href="https://github.com/KusionStack/karpor/blob/main/README-zh.md" target="_blank">[中文]</a> 
  <a href="https://github.com/KusionStack/karpor/blob/main/README-pt.md" target="_blank">[Português]</a>
</p>

[![Karpor](https://github.com/KusionStack/karpor/actions/workflows/release.yaml/badge.svg)](https://github.com/KusionStack/karpor/actions/workflows/release.yaml)
[![GitHub release](https://img.shields.io/github/release/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/KusionStack/karpor)](https://goreportcard.com/report/github.com/KusionStack/karpor)
[![Coverage Status](https://coveralls.io/repos/github/KusionStack/karpor/badge.svg)](https://coveralls.io/github/KusionStack/karpor)
[![Go Reference](https://pkg.go.dev/badge/github.com/KusionStack/karpor.svg)](https://pkg.go.dev/github.com/KusionStack/karpor)
[![license](https://img.shields.io/github/license/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/blob/main/LICENSE)

</div>

## What is Karpor?

Karpor is Intelligence for Kubernetes. It brings advanced **Search**, **Insight** and **AI** to Kubernetes. It is essentially a **Kubernetes Visualization Tool**. With Karpor, you can gain crucial visibility into your Kubernetes clusters across any clouds.

We hope to become a **small and beautiful, vendor-neutral, developer-friendly, community-driven** open-source project!

**Current Status**: We are iterating [v0.5.0 Milestone](https://github.com/KusionStack/karpor/milestone/4), welcome to join the [discussion](https://github.com/KusionStack/karpor/discussions/528).

https://github.com/KusionStack/karpor/assets/49401013/7cf31cc0-7123-42f6-8543-5addcbf4975c

## Why Karpor?

<h3 align="center">🔍 Search</h3>

<table>
  <tr>
    <td>
      <b>Automatic Syncing</b><br />Automatically synchronize your resources across any clusters managed by the multi-cloud platform.<br /><br/>
      <b>Powerful, flexible queries</b><br />Effectively retrieve and locate resources across multi clusters that you are looking for in a quick and easy way.
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/search/search-auto-complete-raw.jpg" /></kbd>
    </td>
  </tr>
</table>

<h3 align="center">💡 Insight</h3>

<table>
  <tr>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/insight/insight-home-raw.jpg" /></kbd>
    </td>
    <td>
      <b>Compliance Governance</b><br />Understand your compliance status across multiple clusters and compliance standards.<br /><br/>
      <b>Resource Topology</b><br />Logical and topological views of relevant resources within their operational context.<br /><br/>
      <b>Cost Optimization</b><br />Coming soon.
    </td>
  </tr>
</table>

<h3 align="center">✨ AI</h3>

<table>
  <tr>
    <td>
      <b>Natural Language Operations</b><br />Interact with Kubernetes using plain language for more intuitive operations.<br /><br />
      <b>Contextual AI Responses</b><br />Get smart, contextual assistance that understands your needs.<br /><br />
      <b>AIOps for Kubernetes</b><br />Automate and optimize Kubernetes management with AI-powered insights.
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/misc/coming-soon.jpeg" /></kbd>
    </td>
  </tr>
</table>

</br>

## 🌈 Our Vision

The increasing complexity of the kubernetes ecosystem is an undeniable trend that is becoming more and more difficult to manage. This complexity not only entails a heavier burden on operations and maintenance but also slows down the adoption of new technologies by users, limiting their ability to fully leverage the potential of kubernetes.

In general, we wish Karpor to focus on search, insights, and AI, to **break through the increasingly complex maze of kubernetes**, achieving the following **value proposition**:

![](https://kusionstack.io/karpor/assets/overview/vision.png)

## ⚙️ Installation

### Install with Helm

Karpor can be simply installed by helm v3.5+, which is a simple command-line tool and you can get it from [here](https://helm.sh/docs/intro/install/).

If you are interested, you can also directly view the [Karpor Chart Repo](https://github.com/KusionStack/charts).

```bash
$ helm repo add kusionstack https://kusionstack.github.io/charts
$ helm repo update
$ helm install karpor kusionstack/karpor
```

For more information about installation, please check the [Installation Guide](https://kusionstack.io/karpor/getting-started/installation) on official website.

## 📖 Documentation

Detailed documentation is available at [Karpor Website](https://kusionstack.io/karpor).

## 🤝 How to contribute

Karpor is still in the initial stage, and there are many capabilities that need to be made up, so we welcome everyone to participate in construction with us.

- If you don't know how to **start contributing**, you can read the [Contribution Guide](https://kusionstack.io/karpor/developer-guide/contribution-guide), you will know all the details.
- If you don’t know **what issues start**, we have prepared a [Community tasks | 新手任务清单 🎖︎](https://github.com/KusionStack/karpor/issues/463), you can choose the issue you like.
- If you have **any questions**, please [Submit the Issue](https://github.com/KusionStack/karpor/issues) or [Post on the discussions](https://github.com/KusionStack/karpor/discussions/new/choose), we will answer as soon as possible.

## 🎖︎ Contributors

Thanks to these wonderful people! Come and [join us](https://kusionstack.io/karpor/developer-guide/contribution-guide)! 

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/elliotxx"><img src="https://avatars.githubusercontent.com/u/9360247?v=4?s=80" width="80px;" alt="elliotxx"/><br /><sub><b>elliotxx</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=elliotxx" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=elliotxx" title="Documentation">📖</a> <a href="#design-elliotxx" title="Design">🎨</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/panshuai-ps"><img src="https://avatars.githubusercontent.com/u/49754046?v=4?s=80" width="80px;" alt="panshuai-ps"/><br /><sub><b>panshuai-ps</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=panshuai-ps" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=panshuai-ps" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ffforest"><img src="https://avatars.githubusercontent.com/u/5624244?v=4?s=80" width="80px;" alt="Forest"/><br /><sub><b>Forest</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=ffforest" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=ffforest" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/hai-tian"><img src="https://avatars.githubusercontent.com/u/20057132?v=4?s=80" width="80px;" alt="hai-tian"/><br /><sub><b>hai-tian</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=hai-tian" title="Code">💻</a> <a href="#design-hai-tian" title="Design">🎨</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/weieigao"><img src="https://avatars.githubusercontent.com/u/2090295?v=4?s=80" width="80px;" alt="weieigao"/><br /><sub><b>weieigao</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=weieigao" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/shaofan-hs"><img src="https://avatars.githubusercontent.com/u/133250733?v=4?s=80" width="80px;" alt="shaofan-hs"/><br /><sub><b>shaofan-hs</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=shaofan-hs" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/healthjyk"><img src="https://avatars.githubusercontent.com/u/68334452?v=4?s=80" width="80px;" alt="KK"/><br /><sub><b>KK</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=healthjyk" title="Documentation">📖</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/adohe"><img src="https://avatars.githubusercontent.com/u/71679464?v=4?s=80" width="80px;" alt="TonyAdo"/><br /><sub><b>TonyAdo</b></sub></a><br /><a href="#ideas-adohe" title="Ideas, Planning, & Feedback">🤔</a> <a href="#fundingFinding-adohe" title="Funding Finding">🔍</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://blog.wu8685.com/"><img src="https://avatars.githubusercontent.com/u/10124459?v=4?s=80" width="80px;" alt="Kan Wu"/><br /><sub><b>Kan Wu</b></sub></a><br /><a href="#ideas-wu8685" title="Ideas, Planning, & Feedback">🤔</a> <a href="#fundingFinding-wu8685" title="Funding Finding">🔍</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Paradiesvogel7"><img src="https://avatars.githubusercontent.com/u/96288496?v=4?s=80" width="80px;" alt="Paradiesvogel7"/><br /><sub><b>Paradiesvogel7</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=Paradiesvogel7" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/virtually-stray"><img src="https://avatars.githubusercontent.com/u/154653861?v=4?s=80" width="80px;" alt="Stray"/><br /><sub><b>Stray</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=virtually-stray" title="Documentation">📖</a> <a href="https://github.com/KusionStack/karpor/commits?author=virtually-stray" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ruquanzhao"><img src="https://avatars.githubusercontent.com/u/49401013?v=4?s=80" width="80px;" alt="ZhaoRuquan"/><br /><sub><b>ZhaoRuquan</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=ruquanzhao" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=ruquanzhao" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/SparkYuan"><img src="https://avatars.githubusercontent.com/u/4793557?v=4?s=80" width="80px;" alt="Dayuan"/><br /><sub><b>Dayuan</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=SparkYuan" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/wolfcode111"><img src="https://avatars.githubusercontent.com/u/68718623?v=4?s=80" width="80px;" alt="huadongxu"/><br /><sub><b>huadongxu</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=wolfcode111" title="Documentation">📖</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://www.cnblogs.com/sting2me/"><img src="https://avatars.githubusercontent.com/u/3829504?v=4?s=80" width="80px;" alt="Peter Wang"/><br /><sub><b>Peter Wang</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=peter-wangxu" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://blog.solarhell.com/"><img src="https://avatars.githubusercontent.com/u/10279583?v=4?s=80" width="80px;" alt="jiaxin"/><br /><sub><b>jiaxin</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=solarhell" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## ☎️ Contact

If you have any questions, feel free to reach out to us in the following ways:

- [Slack](https://kusionstack.slack.com) | [Join](https://join.slack.com/t/kusionstack/shared_invite/zt-2drafxksz-VzCZZwlraHP4xpPeh_g8lg)
- [DingTalk Group](https://page.dingtalk.com/wow/dingtalk/act/en-home): `42753001`  (Chinese)
- WeChat Group (Chinese): Add the WeChat assistant to bring you into the user group.

  <img src="assets/img/wechat.png" width="200" height="200"/>
