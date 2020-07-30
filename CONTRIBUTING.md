# Contributing
When contributing to this repository, please first discuss the change you wish to make via an issue.

This project uses [EditorConfig](https://editorconfig.org) to set some style rules. If your IDE does not support it natively, please install [a plugin](https://editorconfig.org/#download).

## Learning
Developer documentation may be found under `docs/dev/`.

## Pull Request Process
1. Make sure to update `CHANGES.md` to list your changes. For non-release commits, add a heading
   of the form `## x.y.z (TBD)` under which to add your changes.
2. When changing, adding or removing functionality (i.e. anything other than bug-fixes and improving
   code quality) make sure to reflect this in the documentation under `docs/`.
3. Before committing anything, make sure the following call of `golangci-lint` reports no problems:
   `golangci-lint run --enable-all --disable stylecheck --disable gochecknoglobals --disable gochecknoinits --disable lll --disable gocritic --disable dupl --disable golint --disable gofmt --disable goimports ./...`

## Building
Please see [building.md](docs/building.md) for general building instructions.

## Branches
Main release branch: **master** (protected)
  - Each commit is one release
  - Each commit is tagged with the version
  - This is always the latest, stable and released version
  - Each tag creates a new release

Integration branch: **develop** (protected, default)
  - Finished features, fixes etc. are merged into this branch
  - If possible, it should at any time build and pass tests
  - Tags can be used for test releases and should be distinct from existing version tags
  - Tags will create full releases, incl. docker images

Feature-, bugfix and development branches:
  - Any other branch

## Development cycle
Development takes place on **develop**. It should always build and any larger features should be done in feature branches. Docker images are built for tags and branches matching `develop` or the regular expression `v[0-9]+\..*/`.

Releases are tagged on **master** using v+[SemVer](https://semver.org/). E.g.: `vX.Y.Z`

To create a release for the latest version:
  * Finish all pending changes to develop
  * Make sure CHANGES.md is complete and has the correct publication date and version number
  * Run these commands:
    ```
    git checkout develop
    git pull
    # check with git log that the branch is in the state which you want to release
    git checkout master
    git pull
    git merge develop --ff-only
    git push
    # check if pipeline succeeds - should always be the case!
    ```
Finally, make a release through [the release-overview of the project's GitHub page](https://github.com/setlog/trivrost/releases), creating a new tag against `master`. Tag should be `vX.Y.Z`, title should be `vX.Y.Z (YYYY-MM-DD)` and message should be the markdown-formatted list of fixes, features and changes from `CHANGES.md`.

Alternatively:
```
git tag "vX.Y.Z"
git push origin "vX.Y.Z"
# you can also push all tags with 'git push --tags' if you made sure not to have pointless tags locally
```

Finally, make sure that you do not accidentally continue work on master:
```
git checkout develop
```


To create a release for a maintenance version, do not merge it into `master`. Instead, tag the release on the maintenance branch.
