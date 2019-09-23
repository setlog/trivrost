# Glossary

## Bundle
A bundle simply is a subfolder created inside a `bundles` directory somewhere in the file system, which any files of your choosing are downloaded to. The downloaded files may be nested into further subdirectories. trivrost downloads files that are missing or have changed and deletes files which are no longer wanted based on their current and desired hash values, the latter of which are read from a [bundle information file](#bundle-info). You can have as many bundles as you need and they can contain anything as long as it is a plain file in a normal folder hierarchy. (symbolic links and such are not supported.)

See [file_locations.md](file_locations.md) to learn where trivrost stores the `bundles` directory and other files.

## Program
With *Program*, we usually refer to the trivrost executable. On Linux and Windows, this is equivalent to its binary. On MacOS, however, this will refer to its `.app`-folder, which is an executable folder which contains the binary at `Contents/MacOS/launcher`.

## Application
With *Application* we usually refer to the last binary launched by trivrost: the application which you want to deploy and run.

## launcher-config
A config-file embedded into your build of trivrost which, most importantly, tells trivrost where to find the [deployment-config](#deployment-config). See [launcher-config.md](launcher-config.md).

## deployment-config
The deployment config tells trivrost where to check for updates to itself and its bundles, as well as what programs to execute once everything is in place. See [deployment-config.md](deployment-config.md).

## Bundle info
A bundle info file describes a single bundle with a unique name, timestamp, relative file paths (which resolve to URLs to said files by joining to the URLs specified by `LauncherUpdate.BundleInfo` and `Bundles.BundleInfo` in the [deplyoment-config](#deployment-config)) as well as hash values and sizes of said files. See [bundleinfo.md](bundleinfo.md).

## System mode
If trivrost finds a directory called `systembundles` next to its binary, it will enter *system mode* causing trivrost to consider itself installed and interpret the contained folders as bundles. These bundles will still be validated, but not attempted to be updated if changes are required. The `bundles` folder in the user files will then only be used to keep bundles not already contained in `systembundles`. During execution, trivrost will then first look for executable files relative to `systembundles`. The working directory will however always be the `bundles` folder. This mechanism allows trivrost to, e.g., be installed under `C:\Program Files (x86)` on Windows using our support for `.msi`-files and uninstalled using the OS's control panel.
