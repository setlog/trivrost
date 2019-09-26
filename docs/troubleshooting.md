# Troubleshooting
Here is a list of challenges to expect and how we were able to manage them.

- The `%LOCALAPPDATA%` folder is purposely deleted every time the user logs on or off.
  - Solution: Add the `-roaming` argument.
- Users may not store or run `.exe`-files outside of `%PROGRAMFILES%`.
  - Solution: If your application files do not contain native executables, use a [system mode](lifecycle.md#system-mode) installation. Otherwise, there is no reason to use trivrost: supply a new `.msi`-installer of your own for every new version of your application.
- Self-update fails because binary is already in use by another process.
  - Solution: Use a [system mode](lifecycle.md#system-mode) installation.
- trivrost cannot download files.
  - Solution: Make sure you have communicated what URLs need to be whitelisted. Alternatively, an educative e-mail about the normality of port 443 may be in order.
- No write privileges under `%APPDATA%` or `%LOCALAPPDATA%`.
  - Solution: The user's privileges need to be corrected by an administrator.
- A security application alleges the file would not be secure.
  - Solution: trivrost needs to be whitelisted in the security application by an administrator.
