# How to panic
In trivrost, panicking in the correct manner is important for the logs to read nicely as well as allowing the panic to be reported to the user correctly and to an appropriate extent. To panic, use `panic()` and supply either a `string` or an `error` as the parameter. A common construct is `panic(fmt.Sprintf())`. Do not panic using `logrus.Panic()` or `logrus.Panicf()`, as these will print the current stack themselves, which, however, is already done by trivrost's top-level panic handling function `cmd/launcher/main.go:handlePanic()`.

## Display error message to user
If you panic using a `misc.IUserError` (such as `*misc.NestedError`), the string returned by its `UserError()` function will be displayed to the user in an error message box instead of a generic error message. In either case, the user is given the option of opening the log folder before terminating trivrost.

## Terminate gracefully
Apart from the notable exception in `cmd/launcher/locking/restart.go:RestartWithBinary()`, instead of terminating with `log.Exit(0)`, you should do so by calling `panic(context.Canceled)` instead. The GUI calls the `cancelFunc` of trivrost's top-level `context.Context` when the user tries to close trivrost's main window. Various long-running code then panics with the current context's `error` value (`.Err()`) – namely `context.Canceled` – which leads `handlePanic()` to terminate trivrost normally.

# When to panic
Only `panic()` whenever there is no good way to deal with an `error` and trivrost should terminate. We often use "`Must()`"-functions when both variants are available: one which panics and one which returns an `error`.
