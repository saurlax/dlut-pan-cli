## dlut-pan completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	dlut-pan completion fish | source

To load completions for every new session, execute once:

	dlut-pan completion fish > ~/.config/fish/completions/dlut-pan.fish

You will need to start a new shell for this setup to take effect.


```
dlut-pan completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --auth-url string    Login endpoint (DLUT_PAN_AUTH_URL or http://202.118.66.144:8088/v1/auth/sign_in)
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --device string      Login device, link name and LDAP identifier (default "dlut-pan-cli")
      --dry-run            Preview method, URL and body without credentials or network access
      --json               Print API JSON with secret fields redacted
      --password-stdin     Read the login password from stdin instead of environment or terminal prompt
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
      --token string       Access token; prefer DLUT_PAN_TOKEN or --token-stdin to avoid shell history
      --token-stdin        Read a token from stdin instead of the flag or environment
      --username string    Student/staff account (DLUT_PAN_USERNAME); authenticate for this command
```

### SEE ALSO

* [dlut-pan completion](dlut-pan_completion.md)	 - Generate the autocompletion script for the specified shell
