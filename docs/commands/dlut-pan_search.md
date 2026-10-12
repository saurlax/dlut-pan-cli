## dlut-pan search

Search file names, file contents or groups.

### Synopsis

Each command returns one page. Use --offset and --limit to control pagination.
File searches require a space root; group search does not. No results are automatically downloaded.

```
dlut-pan search [flags]
```

### Options

```
  -h, --help   help for search
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

* [dlut-pan](dlut-pan.md)	 - Operate DLUT Pan from your terminal.
* [dlut-pan search files](dlut-pan_search_files.md)	 - Search file and directory names.
* [dlut-pan search full-text](dlut-pan_search_full-text.md)	 - Search indexed file contents.
* [dlut-pan search groups](dlut-pan_search_groups.md)	 - Search discoverable groups.
