# 无限画布 (infinite-canvas)

开源 AI 创作工作台，集成无限画布、Agent、导演台、全景图、AI 生图、图片编辑、视频生成、画布编排等，兼容 OpenAI 接口及多种 API 服务。

把画布编排、AI 图片 / 视频 / 音频生成、参考图编辑、对话助手、提示词库和素材沉淀放在同一个界面里，适合用来探索视觉方案并连续迭代结果。

## 核心功能

- 无限画布：节点式编排图片、视频、音频素材，自由缩放、连线与分组
- AI 生成：接入 OpenAI 兼容接口及多种 API 渠道，支持生图、图片编辑、视频、音频
- 导演台与全景图：时间轴编排与全景图生成
- 画布 Agent：通过 Canvas Agent / Codex 直接读取和操作当前画布
- 提示词库：内置提示词仓库同步，支持第三方 GitHub 提示词源
- 账号体系：登录同步、管理后台、渠道与令牌管理；未登录时数据保存在浏览器本地

## 快速开始

### Docker 启动

```bash
git clone https://github.com/vow132/infinite-canvas.git
cd infinite-canvas
cp .env.example .env
docker compose -f docker-compose.local.yml up -d --build
```

启动后访问 <http://localhost:3000>，默认管理员账号 `admin`，密码为 `.env` 中的 `ADMIN_PASSWORD`。

### 本地开发

```bash
cp .env.example .env
go run .        # 启动后端
cd web && bun install && bun run dev   # 启动前端
```

首次使用建议：打开右上角配置弹窗，填入自己的 `Base URL`、`API Key` 和模型名；如使用后台渠道模式，再到管理后台补充系统模型与渠道配置。

## CI/CD

GitHub Actions 的 CI 检查 Go 后端、Comfy Bridge 和前端构建。推送 `main` 且 CI 成功后，CD 才通过 SSH 部署对应提交到现有服务：<https://ac.91i.asia>。

服务器部署源码位于 `~/infinite-canvas-deploy`，复用 `~/pi-web-deploy/docker-compose.yml` 的 `infinite-canvas` 服务、环境配置和数据卷，不重启 Caddy 或其他服务。部署后检查前后端及本机反代，失败时尝试恢复上一版镜像。

SSH 连接使用 GitHub Secrets：`SSH_HOST`、`SSH_PORT`、`SSH_USERNAME`、`SSH_PASSWORD`、`SSH_FINGERPRINT`。服务器密码不写入仓库；如需重新部署，可在 Actions 页面重新运行对应提交的 CI。

## 文档

完整文档见 [docs/index.md](docs/index.md)：

- [快速开始](docs/overview/quick-start.md)
- [功能介绍](docs/overview/features.md)
- [Docker 部署](docs/overview/docker.md)
- [画布节点操作手册](docs/canvas/canvas-node-manual.md)
- [画布快捷键](docs/canvas/canvas-shortcuts.md)
- [本地开发](docs/backend/local-development.md)

## 说明

- 未登录时画布项目和“我的素材”保存在浏览器本地；登录且账号同步可用时，会同步保存到账号 / 云端。
- 本地直连模式下，AI API Key 保存在浏览器本地，并由前端直接请求 OpenAI 兼容接口，请注意保管。
- 本项目基于 AGPL-3.0 协议开源，详见 [LICENSE](LICENSE)。
