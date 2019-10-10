# Locking
trivrost uses files as locks for synchronisation between multiple trivrost instances and the processes started by them. In this document we explain this mechanism for developers who have an interest in working on the trivrost source code.

## The problem
- There is a conflict in the file system, when:
  - (a) trivrost attempts to update bundles while a process depending on at least one of them is still running.
  - (b) a user starts trivrost, while another trivrost instance is already running an update.

## Desired behaviour
- (a) trivrost defers the update of bundles until all depending processes are terminated. (This prevents changes to bundles from conflicting with running processes)
- (b) Execution is halted, until the update process is completed, and then continues. (This allows for starting multiple instances of the application)

## Proposed solutions
Of the following proposals, those written in **bold text** were implemented.

For (a)
1. trivrost keeps running after starting the application and watches for its termination. Behaviour with other trivrost instances is communicated via IPC (Inter Process Communication).
    - This can only work for as long as trivrost is not being executed over the network.
    - Go has no feature-complete libraries for this.
2. **We manage a list of process ids in a list and only update bundles when none of the listed processes are running.**

For (b)
1. Synchronize update phase via IPC.
    - Same problems as above.
2. **Create a lock file which contains the process id of trivrost and then only the trivrost instance listed in it may perform bundle updates.**

## Process signatures
- Under Windows, process ids (PIDs) are reused often. (About once every 150 times when restarting a process repeatedly)
  - Because of this, a PID cannot be used to reliably detect whether a process is still running or not.
  - To fix this, we use a process signature which in addition to the PID contains the time in at least
  millisecond precision at which the process was created. All major operating systems provide APIs for this.

## Introducing: the lock file and the process signature files
- `.lock`
  - An actual lock-file in the sense that it is guarded by an exclusive file lock obtained through OS API calls. Never contains any data and is obtained before reading, writing or creating any file except `.lock` itself.
- `.launcher-lock`
  - A Json file containing the process signature of the trivrost instance which is currently allowed to update itself, update bundles and launch the application. Should have been called `launcher-signature.json`.
- `.execution-lock`
  - A Json file containing a list of process signatures identifying the applications started by trivrost. Should have been called `application-signatures.json`.

## Execution behavior
A trivrost instance won't perform any actions on the file system nor launch any other applications unless it owns the `.lock`, whereafter it will ascertain that its process signature is stored in the `.launcher-lock` file. At launch, trivrost will attempt to obtain the `.lock` first. This process can be seen in `AcquireLock()`, triggered by `LauncherMain()`. The self-restarting of trivrost through `Restart()` in the case of `LockClaimed` is happening because in all cases except `LockOwned` the trivrost binary may have changed through an update by another trivrost instance, so we restart to guarantee we are on the latest version. In the case of `LockClaimed` the `Restart()` function writes the process signature of the restarted trivrost process into `.launcher-lock` before releasing the `.lock` and terminating.

The `.execution-lock` is less complicated, because it is only accessed by the one trivrost instance which already holds the `.lock`. In `executeCommands()` the process signature of the started application is added to a list in the `.execution-lock` file. In `Run()` we wait for all processes in the `.execution-lock` file to stop running – using `AwaitApplicationsTerminated()` – in case we have an update to apply while depending processes are still running.
