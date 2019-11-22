# Lifecycle
This document will explain what trivrost actually does when it runs. Whenever we mention or hint at file system operations here, you may refer to [file_locations.md](file_locations.md) to find out where a file or folder is actually located.

## Goal determination
When trivrost is started, it will try to find out what it should do. First, it will check its location in the file system. If it finds that it is not where it wants itself to be, it will check if its desired path is already occupied by another file or folder. If it is not, it will [install](#install). If it is, it will make the educated guess that it is dealing with a previous installation of itself, and try to run the file with the `--build-time` argument. If that fails or it is found that the executable was built at a sooner time than what is currently running, trivrost will also [install](#install). Otherwise, trivrost will act as a shortcut for the installed executable. If trivrost finds that it is in the right location, it will proceed to [update](#update).

Note that whenever trivrost runs, it makes sure that it does not race with itself, in case it is running multiple times; for example when the user wants to launch the application multiple times, which is supported. For more info on this, see [locking.md](dev/locking.md).

## Install
When trivrost has decided that it should install itself, it will copy [itself](glossary.md#trivrost-deployment-artifact) to its desired path, overwriting any previous installation of itself, and create a desktop shortcut as well as two Start Menu entries (MacOS excluded): one for starting and one for uninstallation. After that, it will restart in its new location, upon which it will find that it should [update](#update).

## Update
When trivrost finds that it is installed, it will go through the following update-cycle until everything is up to date:
1. Download the [deployment-config](glossary.md#deployment-config) from the URL specified in the embedded [launcher-config](glossary.md#launcher-config) into memory.
2. If the deployment-config specifies a launcher update for the current platform...
   1. Determine the SHA-256 hash(es) of the running deployment artifact.
   2. Retrieve the according bundle info specified in the deployment-config.
   3. Update the deployment artifact and restart with it if there is any hash mismatch.
3. If the deployment-config specifies any bundles for the current platform...
   1. Determine the SHA-256 hash(es) of the existing bundles.
   2. Retrieve the according bundle info files specified in the deployment-config.
   3. If there is any hash mismatch...
      1. Wait for any running commands which may depend on the bundles to terminate.
      2. Update `bundles` to match the state described by the bundle info files.

When this is complete, trivrost will then [launch](#launch) the commands specified in the deployment-config, i.e. your application.

## Launch
When trivrost begins to launch [the programs you have configured](docs/deployment-config.md), it will do so on the basis of treating the `bundles` directory as the working directory. This way, any downloaded executables under `bundles` can be executed by using a relative path, and any relative file paths contained in program arguments are relative to the `bundles` directory as well.

## Uninstall
Because trivrost is designed to be able to install without administrative privileges, it does not attempt to register typical uninstallation routines, such as an entry under `Add or remove programs` in the control panel of Microsoft Windows. Instead, a Start Menu shortcut is created which runs the program with an `-uninstall` parameter. ([TODO: Figure out where this should go on MacOS](https://github.com/setlog/trivrost/issues/11))

## System mode
System mode is a special behavior engaged by the presence of a folder called `systembundles` next to the trivrost executable. When such a folder is present, trivrost will consider itself installed, and treat the `systembundles` folder as an additional, read-only `bundles` directory. Whether a bundle is interpreted as a system bundle is determined purely by whether it is present under `.\systembundles`; the [deployment-config](deployment-config.md) provides no information on what type of bundle a bundle should be. During [launch](#launch), trivrost will first look for executable files relative to `.\systembundles`. The working directory of executed commands will however always be the `bundles` folder. See also the `IsUpdateMandatory` field in [deployment-config](deployment-config.md).

This feature was introduced to be able to deploy in Windows environments where policy requires native binaries to exist only inside of protected system folders such as `C:\Program Files (x86)`, and is [intended to be preceded by the execution of an `.msi`-installer](building.md#msi), so that it can be uninstalled using the OS's control panel. In system mode, trivrost cannot update itself nor bundles under `systembundles`. If remote bundleinfo files signal any changes to bundles contained in the `systembundles` folder, while also signalling changes to at least one user bundle, trivrost will reject the launch and prompt the user to reinstall the application.
