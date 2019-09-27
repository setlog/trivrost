# Building trivrost
trivrost needs certain resource files generated with `go generate` and then can be built using `go build`. Depending on the operating system and for extra debug information, additional flags are required. For convenience a [Makefile](../Makefile) is provided with several make targets. The default target detects the OS and will build all project components and can be run using just `make` inside the project root.
The resulting binaries are placed under the `out/update_files/<windows|linux|darwin>/` subdirectory.

In this document, we will focus on how to get the project to build, i.e. successfully run just `make`, disregarding its configuration for now. Additional make-targets are mentioned further below.

# Unix

## Prerequisites
- [`git`](https://git-scm.com/)
- [`go`](https://golang.org/) >= v1.13
- [`goversioninfo`](https://github.com/josephspurrier/goversioninfo)
- [`make`](https://superuser.com/questions/352000/whats-a-good-way-to-install-build-essentials-all-common-useful-commands-on)
- `libgtk-3-dev`

## Building
1. `go get -d -u github.com/setlog/trivrost`
    * Ignore warnings about missing `.go`-files.
2. `cd ${GOPATH}/github.com/setlog/trivrost`
3. `make copy-test-files`
4. `make`

# Windows

## Prerequisites
- [`git`](https://git-scm.com/)
- [`go`](https://golang.org/) >= v1.13
- [`goversioninfo`](https://github.com/josephspurrier/goversioninfo)

## Further instructions
1. Install [mingw-w64](https://mingw-w64.org) (specifically the [Mingw-builds](https://mingw-w64.org/doku.php/download/mingw-builds) package) as a recent GCC-compatible compiler.
    * For *Architecture*, choose `x86_64`.
    * For *Threads*, choose `win32`.
    * For *Exception*, choose `seh`.
    * Add the `bin` folder of MinGW to your `%PATH%`.
2. Install Cygwin. (Cygwin effectively simulates a Linux environment with typical Linux developer tools, which we currently need to run `make` on Windows; getting away from this is [an open issue](https://github.com/setlog/trivrost/issues/12))
    * When prompted what packages to install, add `make`.
    * Do not use Cygwin GIT, as it is currently [outdated](https://github.com/me-and/Cygwin-Git/issues/40) and [broken](https://github.com/golang/go/issues/23155).

## Building
1. `go get -d -u github.com/setlog/trivrost`
    * Ignore warnings about missing `.go`-files.
2. Open the Cygwin Terminal and enter the following commands. Note that you can run the Unix command `ls -la` at any time to see what's in the current directory.
3. `cd ${GOPATH}/github.com/setlog/trivrost` (Note how after entering the command, the working directory starts with `/cygdrive/`, under which Cygwin mimics your Windows drives.)
4. `make copy-test-files`
5. `make`

# Additional make targets
Run `make help` to see all available targets.
To force trivrost to spawn a console which shows log output when starting on Windows, set environment variable `TRIVROST_FORCECONSOLE` to anything non-empty.

You can run `make bundle` to build an archive of the program for distribution. Requires `zip` on MacOS and `tar` everywhere else.

On Windows, you can sign the binary by running `make sign`. For this to work, install `signtool` from the Windows SDK and put a base64-encoded p12 certificate into the environment variable `CERT_FILE` and its password into `CERT_KEY`.

# MSI
It is possible to create an MSI package with bundles which get installed system-wide and cannot be updated. The bundles that should be prebundled need a JSON-array under the key `Tags` containing an element called `"msi"` in the deployment-config. The launcher will no longer be able to update itself and store those bundles under a directory called `systembundles` under the 'Program Files' directory. (Usually: `C:\Program Files\Vendor\Product`). We call this a [system mode](lifecycle.md#system-mode) installation.
  - Install WiX installer set.
  - Run `make bundle-msi ARCH=368 DEPLOYMENT_CONFIG=<path to file>` or `make bundle-msi ARCH=amd64 DEPLOYMENT_CONFIG=<path to file>`
  - Note: Creating an MSI package without a package configured to bundle will fail when creating the installer in a 'harvest' phase.
  - The resulting .msi file is placed in the `release_files` directory and can be signed using `make sign-msi`.
  - Signing is possible with the `make sign-msi` target. It will always sign the 32bit and 64bit version and requires a 64bit system.
