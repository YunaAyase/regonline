# 在线报名系统 (RegOnline)

> 基于 Go + Nuxt 4 + SQLite 构建的在线报名管理系统

## 项目简介

RegOnline 是一个前后端分离的在线报名管理系统，支持学生报名、班级管理、数据导出等功能。

| 层级 | 技术 |
|------|------|
| **后端** | Go + Gin + GORM + SQLite（纯 Go 实现，无 CGO） |
| **前端** | Nuxt 4 + Nuxt UI + Tailwind CSS |
| **数据库** | SQLite（无需额外安装） |
| **OCR** | Tesseract-OCR（可选，身份证识别） |

---

## 推荐运行环境

> **强烈建议在 Linux 平台上部署运行**，效果最佳：
> - SQLite 在 Linux 上性能最优
> - Tesseract OCR 识别准确率更高
> - 系统资源占用更低
> - 更适合生产环境长期运行
>
> **推荐配置**：1 核 CPU / 512MB 内存 / 2GB 磁盘

---

## 快速开始

本项目提供 **Docker 一键部署** 和 **本地运行** 两种方案。

> **强烈推荐使用 Docker 方案**——无需安装 Go、Node.js、Tesseract 等任何依赖，三平台操作完全一致。

### 方案一：Docker 部署（推荐，零依赖）

> **Windows 用户**：强烈建议安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/) 并使用容器运行，以获得最佳体验。
>
> **没有 Docker 环境？**：推荐在电脑上安装 Linux 操作系统（虚拟机或物理机）后使用 Docker 方案。

#### 1. 项目根目录创建 `docker-compose.yml`

```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    ports: ["8080:8080"]
    volumes:
      - ./backend/data:/app/data
      - ./photos:/app/photos
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports: ["3000:3000"]
    environment:
      - NUXT_PUBLIC_API_BASE=http://backend:8080/api
      - TZ=Asia/Shanghai
    depends_on: [backend]
    restart: unless-stopped
```

#### 2. 后端 Dockerfile（已提供 `backend/Dockerfile`）

已内置 Tesseract OCR + 中文语言包，无需额外操作。

#### 3. 前端 Dockerfile（已提供 `frontend/Dockerfile`）

使用多阶段构建，最终镜像仅包含运行时产物。

#### 4. 启动

```bash
docker compose up -d
```

访问 `http://localhost:3000` 即可使用。

停止：`docker compose down`
查看日志：`docker compose logs -f`

---

### 方案二：本地运行

#### Linux / macOS

```bash
# 进入项目目录
cd regonline

# 启动后端
cd backend && go run ./cmd/server &

# 启动前端
cd ../frontend && npm install && npm run dev &

# 访问 http://localhost:3000
```

#### Windows

```powershell
# 终端 1：启动后端
cd regonline\backend
go run .\cmd\server

# 终端 2：启动前端
cd regonline\frontend
npm install
npm run dev

# 访问 http://localhost:3000
```

---

## 环境要求

| 组件 | Docker 方案 | 本地方案 |
|------|:-----------:|:--------:|
| **Go** | ❌ 不需要 | ≥ 1.25 |
| **Node.js** | ❌ 不需要 | ≥ 22 |
| **Tesseract OCR** | ❌ 不需要 | ≥ 5.x（可选） |
| **SQLite** | ❌ 不需要 | 不需要（纯 Go 实现） |

> **SQLite 使用 `glebarez/sqlite`（纯 Go 实现），任何平台都无需额外安装 SQLite。**

---

## OCR 身份证识别

OCR 功能用于自动从户口本照片中识别 18 位身份证号。

- **Docker 方案**：已内置 Tesseract + 中文语言包，开箱即用
- **本地方案**：如未安装 Tesseract，自动降级为手动输入

### 本地安装 Tesseract（可选）

**Ubuntu / Debian**
```bash
sudo apt install tesseract-ocr tesseract-ocr-chi-sim
```

**CentOS / RHEL**
```bash
sudo yum install tesseract tesseract-langpack-chi
```

**macOS**
```bash
brew install tesseract tesseract-lang
```

**Windows**
从 [CyaniAgent 服务器分享](https://resources.imikufans.cn/@s/0cIfvRpC) 下载安装包，将 `tesseract.exe` 所在目录加入 PATH。

### 验证安装
```bash
tesseract --version
```

启动后端后，日志显示：
- `Tesseract found, OCR service available` — 可用
- `Tesseract not found in PATH, OCR will fall back to manual input` — 降级为手动输入

---

## 管理员账户

| 用户名 | `admin` |
| 密码 | `admin` |

管理页面：`http://localhost:3000/admin/login`

> 首次登录后建议在「网站设置 → 管理员账号」中修改密码。

---

## 生产部署建议

1. **使用 Docker Compose** 部署，配置 `restart: unless-stopped`
2. **配置 Nginx 反向代理** 处理 HTTPS 证书
3. **定期备份** `backend/data/regonline.db` 和 `photos/` 目录
4. **修改默认管理员密码**
5. **将 `server.mode` 改为 `release`**
6. **配置防火墙**，仅开放 80/443 端口

---

## 项目结构

```
regonline/
├── backend/                  # Go 后端
│   ├── Dockerfile
│   ├── cmd/server/           # 入口文件
│   ├── internal/             # 业务代码
│   │   ├── handler/          # HTTP 处理器
│   │   ├── model/            # 数据模型
│   │   ├── service/          # 业务逻辑
│   │   ├── repository/       # 数据访问层
│   │   └── util/             # 工具函数
│   ├── config.yaml           # 配置文件
│   └── data/                 # SQLite 数据库
├── frontend/                 # Nuxt 4 前端
│   ├── Dockerfile
│   ├── app/
│   │   ├── assets/css/       # 全局样式 & Design Tokens
│   │   ├── components/       # 公共组件
│   │   ├── composables/      # 组合式函数
│   │   ├── pages/            # 页面组件
│   │   └── types/            # TypeScript 类型
│   └── nuxt.config.ts
└── photos/                   # 户口本照片存储
```

---

## 常见问题

**Q: 后端启动后前端无法连接？**
A: 确保 `nuxt.config.ts` 中 `devProxy` 正确指向 `http://localhost:8080/api/`。

**Q: 班级列表为空？**
A: 检查 `config.yaml` 中 `classes.seed_enabled` 是否为 `true`。

**Q: OCR 识别不准确？**
A: 确保照片清晰、光线充足，且包含完整的身份证号码区域。

**Q: 如何备份数据？**
A: 直接复制 `backend/data/regonline.db` 文件即可完整备份。
