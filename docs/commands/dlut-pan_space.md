## dlut-pan space

Inspect built-in space identifiers.

### Synopsis

Display documented space IDs without contacting the server. Personal space discovery has no documented API;
provide your own root ID with --root-id or DLUT_PAN_ROOT_ID.

```
dlut-pan space [flags]
```

### Options

```
  -h, --help   help for space
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
* [dlut-pan space naosi](dlut-pan_space_naosi.md)	 - Show the NAOSI group and root IDs.
