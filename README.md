## FVM Go
Flutter Version Manager written in Go.

### 🍦 Features
- [x] Manage Multiple Flutter SDKs
- [x] Project Versioning
- [x] **Portable** - No need to install Flutter SDK to install fvm


### 💾 Installation
- Download the latest release from [HERE](https://github.com/lollipopkit/gofvm/releases)
- `go install github.com/lollipopkit/gofvm@latest`


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