# Commandline options

## trivrost
* `uninstall`: Flag to uninstall the launcher and its bundles on the local machine.
* `debug`: Enable debug log level.
* `skipselfupdate`: Never perform a self-update.
* `roaming`: Cause all files which would be written under `%LOCALAPPDATA%` to be written under `%APPDATA%` instead. (Windows only)
* `build-time`: Print the output of 'date -u "+%Y-%m-%d %H:%M:%S UTC"' from the time the binary was built to standard out and exit immediately.
* `deployment-config`: Override the embedded URL of the deployment-config.
* `accept-install`: Accept install prompt when it is dismissed. Use with `-dismiss-gui-prompts`.
* `accept-uninstall`: Accept uninstall prompt when it is dismissed. Use with `-dismiss-gui-prompts`.
* `dismiss-gui-prompts`: Automatically dismiss GUI prompts.
* `nostreampassing`: Do not relay standard streams to executed commands.
* `extra-env`: Pass all arguments to execution as environment variables. Different variables are separated via `;`. Variable name and value are separated by `=`.

## hasher
Hasher is a utility which generates [bundle info files](walkthrough.md#Bundle-info) given a directory path as an input. Usage:  
`hasher unique_bundle_name path/to/bundle/folder`

## bundown
Bundown is a utility which can download bundles for a desired OS/Arch combination.

* `deployment-config`: Path to a trivrost deployment-config to download bundles for. (default "trivrost/deployment-config.json")
* `os`: GOOS-style name of the operating system to download bundles for. (default "linux")
* `arch`: GOARCH-style name of the architecture to download bundles for. (default "amd64")
* `out`: Path to the directory to download files to. Will be created if missing. (default "bundles")
* `tags`: Only download bundles with one of these comma-separated tags. The special tag `untagged` implicitly exists on all bundles without tags. The special tag `all` will instruct bundown to download all bundles regardless of tags. (default "untagged")
* `pub`: Path to a custom public key file to verify signatures of downloaded bundle info files. (optional)

## validator
Validator is a utility which can validate your deployment-config as well as whether required resources are actually available at the URLs it defines.

Usage: `validator [-skipurlcheck] [-skipjarcheck] path/to/deployment-config.json`

* `-skipurlcheck`: Disable checking of availability of all URLs in the config.
* `-skipjarcheck`: When checking URLs, disable checking of availability of `.jar`-files given to a `java`, `java.exe` or `javaw.exe` binary with the `-jar`-argument in defined commands.
* `-act-as-service`: Validate deployment-config for HTTP GET requests on :80/validate.
* `-port`: Override port for --act-as-service.
