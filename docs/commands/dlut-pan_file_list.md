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
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --dry-run            Preview method, URL and body without credentials or network access
      --json               Print API JSON with secret fields redacted
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
```

### SEE ALSO

* [dlut-pan file](dlut-pan_file.md)	 - List directories and download files.
