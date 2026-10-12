## dlut-pan

Operate DLUT Pan from your terminal.

### Synopsis

Unofficial CLI for Dalian University of Technology's cloud drive.
Run 'dlut-pan auth login' once to save your session in the system credential store.
File, search and group commands reuse the saved session without asking for credentials.
The default root is NAOSI; use --root-id for your personal or another group space.
Run any command group without a subcommand to see its help.

```
dlut-pan [flags]
```

### Options

```
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --dry-run            Preview method, URL and body without credentials or network access
  -h, --help               help for dlut-pan
      --json               Print API JSON with secret fields redacted
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
```

### SEE ALSO

* [dlut-pan auth](dlut-pan_auth.md)	 - Log in, inspect the saved session or log out.
* [dlut-pan completion](dlut-pan_completion.md)	 - Generate the autocompletion script for the specified shell
* [dlut-pan file](dlut-pan_file.md)	 - List directories and download files.
* [dlut-pan group](dlut-pan_group.md)	 - Join a known group.
* [dlut-pan search](dlut-pan_search.md)	 - Search file names, file contents or groups.
* [dlut-pan space](dlut-pan_space.md)	 - Inspect built-in space identifiers.
