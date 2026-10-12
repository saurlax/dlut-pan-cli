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
      --base-url string    API base URL (DLUT_PAN_BASE_URL or http://pan.dlut.edu.cn/v1)
      --dry-run            Preview method, URL and body without credentials or network access
      --json               Print API JSON with secret fields redacted
      --root-id string     Space root ID (DLUT_PAN_ROOT_ID or NAOSI ak7an4chdvtq)
      --timeout duration   HTTP request timeout, including download body (e.g. 30s or 10m) (default 30s)
```

### SEE ALSO

* [dlut-pan](dlut-pan.md)	 - Operate DLUT Pan from your terminal.
* [dlut-pan group join](dlut-pan_group_join.md)	 - Join a group or submit a membership request.
