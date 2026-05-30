# 遂平县青少年活动中心特长班报名系统 — Code Wiki

> **版本**: v1.0.0 | **语言**: Python 3.13+ / HTML5 / JavaScript | **最后更新**: 2026-05-23

---

## 目录

1. [项目概述](#1-项目概述)
2. [系统架构](#2-系统架构)
3. [目录结构](#3-目录结构)
4. [核心模块详解](#4-核心模块详解)
   - [4.1 Web 服务器层](#41-web-服务器层)
   - [4.2 桌面应用层](#42-桌面应用层)
   - [4.3 系统启动与编排层](#43-系统启动与编排层)
   - [4.4 运维与工具层](#44-运维与工具层)
   - [4.5 前端模板层](#45-前端模板层)
5. [数据模型与存储](#5-数据模型与存储)
6. [API 接口文档](#6-api-接口文档)
7. [数据流与交互逻辑](#7-数据流与交互逻辑)
8. [依赖关系](#8-依赖关系)
9. [运行方式](#9-运行方式)
10. [部署方案](#10-部署方案)
11. [安全与限制说明](#11-安全与限制说明)

---

## 1. 项目概述

本项目是一套面向 **遂平县青少年活动中心** 的在线特长班报名管理系统。系统支持家长通过手机扫码在线报名、教务人员通过桌面管理后台查看和删除报名记录，并内置 OCR 身份证识别、二维码扫码访问、年龄校验、班级人数上限控制等核心功能。

### 核心功能

| 功能 | 描述 |
|------|------|
| 在线报名 | 响应式 Web 表单，支持手机/PC，含 11 种特长班可选 |
| OCR 身份证识别 | 上传户口本照片自动识别身份证号（基于 Tesseract） |
| 二维码扫码访问 | 自动生成局域网访问二维码 |
| 年龄校验 | 根据出生日期计算年龄，自动校验是否符合班级年龄限制 |
| 班级人数控制 | 按班级设置最大人数，满额自动拒绝 |
| 查重校验 | 同一学生同班级不可重复报名 |
| 管理后台 | PyQt5 桌面 UI，支持查看/搜索/删除报名记录，修改班级配置 |
| 数据持久化 | JSON 文件本地存储，无需数据库 |

---

## 2. 系统架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                        部署方式选择                                   │
├───────────────────┬──────────────────────┬───────────────────────────┤
│  方式A: 一体机模式  │  方式B: Web Server    │  方式C: 桌面单机模式        │
│  (start_system.py) │  (production_server.py)│ (网上报名系统.py)         │
└───────┬───────────┴─────────┬────────────┴───────────┬───────────────┘
        │                     │                        │
        ▼                     ▼                        ▼
┌───────────────┐   ┌─────────────────┐   ┌───────────────────────┐
│  Flask Web    │   │  Flask Web      │   │  PyQt5 Desktop App    │
│  Server       │   │  Server         │   │  (注册 + 管理一体化)    │
│  (后台线程)    │   │  (前台运行)      │   │                       │
├───────────────┤   ├─────────────────┤   ├───────────────────────┤
│  PyQt5 Admin  │   │  家长手机扫码    │   │  OCR 身份证识别        │
│  UI (主线程)   │   │  访问注册页面    │   │  JSON 本地存储         │
├───────────────┤   └─────────────────┘   └───────────────────────┘
│  QRCode 生成   │
└───────┬───────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────────┐
│                   数据持久化层                                     │
│         registration_data.json  +  photos/ 目录                   │
└─────────────────────────────────────────────────────────────────┘
```

### 三种运行模式对比

| 模式 | 启动方式 | 适用场景 | 特点 |
|------|---------|---------|------|
| 一体机模式 | `启动系统.bat` → `start_system.py` | 单台电脑布展现场 | Web 服务 + 管理 UI + 二维码 全自动启动 |
| Web Server 模式 | `一键布置服务器.bat` → `production_server.py` | 局域网多终端 | 纯 Web 服务，家长手机访问 |
| 桌面单机模式 | `python 网上报名系统.py` | 单人操作 | PyQt5 桌面应用，无需浏览器 |

---

## 3. 目录结构

```
regonline/
├── web_server.py              # Web 服务器（主版本，含 OCR）
├── web_server_fixed.py        # Web 服务器（修复版，更强错误处理）
├── production_server.py       # Web 服务器（生产版，最简稳定）
├── simple_server.py           # Web 服务器（最简测试版）
├── 网上报名系统.py             # PyQt5 桌面报名系统（单机版 main.py）
├── admin_ui.py                # PyQt5 管理后台 UI
├── start_system.py            # 系统编排启动（Web + Admin + QR）
├── start_system_debug.py      # 系统编排启动（调试模式）
├── sever_manager.py           # 命令行服务器管理器
├── diagnose.py                # 系统诊断工具
├── build_installer.py         # 打包工具（PyInstaller + Inno Setup）
├── 生成二维码.py               # 独立二维码生成工具
├── setup.spec                 # PyInstaller 打包配置
│
├── templates/
│   └── registration.html      # Web 报名页面模板
├── registration.html          # 独立报名页面（与 templates/ 内同步）
├── photos/                    # 照片存储目录
│   └── photo_*.jpg            # 上传的户口本照片
├── registration_data.json     # 报名数据持久化文件
├── qrcode.png                 # 生成的二维码图片
├── 报名二维码.png              # 生成的二维码图片（中文名）
├── 报名二维码_大尺寸.png        # 大尺寸二维码
│
├── 一键布置服务器.bat           # Windows 部署脚本
├── 一键打包.bat                # Windows 打包脚本
├── 启动系统.bat                # Windows 快速启动
├── 启动系统_修复后.bat          # Windows 修复版启动
├── 设置固定IP.bat              # 静态 IP 设置指南
├── 配置防火墙.bat              # Windows 防火墙配置
├── install_sever.bat          # 生产环境依赖安装
├── install_service.bat        # Windows 服务安装
└── CODE_WIKI.md               # 本文档
```

---

## 4. 核心模块详解

### 4.1 Web 服务器层

项目包含 **4 个 Flask Web 服务器变体**，功能从简到繁：

#### 4.1.1 `web_server.py` — 主版本 Web 服务器

[web_server.py](file:///d:/code/regonline/web_server.py)

**职责**: 完整的 HTTP API 服务器，提供 Web 报名页面和所有 API 端点。

**关键全局配置**:

```python
CLASSES = {
    '编程班': {'max_students': 30, 'min_age': 7, 'max_age': 16},
    '手工班': {'max_students': 25, 'min_age': 6, 'max_age': 15},
    '国画班': {'max_students': 20, 'min_age': 8, 'max_age': 18},
    '舞蹈班': {'max_students': 25, 'min_age': 5, 'max_age': 14},
    '钢琴班': {'max_students': 15, 'min_age': 6, 'max_age': 18},
    '少儿绘画1班': {'max_students': 25, 'min_age': 5, 'max_age': 12},
    '少儿绘画2班': {'max_students': 25, 'min_age': 5, 'max_age': 12},
    '动漫班': {'max_students': 25, 'min_age': 9, 'max_age': 18},
    '线描班': {'max_students': 20, 'min_age': 7, 'max_age': 16},
    '拼装班': {'max_students': 20, 'min_age': 6, 'max_age': 14},
    '书法班': {'max_students': 25, 'min_age': 7, 'max_age': 17}
}
```

**路由表**:

| 方法 | 路径 | 处理函数 | 描述 |
|------|------|---------|------|
| `GET` | `/` | `index()` | 渲染报名主页 HTML |
| `GET` | `/api/class_info` | `get_class_info()` | 查询指定班级信息（含当前报名人数） |
| `POST` | `/api/submit` | `submit_registration()` | 提交报名表单 |
| `POST` | `/api/recognize_id` | `recognize_id()` | OCR 身份证号识别 |
| `GET` | `/photos/<filename>` | `get_photo()` | 获取上传照片文件 |

**关键函数**:

| 函数 | 签名 | 职责 |
|------|------|------|
| `load_data()` | `() -> dict` | 从 `registration_data.json` 加载数据 |
| `save_data(data)` | `(dict) -> None` | 保存数据到 JSON 文件 |
| `check_duplicate(name, id_number, class_name, registrations)` | `(str, str, str, list) -> bool` | 检查是否重复报名（同姓名或同身份证 + 同班级） |
| `calculate_age(birth_date)` | `(datetime) -> int` | 根据出生日期计算年龄 |

**`submit_registration` 处理流程**:

```
1. 提取表单字段 → 2. 空值校验 → 3. 手机号正则校验 (^1[3-9]\d{9}$)
→ 4. 身份证号正则校验 (^\d{17}[\dXx]$|\d{15}$)
→ 5. 年龄校验 (对比班级 min_age/max_age)
→ 6. 班级人数上限校验
→ 7. 重复报名校验 (姓名+身份证)
→ 8. 照片保存到 photos/ 目录 (MD5 哈希命名)
→ 9. 构建 registration dict → 写入 JSON
→ 10. 返回 {success: True, message: '报名成功！'}
```

#### 4.1.2 `web_server_fixed.py` — 修复版

[web_server_fixed.py](file:///d:/code/regonline/web_server_fixed.py)

**与主版本差异**:
- 添加 Windows UTF-8 编码处理 (Stdout I/O 包装)
- 更强的异常捕获和 `traceback.print_exc()` 调试输出
- 移除了 OCR 的 `pytesseract` 和 `PIL` 顶层导入（改为延迟导入）
- 添加 `@app.errorhandler(404)` 和 `@app.errorhandler(500)` 错误处理
- 启动前检查 `templates/registration.html` 是否存在
- 使用 `debug=True, use_reloader=False` 模式

#### 4.1.3 `production_server.py` — 生产版

[production_server.py](file:///d:/code/regonline/production_server.py)

**与主版本差异**:
- 移除了 OCR 识别 API (`/api/recognize_id`)
- 移除了年龄校验和重复报名校验（简化版，信任前端校验）
- 移除了 `pytesseract` 和 `PIL` 依赖
- 新增 `GET /api/classes` 接口返回全部班级列表
- 编程班名额调整为 10 人
- 使用 `threaded=True` 多线程模式

#### 4.1.4 `simple_server.py` — 测试版

[simple_server.py](file:///d:/code/regonline/simple_server.py)

**职责**: 最简 Flask 服务器，仅用于诊断网络连接是否正常。

提供 `GET /` 返回一个纯 HTML 测试页面，`GET /api/test` 返回 JSON 测试响应。不含任何业务逻辑。

---

### 4.2 桌面应用层

#### 4.2.1 `网上报名系统.py` (main.py) — PyQt5 桌面报名系统

[网上报名系统.py](file:///d:/code/regonline/网上报名系统.py)

**主类**: `RegistrationSystem(QMainWindow)` — 全功能桌面报名应用

**UI 结构** (3 个 Tab 页):

| Tab | 内容 | 职责 |
|-----|------|------|
| 报名登记 | 表单 + 照片上传 + 班级信息 | 录入新报名 |
| 报名管理 | 搜索栏 + 报名列表 + 删除按钮 | 查询/管理已有记录 |
| 班级设置 | 班级配置表 (人数/年龄) + 保存按钮 | 修改班级参数 |

**关键方法**:

| 方法 | 职责 |
|------|------|
| `init_ui()` | 构建 UI，设置 3 个 Tab 页 |
| `setup_registration_tab(tab)` | 构建报名登记页 UI |
| `setup_manage_tab(tab)` | 构建报名管理页 UI |
| `setup_settings_tab(tab)` | 构建班级设置页 UI |
| `submit_registration()` | 提交报名（含校验逻辑） |
| `upload_photo()` | 上传照片并触发 OCR 识别 |
| `recognize_id_number(image_path)` | OCR 识别身份证号 |
| `check_duplicate(name, id_number, class_name)` | 查重校验 |
| `calculate_age(birth_date)` | 年龄计算 |
| `search_registrations()` | 按姓名/班级筛选 |
| `delete_registration(registration)` | 删除报名记录 |
| `save_data()` / `load_data()` | 数据持久化 |
| `save_class_settings()` | 保存班级配置 |

**OCR 路径自动发现**: `setup_tesseract()` 函数按优先级尝试多个安装路径：
```
1. C:\Program Files\Tesseract-OCR\tesseract.exe
2. C:\Program Files (x86)\Tesseract-OCR\tesseract.exe
3. 打包目录下的 tesseract.exe
4. PATH 环境变量中的 tesseract
```

**PyInstaller 兼容**: `resource_path(relative_path)` 通过 `sys._MEIPASS` 支持打包后路径解析。

#### 4.2.2 `admin_ui.py` — PyQt5 管理后台

[admin_ui.py](file:///d:/code/regonline/admin_ui.py)

**主类**: `AdminUI(QMainWindow)` — 独立的教务管理后台

**UI 组成**:
- 顶部标题: "遂平县青少年活动中心报名系统 - 管理后台"
- 统计栏: 总报名人数 + 各班级人数统计
- 报名列表: 12 列 Table (ID/姓名/性别/出生年月/年级/特长班/家长姓名/家长电话/家庭住址/身份证号/报名时间/操作)
- 每行带"删除"按钮

**关键方法**:

| 方法 | 职责 |
|------|------|
| `load_data()` | 加载 JSON 数据 |
| `save_data()` | 保存 JSON 数据 |
| `update_stats()` | 更新统计信息显示 |
| `refresh_table()` | 刷新报名表格 |
| `delete_registration(registration)` | 确认后删除报名记录 |

**入口函数**: `run_admin()` 创建 `QApplication` 并启动窗口。

---

### 4.3 系统启动与编排层

#### 4.3.1 `start_system.py` — 主启动器

[start_system.py](file:///d:/code/regonline/start_system.py)

**职责**: 编排整个系统的启动流程。

**执行流程**:
```
1. get_local_ip()           → 获取本机局域网 IP
2. generate_qrcode(url)     → 生成二维码图片并打开
3. start_web_server()       → 后台线程启动 Flask
4. start_admin_ui()         → 主线程启动 PyQt5 管理 UI
```

**关键函数**:

```python
def get_local_ip():
    """通过 UDP 连接 8.8.8.8 获取本机实际局域网 IP"""
    
def generate_qrcode(url, filename='qrcode.png'):
    """使用 qrcode 库生成二维码 PNG"""
    
def start_web_server():
    """从 web_server 模块启动 Flask app (debug=False)"""
    
def start_admin_ui():
    """等待 2 秒后启动 PyQt5 管理界面"""
```

#### 4.3.2 `start_system_debug.py` — 调试启动器

[start_system_debug.py](file:///d:/code/regonline/start_system_debug.py)

**与主启动器差异**:
- 启动前检查必要文件 (`web_server.py`, `admin_ui.py`, `templates/registration.html`)
- Web 服务器使用 `debug=True, threaded=True`
- 启动后自动执行 `test_server()` 验证服务器可达性
- 所有异常均 `traceback.print_exc()` 输出完整堆栈
- 管理界面启动失败后等待用户按键退出

---

### 4.4 运维与工具层

#### 4.4.1 `sever_manager.py` — 服务器管理器

[sever_manager.py](file:///d:/code/regonline/sever_manager.py)

**主类**: `ServerManager` — CLI 菜单式服务器管理

| 方法 | 职责 |
|------|------|
| `start()` | 通过 `subprocess.Popen` 启动 `production_server.py` |
| `stop()` | 先 `terminate()`，2 秒后若未退出则 `kill()` |
| `restart()` | stop → wait 2s → start |
| `status()` | 通过 `requests.get('http://127.0.0.1:5000')` 检查存活 |

**交互菜单**: 5 个选项 (启动/停止/重启/查看状态/退出)

#### 4.4.2 `diagnose.py` — 系统诊断工具

[diagnose.py](file:///d:/code/regonline/diagnose.py)

**检查项**:

| 序号 | 检查项 | 说明 |
|------|--------|------|
| 1 | Python 环境 | 版本号 + 路径 |
| 2 | 依赖包 | flask, flask_cors, PIL, pytesseract, PyQt5, qrcode |
| 3 | 必要文件 | web_server.py, admin_ui.py, templates/registration.html |
| 4 | 5000 端口 | socket 连接测试是否被占用 |

#### 4.4.3 `build_installer.py` — 打包工具

[build_installer.py](file:///d:/code/regonline/build_installer.py)

**打包流程**:

```
1. clean_build()                    → 清理旧 build/dist 目录
2. create_version_file()            → 生成版本信息
3. create_readme() / create_license() → 生成文档
4. build_with_pyinstaller()         → PyInstaller 打包为单 exe
5. create_installer_script()        → 生成 Inno Setup .iss 脚本
6. build_installer()                → 编译安装包
7. create_portable_version()        → 生成便携版 ZIP
```

**PyInstaller 参数**: `--onefile --windowed --name=ABCDEFG报名系统 --add-data=photos;photos --hidden-import=PIL,pytesseract,PyQt5`

#### 4.4.4 `生成二维码.py` — 独立二维码生成器

[生成二维码.py](file:///d:/code/regonline/生成二维码.py)

**职责**: 独立获取本机 IP 并生成两种尺寸的二维码 (`报名二维码.png` + `报名二维码_大尺寸.png`)。

---

### 4.5 前端模板层

#### 4.5.1 `templates/registration.html` — 报名页面

[templates/registration.html](file:///d:/code/regonline/templates/registration.html)

**技术栈**: 纯 HTML5 + CSS3 + Vanilla JavaScript (无框架)

**设计特征**:
- 渐变背景 (`linear-gradient(135deg, #667eea, #764ba2)`)
- 圆角卡片式容器 (max-width: 600px, 居中)
- 响应式布局 (支持 Mobile/Tablet/Desktop)
- CSS 动画 (spinner 加载动画, 按钮 hover 效果)

**表单字段**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| 学生姓名 | text | ✓ | - |
| 性别 | select (男/女) | ✓ | - |
| 出生年月 | date | ✓ | 限制 5-14 岁范围 |
| 年级 | select (幼儿园~初三) | ✓ | - |
| 特长班 | select (动态加载) | ✓ | 随选择实时显示班级信息 |
| 家长姓名 | text | ✓ | - |
| 家长电话 | tel | ✓ | 前端未校验（后端校验） |
| 家庭住址 | textarea | ✓ | - |
| 户口本照片 | file (image/*) | ✓ | 点击区域上传，支持预览 |
| 身份证号 | hidden + manual input | ✓ | 照片上传后显示输入框 |

**JavaScript 交互逻辑**:

```
1. DOMContentLoaded → loadClasses() → fetch /api/classes → 填充班级下拉
2. classSelect.onchange → fetch /api/class_info → 显示人数/年龄信息
3. photoInput.onchange → FileReader 预览 + 显示身份证手动输入框
4. form.onsubmit → FormData → fetch POST /api/submit → 显示结果
```

**日期范围限制**: 出生日期自动限制在 `[today-14年, today-5年]` 范围。

---

## 5. 数据模型与存储

### 5.1 存储机制

使用**单一 JSON 文件** + **文件系统照片存储**，无需数据库。

```
registration_data.json          ← 结构化报名数据
photos/                          ← 上传的照片文件
    └── {timestamp}_{md5hash}.jpg
```

### 5.2 JSON 数据结构

```json
{
  "registrations": [
    {
      "id": 1,
      "name": "张三",
      "gender": "男",
      "birth_date": "2013-08-12",
      "grade": "五年级",
      "class": "编程班",
      "parent_name": "张某某",
      "parent_phone": "13800138000",
      "address": "XX小区X栋",
      "id_number": "411728201308120053",
      "photo_path": "photos/20260518_164440_abc12345.jpg",
      "registration_time": "2026-05-18 16:44:40"
    }
  ],
  "classes": {
    "编程班": {
      "max_students": 30,
      "min_age": 7,
      "max_age": 16
    }
  }
}
```

### 5.3 班级配置模型

```python
{
    '班级名称': {
        'max_students': int,   # 最大人数上限 (1-100)
        'min_age': int,        # 最小年龄 (3-20)
        'max_age': int         # 最大年龄 (5-25)
    }
}
```

**11 个预设班级**: 编程班、手工班、国画班、舞蹈班、钢琴班、少儿绘画1班、少儿绘画2班、动漫班、线描班、拼装班、书法班。

### 5.4 照片命名规则

```
{timestamp}_{md5(student_name)[:8]}.jpg
```

示例: `20260518_164440_abc12345.jpg` — 时间戳确保唯一性，MD5 前 8 位关联学生姓名。

---

## 6. API 接口文档

### 6.1 基础信息

| 属性 | 值 |
|------|-----|
| 协议 | HTTP/1.1 |
| 数据格式 | JSON (API) / FormData (提交) |
| 编码 | UTF-8 |
| 默认端口 | 5000 |
| CORS | 已启用 (`flask-cors`) |

### 6.2 接口列表

#### `GET /` — 报名主页

返回渲染后的 HTML 报名页面。

**响应**: `text/html`

---

#### `GET /api/classes` — 获取班级列表

> 仅 `production_server.py` 提供

**响应示例**:
```json
{
  "编程班": {"max_students": 30, "min_age": 7, "max_age": 16},
  "手工班": {"max_students": 25, "min_age": 6, "max_age": 15}
}
```

---

#### `GET /api/class_info?class_name={name}` — 获取班级详情

**参数**: `class_name` (query string, required)

**成功响应 (200)**:
```json
{
  "max_students": 30,
  "min_age": 7,
  "max_age": 16,
  "current_count": 6
}
```

**错误响应 (404)**:
```json
{"error": "班级不存在"}
```

---

#### `POST /api/submit` — 提交报名

**Content-Type**: `multipart/form-data`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | ✓ | 学生姓名 |
| gender | string | ✓ | 男/女 |
| birth_date | string | ✓ | 格式: YYYY-MM-DD |
| grade | string | ✓ | 幼儿园~初三 |
| class_name | string | ✓ | 班级名称 |
| parent_name | string | ✓ | 家长姓名 |
| parent_phone | string | ✓ | 11 位手机号 |
| address | string | ✓ | 家庭住址 |
| id_number | string | ✓ | 15 或 18 位身份证号 |
| photo | file | - | 户口本照片 (jpg/png) |

**成功响应**:
```json
{"success": true, "message": "报名成功！"}
```

**失败响应**:
```json
{"success": false, "message": "请填写所有必填项"}
```

---

#### `POST /api/recognize_id` — OCR 识别身份证号

> 仅 `web_server.py` / `web_server_fixed.py` 提供

**Content-Type**: `multipart/form-data`

**参数**: `photo` (file, required) — 带有身份证号的图片

**成功响应**:
```json
{"success": true, "id_number": "411728201308120053"}
```

---

#### `GET /photos/<filename>` — 获取照片

返回 `photos/` 目录下的文件。用于报名提交后查看上传的照片。

---

## 7. 数据流与交互逻辑

### 7.1 一体机模式数据流

```
┌─────────────┐     HTTP GET /     ┌──────────────┐
│  家长手机     │ ─────────────────→ │  Flask Web    │
│  (浏览器)     │ ←─── HTML 页面 ── │  Server       │
└──────┬──────┘                    └──────┬───────┘
       │ POST /api/submit                │
       │ (FormData + photo)              │ read/write
       └─────────────────────────────────┤
                                         ▼
                                  ┌──────────────┐
                                  │  JSON 文件    │
                                  │  + photos/    │
                                  └──────┬───────┘
                                         │ read
                                         ▼
                                  ┌──────────────┐
                                  │  PyQt5 Admin  │
                                  │  UI (桌面)     │
                                  └──────────────┘
```

### 7.2 报名提交流程 (完整)

```
家长填写表单
    │
    ├─ 选择班级 → GET /api/class_info → 显示人数/年龄信息
    ├─ 上传照片 → FileReader 预览 → 手动输入身份证号
    │
    ▼
点击"立即报名"
    │
    ▼
JavaScript FormData 构造 → fetch POST /api/submit
    │
    ▼
Flask submit_registration()
    │
    ├─ [1] 提取字段 .strip()
    ├─ [2] all() 空值校验
    ├─ [3] re.match 手机号校验
    ├─ [4] re.match 身份证号校验
    ├─ [5] datetime.strptime 解析出生日期
    ├─ [6] calculate_age() + 班级年龄范围校验
    ├─ [7] 班级人数上限校验 (遍历统计)
    ├─ [8] check_duplicate() 查重
    ├─ [9] 保存照片到 photos/ (md5 哈希命名)
    ├─ [10] 构建 registration dict
    ├─ [11] json.dump 写入文件
    │
    ▼
返回 {"success": true} → 前端显示"报名成功！" → 2 秒后清空表单
```

### 7.3 查重逻辑

```python
def check_duplicate(name, id_number, class_name, registrations):
    for reg in registrations:
        if (reg['name'] == name or reg['id_number'] == id_number) \
           and reg['class'] == class_name:
            return True
    return False
```

规则: **(姓名相同 或 身份证号相同) 且 班级相同** → 视为重复

---

## 8. 依赖关系

### 8.1 Python 包依赖

| 包名 | 版本要求 | 用途 | 必需? |
|------|---------|------|-------|
| `flask` | ≥2.0 | Web 框架 | ✓ 核心 |
| `flask-cors` | ≥3.0 | 跨域支持 | ✓ 核心 |
| `PyQt5` | ≥5.15 | 桌面管理 UI | ✓ 核心 |
| `Pillow (PIL)` | ≥9.0 | 图片处理 + OCR 图像预处理 | ✓ 核心 |
| `pytesseract` | ≥0.3 | OCR 身份证识别 | 可选 |
| `qrcode` | ≥7.0 | 生成访问二维码 | 可选 |
| `waitress` | ≥2.0 | 生产级 WSGI 服务器 | 可选 |
| `requests` | ≥2.0 | 服务器状态检测 | 可选 |
| `pyinstaller` | ≥5.0 | 打包为 .exe | 仅开发 |

### 8.2 系统依赖

| 依赖项 | 用途 | 说明 |
|--------|------|------|
| Tesseract OCR | 身份证识别引擎 | 需单独安装并配置中文语言包 (`chi_sim`) |
| Inno Setup 6 | 构建 Windows 安装包 | 仅打包时需要 |
| Windows 防火墙 | 放行 5000 端口 | 局域网访问必需 |

### 8.3 模块依赖图

```
start_system.py ─────import───→ web_server.py
      │                              │
      ├──import──→ admin_ui.py       ├──import──→ flask
      │                              ├──import──→ flask_cors
      ├──import──→ qrcode            ├──import──→ json
      │                              ├──import──→ re
      └──import──→ socket            ├──import──→ hashlib
                                     ├──import──→ pytesseract
                                     └──import──→ PIL

网上报名系统.py (main.py) ──import──→ PyQt5, PIL, pytesseract, json, re
admin_ui.py ───────────────import──→ PyQt5, json

production_server.py ─────import──→ flask, flask_cors, json, re, hashlib
sever_manager.py ─────────import──→ subprocess, requests, socket
diagnose.py ──────────────import──→ (无外部依赖，仅 sys/os)
build_installer.py ───────import──→ subprocess, shutil
```

---

## 9. 运行方式

### 9.1 环境要求

- **操作系统**: Windows 7/8/10/11 (64位)
- **Python**: 3.8+
- **Tesseract OCR**: 可选（用于自动识别身份证）
- **网络**: 手机和电脑在同一局域网

### 9.2 快速启动 (推荐)

```batch
# 一键安装依赖 + 启动
双击: 启动系统_修复后.bat
```

该脚本自动执行:
1. `pip install flask flask-cors Pillow pytesseract PyQt5 qrcode requests`
2. 创建 `templates/` 和 `photos/` 目录
3. 运行 `diagnose.py` 诊断
4. 启动 `web_server_fixed.py`

### 9.3 手动启动

```powershell
# 1. 安装依赖
pip install flask flask-cors Pillow pytesseract PyQt5 qrcode

# 2. 启动一体机模式 (Web + Admin + QR)
python start_system.py

# 3. 或仅启动 Web 服务器
python web_server.py

# 4. 或仅启动管理后台
python admin_ui.py

# 5. 或启动桌面单机版
python 网上报名系统.py
```

### 9.4 服务部署

```batch
# 一键服务器部署
双击: 一键布置服务器.bat

# 配置防火墙 (允许外部访问)
双击: 配置防火墙.bat

# 设为 Windows 服务 (开机自启)
双击: install_service.bat
```

### 9.5 打包分发

```batch
双击: 一键打包.bat
```

产物:
- `dist/ABCDEFG报名系统.exe` — 单文件可执行程序
- `installer/ABCDEFG报名系统_Setup.exe` — Windows 安装包
- `ABCDEFG报名系统_便携版.zip` — 绿色免安装版

### 9.6 系统诊断

```powershell
python diagnose.py
```

检查 Python 环境、依赖包、必要文件、5000 端口状态。

---

## 10. 部署方案

### 10.1 局域网报名场景 (推荐)

```
┌──────────────────────────────────────────┐
│  电脑 (Windows, 运行 Flask + Admin UI)    │
│  IP: 192.168.1.100:5000                  │
│  连接到 WiFi                              │
└──────────────┬───────────────────────────┘
               │ 局域网 (同一 WiFi)
    ┌──────────┼──────────┐
    ▼          ▼          ▼
 ┌──────┐  ┌──────┐  ┌──────┐
 │手机 1│  │手机 2│  │手机 3│  ... (家长扫码报名)
 └──────┘  └──────┘  └──────┘
```

**配置步骤**:
1. 电脑连接 WiFi，设置固定 IP (参考 `设置固定IP.bat`)
2. 运行 `配置防火墙.bat` 放行 5000 端口
3. 运行 `一键布置服务器.bat` 启动服务
4. 展示生成的二维码供家长扫码

### 10.2 单机报名场景

直接运行 `python 网上报名系统.py`，在一台电脑上完成报名登记和管理。

### 10.3 打包分发场景

使用 `一键打包.bat` 生成安装包，分发到其他电脑使用。

---

## 11. 安全与限制说明

### 11.1 已知限制

| 限制 | 说明 | 影响 |
|------|------|------|
| 无数据库 | 使用 JSON 文件存储 | 不支持并发写入，高并发场景可能数据丢失 |
| 无认证机制 | 管理后台无登录验证 | 任何能访问 5000 端口的用户都可以提交报名 |
| 无 HTTPS | HTTP 明文传输 | 身份证号等敏感信息在局域网内明文传输 |
| 单文件存储 | 单一 JSON 文件 | 文件损坏则所有数据丢失 |
| 无分页 | 报名列表全量加载 | 大量报名记录时性能下降 |
| OCR 不稳定 | 依赖 Tesseract 中文识别 | 需手动输入身份证号作为备选方案 |

### 11.2 数据安全建议

1. **定期备份** `registration_data.json` 和 `photos/` 目录
2. **不在公网暴露** 此服务，仅在局域网内使用
3. **不要将 `registration_data.json` 提交到公开仓库** (含真实身份证号和手机号)
4. 建议在班级报名结束后**清理 JSON 数据文件中的敏感字段**

---

> **文档维护**: 本文档基于 2026-05-23 代码状态生成，项目持续迭代中。  
> **联系方式**: 遂平县青少年活动中心  
> **版权**: 编程班 © 2024-2026 All Rights Reserved.