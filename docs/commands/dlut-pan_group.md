## dlut-pan group

Join a known group.

### Synopsis

Use search groups to find group IDs. A group ID differs from its space root ID.
NAOSI group: ej7f8s176d8f; NAOSI root: ak7an4chdvtq.
Membership is changed only when you explicitly run join with --yes.

```
dlut-pan group [flags]
```

### Options

```
  -h, --help   help for group
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

* [dlut-pan](dlut-pan.md)	 - Operate DLUT Pan from your terminal.
* [dlut-pan group join](dlut-pan_group_join.md)	 - Join a group or submit a membership request.
