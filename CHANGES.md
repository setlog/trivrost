# Release-Changelog

## TBD (TBD)
### Changes
* Every TLS Certificate fingerprint will be logged once with the host name it has first been seen on.
* DWARF symbols are now stripped from the trivrost binary to reduce file size. This can save a few bytes on some platforms.
* The binary is now compressed with UPX when using `make compress`. Reduces the final filesize to less than 50%.
* Shorter log-output for proxy detection. Reduces average size of the log output by 5–15%.
* Shorter log-output for HTTP errors, reduces size of log output by a few percent.
* Update dependencies to recent versions: gopsutils, testify, gojsonschema, logrus, prometheus/client_golang, go_ieproxy, fatih/color, golang/x/sys, golang/x/net
* Do not hide the download speed label, even if the speed is zero.
* The download-speed label now shows a 3 second average.
* The same download-related log messages will now be printed at most 5 times (with information about this limit in the last message).
* `hasher` will no longer blindly overwrite an existing bundleinfo.json but instead error out.
* `hasher` will now exit with an error when the `pathToHash` has no files to hash.
* `timestamps.json` is ignored, if it is corrupt.
### Features
* trivrost will log the progress of downloads if the connection was interrupted for any reason.
### Fixes
* `hasher` will no longer create a directory if a non-existing one is passed as an argument.
* trivrost will no longer attempt to repeat range requests to a host after it has failed to conformly respond while displaying the confusing message `Taking longer than usual: HTTP Status 200` and will now fail immediately in such cases instead.
* trivrost will no longer fail to comply with HTTP 2 strictly using lower-case HTTP Header names. This had been caused by methods of `http.Header` still being oriented around HTTP 1 canonical header names due to Go's backwards compatibility promise.
* Instead of always showing 'Cannot reach server' to the user, show more precise/useful messages on connection issues.

## 1.4.6 (2021-01-25)
### Fixes
* Windows binary signing: Use RFC-3161 timestamp server with sha 256 config. SHA-1 ciphers are considered deprecated. Nothing should change for the enduser.

## 1.4.5 (2021-01-04)
### Fixes
* Switch timestamp server for signing from Verisign to Globalsign.

## 1.4.4 (2020-01-17)
### Changes
* Validator's `/validate` endpoint now always answers with status code 200 OK. Endpoint `/metrics` should be tested for value of `trivrost_validation_ok` as a healthcheck instead.

## 1.4.3 (2020-01-17)
### Fixes
* Fix validator displaying wrong URL when using `configurl`.

## 1.4.2 (2020-01-17)
### Fixes
* Fix validator saying `configurl` is not set when it is.

## 1.4.1 (2020-01-17)
### Features
* `cmd/validator`: Run as service for use as healthcheck when started with `--act-as-service` command line argument; send HTTP GET requests to `:80/validate`. Customize port with `--port`. Override deployment-config URL with `configurl` query parameter.

## 1.4.0 (2019-11-22)
### Features
* `cmd/validator` will now check that command binary's *will be* downloaded per platform, where previously it would only check if they are available under their respective bundle URL.
* Added new field `IsUpdateMandatory` to [deployment-config](docs/deployment-config.md), which allows to deny the user the launch of the application when system bundles require changes.
### Fixes
* Fixed changes to the launcher icon not being applied after a self-update on Linux.
* `cmd/validator`: fix only the first missing bundle URL being reported.
* Fix missing `LauncherUpdate` in deployment-config leading to error when documentation explicitly allows it.
* Fix false positives in log warning about program name on disk having diverged from what is configured in the embedded launcher-config on Windows and MacOS.
### Changes
* trivrost now allows users to continue past failed self-updates as well as omitted changes to system bundles in [system mode](docs/lifecycle.md#system-mode) installations, unless [further configuration](docs/deployment-config.md) is made; see the new field `IsUpdateMandatory` in deployment-config.

## 1.3.5 (2019-11-11)
### Fixes
* Ignore `--psn_X_YYYYY` argument which is appended by MacOS to the command line arguments when the programm was launched through a Gatekeeper context. Regression found in v1.3.2, v1.3.3, v1.3.4.

## 1.3.4 (2019-10-29)
### Fixes
* Fix failure to launch on Windows if trivrost binary is located on a mounted volume for which no unique path exists because it has no drive letter assigned. See [this Go issue](https://github.com/golang/go/issues/20506#issuecomment-318514515) for details.

## 1.3.3 (2019-10-18)
### Features
* Users can now always open the log folder through a `Show logs...`-link in the lower-right corner of trivrost's progress window.
* trivrost can now be closed by the user during the `DetermineLocalLauncherVersion` and `DetermineLocalBundleVersions` stages, where it would previously block until the stage completed.
### Changes
* trivrost will now retry the launch of programs in the execution phase every three seconds on error.
* When informing about bad command line arguments via GUI, display the most common arguments as a hint.
* Build-time-difference of the launcher is now compared using the time package instead of a plain string comparison, bringing more robustness should the string ever be garbled.

## 1.3.2 (2019-10-14)
### Features
* trivrost will now handle the following signals, logging the stack trace of all goroutines before terminating: `SIGINT`, `SIGQUIT`, `SIGABRT`, `SIGTERM`, `SIGHUP`.
### Fixes
* `cmd/signer` can now parse PKCS1-formatted private keys in the same way that `scripts/signer` could.
### Changes
* trivrost will now inform about bad command line arguments via GUI.

## 1.3.1 (2019-10-01)
### Changes
* Missing `Content-Length`-headers will now cause size-checks to be skipped instead of failing. In those cases, bad files will only be detected by SHA-mismatch as soon as they have been downloaded entirely.

## 1.3.0 (2019-09-27)
### Features
* trivrost will now hint the final size of bundle files to the operating system, eliminating file fragmentation.
* Added field `IgnoreLauncherBundleInfoHashes` to launcher-config, where you can specify the SHA-256 hash values of bundleinfo files to be ignored by trivrost when it checks for a self-update. This behaviour can be used to hand out specialized builds to specific users for hotfixing purposes without having to worry about the need to add (and later remove) the `-skipselfupdate` argument.
* Added command line argument `-deployment-config` to override the embedded deployment-config URL.
### Fixes
* Fix signing of MSI bundles by no longer setting the SN.
* Fixed download progress of previous stages adding to the displayed progress of later stages.
* Fixed regression issue which would cause Kubuntu's UI to become unresponsive while trivrost is running because apparently it cannot deal with 10 window title text changes per second.
### Changes
* Failure during self-update will now be reported to the user in a less generic, more detailed message.
* Remove debug symbols from final binaries. This saves up to 30% for the Linux binary's and 60% for the Windows binary's filesize.
* Removed the reinstall dialog. trivrost will now launch installed trivrost binaries which report the same or a more recent [build time](docs/cmdline.md) as itself (effectively acting as a shortcut) and install over those which don't. If launching an installed binary fails, trivrost will fall back to reinstalling.
* The messages shown when trivrost hashes local files are now distinct from the messages shown when trivrost retrieves remote hashes/bundle info.
* Building the project using the provided Makefile now allows the `icon.png` resource to be missing for Linux builds.
* Project now uses GoLang 1.13.
* Renamed project to "trivrost" for open source release.
* `scripts/signer` has been reimplemented in Go (under `cmd/signer`) and will be built when running `make tools`. `scripts/signer` has been removed.

## 1.2.2 (2019-09-12)
### Changes
* Backported from 1.3.0: Failure during self-update will now be reported to the user in a less generic, more detailed message.

## 1.2.1 (2019-09-11)
### Fixes
* Backported from 1.3.0: Fix signing of MSI bundles by no longer setting the SN.

## 1.2.0 (2019-09-09)
### Features
* trivrost will now download a number of files concurrently to reduce delays introduced by round-trip times.
* trivrost will now stream file downloads to disk directly, without holding a copy in RAM first.
* trivrost now only shows a single progress bar which sums up the total progress of all launcher stages.
* trivrost now shows the download speed during stages which connect to the internet. Download size is no longer shown.
* Resolvable problems during downloading (such as an unstable internet connection) will now be unintrusively reported in the progress window for as long as they persist.
* trivrost now first determines the state of **all** local bundles before downloading **all** remote bundle info files instead of doing both tasks alternatingly for every single bundle.
  * This also profits from the concurrent download feature, further reducing round-trip times.
* You can now customize the status messages displayed by trivrost through the launcher-config.
* trivrost will now correctly abort when a remote file changes during download.
* New flag `-roaming` which, when set, will cause all files which would be written under `%LOCALAPPDATA%` to be written under `%APPDATA%` instead. (Windows only)
  * If trivrost is installed with this flag, all shortcuts will have it set as well.
### Fixes
* Fixed regression issue which would have caused false positives in system bundles change detection.
* Fix possible issue related to the GUI event queue becoming full.
* Fixed old binaries not being removed due to incorrect name being checked.
  * trivrost now also removes binary files from aborted self updates and failed uninstallations on sight.
* Fixed documentation for bundle-paths containing some incorrect paths.
### Changes
* trivrost will now only deny system mode installations which require changes to system bundles from executing if there are also changes to user bundles.
* Removed delay before showing progress window after starting. The intention was to not show the progress window at all if trivrost would finish quickly, but in practice, it never finishes quickly (i.e. in less than half a second).
* `hasher` now adds file-size information to the bundleinfo file which is used draw a correct progress bar in the GUI starting in this version.
* Updated `CONTRIBUTING.md` to reflect new fast-forward workflow.
* Set `Vendor` and `OriginalFilename` correctly in Windows manifest.
* trivrost now updates the bundle folders directly, instead of downloading to a temporary directory first. The original idea was to prevent inconsistent states if trivrost should crash. However, since trivrost never attempts to launch the application in a state which does not match the one described in `bundleinfo.json`, this behavior is no longer required. It also enables trivrost to resume an update even if the launcher process was killed.
* Less verbose linker output when building.

## 1.1.8 (2019-07-17)
### Fixes
* Validator tool now correctly concatenates package paths.
### Changes
* Validator tool now logs the URLs it checked with a leading "OK:", as well as why it checked them.

## 1.1.7 (2019-07-12)
### Fixes
* More verbose logging output when running a self-update.

## 1.1.6 (2019-07-10)
### Features
* trivrost will now reuse the http client for slightly faster downloads of many small files.
### Fixes
* Fix progress bar indicating wrong progress when trivrost issues HTTP range requests.
* Fix trivrost being unresponsive when trying to close it while downloading.
### Changes
* Window resizing and maximizing is now disabled under Windows.
* File system errors related to missing privileges will now display an error message instructing the user to contact their system administrator, instead of the more generic error message.

## 1.1.5 (2019-06-28)
### Features
* trivrost will now perform range requests on the remainder of failed transfers instead of starting from scratch.
### Fixes
* If automatic Windows proxy was configured but not set, retrieving the proxy config would fail on windows. trivrost checks the automatic proxy config and if it is enabled but does not reutrn anything, it falls back to scripted or static proxy configuration. trivrost logs the detected proxy settings when starting. When started with -debug, it even logs the active proxy for each request.

## 1.1.4 (2019-06-27)
### Fixes
* Idle connections now time out.
* trivrost now respects Windows system proxy settings.

## 1.1.3 (2019-05-31)
### Fixes
* Fix two trivrost instances being able to lock each other out from starting.
### Changes
* Log file timestamps are now always in UTC.

## 1.1.2 (2019-05-29)
### Features
* New tool `validator` which can detect several problems in a given deployment-config.
* Support for the `file://`-scheme in URLs using absolute local file paths.
### Fixes
* Fix race condition in file locking mechanism by using OS file locking APIs. Resolves 'Could not unmarshal process signature json' error message.
* trivrost will no longer give up retrying on read errors during file download. Should alleviate troubles with trying to get through the chinese wall.
* Show warning when environment-variables for signing are not set.
### Changes
* `.msi`-installer now displays correct `BrandingName` instead of a random temporary file name generated by Windows in the UAC confirmation dialog.

## 1.1.1 (2019-05-17)
### Fixes
* Workaround for potential MSI validation bugs when creating MSI installer.

## 1.1.0 (2019-05-16)
### Features
* trivrost will now remove unknown bundles; i.e. if the `bundles` folder contains a folder which does not match with the `LocalDirectory` of one of the `Bundles` in the deployment-config, then that folder will be deleted. This should keep users' hard drives clean if bundles are ever removed or renamed in the deployment-config.
* The OS and architecture placeholders `{{.OS}}` and `{{.Arch}}` for deployment-config are back and will expand into `windows`/`linux`/`darwin` and `386`/`amd64`.
* All windows spawned by trivrost now appear centered on screen on Windows.
* All windows spawned by trivrost now appear on top of all other windows – not grabbing focus – on Windows.
* All windows spawned by trivrost now use the 16x16 and 32x32 variants of the embedded icon in the window titlebar and taskbar on Windows, when available.
* Introduced new field `LingerTimeMilliseconds` under `Execution` in deployment-config: set a time in milliseconds (as a JSON number, i.e. no surrounding quotes) for the trivrost window to remain on screen after kicking off the application. Useful for applications which always take some time to start, so the user can tell that their computer is still busy instead of thinking that trivrost crashed.
* Added new flag `--nostreampassing` which disables passing of standard streams (stdout, stderr and stdin) to programs started through trivrost execution.
* New helper-tool `bundown`, provided a deployment-config, can download bundles for a desired OS and Arch into a desired folder.
* `Bundles` in deployment-config can now have `Tags`, which is an array of strings which can be used to limit what `bundown` downloads.
* Added support for preinstalled bundles in write-protected folders: if trivrost finds a folder called `systembundles` next to [itself](docs/glossary.md#trivrost-deployment-artifact), it will consider itself installed and interpret the contained folders as bundles. These bundles will still be validated, but not attempted to be updated if changes are required. The `bundles` folder in the user files will then only keep bundles not already contained in `systembundles`. For execution, trivrost will then first look for executable files relative to `systembundles`. The working directory will however always be the `bundles` folder. This mechanism allows trivrost to, e.g., be installed under `C:\Program Files (x86)` on Windows.
* Two new jobs called `bundle-msi` and `sign-msi` for handling msi building of trivrost. `bundle-msi` can be given `DEPLOYMENT_CONFIG=foo` and `arch=amd64/368`. Calling this job after `bundle` will create a 32bit or 64bit MSI Windows installer file. It will install the binary plus all bundles with the tag `msi` from the deployment config. A desktop shortcut and start menu shortcut will be created. Requires WIX Toolkit to be installed.
### Fixes
* The `signer` script can be run in parallel now.
* Fixed path to the downloaded binary file (during update) being relative to the working directory.
* Fixed several misuses of `len()` on strings which would cause problems with non-ASCII text.
* trivrost's progress window remains open after execution if the user was dragging it while trivrost tried to terminate the UI. (Caused by a bug in libui)
* The correct git hash, branch and tag is now printed in the log.
### Changes
* Removed the install prompt. trivrost installs immediately, unless its desired target path is already occupied by a file. In that case, a reinstall dialog is displayed.
* The trivrost binary will no longer remove itself upon install.
* Old binaries will now be removed based on their name ending in `.old.exe` on Windows, instead of using the `-remove` command line flag.
* Removed the `-remove` command line flag.
* The trivrost binary will now use the `BinaryName` configured in the launcher-config on install, instead of retaining whatever the executable file's name was.

## 1.0.1 (2019-04-17)
### Features
* Added support for Windows systems with broken `%APPDATA%` and `%LOCALAPPDATA%` environment variables.

## 1.0.0 (2019-04-17)
* First internal release.
