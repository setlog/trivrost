# What files and folders does trivrost create?
* [Itself](glossary.md#trivrost-deployment-artifact).
* All bundles you define, with their contained files, stored in a folder called `bundles`.
* A lock-file `.lock` which is locked using the OS's file system API, to [prevent trivrost from racing with other instances of itself](dev/locking.md).
* A file `.launcher-lock` which contains information on the currently locking trivrost instance.
* A file `.execution-lock` which prevents trivrost from updating bundles while your application is running.
* A `timestamps.json` file used to protect against attacks.
* `.log`-files in a `log`-folder.
* A desktop shortcut to its binary.
* A Start menu shortcut to its binary.
* A Start menu shortcut to its binary with the `--uninstall` parameter.
* The icon you defined. (Linux only; required for shortcut to display icon)

# Where does trivrost write files?
trivrost uses the following user- and platform-specific folders to store files. `<VendorName>` and `<ProductName>` are resolved to their values in `launcher-config.json`.

## Windows
### Default
Deployment artifact:  
`%APPDATA%\<VendorName>\<ProductName>\`

`bundles`-folder, lock-files and `timestamps.json`:  
`%LOCALAPPDATA%\<VendorName>\<ProductName>\`  
If trivrost is started with the [`--roaming` parameter](cmdline.md#trivrost), the path changes to:
`%APPDATA%\<VendorName>\<ProductName>\`

Desktop shortcut:  
`%USERPROFILE%\Desktop\`

Start menu shortcuts:  
`%APPDATA%\Microsoft\Windows\Start Menu\<VendorName>\`  
`%APPDATA%\Microsoft\Windows\Start Menu\<VendorName>\Uninstall\`

Log-files:  
`%LOCALAPPDATA%\Temp\<VendorName>\<ProductName>\log\`  
If trivrost is started with the [`--roaming` parameter](cmdline.md#trivrost), the path changes to:
`%APPDATA%\Temp\<VendorName>\<ProductName>\log\`

### System mode
As **Default**, but with the following changes/additions.

Deployment artifact and `systembundles`-folder:  
`%ProgramFiles%\<VendorName>\<ProductName>\`

Desktop shortcut:  
`%PUBLIC%\Desktop\`

Start menu shortcuts:  
`%ALLUSERSPROFILE%\Microsoft\Windows\Start Menu\<VendorName>\`  
(Uninstall shortcut not installed by system mode-`.msi`)

## MacOS
Deployment artifact, `bundles`-folder, lock-files and `timestamps.json`:  
`$HOME/Library/Application Support/<VendorName>/<ProductName>/`

Desktop shortcut:  
`$HOME/Desktop/`

Start menu shortcuts: N/A

Log-files:  
`$HOME/Library/Caches/<VendorName>/<ProductName>/log/`

## Linux
### Default
Deployment artifact, `bundles`-folder, icon, lock-files and `timestamps.json`:  
`$HOME/.local/share/<VendorName>/<ProductName>/`

Desktop shortcut:  
`$HOME/Desktop/`

Start menu shortcuts:  
`$HOME/.local/share/applications/<VendorName>/<ProductName>/`
`$HOME/.local/share/applications/<VendorName>/<ProductName>/Uninstall/`

Log-files:  
`$HOME/.cache/<VendorName>/<ProductName>/log/`

### XDG
If set/possible, the following [XDG](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)-related configurations will supercede the above:

Desktop shortcut:  
`$(xdg-user-dir DESKTOP)/`

Log-files:  
`$XDG_CACHE_HOME/<VendorName>/<ProductName>/log`
