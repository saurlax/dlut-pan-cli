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
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --dry-run            Preview method, URL and body without credentials or network access
      --json               Print API JSON with secret fields redacted
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
```

### SEE ALSO

* [dlut-pan file](dlut-pan_file.md)	 - List directories and download files.
