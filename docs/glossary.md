# Glossary

## Bundle
A bundle simply is a subfolder created inside a `bundles` directory somewhere in the file system, which any files of your choosing are downloaded to. The downloaded files may be nested into further subdirectories. trivrost downloads files that are missing or have changed and deletes files which are no longer wanted based on their current and desired hash values, the latter of which are read from a [bundle information file](#bundle-info). You can have as many bundles as you need and they can contain anything as long as it is a plain file in a normal folder hierarchy. (symbolic links and such are not supported.)

### User bundle
A typical bundle as described above, located in a `bundles` directory located somewhere under the user's home directory. Specifically, *not* a [system bundle](#system-bundle).

### System bundle
A bundle located under `.\systembundles`. See [system mode](lifecycle.md#system-mode).

## Bundle info
A bundle info file describes a single bundle with a unique name, timestamp, relative file paths (which resolve to URLs to said files by joining to the URLs specified by `LauncherUpdate.BundleInfo` and `Bundles.BundleInfo` in the [deplyoment-config](#deployment-config)) as well as hash values and sizes of said files. See [bundleinfo.md](bundleinfo.md).

## launcher-config
A config-file embedded into your build of trivrost which, most importantly, tells trivrost where to find the [deployment-config](#deployment-config). See [launcher-config.md](launcher-config.md).

## deployment-config
The deployment config tells trivrost where to check for updates to itself and its bundles, as well as what programs to execute once everything is in place. See [deployment-config.md](deployment-config.md).

## trivrost deployment artifact
We occasionally need to differ between just trivrost's executable binary and its entire deployment artifact as a whole. On Windows and Linux these two concepts are the same thing. On MacOS however, the deployment artifact of trivrost actually is a `.app`-folder, which contains the binary at `Contents/MacOS/launcher`.

## Transmitting flags
Whenever trivrost restarts itself as part of its exclusive lock and self-update mechanisms, it passes most of the command line arguments it was run with (such as `-skipselfupdate`) to the new instance as well. We refer to arguments affected by this as "transmitting flags".
