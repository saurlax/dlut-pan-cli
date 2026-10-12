## dlut-pan completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(dlut-pan completion zsh)

To load completions for every new session, execute once:

#### Linux:

	dlut-pan completion zsh > "${fpath[1]}/_dlut-pan"

#### macOS:

	dlut-pan completion zsh > $(brew --prefix)/share/zsh/site-functions/_dlut-pan

You will need to start a new shell for this setup to take effect.


```
dlut-pan completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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
