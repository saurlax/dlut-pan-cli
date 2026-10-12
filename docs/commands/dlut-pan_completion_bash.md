## dlut-pan completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(dlut-pan completion bash)

To load completions for every new session, execute once:

#### Linux:

	dlut-pan completion bash > /etc/bash_completion.d/dlut-pan

#### macOS:

	dlut-pan completion bash > $(brew --prefix)/etc/bash_completion.d/dlut-pan

You will need to start a new shell for this setup to take effect.


```
dlut-pan completion bash
```

### Options

```
  -h, --help              help for bash
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
