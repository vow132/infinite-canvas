---
title: 待测试
description: 当前版本已实现但仍需人工验证的变更项
---

# 待测试

- 仓库地址整体迁移到 `vow132/infinite-canvas`：Go module 及 import 路径、页面 GitHub 图标跳转链接、版本检查 raw 地址、画布 Codex 连接命令、docker-compose 镜像名、canvas-agent 仓库字段、文档克隆地址均已替换，需确认编译和页面跳转正常。
- 新增 CD 工作流（`.github/workflows/deploy.yml`）：推送 main 后自动 SSH 到目标服务器，克隆 / 拉取源码并 `docker compose -f docker-compose.local.yml up -d --build` 部署，应用运行在服务器 3000 端口；SSH 凭据存放在 GitHub Encrypted Secrets（SSH_HOST / SSH_PORT / SSH_USERNAME / SSH_PASSWORD），首次部署与服务器上的 `.env` 默认密码需人工确认。
