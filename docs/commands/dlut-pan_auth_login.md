## dlut-pan auth login

Log in and save the session for subsequent commands.

### Synopsis

Enter your account and password at terminal prompts, or supply --username and --password-stdin.
The returned token and API endpoint are saved in the system credential store; passwords are discarded.
For a browser CAS session, --token-stdin imports its token without remote validation.
A successful login replaces the previous saved session. Later commands require no credential input.

```
dlut-pan auth login [flags]
```

### Examples

```
  dlut-pan auth login
  dlut-pan auth login --username YOUR_ACCOUNT
  dlut-pan file list /
```

### Options

```
      --auth-url string   Login endpoint (DLUT_PAN_AUTH_URL or http://202.118.66.144:8088/v1/auth/sign_in)
      --device string     Login device, link name and LDAP identifier (default "dlut-pan-cli")
  -h, --help              help for login
      --password-stdin    Read login password from stdin instead of environment or terminal prompt
      --token-stdin       Import a browser session token from stdin and save it (not remotely validated)
      --username string   Student/staff account (DLUT_PAN_USERNAME or terminal prompt)
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

* [dlut-pan auth](dlut-pan_auth.md)	 - Log in, inspect the saved session or log out.
