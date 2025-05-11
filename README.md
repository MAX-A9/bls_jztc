## 前后端全部开源同城微信小程序（Goframe + uniapp）

### QQ交流群
<a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=cqdDWhjH05oC8FkSsejNeMbirHWO3Ng_&jump_from=webapi&authKey=3M91jbO/ceHsNcxtcyGMbzVS1I3kQTb/Sg/HUeXmBv9GSyWuxCEo74vkyfeN6vb2"><img border="0" src="https://img.shields.io/badge/点击加入-9518324-green.svg" alt="布鲁斯社区①群" title="布鲁斯社区①群"></a>

## 重要信息
1. 项目合作洽谈，请联系客服微信（使用微信扫码添加好友，请注明来意）。
2. 如需二开，定制开发，请联系客服。<br>
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/9GEQOi3pI9s9NQZU.png "微信")

## 使用须知
### ✅允许
- 个人学习使用
- 允许用于学习、毕设等
- 允许进行商业使用，请自觉遵守使用协议，如需要二开增加功能请联系我
- 请遵守 Apache License2.0 协议，再次开源请注明出处
- 推荐Watch、Star项目，获取项目第一时间更新，同时也是对项目最好的支持
- 希望大家多多支持原创作品
- 禁止添加授权二次销售，如若发现，后果自负

## 安装教程

* 配置环境（安装Go环境，并安装依赖）
    * 进入demo目录，执行go mod tidy
    * 进入jztc目录，执行npm install
* 各项命令
    * 启动服务 go run main.go
    * 启动前端 npm run dev
    * 打包后端 go run build
    * 打包前端 npm run build
* 打包相关问题（支持前后端分离部署，也可以一起部署）
    *前后端一起部署，将前端打包后的dist目录，放到resource目录中即可。
    *前后端分离部署，需要在前端和后端中进行配置，前端在.env文件，后端在manifest/config目录中
* 本地环境
    * http://localhost:8000/swagger（文档地址，可在config中关闭）
    * 小程序constants.js中配置域名

## 页面展示
### 接口文档
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/PTmqxTgZo1IvTFMU.jpg "接口文档")
### 首页
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/SCfng8PgXD3zI91L.jpg "首页")
### 闲置
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/ZGzh4OhtBagjFJZy.jpg "闲置")
### 我的
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/JdkRJT70xxd5zZag.jpg "我的页面")
### 发布
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/4EFWR9xjou7KcaQ3.jpg "发布页面")
### 消息
![](https://haowen.oss-cn-beijing.aliyuncs.com/test/uploads/image/20250511/PG78fq1MNLuhOSss.jpg "消息页面")