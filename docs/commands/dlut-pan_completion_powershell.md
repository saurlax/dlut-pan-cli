## dlut-pan completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	dlut-pan completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
dlut-pan completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --dry-run            Preview method, URL and body without credentials or network access
      --json               Print API JSON with secret fields redacted
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
```

### SEE ALSO

* [dlut-pan completion](dlut-pan_completion.md)	 - Generate the autocompletion script for the specified shell
