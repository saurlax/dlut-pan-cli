## dlut-pan auth status

Show whether a local login session is saved.

### Synopsis

Inspect the system credential store without network requests or displaying the token.
Session presence does not prove the server still accepts it; expired sessions require auth login again.

```
dlut-pan auth status [flags]
```

### Examples

```
  dlut-pan auth status
```

### Options

```
  -h, --help   help for status
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
