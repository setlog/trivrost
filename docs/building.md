# Building trivrost for your project
To distribute your project using trivrost, it is recommended to use a CI system. First you create the necessary configuration files inside your project (e.g. in a `/trivrost` subdirectory). Then you configure your pipeline to clone `trivrost`, copy all files to the right directories and build trivrost. In this document, we will focus on how to just get the project to build.

Requirements:
  - golang >= 1.13
  - git
  - make
  - zip
  - goversioninfo
  - libgtk-3-dev (under Linux/MacOS only)
  - signtool (if signing executable or MSI under Windows)
  - candle, light (WiX toolkit) (if building/signing MSIs under Windows)

trivrost needs certain resource files generated with `go generate` and then can be built using `go build`. Depending on the operating system and for extra debug information, additional flags are required. For convenience a [Makefile](../Makefile) is provided with several make targets. The default target detects the OS and will build all project components and can be run using just `make` inside the project root.
The resulting binaries are placed under the `out/update_files/<windows|linux|darwin>/` subdirectory.

# Make targets
Run `make help` to see all available targets.
To force trivrost to spawn a console which shows log output when starting on Windows, set environment variable `TRIVROST_FORCECONSOLE` to anything non-empty.

# Windows
It uses a go port of [libui](https://github.com/andlabs/libui) using cgo which requires Windows libraries to compile for Windows. Thus cross-compiling causes problems (and is not supported with the provided Makefile). It is highly recommended to build trivrost under Windows using a native Windows.
To build under Windows using the provided Makefile:
  1. Install the native [GIT for Windows](https://git-scm.com/download/win) package (cygwin GIT is currently [outdated](https://github.com/me-and/Cygwin-Git/issues/40) and [broken](https://github.com/golang/go/issues/23155)).
  2. Install [mingw-w64](https://mingw-w64.org) (specifically the [Mingw-builds](https://mingw-w64.org/doku.php/download/mingw-builds) package) as a recent GCC compatible compiler.
  3. Install golang.
  4. Install Cygwin. (Cygwin effectively simulates a Linux environment with typical Linux developer tools, which we currently need; getting away from this is an open issue)
  5. Open the Cygwin Terminal and type `cd /cygdrive`. In there, your drives will appear as if they were folders. From there, `cd` to your project root. Run `ls` at any time to see what's in the current folder.
  6. Run `make` from inside the Cygwin terminal.

## Signing binaries on Windows (optional)
Optionally, sign the final binary. Requires `signtool` from the Windows SDK.
  - Install `signtool` from the Windows SDK.
  - Put a base64 encoded p12 certificate into the environment variable `CERT_FILE` and its password into `CERT_KEY`
  - Run `make sign`

## Building and signing MSIs on Windows (optional)
It is possible to create an MSI package with prebundled bundles. The bundles that should be prebundled need a JSON-array under the key `Tags` containing an element called `"msi"` in the deployment-config. The launcher will no longer be able to update itself and store those bundles under a directory called `systembundles` under the 'Program Files' directory. (Usually: `C:\Program Files\Vendor\Product`).
  - Install WiX installer set.
  - Run `make bundle-msi ARCH=368 DEPLOYMENT_CONFIG=<path to file>` or `make bundle-msi ARCH=amd64 DEPLOYMENT_CONFIG=<path to file>`
  - Note: Creating an MSI package without a package configured to bundle will fail when creating the installer in a 'harvest' phase.
  - The resulting .msi file is placed in the `release_files` directory and can be signed using `make sign-msi`.
  - Signing is possible with the `make sign-msi` target. It will always sign the 32bit and 64bit version and requires a 64bit system.

# Linux
Install golang, GTK headers and build-essentials and just run `make`.

# Darwin
Install golang, the Developer Tools and just run `make`.

# Developing trivrost
For easier local developing and testing, the following make targets `make copy-test-files`, `make test` and `make test-integration` exist.
