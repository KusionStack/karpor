<div align="center">
<img width="60%" src="https://kusionstack.io/karpor/assets/logo/logo-full.svg" style="max-width:100%;">

<h2 style="font-size: 1.5em;">
  Intuitive Discovery 🔍 ,<br>
  Limitless Insight 📊.<span style="color: gray; font-weight: normal;"> With AI.<sup style="color: gray; font-weight: normal; font-size: 0.5em"> coming soon</sup></span>✨
</h2>

<p align="center">
  <a href="https://karpor-demo.kusionstack.io" target="_blank"><b>🎮 现场演示</b></a> •
  <a href="https://kusionstack.io/karpor/" target="_blank"><b>🚀 概览</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/installation" target="_blank"><b>⚙️ 安装</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/quick-start" target="_blank"><b>⚡️ 快速开始</b></a> •
  <a href="https://kusionstack.io/karpor" target="_blank"><b>📚 文档</b></a>
</p>

[![Karpor](https://github.com/KusionStack/karpor/actions/workflows/release.yaml/badge.svg)](https://github.com/KusionStack/karpor/actions/workflows/release.yaml) 
[![GitHub release](https://img.shields.io/github/release/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/releases) 
[![Go Report Card](https://goreportcard.com/badge/github.com/KusionStack/karpor)](https://goreportcard.com/report/github.com/KusionStack/karpor) 
[![Coverage Status](https://coveralls.io/repos/github/KusionStack/karpor/badge.svg)](https://coveralls.io/github/KusionStack/karpor) 
[![Go Reference](https://pkg.go.dev/badge/github.com/KusionStack/karpor.svg)](https://pkg.go.dev/github.com/KusionStack/karpor) 
[![license](https://img.shields.io/github/license/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/blob/main/LICENSE) 
</div>

## Karpor 是什么？

Karpor 是一个现代化的 **Kubernetes 资源探索器**，专注于 **🔍 搜索**、**💡 洞察** 和 **🤖 智能**。它具备 `自托管`、`非侵入式`、`只读`、`安全合规`、`多云/多集群支持`、`自定义逻辑视图` 等特性，并且可以作为 **Kubernetes 数据面**，降低发现和理解 Kubernetes 资源的成本。

## 为什么选择 Karpor？

<h3 align="center">🤝 用户友好</h3>

<table>
  <tr>
    <td style="vertical-align: middle;">
      <b>⚡️ 轻量级，易于安装</b><br />仅需一个 `helm` 命令即可完成安装。<br /><br/>
      <b>📦 开箱即用</b><br />内置安全和合规策略、资源同步策略、资源转换规则和拓扑关系定义。<br /><br />
      <b>🔍 快速搜索和定位资源</b><br />以用户友好的方式快速定位跨集群资源。
    </td>
    <td style="width: 60%; vertical-align: middle;">
      <img src="https://kusionstack.io/karpor/assets/overview/user-friendly.png" alt="用户友好" />
    </td>
  </tr>
</table>

<h3 align="center">✨ 智能</h3>

<table>
  <tr>
    <td style="width: 60%; vertical-align: middle;">
      <img src="https://kusionstack.io/karpor/assets/overview/intelligent.png" alt="智能" />
    </td>
    <td style="vertical-align: middle;">
      <b>🔒 合规保护</b><br />自动识别潜在风险并生成 AI 整改建议。<br /><br/>
      <b>📊 逻辑和拓扑视图</b><br />在其运行上下文中展示相关资源的逻辑和拓扑视图。
    </td>
  </tr>
</table>

<h3 align="center">⚡️ 低负担</h3>

<table>
  <tr>
    <td style="vertical-align: middle;">
      <b>🔒 只读，非侵入式</b><br />只读数据平面，对用户集群无侵入。<br /><br />
      <b>⚙️ 兼容 Kubernetes 原生 API</b><br />与现有 Kubernetes 工具链的无缝集成。<br /><br />
      <b>☁️ 多集群和多云/混合云</b><br />原生支持多集群和多云/混合云。
    </td>
    <td style="width: 60%; vertical-align: middle;">
      <img src="https://kusionstack.io/karpor/assets/overview/low-burden.png" alt="低负担" />
    </td>
  </tr>
</table>

</br>

## 安装

### 使用 Helm 安装

Karpor 可以通过 Helm v3.5+ 简单安装。你可以从 [这里](https://helm.sh/docs/intro/install/) 获取 Helm。

如果你感兴趣，也可以直接查看 [Karpor Chart Repo](https://github.com/KusionStack/charts)。

安装命令:

```bash
helm repo add kusionstack https://kusionstack.github.io/charts
helm repo update
helm install karpor kusionstack/karpor
```

更多安装信息，请查看 [官方安装指南](https://kusionstack.io/karpor/getting-started/installation)。

## 文档

详细文档可在 [Karpor 官网](https://kusionstack.io/karpor) 查阅。

## 如何贡献

Karpor 仍处于初期阶段，我们欢迎每个人参与贡献。

- 如果你不知道如何 **开始贡献**，可以阅读 [贡献指南](https://kusionstack.io/karpor/developer-guide/contribution-guide) 了解所有细节。
- 如果你不知道 **从哪开始**，我们准备了 [社区任务 | 新手任务清单 🎖︎](https://github.com/KusionStack/karpor/issues/463)，你可以选择感兴趣的开始。
- 如果你有任何 **问题**，请 [提交问题](https://github.com/KusionStack/karpor/issues) 或 [参与讨论](https://github.com/KusionStack/karpor/discussions/new/choose)，我们会尽快回答。

## 贡献者

感谢所有贡献者！欢迎[加入我们](https://kusionstack.io/karpor/developer-guide/contribution-guide)。

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/elliotxx"><img src="https://avatars.githubusercontent.com/u/9360247?v=4?s=70" width="70px;" alt="elliotxx"/><br /><sub><b>elliotxx</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=elliotxx" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=elliotxx" title="Documentation">📖</a> <a href="#design-elliotxx" title="Design">🎨</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/panshuai-ps"><img src="https://avatars.githubusercontent.com/u/49754046?v=4?s=70" width="70px;" alt="panshuai-ps"/><br /><sub><b>panshuai-ps</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=panshuai-ps" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=panshuai-ps" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ffforest"><img src="https://avatars.githubusercontent.com/u/5624244?v=4?s=70" width="70px;" alt="Forest"/><br /><sub><b>Forest</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=ffforest" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=ffforest" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/hai-tian"><img src="https://avatars.githubusercontent.com/u/20057132?v=4?s=70" width="70px;" alt="hai-tian"/><br /><sub><b>hai-tian</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=hai-tian" title="Code">💻</a> <a href="#design-hai-tian" title="Design">🎨</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/weieigao"><img src="https://avatars.githubusercontent.com/u/2090295?v=4?s=70" width="70px;" alt="weieigao"/><br /><sub><b>weieigao</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=weieigao" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/shaofan-hs"><img src="https://avatars.githubusercontent.com/u/133250733?v=4?s=70" width="70px;" alt="shaofan-hs"/><br /><sub><b>shaofan-hs</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=shaofan-hs" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/healthjyk"><img src="https://avatars.githubusercontent.com/u/68334452?v=4?s=70" width="70px;" alt="KK"/><br /><sub><b>KK</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=healthjyk" title="Documentation">📖</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/adohe"><img src="https://avatars.githubusercontent.com/u/71679464?v=4?s=70" width="70px;" alt="TonyAdo"/><br /><sub><b>TonyAdo</b></sub></a><br /><a href="#ideas-adohe" title="Ideas, Planning, & Feedback">🤔</a> <a href="#fundingFinding-adohe" title="Funding Finding">🔍</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://blog.wu8685.com/"><img src="https://avatars.githubusercontent.com/u/10124459?v=4?s=70" width="70px;" alt="Kan Wu"/><br /><sub><b>Kan Wu</b></sub></a><br /><a href="#ideas-wu8685" title="Ideas, Planning, & Feedback">🤔</a> <a href="#fundingFinding-wu8685" title="Funding Finding">🔍</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->
