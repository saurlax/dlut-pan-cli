## dlut-pan group join

Join a group or submit a membership request.

### Synopsis

POST a membership request for the authenticated user's ID. Supply the user ID from your browser session
(userid), which can differ from your student/staff account. --yes confirms the membership change.
relation.is_activated=true means joined; false means awaiting approval. Joining is never automatic.

```
dlut-pan group join GROUP_ID --user-id USER_ID [flags]
```

### Examples

```
  dlut-pan group join ej7f8s176d8f --user-id YOUR_USER_ID --yes
  dlut-pan group join GROUP_ID --user-id YOUR_USER_ID --remarks '课程学习' --dry-run
```

### Options

```
  -h, --help             help for join
      --remarks string   Membership request reason
      --user-id string   Authenticated user's internal ID, not student/staff number (required)
  -y, --yes              Confirm the group membership change
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

* [dlut-pan group](dlut-pan_group.md)	 - Join a known group.
