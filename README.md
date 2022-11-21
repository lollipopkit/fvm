<h1 align="center">Go FVM</h1>
<p align="center">Flutter SDK Version Manager written in Go.</p>

<p align="center">
    <img alt="badge-lang" src="https://badgen.net/badge/FVM/0.0.7/cyan">
    <img alt="badge-lang" src="https://badgen.net/badge/Go/1.19/purple">
</p>


English | [简体中文](README_CN.md)

### 🍦 Features
- [x] Manage Multiple Flutter SDKs
- [x] Project Versioning
- [x] **Portable** - No need to install Flutter SDK to install fvm
- [x] Auto config PATH and IDE settings
- [x] Alias. e.g. `fvm dart format .` or `dart format .`


### 💾 Installation
- Download the latest release from [HERE](https://github.com/lollipopkit/fvm/releases)
- `go install github.com/lollipopkit/fvm@latest`

### 🔖 Attention
- Under normal conditions: 
   1. Install `fvm` through `go install` or download the binary file. 
   2. It's highly recommended to set env `FVM_HOME`. If you don't set `FVM_HOME`, fvm will use `$HOME/.fvm` as default.
   3. Run `fvm install <version>` to install flutter sdk.
   4. Enter project folder, run `fvm use <version>` to use flutter sdk.
- This tool is partly compatible with `fvm-dart`, but there are still some differences. You may need to reconfigure the environment as described above.
- This tool is still in the early stage of development. If you encounter any problems, please submit an issue.

### ⚒️ CLI
```
NAME:
   fvm - Flutter Version Manager written in Go

USAGE:
   fvm [global options] command [command options] [arguments...]

COMMANDS:
   config, c   Config something
   dart, d     Proxy dart commands
   flutter, f  Proxy flutter commands
   global, g   Manage global version of Flutter
   install, i  Install a specific version of Flutter
   list, l     List all installed versions of Flutter
   release, r  List all releases of Flutter
   use, u      Use a specific version of Flutter
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### 📝 License
```
lollipopkit LPGL-3.0
```