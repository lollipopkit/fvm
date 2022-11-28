<h1 align="center">Go FVM</h1>
<p align="center">Flutter SDK 版本管理工具 - Golang实现</p>

<p align="center">
    <img alt="badge-lang" src="https://badgen.net/badge/FVM/0.1.0/cyan">
    <img alt="badge-lang" src="https://badgen.net/badge/Go/1.19/purple">
</p>


简体中文 | [English](README_en.md)

### 🍦 特点
- [x] 管理多个Flutter SDK
- [x] 支持为不同项目设置不同的Flutter SDK版本
- [x] **便携** - 不需要安装Flutter SDK即可安装fvm
- [x] 自动配置环境和IDE设置
- [x] 支持设置别名。例如 `fvm dart format .` -> `dart format .`


### 💾 使用
1. 安装 `fvm`。
   - 一键安装（需要Go环境）： `go install github.com/lollipopkit/fvm@latest`
   - 从 [这里](https://github.com/lollipopkit/fvm/releases) 下载最新二进制文件。 记得将 `fvm` 添加到环境变量中。
2. 强烈建议设置环境 `FVM_HOME`。 如果你没有设置环境 `FVM_HOME`，fvm 会使用默认的路径 `$HOME/.fvm` 。
3. 运行 `fvm install <version>` 来安装某个版本的SDK。
4. 通过 `fvm global <version>` 来设定全局默认版本。
5. 进入项目文件夹, 运行 `fvm use <version>` 来设定项目特定版本（仅在该文件夹生效）。
6. 通过 `fvm config alias` 来设置别名。 这样，你可以省略 `fvm` 命令前缀。 例如：格式化项目，使用 `dart format .` 而不是 `fvm dart format .` 命令。

### 🔖 注意
- 这个工具与 `fvm-dart` 部分兼容, 但存在许多差异。 你可能需要根据上述步骤重新配置环境。
- 这个工具目前处于开发阶段，可能会有一些不稳定的地方。 请谨慎使用，遇到问题请发送Issue。

### ⚒️ 命令行
```
名称:
   fvm - Flutter SDK 版本管理工具 - Golang实现

使用:
   fvm [全局选项] 命令 [选项] [参数...]

命令:
   config, c   更改某些配置
   dart, d     代理dart命令
   flutter, f  代理flutter命令
   global, g   设置全局默认Flutter SDK版本
   install, i  安装特定版本的Flutter SDK
   list, l     列出所有已安装的Flutter SDK版本
   release, r  列出所有可用的Flutter SDK版本
   use, u      在某个文件夹中设置使用特定版本的Flutter SDK
   help, h     显示帮助信息

全局选项:
   --help, -h  显示帮助 (默认: 否)
```

### 📝 License
```
lollipopkit LPGL-3.0
```