# Bundle info
Bundle info files are JSON-files which contain information about a [bundle](walkthrough.md#Bundles). Their filename typically is `bundleinfo.json`. Bundle info files contain file paths relative to `LauncherUpdate.BaseURL` and `Bundles.BaseURL` from the [deployment-config](walkthrough.md#deployment-config), respectively, as well as the SHA-256 hash values of the files at those paths, enabling trivrost to validate whether some files in a folder structure represent the described bundle, or not. Bundle information files should be generated using the `hasher` tool contained under `cmd/hasher`.

## Fields
* **`Timestamp`** (string): See [Timestamps](security.md#Timestamps).
* **`UniqueBundleName`** (string): See [Timestamps](security.md#Timestamps).
* **`BundleFiles`** (object): An object where each key describes a file with a relative file path and each value is another object with further file information.
  * **`SHA256`** (string): The cryptographically secure SHA-256 hash value of the file.
  * **`Size`** (int): The size of the file in bytes. Used for accurate download progress reporting in trivrost's GUI.

## Examples
A bundle info file may look something like the following, though real-world examples are likely to be longer:
```
{
  "Timestamp": "2019-04-08 15:13:37",
  "UniqueBundleName": "MediaFilesForHelp",
  "BundleFiles": {
    "README.txt": {
      "SHA256": "62ebd700a6200c8ceba54d8e9af87cf062336cec4fa9df910527a3b67723d779",
      "Size": 16403
    },
    "somedir/somefile": {
      "SHA256": "941b82840ec99a802c9708f8d43bcb458a7ab4bf6133d0ed15a258d2cd45dab1",
      "Size": 13454
    }
}
```