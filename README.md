# vue-evernote-api

backend of vue-evernote

## 技术选型

golang + jwt + echo rest api + mongoDB + docker

## 接口文档

[接口文档](https://github.com/go-study-projs/vue-evernote-api/wiki)

## 功能

- 完成接口

    - [x] 用户注册
    - [x]  用户登录
    - [x] 获取笔记本列表
    - [x] 创建笔记本
    - [x] 修改笔记本
    - [x] 删除笔记本
    - [x] 创建笔记
    - [x] 获取笔记列表
    - [x] 将笔记放入回收站
    - [x] 修改笔记
    - [x] 将记从回收站彻底删除
    - [x] 回收站恢复笔记
    - [x] 获取回收站笔记列表

- 系统功能
    - [ ] refreshToken过期续签
    - [ ] 接入日志
    - [x] 容器化
    - [ ] 单元测试
    - [ ] CI : Github Action

## 运行

```shell
git clone git@github.com:go-study-projs/vue-evernote-api.git
cd vue-evernote-api
docker-compose up -d
```
默认端口8080
