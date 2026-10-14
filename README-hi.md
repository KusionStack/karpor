> 🇮🇳 यह फ़ाइल मूल README का हिन्दी अनुवाद है।


<div align="center">
<p></p><p></p>
<p>
    <img width="60%" src="https://kusionstack.io/karpor/assets/logo/logo-full.svg">
</p>

<h1 style="font-size: 1.5em;">
    Kubernetes के लिए इंटेलिजेंस
</h1>

<p align="center">
  <a href="https://karpor-demo.kusionstack.io" target="_blank"><b>🎮 लाइव डेमो</b></a> •
  <a href="https://kusionstack.io/karpor/" target="_blank"><b>🌐 वेबसाइट</b></a> •
  <a href="https://kusionstack.io/karpor/getting-started/quick-start" target="_blank"><b>⚡️ क्विक स्टार्ट</b></a> •
  <a href="https://kusionstack.io/karpor" target="_blank"><b>📚 डॉक्यूमेंटेशन</b></a> •
  <a href="https://github.com/KusionStack/karpor/discussions" target="_blank"><b>💬 चर्चाएँ</b></a><br>
  <a href="https://github.com/KusionStack/karpor/blob/main/README.md" target="_blank">[English]</a>
  <a href="https://github.com/KusionStack/karpor/blob/main/README-zh.md" target="_blank">[中文]</a>
  <a href="https://github.com/KusionStack/karpor/blob/main/README-pt.md" target="_blank">[Português]</a>
  <a href="https://github.com/KusionStack/karpor/blob/main/README-ja.md" target="_blank">[日本語]</a>
  <a href="https://github.com/KusionStack/karpor/blob/main/README-ko.md" target="_blank">[한국어]</a>
  [हिन्दी]
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

## Karpor क्या है?

Karpor, Kubernetes के लिए इंटेलिजेंस है। यह Kubernetes में उन्नत **खोज (Search)**, **इनसाइट (Insight)** और **AI** लाता है। मूल रूप से यह एक **Kubernetes विज़ुअलाइज़ेशन टूल** है। Karpor की मदद से आप किसी भी क्लाउड पर चल रहे अपने Kubernetes क्लस्टर की ज़रूरी जानकारी एक जगह देख सकते हैं।

हमारा लक्ष्य एक **छोटा और सुंदर, वेंडर-न्यूट्रल, डेवलपर-फ़्रेंडली और समुदाय-संचालित** ओपन-सोर्स प्रोजेक्ट बनना है!

**मौजूदा स्थिति**: हम [v0.7.0](https://github.com/KusionStack/karpor/issues/723) पर काम कर रहे हैं। [चर्चा](https://github.com/KusionStack/karpor/discussions/737) में आपका स्वागत है।

https://github.com/user-attachments/assets/44ecb3b6-e76b-4003-9546-25c6b5f8b949

## Karpor क्यों?

<h3 align="center">🔍 खोज</h3>

<table>
  <tr>
    <td>
      <b>ऑटोमैटिक सिंकिंग</b><br />मल्टी-क्लाउड प्लेटफ़ॉर्म से प्रबंधित किसी भी क्लस्टर के रिसोर्स अपने-आप सिंक करें।<br /><br/>
      <b>शक्तिशाली और लचीली क्वेरी</b><br />कई क्लस्टर में फैले जिन रिसोर्स को आप ढूँढ रहे हैं, उन्हें जल्दी और आसानी से खोजें।
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/search/search-auto-complete-raw.jpg" /></kbd>
    </td>
  </tr>
</table>

<h3 align="center">💡 इनसाइट</h3>

<table>
  <tr>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/insight/insight-home-raw.jpg" /></kbd>
    </td>
    <td>
      <b>कंप्लायंस गवर्नेंस</b><br />कई क्लस्टर और कंप्लायंस मानकों में अपनी कंप्लायंस स्थिति समझें।<br /><br/>
      <b>रिसोर्स टोपोलॉजी</b><br />संबंधित रिसोर्स को उनके ऑपरेशनल संदर्भ में लॉजिकल और टोपोलॉजिकल रूप से देखें।<br /><br/>
      <b>लागत अनुकूलन</b><br />जल्द आ रहा है।
    </td>
  </tr>
</table>

<h3 align="center">✨ AI</h3>

<table>
  <tr>
    <td>
      <b>सामान्य भाषा में ऑपरेशन</b><br />सरल भाषा में Kubernetes से बात करें, ताकि काम करना और आसान हो।<br /><br />
      <b>संदर्भ के अनुसार AI जवाब</b><br />आपकी ज़रूरत समझने वाली स्मार्ट और संदर्भ-आधारित मदद पाएँ।<br /><br />
      <b>Kubernetes के लिए AIOps</b><br />AI-आधारित इनसाइट से Kubernetes प्रबंधन को ऑटोमेट और बेहतर करें।
    </td>
    <td width="60%">
      <kbd><img src="https://kusionstack.io/karpor/assets/ai/event-ai-diagnosis.png" /></kbd>
    </td>
  </tr>
</table>

</br>

## 🌈 हमारा विज़न

Kubernetes इकोसिस्टम की बढ़ती जटिलता एक सच्चाई है, और इसे संभालना लगातार मुश्किल होता जा रहा है। इस जटिलता से ऑपरेशन और मेंटेनेंस का बोझ तो बढ़ता ही है, साथ ही उपयोगकर्ता नई तकनीकें धीरे अपनाते हैं और Kubernetes की पूरी क्षमता का लाभ नहीं उठा पाते।

कुल मिलाकर, हम चाहते हैं कि Karpor खोज, इनसाइट और AI पर ध्यान दे, ताकि **Kubernetes की लगातार जटिल होती भूलभुलैया से बाहर निकला जा सके**, और नीचे दिया गया **वैल्यू प्रपोज़िशन** हासिल हो:

![](https://kusionstack.io/karpor/assets/overview/vision.png)

## ⚙️ इंस्टॉलेशन

### Helm से इंस्टॉल करें

Karpor को helm v3.5+ से आसानी से इंस्टॉल किया जा सकता है। helm एक सरल कमांड-लाइन टूल है, जिसे आप [यहाँ](https://helm.sh/docs/intro/install/) से पा सकते हैं।

अगर आप चाहें, तो सीधे [Karpor Chart Repo](https://github.com/KusionStack/charts) भी देख सकते हैं।

```bash
$ helm repo add kusionstack https://kusionstack.github.io/charts
$ helm repo update
$ helm install karpor kusionstack/karpor
```

इंस्टॉलेशन के बारे में ज़्यादा जानकारी के लिए आधिकारिक वेबसाइट पर [इंस्टॉलेशन गाइड](https://kusionstack.io/karpor/getting-started/installation) देखें।

## 📖 डॉक्यूमेंटेशन

विस्तृत डॉक्यूमेंटेशन [Karpor वेबसाइट](https://kusionstack.io/karpor) पर उपलब्ध है।

## 🤝 योगदान कैसे करें

Karpor अभी शुरुआती दौर में है और कई क्षमताएँ अभी बननी बाकी हैं, इसलिए इसे हमारे साथ मिलकर बनाने में सभी का स्वागत है।

- अगर आपको नहीं पता कि **योगदान कैसे शुरू करें**, तो [योगदान गाइड](https://kusionstack.io/karpor/developer-guide/contribution-guide) पढ़ें। वहाँ सारी जानकारी मिलेगी।
- अगर आपको नहीं पता कि **किस issue से शुरुआत करें**, तो हमने [Community tasks | 新手任务清单 🎖︎](https://github.com/KusionStack/karpor/issues/463) तैयार की है। उसमें से अपनी पसंद का issue चुनें।
- अगर आपके **कोई सवाल** हैं, तो [issue बनाएँ](https://github.com/KusionStack/karpor/issues) या [चर्चा में पोस्ट करें](https://github.com/KusionStack/karpor/discussions/new/choose)। हम जल्द से जल्द जवाब देंगे।

## 🎖︎ योगदानकर्ता

इन सभी शानदार लोगों का धन्यवाद! आप भी [हमसे जुड़ें](https://kusionstack.io/karpor/developer-guide/contribution-guide)!

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

## ☎️ संपर्क

अगर आपके कोई सवाल हैं, तो इन तरीकों से हमसे संपर्क करें:

- **Slack चैनल:**
  - [#kusion](https://cloud-native.slack.com/archives/C07U0395UG0) - Karpor और KusionStack पर तकनीकी चर्चा
  - [#kusion-general](https://cloud-native.slack.com/archives/C07T4LBDB7G) - सामान्य चर्चा, घोषणाएँ और कम्युनिटी अपडेट
- [DingTalk ग्रुप](https://page.dingtalk.com/wow/dingtalk/act/en-home): `42753001` (चीनी भाषा)
- WeChat ग्रुप (चीनी भाषा): यूज़र ग्रुप में जुड़ने के लिए WeChat असिस्टेंट को जोड़ें।

  <img src="assets/img/wechat.png" width="200" height="200"/>


## लाइसेंस
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FKusionStack%2Fkarpor.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FKusionStack%2Fkarpor?ref=badge_large)
