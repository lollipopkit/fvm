<h1 align="center">Go FVM</h1>
<p align="center">Flutter SDK Version Manager written in Go.</p>

<p align="center">
    <img alt="badge-lang" src="https://badgen.net/badge/FVM/0.0.8/cyan">
    <img alt="badge-lang" src="https://badgen.net/badge/Go/1.19/purple">
</p>


English | [简体中文](README_CN.md)

### 🍦 Features
- [x] Manage Multiple Flutter SDKs
- [x] Project Versioning
- [x] **Portable** - No need to install Flutter SDK to install fvm
- [x] Auto config PATH and IDE settings
- [x] Alias. e.g. `fvm dart format .` -> `dart format .`


### 💾 Uasge
1. Install `fvm` through `go install github.com/lollipopkit/fvm@latest` or download the latest release from [HERE](https://github.com/lollipopkit/fvm/releases). It's highly recommended to add `fvm` to PATH.
2. It's highly recommended to set env `FVM_HOME`. If you don't set `FVM_HOME`, fvm will use`$HOME/.fvm` as default.
3. Run `fvm install <version>` to install flutter sdk.
4. Set default global version by `fvm global <version>`.
5. Enter project folder, run `fvm use <version>` to use flutter sdk only in this directory.
6. Config alias by `fvm config alias`. So, you can omit `fvm` and use `dart` or `flutter` directly. eg. `dart format .` instead of `fvm dart format .`.


### 🔖 Attention
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