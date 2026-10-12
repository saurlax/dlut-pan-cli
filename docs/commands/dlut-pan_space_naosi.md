## dlut-pan space naosi

Show the NAOSI group and root IDs.

### Synopsis

Print known NAOSI identifiers. This does not check or modify your membership.

```
dlut-pan space naosi [flags]
```

### Examples

```
  dlut-pan space naosi --json
```

### Options

```
  -h, --help   help for naosi
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

* [dlut-pan space](dlut-pan_space.md)	 - Inspect built-in space identifiers.
