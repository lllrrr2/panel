<p align="right">
[简体中文] | [<a href="README_EN.md">English</a>]
</p>

<h1 align="center" style="font-size: 40px">耗子面板</h1>

<div align="center">

[![Go](https://img.shields.io/github/go-mod/go-version/tnborg/panel)](https://go.dev/)
[![Release](https://img.shields.io/github/release/tnborg/panel.svg)](https://github.com/tnborg/panel/releases)
[![Test](https://github.com/tnborg/panel/actions/workflows/test.yml/badge.svg)](https://github.com/tnborg/panel/actions)
[![Report Card](https://goreportcard.com/badge/github.com/tnborg/panel)](https://goreportcard.com/report/github.com/tnborg/panel)
[![Stars](https://img.shields.io/github/stars/tnborg/panel?style=flat)](https://github.com/tnborg/panel)
[![License](https://img.shields.io/github/license/tnborg/panel)](https://www.gnu.org/licenses/agpl-3.0.html)

</div>

新一代全能服务器运维管理面板。简单轻量，高效运维。

官网：[panel.haozi.net](https://panel.haozi.net) | QQ群：[12370907](https://jq.qq.com/?_wv=1027&k=I1oJKSTH) | 微信群：[复制此链接](https://work.weixin.qq.com/gm/d8ebf618553398d454e3378695c858b6)

## 优势

1. **极低占用:** Go 语言开发，安装包小，占用低，单文件运行，不会对系统性能造成影响
2. **低破坏性:** 设计为尽可能减少对系统的额外修改，在同类产品中，我们对系统的修改最少
3. **追随时代:** 整体设计走在时代前沿，对新系统兼容性好，在同类产品中处于领先地位
4. **高效运维:** 功能完善，自定义能力强，既可快速部署小型网站，也可基于定制化需求部署复杂应用
5. **离线运行:** 支持离线模式，甚至可以在部署完成后停止面板进程，不会对已有服务造成任何影响
6. **安全稳定:** 面板采用业界多种技术保障本体安全性，已在我们的多个生产环境中长期稳定运行
7. **全面开源:** 少有的全开源面板，您可以在遵守开源协议的前提下对面板自由修改、二次开发
8. **永久免费:** 承诺面板本体未来不会引入任何收费/授权功能，永久免费使用

## 快速安装

支持 `amd64` | `arm64` 架构下的干净的主流系统，具体支持的系统请参考[安装文档](https://ratpanel.github.io/zh_CN/quickstart/install)。

```shell
curl -fsLm 10 -o install.sh https://dl.cdn.haozi.net/panel/install.sh && bash install.sh
```

## UI 截图

![UI 截图](.github/assets/ui.png)

## 合作伙伴

如果耗子面板对您有帮助，欢迎[赞助我们](https://github.com/tnborg/panel/issues/90)，同时感谢以下支持者/赞助商的支持：

<p align="center">
  <a href="https://www.weixiaoduo.com/">
    <img height="60" src=".github/assets/wxd.png" alt="微晓朵">
  </a>
  <a href="https://www.dkdun.cn/aff/MQZZNVHQ">
    <img height="60" src=".github/assets/dk.png" alt="林枫云">
  </a>
  <a href="https://waf.pro/">
    <img height="60" src=".github/assets/wafpro.png" alt="WAFPRO">
  </a>
  <a href="https://scdn.ddunyun.com/">
    <img height="60" src=".github/assets/ddunyun.png" alt="盾云SCDN">
  </a>
  <a href="https://1ms.run/">
    <img height="60" src=".github/assets/1ms.svg" alt="毫秒镜像">
  </a>
</p>

<p align="center">
  <a target="_blank" href="https://afdian.com/a/tnblabs">
    <img alt="sponsors" src="https://github.com/tnborg/sponsor/blob/main/sponsors.svg?raw=true"/>
  </a>
</p>

## Star 历史

<a href="https://star-history.com/#tnborg/panel&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=tnborg/panel&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=tnborg/panel&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=tnborg/panel&type=Date" />
 </picture>
</a>
