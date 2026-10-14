> 🇰🇷 이 파일은 원본 README의 한글 번역본입니다


<div align="center">
<p></p><p></p>
<p>
    <img width="60%" src="https://kusionstack.io/karpor/assets/logo/logo-full.svg">
</p>

<h1 style="font-size: 1.5em;">
    쿠버네티스를 위한 인텔리전스
</h1>

<p align="center">
  <a href="https://karpor-demo.kusionstack.io" target="_blank"><b>🎮 라이브 데모</b></a> •
  <a href="https://kusionstack.io/karpor/" target="_blank"><b>🌐 웹사이트</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/quick-start" target="_blank"><b>⚡️ 빠른 시작</b></a> •
  <a href="https://kusionstack.io/karpor" target="_blank"><b>📚 문서</b></a> •
  <a href="https://github.com/KusionStack/karpor/discussions" target="_blank"><b>💬 토론</b></a><br>
  <a href="https://github.com/KusionStack/karpor/blob/main/README.md" target="_blank">[English]</a>
  <a href="https://github.com/KusionStack/karpor/blob/main/README-zh.md" target="_blank">[中文]</a> 
  <a href="https://github.com/KusionStack/karpor/blob/main/README-pt.md" target="_blank">[Português]</a>
  <a href="https://github.com/KusionStack/karpor/blob/main/README-ja.md" target="_blank">[日本語]</a>
  [한국어]
  <a href="https://github.com/KusionStack/karpor/blob/main/README-hi.md" target="_blank">[हिन्दी]</a>
</p>

[![GitHub Release](https://img.shields.io/github/release/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/KusionStack/karpor)](https://goreportcard.com/report/github.com/KusionStack/karpor)
[![Coverage Status](https://coveralls.io/repos/github/KusionStack/karpor/badge.svg)](https://coveralls.io/github/KusionStack/karpor)
[![Go Reference](https://pkg.go.dev/badge/github.com/KusionStack/karpor.svg)](https://pkg.go.dev/github.com/KusionStack/karpor)
[![license](https://img.shields.io/github/license/KusionStack/karpor.svg)](https://github.com/KusionStack/karpor/blob/main/LICENSE)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/karpor)](https://artifacthub.io/packages/helm/kusionstack/karpor)
[![CNCF](https://shields.io/badge/CNCF-Sandbox%20project-blue?logo=linux-foundation&style=flat)](https://landscape.cncf.io/?item=provisioning--automation-configuration--kusionstack)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FKusionStack%2Fkarpor.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FKusionStack%2Fkarpor?ref=badge_shield)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/KusionStack/karpor)
[![Medium](https://img.shields.io/badge/@kusionstack-black?style=flat&logo=medium&logoColor=white&link=https://medium.com/@kusionstack)](https://medium.com/@kusionstack)
[![Slack](https://img.shields.io/badge/slack-kusion-blueviolet?logo=slack)](https://cloud-native.slack.com/archives/C07U0395UG0)

<a href="https://trendshift.io/repositories/11086" target="_blank"><img src="https://trendshift.io/api/badge/repositories/11086" alt="KusionStack%2Fkarpor | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a> <a href="https://www.producthunt.com/posts/karpor" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=770525&theme=light&t=1736447376990" alt="Karpor - &#0032;Effortlessly&#0032;manage&#0032;Kubernetes&#0032;with&#0032;AI&#0045;powered&#0032;insights | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

</div>

## Karpor란 무엇인가요?

Karpor는 쿠버네티스를 위한 인텔리전스입니다. 쿠버네티스에 고급 **검색**, **인사이트**, 그리고 **AI** 기능을 제공합니다. 본질적으로 **쿠버네티스 시각화 도구**입니다. Karpor를 사용하면 어떤 클라우드에서든 쿠버네티스 클러스터 전반에 걸쳐 중요한 가시성을 확보할 수 있습니다.

우리는 **작지만 훌륭하고, 벤더 중립적이며, 개발자 친화적인, 커뮤니티 중심의** 오픈소스 프로젝트가 되기를 희망합니다!

**현재 상태**: [v0.7.0](https://github.com/KusionStack/karpor/issues/723)을 반복 개발 중입니다. [토론](https://github.com/KusionStack/karpor/discussions/737)에 참여해 주세요.

https://github.com/user-attachments/assets/44ecb3b6-e76b-4003-9546-25c6b5f8b949

## 왜 Karpor인가요?

<h3 align="center">🔍 검색</h3>

<table>
  <tr>
    <td>
      <b>자동 동기화</b><br />멀티 클라우드 플랫폼에서 관리되는 모든 클러스터의 리소스를 자동으로 동기화합니다.<br /><br/>
      <b>강력하고 유연한 쿼리</b><br />여러 클러스터에 걸쳐 원하는 리소스를 빠르고 쉽게 검색하고 찾아낼 수 있습니다.
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/search/search-auto-complete-raw.jpg" /></kbd>
    </td>
  </tr>
</table>

<h3 align="center">💡 인사이트</h3>

<table>
  <tr>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/insight/insight-home-raw.jpg" /></kbd>
    </td>
    <td>
      <b>컴플라이언스 거버넌스</b><br />여러 클러스터와 컴플라이언스 기준에 걸친 규정 준수 현황을 파악합니다.<br /><br/>
      <b>리소스 토폴로지</b><br />운영 컨텍스트 내에서 관련 리소스의 논리적·토폴로지적 뷰를 제공합니다.<br /><br/>
      <b>비용 최적화</b><br />출시 예정.
    </td>
  </tr>
</table>

<h3 align="center">✨ AI</h3>

<table>
  <tr>
    <td>
      <b>자연어 운영</b><br />일상 언어를 사용하여 쿠버네티스와 상호작용함으로써 더욱 직관적인 운영을 실현합니다.<br /><br />
      <b>맥락 인식 AI 응답</b><br />사용자의 요구사항을 이해하는 스마트하고 맥락적인 지원을 제공합니다.<br /><br />
      <b>쿠버네티스를 위한 AIOps</b><br />AI 기반 인사이트로 쿠버네티스 관리를 자동화하고 최적화합니다.
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/ai/event-ai-diagnosis.png" /></kbd>
    </td>
  </tr>
</table>

</br>

## 🌈 우리의 비전

쿠버네티스 생태계의 복잡성이 점점 더 증가하는 것은 부정할 수 없는 추세이며, 관리하기가 갈수록 어려워지고 있습니다. 이러한 복잡성은 운영·유지보수 부담을 가중시킬 뿐만 아니라, 사용자들이 새로운 기술을 도입하는 속도를 늦추고 쿠버네티스의 잠재력을 충분히 활용하는 능력을 제한합니다.

전반적으로, 우리는 Karpor가 검색, 인사이트, AI에 집중하여 **점점 복잡해지는 쿠버네티스의 미로를 돌파**하고, 다음과 같은 **가치 제안**을 실현하기를 바랍니다:

![](https://kusionstack.io/karpor/assets/overview/vision.png)

## ⚙️ 설치

### Helm으로 설치

Karpor는 helm v3.5+로 간단하게 설치할 수 있습니다. Helm은 간단한 커맨드라인 도구이며 [여기](https://helm.sh/docs/intro/install/)에서 받을 수 있습니다.

관심이 있으시다면 [Karpor Chart Repo](https://github.com/KusionStack/charts)를 직접 확인해 보실 수도 있습니다.

```bash
$ helm repo add kusionstack https://kusionstack.github.io/charts
$ helm repo update
$ helm install karpor kusionstack/karpor
```

설치에 대한 자세한 내용은 공식 웹사이트의 [설치 가이드](https://kusionstack.io/karpor/getting-started/installation)를 확인하세요.

## 📖 문서

자세한 문서는 [Karpor 웹사이트](https://kusionstack.io/karpor)에서 확인하실 수 있습니다.

## 🤝 기여 방법

Karpor는 아직 초기 단계에 있으며, 구현해야 할 기능이 많습니다. 따라서 모든 분들의 참여를 환영합니다.

- **기여를 시작하는 방법**을 모르신다면 [기여 가이드](https://kusionstack.io/karpor/developer-guide/contribution-guide)를 읽어보세요. 모든 세부 사항을 확인하실 수 있습니다.
- **어떤 이슈부터 시작해야 할지** 모르신다면 [커뮤니티 태스크 | 新手任务清单 🎖︎](https://github.com/KusionStack/karpor/issues/463)를 준비해 두었으니 마음에 드는 이슈를 선택하세요.
- **질문이 있으시면** [이슈 제출](https://github.com/KusionStack/karpor/issues) 또는 [토론에 게시](https://github.com/KusionStack/karpor/discussions/new/choose)해 주시면 최대한 빨리 답변드리겠습니다.

## 🎖︎ 기여자

이 훌륭한 분들께 감사드립니다! 어서 [함께해요](https://kusionstack.io/karpor/developer-guide/contribution-guide)!

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/elliotxx"><img src="https://avatars.githubusercontent.com/u/9360247?v=4?s=80" width="80px;" alt="elliotxx"/><br /><sub><b>elliotxx</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=elliotxx" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=elliotxx" title="Documentation">📖</a> <a href="#design-elliotxx" title="Design">🎨</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/panshuai111"><img src="https://avatars.githubusercontent.com/u/49754046?v=4?s=80" width="80px;" alt="panshuai111"/><br /><sub><b>panshuai111</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=panshuai111" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=panshuai111" title="Documentation">📖</a></td>
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
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/rajeshkio"><img src="https://avatars.githubusercontent.com/u/107089376?v=4?s=80" width="80px;" alt="rajeshkio"/><br /><sub><b>rajeshkio</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=rajeshkio" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://chenchen.link/"><img src="https://avatars.githubusercontent.com/u/30647904?v=4?s=80" width="80px;" alt="Chen Chen"/><br /><sub><b>Chen Chen</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=z1cheng" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/JasonHe-WQ"><img src="https://avatars.githubusercontent.com/u/85824149?v=4?s=80" width="80px;" alt="JasonHe-WQ"/><br /><sub><b>JasonHe-WQ</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=JasonHe-WQ" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/CirillaQL"><img src="https://avatars.githubusercontent.com/u/46399960?v=4?s=80" width="80px;" alt="EduardoQian"/><br /><sub><b>EduardoQian</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=CirillaQL" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/regend"><img src="https://avatars.githubusercontent.com/u/5510284?v=4?s=80" width="80px;" alt="xiao.wu"/><br /><sub><b>xiao.wu</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=regend" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/iamryanchia"><img src="https://avatars.githubusercontent.com/u/41557860?v=4?s=80" width="80px;" alt="iamryanchia"/><br /><sub><b>iamryanchia</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=iamryanchia" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/cheny-alf"><img src="https://avatars.githubusercontent.com/u/71162267?v=4?s=80" width="80px;" alt="cheny-alf"/><br /><sub><b>cheny-alf</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=cheny-alf" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/rajp152k"><img src="https://avatars.githubusercontent.com/u/42045620?v=4?s=80" width="80px;" alt="Raj"/><br /><sub><b>Raj</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=rajp152k" title="Documentation">📖</a> <a href="https://github.com/KusionStack/karpor/commits?author=rajp152k" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://wavee.world/invitation/b96d00e6-b802-4a1b-8a66-2e3854a01ffd"><img src="https://avatars.githubusercontent.com/u/22633385?v=4?s=80" width="80px;" alt="Ikko Eltociear Ashimine"/><br /><sub><b>Ikko Eltociear Ashimine</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=eltociear" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/wujunwei"><img src="https://avatars.githubusercontent.com/u/14191653?v=4?s=80" width="80px;" alt="Adam"/><br /><sub><b>Adam</b></sub></a><br /><a href="#tool-wujunwei" title="Tools">🔧</a> <a href="https://github.com/KusionStack/karpor/commits?author=wujunwei" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Cookiery"><img src="https://avatars.githubusercontent.com/u/33125275?v=4?s=80" width="80px;" alt="Cookie"/><br /><sub><b>Cookie</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=Cookiery" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=Cookiery" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://jinjia.jsisu.com/"><img src="https://avatars.githubusercontent.com/u/28686975?v=4?s=80" width="80px;" alt="jinjia"/><br /><sub><b>jinjia</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=jinjiaKarl" title="Code">💻</a> <a href="https://github.com/KusionStack/karpor/commits?author=jinjiaKarl" title="Documentation">📖</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/LeoSofl"><img src="https://avatars.githubusercontent.com/u/37336398?v=4?s=80" width="80px;" alt="sofl"/><br /><sub><b>sofl</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=LeoSofl" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/nueavv"><img src="https://avatars.githubusercontent.com/u/90682513?v=4?s=80" width="80px;" alt="1102"/><br /><sub><b>1102</b></sub></a><br /><a href="https://github.com/KusionStack/karpor/commits?author=nueavv" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## ☎️ 연락처

질문이 있으시면 아래 방법으로 언제든지 연락 주세요:

- **Slack 채널:**
  - [#kusion](https://cloud-native.slack.com/archives/C07U0395UG0) - Karpor 및 KusionStack에 대한 기술적 토론
  - [#kusion-general](https://cloud-native.slack.com/archives/C07T4LBDB7G) - 일반 토론, 공지 및 커뮤니티 업데이트
- [DingTalk 그룹](https://page.dingtalk.com/wow/dingtalk/act/en-home): `42753001` (중국어)
- WeChat 그룹 (중국어): WeChat 어시스턴트를 추가하면 사용자 그룹에 초대해 드립니다.

  <img src="assets/img/wechat.png" width="200" height="200"/>


## 라이선스
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FKusionStack%2Fkarpor.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FKusionStack%2Fkarpor?ref=badge_large)
