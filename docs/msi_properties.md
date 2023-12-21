## MSI Properties

The Windows .msi can be invoked with `misexec`, where additional properties can be provided:

* `ADDROAMINGARGUMENT=true`: If set, shortcuts will be installed such that they run trivrost with the `--roaming` argument, which is of interest to some Citrix environments.
