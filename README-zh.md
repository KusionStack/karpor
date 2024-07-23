<div align="center">
<p></p><p></p>
<p>
    <img width="60%" src="https://kusionstack.io/karpor/assets/logo/logo-full.svg"> 
</p>

<h1 style="font-size: 1.5em;">
    Intelligence for Kubernetes
</h1>

<p align="center">
  <a href="https://karpor-demo.kusionstack.io" target="_blank"><b>🎮 现场演示</b></a> •
  <a href="https://kusionstack.io/zh/karpor/" target="_blank"><b>🌐 官网</b></a> •
  <a href="https://kusionstack.io/zh/karpor/getting-started/quick-start" target="_blank"><b>⚡️ 快速开始</b></a> •
  <a href="https://kusionstack.io/zh/karpor" target="_blank"><b>📚 文档</b></a> •
  <a href="https://github.com/KusionStack/karpor/discussions" target="_blank"><b>💬 讨论</b></a><br>
  <a href="https://github.com/KusionStack/karpor/blob/main/README.md" target="_blank">[English]</a> 
  [中文]
  <a href="https://github.com/KusionStack/karpor/blob/main/README-pt.md" target="_blank">[Português]</a>
</p>

[![Karpor](https://github.com/KusionStack/karpor/actions/workflows/release.yaml/badge.svg)](https://github.com/KusionStack/karpor/actions/workflows/release.yaml)
[![GitHub release](https://img.shields.io/github/release/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/KusionStack/karpor)](https://goreportcard.com/report/github.com/KusionStack/karpor)
[![Coverage Status](https://coveralls.io/repos/github/KusionStack/karpor/badge.svg)](https://coveralls.io/github/KusionStack/karpor)
[![Go Reference](https://pkg.go.dev/badge/github.com/KusionStack/karpor.svg)](https://pkg.go.dev/github.com/KusionStack/karpor)
[![license](https://img.shields.io/github/license/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/blob/main/LICENSE)

</div>

## Karpor 是什么？

Karpor 是一个现代化的 Kubernetes 可视化工具，核心特性聚焦在 **搜索、洞察、智能** ，它可以从任何云平台上获得对 Kubernetes 跨集群的关键可见性提供给开发者和平台团队，帮助他们更好地理解集群并做出决策。

我们立志成为一个 **小而美、厂商中立、开发者友好、社区驱动** 的开源项目！

**当前状态**: 我们正在迭代 [v0.5.0 里程碑](https://github.com/KusionStack/karpor/milestone/4),  欢迎加入 [讨论](https://github.com/KusionStack/karpor/discussions/528).

https://github.com/KusionStack/karpor/assets/49401013/7cf31cc0-7123-42f6-8543-5addcbf4975c

## 为什么选择 Karpor？

<h3 align="center">🔍 搜索</h3>

<table>
  <tr>
    <td>
      <b>自动同步</b><br />自动同步您在多云平台管理的任何集群中的资源<br /><br/>
      <b>强大灵活的查询</b><br />以快速简单的方式有效地检索和定位跨集群的资源
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/search/search-auto-complete-raw.jpg"  /></kbd>
    </td>
  </tr>
</table>

<h3 align="center">💡 洞察</h3>

<table>
  <tr>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/insight/insight-home-raw.jpg"  /></kbd>
    </td>
    <td>
      <b>安全合规</b><br />了解您在多个集群和合规标准中的合规性状态<br /><br/>
      <b>资源拓扑</b><br />提供包含资源运行上下文信息的关系拓扑和逻辑视图<br /><br/>
      <b>成本优化</b><br />即将推出
    </td>
  </tr>
</table>

<h3 align="center">✨ 智能</h3>

<table>
  <tr>
    <td>
      <b>自然语言操作</b><br />使用自然语言与 Kubernetes 交互，实现更直观的操作<br /><br />
      <b>情境化 AI 响应</b><br />获得智能的、情境化的辅助，满足您的需求<br /><br />
      <b>Kubernetes AIOps</b><br />利用 AI 驱动的洞察，自动化和优化 Kubernetes 管理
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/misc/coming-soon.jpeg"  /></kbd>
    </td>
  </tr>
</table>

</br>

## 🌈 愿景

现如今，Kubernetes 生态系统日益复杂是一个不可否认的趋势，这一趋势越来越难以驾驭。这种复杂性不仅增加了运维的难度，也降低了用户采纳新技术的速度，从而限制了他们充分利用 Kubernetes 的潜力。

我们希望 Karpor 围绕着搜索、洞察和AI，**击穿 Kubernetes 愈演愈烈的复杂性**，达成以下**价值主张**：

![](https://kusionstack.io/karpor/assets/overview/vision.png)

## ⚙️ 安装

### 使用 Helm 安装

Karpor 可以通过 helm v3.5+ 简单安装，这是一个简单的命令行工具，您可以从[这里](https://helm.sh/docs/intro/install/)获取。

如果您感兴趣，您也可以直接查看 [Karpor Chart Repo](https://github.com/KusionStack/charts)。

```bash
$ helm repo add kusionstack https://kusionstack.github.io/charts 
$ helm repo update
$ helm install karpor kusionstack/karpor
```

有关安装的更多信息，请查看官方网站上的 [安装指南](https://kusionstack.io/zh/karpor/getting-started/installation)。

## 📖 文档

详细的文档可在 [Karpor 官网](https://kusionstack.io/zh/karpor) 查阅。

## 🤝 如何贡献

Karpor 仍处于初期阶段，仍有许多功能需要构建，因此我们欢迎每个人与我们共同参与建设。

- 如果您不知道如何 **开始贡献**，您可以阅读[贡献指南](https://kusionstack.io/zh/karpor/developer-guide/contribution-guide)，您将了解所有细节。
- 如果您不知道 **从哪些问题开始**，我们准备了[社区任务 | 新手任务清单 🎖︎](https://github.com/KusionStack/karpor/issues/463)，您可以选择您喜欢的问题。
- 如果您有任何 **问题**，请[提交问题](https://github.com/KusionStack/karpor/issues)或[在讨论中发帖](https://github.com/KusionStack/karpor/discussions/new/choose)，我们将尽快回答。

## 🎖︎ 贡献者

感谢这些了不起的人！来[加入我们](https://kusionstack.io/zh/karpor/developer-guide/contribution-guide)吧！

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

## ☎️ 联系方式

如果您有任何问题，欢迎通过以下方式联系我们：

- [Slack](https://kusionstack.slack.com) | [加入](https://join.slack.com/t/kusionstack/shared_invite/zt-2drafxksz-VzCZZwlraHP4xpPeh_g8lg)
- [钉钉群](https://page.dingtalk.com/wow/dingtalk/act/en-home)：`42753001`（中文）
- 微信群（中文）：添加微信小助手，拉你进用户群

  <img src="assets/img/wechat.png" width="200" height="200"/>
