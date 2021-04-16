# Walkthrough

To build a trivrost based launcher for your application, you need the configuration and build environment. It helps to understand the concepts, file paths and configuration layout, especially if you plan to it in a CI/CD environment.

## Set up your build environment

Depending on your operatingsystem, [different steps are required](building.md), to set up your build environment.

## Compiling the tools

We will need a few tools for creating our bundles. The included Makefile includes a 'tools' target.
When executing `make tools`, the binaries will be placed in the `out/` directory. These need to be executable for the next steps.

## Creating bundle metadata

Determine which 'bundles' your application needs. This could be a runtime, a bunch of application files, resource files or generally 'a bundle of files' that could have a version info attached to them.
Using the tool `hasher` you can create the required metadata for bundles. It takes two arguments: A unique bundle name and a directory.

The unique bundle name will be saved together with a timestamp in the resulting bundleinfo.json. It is a security feature that prevents downgrade attacks if bundle names are reused, even in case of a MITM attacker. A trivrost client will not allow any bundle with the same name and earlier timestamp to be downloaded. But we recommend not reusing bundle names but instead encode version information in this name, e.g. `windows/amd64/openjdk-jre-8u234` or even a GUID. Downgrade attacks are still prevented through timestamped deplyoment-configs and when using a TLS connection.

A `bundleinfo.json` metadata file will be created in the bundle directory.

## Creating the client configuration

trivrost needs the following files to build into the resulting launcher, which you need to provide. They need to be placed into the `cmd/launcher/resources/` directory before building:
* [launcher-config.json](glossary.md#launcher-config) → used during build to configure the launcher. It is embedded into the binary so trivrost and also contains the server URL to find the update infos (deployment config).
* [icon.png](glossary.md#icon) → embedded into the binary as the application icon for Linux. Optional.
* [icon.ico](glossary.md#icon) → embedded into the binary as the application icon for Windows. Optional.
* [icon.icns](glossary.md#icon) → embedded into the application bundle as the application icon for MacOS. Optional
* [public-rsa-keys.pem](security.md) → embedded into the binary to verify signed updates with. May contain multiple public keys to verify bundles signed with different keys. See below.

## Creating the server configuration

The `deployment-config.json` file will be later placed on the server and controls which bundles are downloaded. See [Deployment-config specification](docs/deployment-config.md). This file should NOT be placed into the client's resource directory.

## Signing

To validate the authenticity of bundles, the generated metadata must be signed. You need to create a private key and keep it secret, the public key is added as public-rsa-keys.pem into the launcher.

When retrieving the `bundleinfo.json` files as well as the `deployment-config.json` file, trivrost expects a `….signature`-counterpart to verify the validity of their contents. To generate these signature files, you can use the tool `signer`.

You can generate a private key using the OpenSSL command line utilities:
```sh
openssl genrsa -out private-rsa-key.pem 4096 # Generate private key
openssl rsa -in private-rsa-key.pem -pubout -out public-rsa-keys.pem # Extract public key
```

The tool `signer` takes two arguments: the private key and the files (or list of files) to sign.
```sh
signer private-rsa-key.pem /home/foo/myapp/foobundle1/bundleinfo.json
signer private-rsa-key.pem /home/foo/myapp/fooruntime/bundleinfo.json
signer private-rsa-key.pem /home/foo/myapp/trivrost/deployment-config.json
```

Make sure that noone but you has access to `private-rsa-key.pem`. Copy the public key `public-rsa-keys.pem` to `cmd/launcher/resources/`. See [security.md](security.md#Signing) for more info.

## Compiling

trivrost will create a .exe on Windows, a .tar with an executable on Linux and a .zip on MacOS that is a launchable Mac application. Use `make` to build the binaries to `out/…` and `make bundle` to create the final artifacts.

Launching trivrost now should result in an idling trivrost window with the log showing the missing remote update files. trivrost is extremly robust towards network issues. If you leave the application running like this, it will 'magically work' the moment the remote files are available, not even stumbling over incomplete files while they are still being uploaded or interrupted network connections.

## Setting up a webserver/backend

To use trivrost, you have to operate at least one webserver. This webserver only has to deliver static content, and should support range requests. Even though the config files and packages are cryptographically signed, it is recommended to use TLS for securing connections to this webserver. An exemplary file/folder structure on the webserver could look like this:
```
myapp
├ deployment-config.json
├ deployment-config.json.signature
├ foobundle1
| ├ bundleinfo.json
| ├ bundleinfo.json.signature
| ├ soundfx1.wav
| └ soundfx2.wav
├ fooruntime
| ├ bundleinfo.json
| ├ bundleinfo.json.signature
| ├ myapp.exe
└ launcher
  ├ linux
  | ├ myapp.tar
  | ├ bundleinfo.json
  | └ bundleinfo.json.signature
  └ mac
    ├ myapp.zip
    ├ bundleinfo.json
    └ bundleinfo.json.signature

```

The URL in the `launcher-config.json` for the deployment config should point to the `deployment-conig.json` on this webserver.

## Deploying an update

The recommended way to create an update is to create a new bundle (with a new unique name), make it available and update the `deployment-config.json`. If trivrost is launched with a `deployment-config.json` of a non-existing bundle, or a not yet fully uploaded bundle, it will just sit and try endlessly. Incomplete uploads are detected via the hashes and retried. If locally files did not change, they are not re-downloaded. If the network connections fails in the middle of a download, it is resumed where it stopped.
