# Moments - 极简朋友圈

[![release](https://img.shields.io/badge/release-更新记录-blue)](https://github.com/kingwrcy/moments/releases)
[![docker-release-status](https://img.shields.io/github/actions/workflow/status/kingwrcy/moments/docker-image-release.yml)](https://github.com/kingwrcy/moments/actions/workflows/docker-image-release.yml)
[![docker-pull](https://img.shields.io/docker/pulls/kingwrcy/moments)](https://hub.docker.com/repository/docker/kingwrcy/moments)
[![telegram-group](https://img.shields.io/badge/Telegram-group-blue)](https://t.me/simple_moments)
[![discussion](https://img.shields.io/badge/moments-论坛-blue)](https://discussion.mblog.club)

> 从 v0.2.1 开始，Moments 采用 Golang 重写服务端，包体积更小，功能更强！
>
> 仍需 v0.2.0 版本？[点这里](https://github.com/kingwrcy/moments/blob/master/README.md)

---

## 功能说明

### 用户系统

- 默认管理员：`admin/a123456`，登录后可在后台修改密码
- 多用户模式：可以在后台设置是否允许用户注册

### Memo

- 支持使用 Markdown 语法编写
- 支持修改发布时间，若发布时间晚于当前时间，则对游客不可见
- 支持添加标签，方便分类查找
- 支持上传图片、视频
- 支持引用外部音乐、外部视频、外部链接
- 支持引用豆瓣读书、豆瓣电影
- 支持点赞、评论（区分访客和注册用户），可在后台开启/关闭评论功能
- 支持搜索、置顶、编辑和删除内容
- 支持插入微信 Emoji
- 支持用户主页布局切换（页面相册），可在后台开启/关闭
- 支持内容公开/私密设置
- 支持自定义 css/js 等诸多功能，可在后台配置

### 文件上传

上传到服务器时：

- 支持上传到服务器的 `$UPLOAD_DIR` 目录下
- 上传文件时，会先通过 SHA256 检查文件是否被上传过，若上传过则可以“秒传”
- 上传图片时，会自动创建图片的缩略图，在浏览 Memo 时默认加载缩略图
- 上传文件时，会和 Memo 产生关联，可以在设置中清理无用的文件，无用的文件将会被整理到 `$UPLOAD_DIR/removed` 目录下

上传到 S3 时：

- 支持上传到兼容 S3 的存储服务
- 开启 S3 时，可以设置图片压缩后缀，在浏览 Memo 时，会默认请求带后缀的图片加载缩略图

### 其他功能

以下功能没有详细介绍，请自行探索：

- 友情链接
- 用户管理
- 邮件通知 / 企业微信通知 / 站内消息通知
- OIDC 标准协议用户认证
- AI 润色
- RSS 订阅
- 暗黑模式
- 数据库自动备份，当检测到启动的版本发生变化时，自动备份数据库文件到 `$DB` 目录下

---

## 快速上手

### 应用配置

支持通过环境变量进行配置：

| 变量名         | 说明                   | 默认值                                                                         |
| ----------------- | ---------------------- | ------------------------------------------------------------------------------ |
| PORT           | 监听端口               | 3000                                                                           |
| CORS_ORIGIN    | 允许的跨域 Origin 列表 | 空，多个 Origin 可以使用英文逗号分隔，如 `http://127.0.0.1,http://10.10.10.10` |
| JWT_KEY        | JWT 密钥               | 空，不填写则随机生成，重启后需重新登录                                         |
| DB             | SQLite 数据库存放目录  | /app/data/db.sqlite                                                            |
| UPLOAD_DIR     | 上传文件本地目录       | /app/data/upload                                                               |
| LOG_LEVEL      | 日志级别               | info，可选 debug                                                               |
| ENABLE_SWAGGER | 启用 Swagger 文档      | false，可选 true，通过 `/swagger/index.html` 访问                              |
| ENABLE_SQL_OUTPUT | 启用 SQL 输出          | false，可选 true                                                               |

支持 `.env` 文件加载环境变量，示例：

```env
JWT_KEY=your_secret_key
LOG_LEVEL=info
```

### 使用 Docker Cli 启动

启动容器（需替换 `$JWT_KEY`）：

```bash
docker run -d \
  -e PORT=3000 \
  -e JWT_KEY=$JWT_KEY \
  -p 3000:3000 \
  -v /var/moments:/app/data \
  --name moments \
  kingwrcy/moments:latest
```

镜像标签可选：

- latest：稳定版
- dev：开发版，功能前沿但相对不稳定

### 使用 Docker Compose 启动

```yaml
services:
  moments:
    image: kingwrcy/moments:latest
    container_name: moments
    restart: always
    environment:
      PORT: 3000
      JWT_KEY: $JWT_KEY
    ports:
      - 3000:3000
    volumes:
      - /var/moments:/app/data # 持久化数据到主机的 /var/moments 目录，可以按需修改
```

### 使用可执行文件启动

[下载最新版本](https://github.com/kingwrcy/moments/releases)

示例（Windows 版）：

| 文件名                                         | 说明                         |
| ---------------------------------------------- | ---------------------------- |
| `moments-windows-amd64-x.x.x.exe.zip`          | 压缩包，解压后可直接运行     |
| `moments-windows-amd64-x.x.x.exe-checksum.txt` | `MD5` 校验码，验证文件完整性 |

---

## JWT_KEY 生成方式

方法 1：OpenSSL

```bash
openssl rand -hex 32
```

方法 2：SHA256

```bash
echo $RANDOM | sha256sum
```

方法 3：在线生成（[点这里](https://tool.lu/uuid) 生成 UUID）

---

## 开发

### 依赖环境

后端：Go 1.23.3+
前端：NodeJS 18+，使用 PNPM 管理依赖

VSCode 推荐插件：

- gitlens：Git 扩展
- prettier：代码格式化
- eslint：代码规范检查
- golang：Go 语言支持

### 启动

#### 使用 make（推荐）

后端：

```bash
cd moments
make backend-dev
```

前端（新终端）：

```bash
cd moments
make frontend-install
make frontend-dev
```

#### 手动运行

后端：

```bash
cd moments/backend
go build -ldflags="-X main.version=local -X main.commitId=local" -o ./dist/moments
./dist/moments
```

前端：

```bash
cd moments/front
pnpm install
pnpm run dev
```

启动后访问 `http://localhost:3000`

---

## 其他版本

| 项目                                                            | 演示地址                                                             |
| --------------------------------------------------------------- | -------------------------------------------------------------------- |
| [RandallAnjie/moments](https://github.com/RandallAnjie/moments) | [https://moments.randallanjie.com](https://moments.randallanjie.com) |

---

## 致谢 Contributors

感谢所有贡献者!

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/kingwrcy"><img src="https://avatars.githubusercontent.com/u/1247324?v=4?s=80" width="80px;" alt="kingwrcy"/><br /><sub><b>kingwrcy</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/RandallAnjie"><img src="https://avatars.githubusercontent.com/u/84122428?v=4?s=80" width="80px;" alt="Randall"/><br /><sub><b>Randall</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Jonnyan404"><img src="https://avatars.githubusercontent.com/u/20352705?v=4?s=80" width="80px;" alt="jonny"/><br /><sub><b>jonny</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/akarikun"><img src="https://avatars.githubusercontent.com/u/11921182?v=4?s=80" width="80px;" alt="akari"/><br /><sub><b>akari</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/douseful"><img src="https://avatars.githubusercontent.com/u/52767905?v=4?s=80" width="80px;" alt="yee"/><br /><sub><b>yee</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://www.jschef.com"><img src="https://avatars.githubusercontent.com/u/38160059?v=4?s=80" width="80px;" alt="Chef"/><br /><sub><b>Chef</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://xwsir.cn"><img src="https://avatars.githubusercontent.com/u/17978673?v=4?s=80" width="80px;" alt="小王先森"/><br /><sub><b>小王先森</b></sub></a><br /></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://www.gooth.org"><img src="https://avatars.githubusercontent.com/u/126313?v=4?s=80" width="80px;" alt="Athurg Gooth"/><br /><sub><b>Athurg Gooth</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xuewenG"><img src="https://avatars.githubusercontent.com/u/32838722?v=4?s=80" width="80px;" alt="xuewenG"/><br /><sub><b>xuewenG</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Secretlovez"><img src="https://avatars.githubusercontent.com/u/40491055?v=4?s=80" width="80px;" alt="Secretlovez"/><br /><sub><b>Secretlovez</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/jkjoy"><img src="https://avatars.githubusercontent.com/u/23159890?v=4?s=80" width="80px;" alt="浪子"/><br /><sub><b>浪子</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/lateautumn2"><img src="https://avatars.githubusercontent.com/u/57248164?v=4?s=80" width="80px;" alt="lateautumn2"/><br /><sub><b>lateautumn2</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Jinvic"><img src="https://avatars.githubusercontent.com/u/77521861?v=4?s=80" width="80px;" alt="Jinvic"/><br /><sub><b>Jinvic</b></sub></a><br /></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/dianso"><img src="https://avatars.githubusercontent.com/u/1454808?v=4?s=80" width="80px;" alt="DIANSO"/><br /><sub><b>DIANSO</b></sub></a><br /></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

欢迎贡献！详情见 [all-contributors](https://github.com/all-contributors/all-contributors) 规范。

---

## Star History

[![Star History](https://api.star-history.com/svg?repos=kingwrcy/moments&type=Date)](https://star-history.com/#kingwrcy/moments&Date)

如果你觉得 Moments 还不错，欢迎点个 Star！
