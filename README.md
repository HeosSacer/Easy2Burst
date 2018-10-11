<img align="right" width="120" height="120" title="Burst Logo" src="https://raw.githubusercontent.com/PoC-Consortium/Marketing_Resources/master/BURST_LOGO/PNG/icon_blue.png" />

[![Get Support at https://discord.gg/NKXGM6N](https://img.shields.io/badge/join-discord-blue.svg)](https://discordapp.com/invite/NnVBPJX)
[![MIT](https://img.shields.io/badge/license-GPLv3-blue.svg)](LICENSE)

# Easy2Burst

The easy to start wallet for burst.

### Features
- ***Burst Wallet for windows, macOS and linux***
- ***1-click-to-setup*** philosophy
- fully automated updates of components
- account-manager
- using electron with the newest web technologies

### For Collaborators
1. Install golang version > 1.11

2. Install the go-astilectron packages.
``` shell
$ go get -u github.com/asticode/go-astilectron
$ go get -u github.com/asticode/go-astilectron-bundler/...
```

3. Install the Easy2Burst package.
``` shell
$ go get -u github.com/HeosSacer/Easy2Burst
```

4. Use the *go-astilectron-bundler* to create binarys into *.../Easy2Burst/output*.
``` shell
$ cd .../go/github.com/HeosSacer/Easy2Burst
$ ./go/bin/astilectron-bundler.exe -v
```

5. Change *bundler.json* to add a target os.
``` json
{ "app_name": "Easy2Burst",
  "icon_path_windows": "resources/icon.ico",
  "icon_path_darwin": "resources/icon.icns",
  "icon_path_linux": "resources/icon.png",
  "environments": [
        {"arch": "amd64", "os": "windows"},
        {"arch": "amd64", "os": "darwin"},
        {"arch": "386", "os": "linux"}]
}
```
