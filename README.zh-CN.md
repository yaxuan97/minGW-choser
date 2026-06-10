# MinGW Chooser

[English](README.md)

一个跨平台命令行工具，用于检测你的系统并推荐最合适的 MinGW-w64 构建版本。不再纠结 i686 还是 x86_64、posix 还是 win32、seh 还是 dwarf 还是 sjlj、ucrt 还是 msvcrt。

## 快速开始

```bash
# 从 Releases 下载最新二进制，或从源码构建：
go install github.com/yourusername/mingw-chooser@latest

# 零参数运行 — 自动检测并推荐：
mingw-chooser

# 离线模式（无网络）：
mingw-chooser --offline

# JSON 输出（便于脚本处理）：
mingw-chooser --json
```

## 功能

运行 `mingw-chooser`（无需任何参数）：

1. **检测** — 你的 CPU 架构、操作系统、是否运行在 WoW64 下
2. **获取** — 从 [mingw-builds](https://github.com/niXman/mingw-builds-binaries) 和 [WinLibs](https://github.com/brechtsanders/winlibs_mingw) 拉取最新可用构建
3. **评分** — 按系统匹配度为每个构建打分（posix > win32，x86_64 上 seh > dwarf，ucrt > msvcrt）
4. **推荐** — 给出最佳构建的下载链接、安装说明，以及每个选择的原因

## 输出示例

```
$ mingw-chooser

Detected system:
  CPU: x86_64 (64-bit)
  OS:  Windows 11 Pro

Recommended build:
  winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64ucrt-14.0.0-r2.7z
  https://github.com/brechtsanders/winlibs_mingw/releases/...

How to install:
  1. Extract the .7z archive to C:\mingw64 (or your preferred location)
  2. Add C:\mingw64\bin to your system PATH
  3. Open a new terminal and run: gcc --version

Why this build?
  x86_64  — your CPU is 64-bit
  posix   — best C++11 std::thread support, wider compatibility
  seh     — optimal exception handling performance on x86_64
  ucrt    — modern Windows C runtime, recommended by Microsoft
```

## 参数

| 参数 | 说明 |
|------|------|
| `--arch <arch>` | 手动指定架构（`x86_64`、`i686`、`aarch64`） |
| `--thread <model>` | 手动指定线程模型（`posix`、`win32`） |
| `--exception <type>` | 手动指定异常处理（`seh`、`dwarf`、`sjlj`） |
| `--crt <type>` | 手动指定 CRT（`ucrt`、`msvcrt`） |
| `--json` | 以 JSON 格式输出 |
| `--offline` | 仅使用内置构建快照（不联网） |
| `--list` | 列出所有匹配的构建，而非仅最佳 |
| `--version` | 显示版本号 |

## JSON 输出

```json
{
  "system": {"os": "windows", "os_version": "Windows 11 Pro", "arch": "x86_64"},
  "recommended": {"name": "winlibs-x86_64-...", "url": "https://..."},
  "alternatives": [...],
  "explanation": [
    {"dimension": "arch", "choice": "x86_64", "reason": "你的 CPU 是 64 位", "manual": false}
  ],
  "warning": null
}
```

## 工作原理

```
main.go
  ├── detect/    平台检测（Windows/Linux/macOS）
  ├── fetch/     GitHub Releases API 客户端（mingw-builds + WinLibs）
  ├── match/     评分引擎 — 过滤、打分、排序
  ├── output/    文本 & JSON 格式化器
  └── builds.json (内嵌)  匹配规则 + 回退快照
```

### 匹配算法

1. **过滤** — 只保留与目标架构匹配的构建
2. **评分** — 每个维度按偏好位置给分（第一选择 = 3 分，第二 = 2 分，第三 = 1 分）
3. **排序** — 高分优先。平局决胜：GCC 版本 → 源优先级
4. **解释** — 每个维度的选择都有说明

### 数据源

| 源 | 优先级 | 说明 |
|----|--------|------|
| [WinLibs](https://winlibs.com/) | 高 | 更新频繁，附带额外库 |
| [mingw-builds](https://github.com/niXman/mingw-builds-binaries) | 基准 | 官方独立构建 |

工具从**两个源**同时获取构建，统一评分后选出最佳。当规格完全相同时，WinLibs 略占优势 — 对 Windows 用户来说它通常更新更及时。

### 边界情况处理

- **WoW64** — 64 位 Windows 上运行 32 位进程？检测真实能力，发出警告，仍推荐 64 位
- **ARM64 Windows** — 若无原生 ARM64 构建，建议使用 x86_64 交叉编译
- **离线** — 回退到内嵌构建快照，并引导用户访问发布页面
- **命名变化** — 若 API 响应无法解析，优雅回退

## 从源码构建

```bash
# 需要 Go 1.23+
git clone https://github.com/yourusername/mingw-chooser.git
cd mingw-chooser
go build -o mingw-chooser .

# 交叉编译
GOOS=windows GOARCH=amd64 go build -o mingw-chooser.exe .
GOOS=linux   GOARCH=amd64 go build -o mingw-chooser .
GOOS=darwin  GOARCH=amd64 go build -o mingw-chooser .
```

零外部依赖 — 仅使用 Go 标准库。

## 项目结构

```
mingw-chooser/
├── main.go              命令行入口
├── config.go            内嵌配置加载器
├── builds.json          匹配规则 + 回退构建
├── detect/              平台检测（构建标签）
├── fetch/               GitHub API 客户端 + 资产解析器
├── match/               评分引擎 + 测试
└── output/              文本 & JSON 格式化器
```

## 为什么要用这个工具

MinGW-w64 的构建变体组合有数十种，选择错误会导致奇怪的编译错误或运行时问题。这个工具：

- **零选择困难** — 直接告诉你最佳选择，并解释原因
- **始终最新** — 从两个主流发布源实时获取，新 GCC 版本自动跟进
- **可手动覆盖** — 如果你知道自己需要什么，用 `--thread win32 --crt msvcrt` 就行
- **可编程** — `--json` 输出便于 CI 脚本和 GUI 包装

## 许可证

MIT
