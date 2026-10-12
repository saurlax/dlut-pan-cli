## dlut-pan file list

List a directory (defaults to /).

### Synopsis

List a single directory using POST /fileops/list_folder with an empty body.
Use --json to inspect all returned metadata and find file IDs for downloading.

```
dlut-pan file list [PATH] [flags]
```

### Examples

```
  dlut-pan file list /
  dlut-pan file list /courses --root-id YOUR_ROOT --json
```

### Options

```
  -h, --help   help for list
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

* [dlut-pan file](dlut-pan_file.md)	 - List directories and download files.
