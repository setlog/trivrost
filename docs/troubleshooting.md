# Troubleshooting
Here is a list of challenges to expect and how we were able to manage them.

- trivrost does not start, the system shows an error message saying the launcher cannot be run or is not a valid W32 application.
  - Solution: The downloaded file is corrupt. Tell the user to re-download the file.
- trivrost re-downloads the files every start because the `%LOCALAPPDATA%` folder is purposely deleted every time the user logs on or off (as seen with some Citrix systems).
  - Solution: Add the `-roaming` argument.
- trivrost panics because users may not store or run `.exe`-files outside of `%PROGRAMFILES%`.
  - Solution: If your application files do not contain native executables, use a [system mode](lifecycle.md#system-mode) installation. Otherwise, there is no reason to use trivrost: supply a new `.msi`-installer of your own for every new version of your application.
- Self-update fails because binary is already in use by another process.
  - Solution: Find the process blocking the update. If some mechanism prevents processes from replacing themselves, use a [system mode](lifecycle.md#system-mode) installation.
- trivrost starts, but does not progress beyond 0%.
  - Solution: The logfile will tell you why the download failed. There might be a typo in the URL, a missing file or some firewall blocking the network access. Some networks-firewalls restrict application access to certain IPs. Make sure you have communicated what URLs need to be whitelisted.
  Some firewall appliances manipulate TLS certificates and cause unpredicted failures. trivrost cannot garantuee to work with such broken networks.
- trivrost panics because of no write privileges under `%APPDATA%` or `%LOCALAPPDATA%`.
  - Solution: Whenever this happened, it was cuased by a broken client system. The user's privileges need to be corrected by an administrator.
- A security application alleges the file would not be secure.
  - Solution: Since trivrost can download and execute data, some security applications raise a warning. Use Microsoft Windows executable signing with the `make sign` target and whitelist your company in the security application through an administrator.

## Java specific troubleshooting
- Under Windows, the Java application launched by trivrost always also starts a terminal.
  - Solution: Java provides an extra executable to hide terminal output. Use `javaw.exe` instead of `java.exe`
- Swing or JavaFX fonts under Linux are rendered blurry or with wrong subpixel hinting.
  - Solution: Add the parameter `-Dawt.useSystemAAFontSettings=on` (before `-jar`)
- Swing Application HighDPI scaling issues.
  - Solution: Add the parameter `-Dsun.java2d.dpiaware=false -Dawt.useSystemAAFontSettings=on` (before `-jar`)
