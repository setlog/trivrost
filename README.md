<p align="center"><img width="128px" height="128px" style="width:128px;height:128px" src="docs/res/trivrostlogo.png" alt="trivrost" /></p>

# trivrost
![Test](https://github.com/setlog/trivrost/workflows/Test/badge.svg?branch=master)
[![Renovate](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://renovatebot.com)

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-science.svg)](https://forthebadge.com)

trivrost is a repurposable application-downloader and -launcher in the form of a native executable: it updates some files on a computer and executes a command afterwards, no questions asked. It can also update itself to introduce new features without the need for user interaction. See [reasons.md](docs/reasons.md) for more background information.

## When do I need trivrost?
When you need to deploy an always-online desktop application which always needs to be up to date to many users using all three major OSs and all of them expect your software to *just work*.

## What does it look like?

![Screenshot of trivrost progress window under Windows](docs/res/screenshot.png "Progress window under Windows")

## How does it work?
You release your own build of a trivrost executable to your users. The users start it, causing it to [install and run](docs/lifecycle.md) your software by downloading required files from a webserver administrated by you.

## State of this project
Production-ready, with high confidence for Linux and Windows builds. Has approximately 10.000 active Windows users for one of our builds. MacOS-support [needs input](https://github.com/setlog/trivrost/issues/11).

## Learn more
1. [Background info](docs/reasons.md)
2. [Glossary](docs/glossary.md)
3. [Lifecycle](docs/lifecycle.md)
4. [File locations](docs/file_locations.md)
5. [Building](docs/building.md)
6. [Walkthrough](docs/walkthrough.md)
7. [Launcher-config specification](docs/launcher-config.md)
8. [Deployment-config specification](docs/deployment-config.md)
9.  [Bundle info specification](docs/bundleinfo.md)
10. [Security](docs/security.md)
11. [Command line reference](docs/cmdline.md)
12. [Troubleshooting](docs/troubleshooting.md)

## Contribute to development
See [CONTRIBUTING.md](CONTRIBUTING.md).
