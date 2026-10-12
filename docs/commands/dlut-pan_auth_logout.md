## dlut-pan auth logout

Remove the saved local login session.

### Synopsis

Delete the session from the system credential store. Later remote commands require auth login again.
This is a local logout; no documented endpoint is available to revoke the server token.

```
dlut-pan auth logout [flags]
```

### Examples

```
  dlut-pan auth logout
```

### Options

```
  -h, --help   help for logout
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

* [dlut-pan auth](dlut-pan_auth.md)	 - Log in, inspect the saved session or log out.
