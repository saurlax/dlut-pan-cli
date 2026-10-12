## dlut-pan auth

Log in, inspect the saved session or log out.

### Synopsis

Log in once before using remote commands. Sessions are saved in the system credential store.
Passwords are never saved. Use status to inspect local session presence and logout to remove it.

```
dlut-pan auth [flags]
```

### Options

```
  -h, --help   help for auth
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

* [dlut-pan](dlut-pan.md)	 - Operate DLUT Pan from your terminal.
* [dlut-pan auth login](dlut-pan_auth_login.md)	 - Log in and save the session for subsequent commands.
* [dlut-pan auth logout](dlut-pan_auth_logout.md)	 - Remove the saved local login session.
* [dlut-pan auth status](dlut-pan_auth_status.md)	 - Show whether a local login session is saved.
