# launcher-config
The launcher-config is a JSON-file which is embedded into your build of trivrost and contains build- and runtime-information. Its filename typically is `launcher-config.json`.

## Fields
* **`DeploymentConfigURL`** (string): The URL where the [deployment-config](deployment-config.md) can be retrieved with a HTTP GET-request. HTTPS is supported.
* **`VendorName`** (string): Name of every highest-level folder created by trivrost. Should be the name of the vendor, company or publisher releasing the build you are making.
* **`ProductName`** (string): Name of most second-highest-level folders created by trivrost. Should be the name of the software product you intend trivrost to download/update and launch.
* **`BinaryName`** (string): Name of the built binary (or `.app`-bundle, in the case of MacOS) which is put into the folder named by the value of `ProductName`, omitting extension. This value may also be used to determine the name of the binary as to be released to users.
* **`BrandingName`** (string): Sets the name the program will present itself to the user with. It affects the names of shortcuts created by trivrost and is visible in the title bar and labels of windows shown by trivrost during installation, updates and launching.
* **`BrandingNameShort`** (string): Same as `BrandingName`, but no longer than 15 bytes. Used on MacOS and not usually seen by the user.
* **`ReverseDnsProductId`** (string): **Unique** name of your product in reverse-DNS notation, e.g. `com.example.product.client.environment`. This matches `CFBundleIdentifier` in a MacOS `Info.plist` and `assemblyIdentity.name` in a Windows application manifest file.
* **`ProductVersion`** (object): The [semantic version](https://semver.org/) of your software, as well as an additional build number. You can leave the build number at `0`. The version can be seen e.g. when a user hovers their mouse cursor over the file in Windows Explorer.
  * **`Major`** (int): Major version of your software.
  * **`Minor`** (int): Minor version of your software.
  * **`Patch`** (int): Patch version of your software.
  * **`Build`** (int): Build number of your software.
* **`StatusMessages`** (object): Customize status messages displayed by trivrost during various stages of its execution.
  * **`AcquireLock`** (string): trivrost is waiting for another trivrost instance. (default: `Waiting for other launcher instance to finish...`)
  * **`GetDeploymentConfig`** (string): Deployment-config is being downloaded and parsed. (default: `Retrieving application configuration...`)
  * **`DetermineLocalLauncherVersion`** (string): SHA-256 hash values of current trivrost installation is being determined by reading the installed file(s). (default: `Determining launcher version...`)
  * **`RetrieveRemoteLauncherVersion`** (string): SHA-256 hash values of remote trivrost artifact are being determined by downloading its bundle info file. (default: `Checking for launcher updates...`)
  * **`SelfUpdate`** (string): New launcher binary or application bundle is being downloaded/applied. (default: `Updating launcher...`)
  * **`DetermineLocalBundleVersions`** (string): SHA-256 hash values of local bundle files are being determined by reading the files. (default: `Determining application version...`)
  * **`RetrieveRemoteBundleVersions`** (string): SHA-256 hash values of remote bundle files are being determined by downloading bundle info files. (default: `Checking for application updates...`)
  * **`AwaitApplicationsTerminated`** (string): Waiting for all instances of the application to exit so the update can be applied safely. (default: `Please close all instances of the application to apply the update.`)
  * **`DownloadBundleUpdates`** (string): New bundle files are being downloaded. (default: `Retrieving application update...`)
  * **`LaunchApplication`** (string): Executing commands specified in deployment-config. (default: `Launching application...`)
* **`IgnoreLauncherBundleInfoHashes`** (array): An array of SHA-256 hash values as hex-encoded strings of launcher bundleinfo files which trivrost should ignore, i.e. act as if no update was available, regardless of whether that is the case. This behaviour can be used to hand out specialized builds to specific users for hotfixing purposes without having to worry about the need to add (and later remove) the `-skipselfupdate` argument.

## Remarks
**You should avoid changing `VendorName` and `ProductName` after distributing the trivrost executable of a project. Currently, if you do change either, trivrost will move its installation location and redownload all bundles, without cleaning up after itself, and without updating the shortcuts.**

## Examples
* [Thorough example](../examples/launcher-config.json.example)