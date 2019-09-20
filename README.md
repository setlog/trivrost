# trivrost
trivrost is a reupurposable application-downloader and -launcher: it updates some files on a computer and executes a command afterwards, no questions asked. It can also update itself to introduce new features without the need for user interaction. trivrost is our solution to allow for one-click deployment of an application on a computer, maximizing chances of deployment success while minimizing chances of incoming support tickets. See [reasons.md](docs/reasons.md) for more background information.

## GitHub migration
We are currently in the process of migrating the project as OpenSource to GitHub. The first official OpenSource release will be v1.3.0 in a few days.

## When do I need trivrost?
When you need to deploy a desktop application which always needs to be up to date to many users using all three major OSes and all of them expect your software to *just work*.

## What could go wrong?
See [Troubleshooting](docs/troubleshooting.md).

## What does it look like?

![Screenshot of trivrost progress window](docs/res/screenshot.png "Progress window")

## How does it work?
You release your own build of a trivrost binary to your users. The users execute it, causing it to [install and run](#Install) your software by downloading required files from a webserver administrated by you.

### Install
When trivrost is started, it will check its location in the file system. If it finds that it is not where it wants itself to be, it will copy itself there (all OSs have appropriate folders for this occasion) and create desktop and Start menu shortcuts for quick access. After that, trivrost will [run](#Run). Alternatively, there is [system mode](docs/glossary.md#System-mode) installation.

### Run
When trivrost finds that it is installed, it will run through the following 3 phases:
1. Update itself using files on a webserver which you operate, restarting on success.
2. Download [bundles](docs/glossary.md#Bundle) from said webserver into a `bundles` directory, updating outdated and missing files and deleting unwanted ones.
3. [Execute](#Execute) commands. (e.g.: run Java with the `-jar` argument)

During this, trivrost makes sure that it does not race with itself, in case it is running multiple times; for example when the user wants to launch the application multiple times, which is supported.

### Execute
When trivrost begins to execute [the programs you have configured](docs/deployment-config.md), it will do so on the basis of treating the `bundles` directory as the working directory. This way, any downloaded programs inside `bundles` can be executed by using a relative path, and any relative file paths contained in program arguments are relative to the `bundles` directory as well.

### Uninstall
Because trivrost is designed to be able to install without administrative privileges, it does not attempt to register typical uninstallation routines, such as an entry under `Add or remove programs` in the control panel of Microsoft Windows. Instead, a Start menu shortcut is created which runs the program with an `-uninstall` parameter. On Windows, this shortcut is placed in the start menu. Again, [System mode](docs/glossary.md#System-mode) is different.

## Learn more
1. [Background info](docs/reasons.md)
2. [Glossary](docs/glossary.md)
3. [Walkthrough](docs/walkthrough.md)
4. [File locations](docs/file_locations.md)
5. [Launcher-config specification](docs/launcher-config.md)
6. [Deployment-config specification](docs/deployment-config.md)
7. [Bundle info specification](docs/bundleinfo.md)
8. [Building](docs/building.md)
9. [Security](docs/security.md)
10. [Command line reference](docs/cmdline.md)

## Contribute to development
See [CONTRIBUTING.md](CONTRIBUTING.md).
