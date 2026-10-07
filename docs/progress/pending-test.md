---
title: 待测试
description: 当前版本已实现但仍需人工验证的变更项
---

# 待测试

- 仓库地址整体迁移到 `vow132/infinite-canvas`：Go module 及 import 路径、页面 GitHub 图标跳转链接、版本检查 raw 地址、画布 Codex 连接命令、docker-compose 镜像名、canvas-agent 仓库字段、文档克隆地址均已替换，需确认编译和页面跳转正常。
- CI 保留 Go、Comfy Bridge 测试和前端构建，补充只读权限、运行超时及同分支重复运行取消。
- CD 仅在 main 的 push 对应 CI 成功后执行，部署经过 CI 的精确提交；服务器使用独立的 `~/infinite-canvas-deploy` 源码目录，在构建完成后仅更新 `~/pi-web-deploy` 的 `infinite-canvas` 服务，保留现有环境变量、数据卷、Caddy 和其他容器。
- CD 增加 SSH 主机指纹校验、后端与前端及本机域名反代检查；失败时尝试恢复上一版镜像。凭据只存放在 GitHub Secrets，部署运行结果仍需确认。
