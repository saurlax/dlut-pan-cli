## dlut-pan

Operate DLUT Pan from your terminal.

### Synopsis

Unofficial CLI for Dalian University of Technology's cloud drive.
List folders, search files or groups, join groups and download individual files.
Credentials stay in memory. Use DLUT_PAN_TOKEN or --username with a password prompt.
The default root is NAOSI; use --root-id for your personal or another group space.
Run any command group without a subcommand to see its help.

```
dlut-pan [flags]
```

### Options

```
      --auth-url string    Login endpoint (DLUT_PAN_AUTH_URL or http://202.118.66.144:8088/v1/auth/sign_in)
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --device string      Login device, link name and LDAP identifier (default "dlut-pan-cli")
      --dry-run            Preview method, URL and body without credentials or network access
  -h, --help               help for dlut-pan
      --json               Print API JSON with secret fields redacted
      --password-stdin     Read the login password from stdin instead of environment or terminal prompt
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
      --token string       Access token; prefer DLUT_PAN_TOKEN or --token-stdin to avoid shell history
      --token-stdin        Read a token from stdin instead of the flag or environment
      --username string    Student/staff account (DLUT_PAN_USERNAME); authenticate for this command
```

### SEE ALSO

* [dlut-pan auth](dlut-pan_auth.md)	 - Check authentication without saving credentials.
* [dlut-pan completion](dlut-pan_completion.md)	 - Generate the autocompletion script for the specified shell
* [dlut-pan file](dlut-pan_file.md)	 - List directories and download files.
* [dlut-pan group](dlut-pan_group.md)	 - Join a known group.
* [dlut-pan search](dlut-pan_search.md)	 - Search file names, file contents or groups.
* [dlut-pan space](dlut-pan_space.md)	 - Inspect built-in space identifiers.
