## dlut-pan file download

Download one file to an explicit local path.

### Synopsis

Stream one file to --output. Redirects to storage servers are supported.
An existing file is preserved unless --overwrite is passed. Interrupted downloads discard partial bytes.
The parent directory must exist. Increase --timeout for large files. HTML login pages are rejected.

```
dlut-pan file download FILE_ID --output FILE [flags]
```

### Examples

```
  dlut-pan file download FILE_ID --output lecture.pdf
  dlut-pan file download FILE_ID --output lecture.pdf --timeout 10m --overwrite
  dlut-pan file download FILE_ID --output lecture.pdf --dry-run
```

### Options

```
  -h, --help            help for download
  -o, --output string   Local destination filename (required)
      --overwrite       Replace an existing destination only after a complete download
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
