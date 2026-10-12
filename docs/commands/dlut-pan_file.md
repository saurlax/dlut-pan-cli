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
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --dry-run            Preview method, URL and body without credentials or network access
      --json               Print API JSON with secret fields redacted
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
```

### SEE ALSO

* [dlut-pan](dlut-pan.md)	 - Operate DLUT Pan from your terminal.
* [dlut-pan file download](dlut-pan_file_download.md)	 - Download one file to an explicit local path.
* [dlut-pan file list](dlut-pan_file_list.md)	 - List a directory (defaults to /).
