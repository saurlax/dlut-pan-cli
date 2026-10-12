## dlut-pan search groups

Search discoverable groups.

### Synopsis

Search group metadata without a root ID. Joining a group is a separate explicit command.
Use --json to inspect total, offset, entries and any additional server fields.

```
dlut-pan search groups QUERY [flags]
```

### Examples

```
  dlut-pan search groups '操作系统' --limit 20
  dlut-pan search groups '操作系统' --offset 20 --json
```

### Options

```
  -h, --help         help for groups
      --limit int    Page size (positive; server limits may apply) (default 20)
      --offset int   Pagination offset (zero or greater)
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

* [dlut-pan search](dlut-pan_search.md)	 - Search file names, file contents or groups.
