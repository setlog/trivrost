# What files does trivrost create?
* Its executable file. (On MacOS, actually a `.app`-folder posing as an application, which is canon in the Mac world)
* All files of all bundles you define, stored in a folder called `bundles`.
* A lock-file `.launcher-lock` which prevents trivrost from racing with itself.
* A lock-file `.execution-lock` which prevents trivrost from updating bundles while your application is running.
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
Executable:  
`%APPDATA%\<VendorName>\<ProductName>\`

`bundles`-folder, lock-files and `timestamps.json`:  
`%LOCALAPPDATA%\<VendorName>\<ProductName>\`

Desktop shortcut:  
`%USERPROFILE%\Desktop\`

Start menu shortcuts:  
`%APPDATA%\Microsoft\Windows\Start Menu\<VendorName>\`  
`%APPDATA%\Microsoft\Windows\Start Menu\<VendorName>\Uninstall\`

Log-files:  
`%LOCALAPPDATA%\Temp\<VendorName>\<ProductName>\log\`

### System mode
As **Default**, but with the following changes/additions.

Executable and `systembundles`-folder:  
`%ProgramFiles%\<VendorName>\<ProductName>\`

Desktop shortcut:  
`%PUBLIC%\Desktop\`

Start menu shortcuts:  
`%ALLUSERSPROFILE%\Microsoft\Windows\Start Menu\<VendorName>\`  
(Uninstall shortcut not installed by system mode-`.msi`)

## MacOS
Executable, `bundles`-folder, lock-files and `timestamps.json`:  
`$HOME/Library/Application Support/<VendorName>/<ProductName>/`

Desktop shortcut:  
`$HOME/Desktop/`

Start menu shortcuts: N/A

Log-files:  
`$HOME/Library/Caches/<VendorName>/<ProductName>/log/`

## Linux
### Default
Executable, `bundles`-folder, icon, lock-files and `timestamps.json`:  
`$HOME/.local/share/<VendorName>/<ProductName>/`

Desktop shortcut:  
`$HOME/Desktop/`

Start menu shortcuts:  
`$HOME/.local/share/applications/<VendorName>/<ProductName>/`
`$HOME/.local/share/applications/<VendorName>/<ProductName>/Uninstall/`

Log-files:  
`$HOME/.cache/<VendorName>/<ProductName>/log/`

### XDG
If set/possible, the following [XDG](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)-related configurations will be used:

Desktop shortcut:  
`$(xdg-user-dir DESKTOP)/`

Log-files:  
`$XDG_CACHE_HOME/<VendorName>/<ProductName>/log`
