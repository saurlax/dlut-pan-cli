## dlut-pan search files

Search file and directory names.

### Synopsis

Search names in the selected space; --path optionally restricts the search scope.
Use --json to inspect total, offset, entries and any additional server fields.

```
dlut-pan search files QUERY [flags]
```

### Examples

```
  dlut-pan search files '操作系统' --limit 20
  dlut-pan search files '操作系统' --offset 20 --json
```

### Options

```
  -h, --help          help for files
      --limit int     Page size (positive; server limits may apply) (default 20)
      --offset int    Pagination offset (zero or greater)
      --path string   Restrict to a remote absolute directory path
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

* [dlut-pan search](dlut-pan_search.md)	 - Search file names, file contents or groups.
