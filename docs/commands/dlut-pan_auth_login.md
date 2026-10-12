## dlut-pan auth login

Verify account credentials using the sign-in API.

### Synopsis

Authenticate with --username or DLUT_PAN_USERNAME. The password is read from a hidden terminal prompt,
DLUT_PAN_PASSWORD, or --password-stdin. The returned token is discarded when this process exits.
Use the same account options on a data command to authenticate and perform the operation in one process.

```
dlut-pan auth login [flags]
```

### Examples

```
  dlut-pan auth login --username 2026000000
  dlut-pan file list / --username 2026000000
```

### Options

```
  -h, --help   help for login
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

* [dlut-pan auth](dlut-pan_auth.md)	 - Check authentication without saving credentials.
