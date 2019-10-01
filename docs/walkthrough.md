# Walkthrough

## In short
1. [Set your system up to be able to build trivrost.](building.md)
2. Write a `launcher-config.json` and put it under `cmd/launcher/resources`.
3. Write a `deployment-config.json`.
4. Determine what bundles your project needs and put them in a folder.
5. Generate `bundleinfo.json` files, e.g. using the `hasher` tool.
6. Create an RSA key pair for signing and validation of files.
7. Using the private key, create signatures of `bundleinfo.json` files as well as of `deployment-config.json` using `scripts/signer`.
8. Copy the public key to `cmd/launcher/resources/public-rsa-keys.pem`.
9. Run `make` and `make bundle`.
10. Set up a webserver.
11. Upload required files to webserver.
12. Publish the files under `out/release_files/os/binaryname`[`.exe`/`.app`] to your users.

## Configure resources
trivrost needs the following files to generate the required sources automatically, which you need to provide. They need to be placed into the `cmd/launcher/resources/` directory before building:
* [launcher-config.json](glossary.md#launcher-config) → used during build and embedded into the binary so trivrost knows where to find the deployment-config.
* [icon.png](glossary.md#icon) → embedded into the binary as the application icon for Linux.
* [icon.ico](glossary.md#icon) → embedded into the binary as the application icon for Windows.
* [icon.icns](glossary.md#icon) → embedded into the application bundle as the application icon for MacOS.
* [public-rsa-keys.pem](security.md) → embedded into the binary to verify signed updates with.

## Hashing and signing bundles
When executing `make tools`, a binary called `hasher` will be created that takes a directory as an argument and creates the `bundleinfo.json` file. A utility called `signer` is also built which takes a private key and a list of files to sign.

## Signing files
All `bundleinfo.json` files as well as the `deployment-config.json` file require a `.signature`-counterpart to verify the validity of their contents. A helper-script `scripts/signer` is provided which takes a private key and a list of files to sign using `openssl`. On Windows, you need to run the script via Cygwin. Here is an example for how to generate a key pair and sign a `bundleinfo.json`:
```sh
openssl genrsa -out private-rsa-key.pem 4096 # Generate private key
openssl rsa -in private-rsa-key.pem -pubout -out public-rsa-keys.pem # Extract public key
scripts\signer private-rsa-key.pem D:\bundles\myapp\bundleinfo.json
```
Make sure that noone but you has access to `private-rsa-key.pem`. Copy `public-rsa-keys.pem` to `cmd/launcher/resources/`. See [security.md](security.md#Signing) for more info.

## Backend
To use trivrost, you have to operate at least one webserver. This webserver only has to deliver static content, and should support range requests. Even though the config files and packages are cryptographically signed, it is hightly recommended to use TLS for securing connections to this webserver. An exemplary file/folder structure on the webserver could look like this:
```
.
├ deployment-config.json
├ deployment-config.json.signature
└ myapp
  ├ bundleinfo.json
  ├ bundleinfo.json.signature
  ├ myapp.exe
  ├ soundfx1.wav
  └ soundfx2.wav
```
