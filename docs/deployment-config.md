# deployment-config
The deployment-config is a JSON-file which is supposed to be hosted on a webserver operated by you and downloaded by trivrost every time it starts. Its filename typically is `deployment-config.json`.

## Fields
* **`Timestamp`** (string): A timestamp in the form `YYYY-MM-DD HH:mm:SS` which indicates when the deployment-config was last changed. This field protects trivrost against attacks. See [security.md](security.md) for more information.
* **`LauncherUpdate`** (array): An array of objects which define bundle configurations for how trivrost updates itself. When trivrost runs, this list must boil down to either one single configuration, or no configurations, through filtering by `TargetPlatforms`.
  * **`BundleInfoURL`**, **`BaseURL`**, **`TargetPlatforms`**: See [Common fields](#Common-fields) below.
* **`Bundles`** (array): An array of objects which define the bundles which trivrost should download and keep up to date.
  * **`BundleInfoURL`**, **`BaseURL`**, **`TargetPlatforms`**: See [Common fields](#Common-fields) below.
  * **`LocalDirectory`** (string): Desired name of the bundle's folder in the file system.
  * **`Tags`** (array): An array of strings describing arbitrary tags. Currently only used by bundown to fetch the files required to build `.msi`-installers for Windows for [system mode](walkthrough.md#System-mode).
* **`Execution`** (object): Object which describes trivrost's behavior after having downloaded and updated itself and all bundles.
  * **`Commands`** (array): An array of objects which define individual commands which will be executed in the order they appear. After starting the last command, trivrost will terminate without waiting for it to complete.
    * **`TargetPlatforms`**: See [Common fields](#Common-fields) below.
    * **`Name`** (string): Name of the program to run, or a relative or absolute path to it. Relative paths will be resolved relative to the `bundles` folder. If trivrost is in *system mode*, relative paths will be resolved relative to the `systembundles` folder first. If you provide only a name, without path separators, the system's `PATH` environment variable will be consulted. Note that this is not a shell command. Syntax such as `echo foo > bar` will not work.
    * **`Arguments`** (array): An array of strings, defining program arguments, e.g. `[ "-jar", "myapp.jar" ]`.
    * **`Env`** (object): Set environment variables for the executed program. Keys represent environment variable names. The value then must be either of type string (set/override variable) or `null` (clear variable).
  * **`LingerTimeMilliseconds`** (int): A time, in milliseconds, that the trivrost progress window should remain open after having executed the last command. Useful if you know that the launched application takes some time to become responsive and want to keep the user entertained.

## Common fields
* **`BundleInfoURL`** (string): URL to a [bundle information file](walkthrough.md#Bundle-info) describing this bundle.
* **`BaseURL`** (string): URL, to which bundle info file paths get joined to to determine download URLs for all files. If omitted, it will be inferred by taking `BundleInfoURL` and stripping the last path element from it.
* **`TargetPlatforms`** (array): Array of strings specifying allowed OS/architecture combinations ("platforms") which this element applies to, using the `GOOS` and `GOARCH` [naming scheme](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) in one of the forms `GOOS`, `GOARCH` or `GOOS-GOARCH`, e.g. `windows-amd64`. If omitted, the element applies to all platforms. See also: [Placeholders](#placeholders).
* **`IsUpdateMandatory`** (bool): For `LauncherUpdate`, if set to true, specifies that the launcher should attempt a self-update even if it is in system mode. For `Bundles`, if set to true, specifies that the user cannot choose to ignore when required changes to a bundle are omitted due to it being a [system bundle](glossary.md#system-bundle); the user will however be informed about omitted updates either way. This has no effect on [user bundles](glossary.md#user-bundle).

## Placeholders
* **`{{.OS}}`**: Identifier for the operating system the running trivrost binary was built for. (`darwin`, `windows` or `linux`)
* **`{{.Arch}}`**: Identifier for the operating system architecture which trivrost is running on. (`386` or `amd64`)

## Examples
* [Basic example](../examples/deployment-config.json.simple.example)
* [Exhaustive example](../examples/deployment-config.json.complex.example)
