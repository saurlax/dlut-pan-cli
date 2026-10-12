## dlut-pan file

List directories and download files.

### Synopsis

Paths are remote absolute paths, such as / or /courses. IDs come from list or search results.
The default root is NAOSI; permission errors require membership or an explicit --root-id.
Downloads save one specified file; there is no recursive harvesting.

```
dlut-pan file [flags]
```

### Options

```
  -h, --help   help for file
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
* [dlut-pan file download](dlut-pan_file_download.md)	 - Download one file to an explicit local path.
* [dlut-pan file list](dlut-pan_file_list.md)	 - List a directory (defaults to /).
