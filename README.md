## Go FVM
Flutter SDK Version Manager written in Go.

### 🍦 Features
- [x] Manage Multiple Flutter SDKs
- [x] Project Versioning
- [x] **Portable** - No need to install Flutter SDK to install fvm
- [x] Auto config PATH and IDE settings


### 💾 Installation
- Download the latest release from [HERE](https://github.com/lollipopkit/gofvm/releases)
- `go install github.com/lollipopkit/gofvm@latest`

### 🔖 Attention
- Under normal conditions: 
   1. Install `fvm` through `go install` or download the binary file. 
   2. It's highly recommended to set env `FVM_HOME`. If you don't set `FVM_HOME`, fvm will use `$HOME/.fvm` as default.
   3. Run `fvm install <version>` to install flutter sdk.
   4. Enter project folder, run `fvm use <version>` to use flutter sdk.
- This tool is largely compatible with `fvm-dart`, but there are still some differences. You may need to reconfigure the environment as described above.
- This tool is still in the early stage of development. If you encounter any problems, please submit an issue.

### ⚒️ CLI
```
NAME:
   gofvm - Flutter Version Manager written in Go

USAGE:
   gofvm [global options] command [command options] [arguments...]

COMMANDS:
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