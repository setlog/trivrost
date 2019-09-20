# Troubleshooting
Here is a list of challenges to expect and how we were able to manage them.

- No write privileges under `%APPDATA%` or `%LOCALAPPDATA%`.
  - Solution: Tell their IT that they need to fix the computer.
- Executable file gets deleted randomly.
  - Solution: Not worth anyone's time. Tell them to apply for a new machine.
- `<security application>` alleges the file would not be secure.
  - Solution: Tell their IT to whitelist the file in `<security application>`.
- trivrost cannot download files.
  - Solution: An educative e-mail about the normality of port 443 is in order.
- For reasons, the `%LOCALAPPDATA%` folder is purposely deleted every time the user logs on or off.
  - Solution: Suggest to add the `-roaming` argument.
- For reasons, users may not store or run `.exe` files outside of `%PROGRAMFILES%`.
  - Solution: Provide a [system mode](glossary.md#System-mode) installation.
- Citrix
  - Solution: Provide a [system mode](glossary.md#System-mode) installation.

The list will go on. We recommend patience, resilience, politeness and good humor.
