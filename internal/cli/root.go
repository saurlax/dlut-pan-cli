// Package cli defines an independently constructed command tree for each run.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/saurlax/dlut-pan-cli/internal/api"
	"github.com/saurlax/dlut-pan-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type usageError struct{ error }

func invalid(message string) error { return usageError{errors.New(message)} }

type options struct {
	baseURL, authURL, rootID, token, username, device string
	tokenStdin, passwordStdin, json, dryRun           bool
	timeout                                           time.Duration
}

// Execute returns 0 for success, 2 for usage errors, and 1 for runtime errors.
func Execute(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer, version string) int {
	root := New(version)
	root.SetArgs(args)
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	cmd, err := root.ExecuteContextC(ctx)
	if err == nil {
		return 0
	}
	fmt.Fprintln(errOut, "Error:", err)
	var usage usageError
	if errors.As(err, &usage) || strings.HasPrefix(err.Error(), "unknown command") {
		if cmd == nil {
			cmd = root
		}
		fmt.Fprintln(errOut, "Run '"+cmd.CommandPath()+" --help' for usage.")
		return 2
	}
	return 1
}

func args(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, values []string) error {
		if err := cobra.ExactArgs(n)(cmd, values); err != nil {
			return usageError{err}
		}
		for _, value := range values {
			if strings.TrimSpace(value) == "" {
				return invalid("arguments must not be empty")
			}
		}
		return nil
	}
}

func group(use, short, long string) *cobra.Command {
	return &cobra.Command{Use: use, Short: short, Long: long, Args: args(0), RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() }}
}

func New(version string) *cobra.Command {
	o := &options{}
	root := group("dlut-pan", "Operate DLUT Pan from your terminal.", "Unofficial CLI for Dalian University of Technology's cloud drive.\nList folders, search files or groups, join groups and download individual files.\nCredentials stay in memory. Use DLUT_PAN_TOKEN or --username with a password prompt.\nThe default root is NAOSI; use --root-id for your personal or another group space.\nRun any command group without a subcommand to see its help.")
	root.Version = version
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return usageError{err} })
	f := root.PersistentFlags()
	f.StringVar(&o.baseURL, "base-url", "", "API base URL (DLUT_PAN_BASE_URL or "+config.BaseURL+")")
	f.StringVar(&o.authURL, "auth-url", "", "Login endpoint (DLUT_PAN_AUTH_URL or "+config.AuthURL+")")
	f.StringVar(&o.rootID, "root-id", "", "Space root ID (DLUT_PAN_ROOT_ID or NAOSI "+config.NAOSIRoot+")")
	f.StringVar(&o.token, "token", "", "Access token; prefer DLUT_PAN_TOKEN or --token-stdin to avoid shell history")
	f.BoolVar(&o.tokenStdin, "token-stdin", false, "Read a token from stdin instead of the flag or environment")
	f.StringVar(&o.username, "username", "", "Student/staff account (DLUT_PAN_USERNAME); authenticate for this command")
	f.BoolVar(&o.passwordStdin, "password-stdin", false, "Read the login password from stdin instead of environment or terminal prompt")
	f.StringVar(&o.device, "device", "dlut-pan-cli", "Login device, link name and LDAP identifier")
	f.DurationVar(&o.timeout, "timeout", 30*time.Second, "HTTP request timeout, including download body (e.g. 30s or 10m)")
	f.BoolVar(&o.json, "json", false, "Print API JSON with secret fields redacted")
	f.BoolVar(&o.dryRun, "dry-run", false, "Preview method, URL and body without credentials or network access")
	root.MarkFlagsMutuallyExclusive("token", "token-stdin")
	root.MarkFlagsMutuallyExclusive("token-stdin", "password-stdin")
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if err := cmd.ValidateRequiredFlags(); err != nil {
			return usageError{err}
		}
		if err := cmd.ValidateFlagGroups(); err != nil {
			return usageError{err}
		}
		if o.timeout <= 0 {
			return invalid("--timeout must be positive")
		}
		for _, raw := range []string{o.base(), o.auth()} {
			if err := api.ValidateURL(raw); err != nil {
				return usageError{err}
			}
		}
		return nil
	}

	auth := group("auth", "Check authentication without saving credentials.", "Use an existing browser token for data commands, or authenticate with your account per invocation.\nNo passwords or tokens are written to disk or printed. See docs/authentication.md for CAS token input.")
	login := &cobra.Command{Use: "login", Short: "Verify account credentials using the sign-in API.", Long: "Authenticate with --username or DLUT_PAN_USERNAME. The password is read from a hidden terminal prompt,\nDLUT_PAN_PASSWORD, or --password-stdin. The returned token is discarded when this process exits.\nUse the same account options on a data command to authenticate and perform the operation in one process.", Example: "  dlut-pan auth login --username 2026000000\n  dlut-pan file list / --username 2026000000", Args: args(0), RunE: func(cmd *cobra.Command, _ []string) error {
		if o.dryRun {
			return o.preview(cmd, "POST", o.auth(), url.Values{"username": {"[REDACTED]"}, "password": {"[REDACTED]"}, "link_device": {o.device}, "link_name": {o.device}, "ldap_name": {o.device}}, nil)
		}
		client := o.client()
		result, err := o.login(cmd, client)
		if err != nil {
			return err
		}
		return o.render(cmd, result)
	}}
	auth.AddCommand(login)

	file := group("file", "List directories and download files.", "Paths are remote absolute paths, such as / or /courses. IDs come from list or search results.\nThe default root is NAOSI; permission errors require membership or an explicit --root-id.\nDownloads save one specified file; there is no recursive harvesting.")
	list := &cobra.Command{Use: "list [PATH]", Aliases: []string{"ls"}, Short: "List a directory (defaults to /).", Long: "List a single directory using POST /fileops/list_folder with an empty body.\nUse --json to inspect all returned metadata and find file IDs for downloading.", Example: "  dlut-pan file list /\n  dlut-pan file list /courses --root-id YOUR_ROOT --json", Args: func(cmd *cobra.Command, values []string) error {
		if len(values) > 1 {
			return invalid("expected at most one remote path")
		}
		if len(values) == 1 {
			return remotePath(values[0])
		}
		return nil
	}, RunE: func(cmd *cobra.Command, values []string) error {
		path := "/"
		if len(values) == 1 {
			path = values[0]
		}
		if err := segment(o.space()); err != nil {
			return err
		}
		return o.request(cmd, "POST", "/fileops/list_folder", url.Values{"root_id": {o.space()}, "path": {path}, "locale": {"zh_CN"}}, "")
	}}
	var output string
	var overwrite bool
	download := &cobra.Command{Use: "download FILE_ID --output FILE", Short: "Download one file to an explicit local path.", Long: "Stream one file to --output. Redirects to storage servers are supported.\nAn existing file is preserved unless --overwrite is passed. Interrupted downloads discard partial bytes.\nThe parent directory must exist. Increase --timeout for large files. HTML login pages are rejected.", Example: "  dlut-pan file download FILE_ID --output lecture.pdf\n  dlut-pan file download FILE_ID --output lecture.pdf --timeout 10m --overwrite\n  dlut-pan file download FILE_ID --output lecture.pdf --dry-run", Args: args(1), RunE: func(cmd *cobra.Command, values []string) error {
		if err := segment(values[0]); err != nil {
			return err
		}
		if err := segment(o.space()); err != nil {
			return err
		}
		if strings.TrimSpace(output) == "" || output == "-" {
			return invalid("--output must name a local file (stdout streaming is not supported)")
		}
		if o.dryRun {
			return o.preview(cmd, "GET", o.base()+"/roots/"+url.PathEscape(o.space())+"/files/"+url.PathEscape(values[0]), nil, map[string]any{"output": output, "overwrite": overwrite})
		}
		client, err := o.authenticated(cmd)
		if err != nil {
			return err
		}
		n, err := client.Download(cmd.Context(), o.space(), values[0], output, overwrite)
		if err != nil {
			return err
		}
		return o.render(cmd, map[string]any{"output": output, "bytes": n})
	}}
	download.Flags().StringVarP(&output, "output", "o", "", "Local destination filename (required)")
	download.Flags().BoolVar(&overwrite, "overwrite", false, "Replace an existing destination only after a complete download")
	_ = download.MarkFlagRequired("output")
	_ = download.MarkFlagFilename("output")
	file.AddCommand(list, download)

	search := group("search", "Search file names, file contents or groups.", "Each command returns one page. Use --offset and --limit to control pagination.\nFile searches require a space root; group search does not. No results are automatically downloaded.")
	for _, spec := range []struct{ name, endpoint, short, long string }{
		{"files", "files", "Search file and directory names.", "Search names in the selected space; --path optionally restricts the search scope."},
		{"full-text", "full_text_files", "Search indexed file contents.", "Search the server's full-text index in the selected space. Index coverage depends on the server."},
		{"groups", "groups", "Search discoverable groups.", "Search group metadata without a root ID. Joining a group is a separate explicit command."},
	} {
		var offset, limit int
		var path string
		command := &cobra.Command{Use: spec.name + " QUERY", Short: spec.short, Long: spec.long + "\nUse --json to inspect total, offset, entries and any additional server fields.", Example: "  dlut-pan search " + spec.name + " '操作系统' --limit 20\n  dlut-pan search " + spec.name + " '操作系统' --offset 20 --json", Args: args(1), RunE: func(cmd *cobra.Command, values []string) error {
			if offset < 0 || limit <= 0 {
				return invalid("--offset must be nonnegative and --limit must be positive")
			}
			q := url.Values{"query": {values[0]}, "offset": {fmt.Sprint(offset)}, "limit": {fmt.Sprint(limit)}}
			if spec.name != "groups" {
				if err := segment(o.space()); err != nil {
					return err
				}
				q.Set("root_id", o.space())
				if path != "" {
					if err := remotePath(path); err != nil {
						return err
					}
					q.Set("path", path)
				}
			}
			return o.request(cmd, "GET", "/search/"+spec.endpoint, q, "")
		}}
		command.Flags().IntVar(&offset, "offset", 0, "Pagination offset (zero or greater)")
		command.Flags().IntVar(&limit, "limit", 20, "Page size (positive; server limits may apply)")
		if spec.name != "groups" {
			command.Flags().StringVar(&path, "path", "", "Restrict to a remote absolute directory path")
		}
		search.AddCommand(command)
	}

	groups := group("group", "Join a known group.", "Use search groups to find group IDs. A group ID differs from its space root ID.\nNAOSI group: "+config.NAOSIGroup+"; NAOSI root: "+config.NAOSIRoot+".\nMembership is changed only when you explicitly run join with --yes.")
	var userID, remarks string
	var yes bool
	join := &cobra.Command{Use: "join GROUP_ID --user-id USER_ID", Short: "Join a group or submit a membership request.", Long: "POST a membership request for the authenticated user's ID. Supply the user ID from your browser session\n(userid), which can differ from your student/staff account. --yes confirms the membership change.\nrelation.is_activated=true means joined; false means awaiting approval. Joining is never automatic.", Example: "  dlut-pan group join " + config.NAOSIGroup + " --user-id YOUR_USER_ID --yes\n  dlut-pan group join GROUP_ID --user-id YOUR_USER_ID --remarks '课程学习' --dry-run", Args: args(1), RunE: func(cmd *cobra.Command, values []string) error {
		if err := segment(values[0]); err != nil {
			return err
		}
		if err := segment(userID); err != nil {
			return err
		}
		if !yes && !o.dryRun {
			return invalid("joining changes group membership; pass --yes to proceed")
		}
		body, _ := json.Marshal(map[string]string{"remarks": remarks})
		return o.request(cmd, "POST", "/users/"+url.PathEscape(userID)+"/groups", url.Values{"group_id": {values[0]}}, string(body))
	}}
	join.Flags().StringVar(&userID, "user-id", "", "Authenticated user's internal ID, not student/staff number (required)")
	join.Flags().StringVar(&remarks, "remarks", "", "Membership request reason")
	join.Flags().BoolVarP(&yes, "yes", "y", false, "Confirm the group membership change")
	_ = join.MarkFlagRequired("user-id")
	groups.AddCommand(join)

	space := group("space", "Inspect built-in space identifiers.", "Display documented space IDs without contacting the server. Personal space discovery has no documented API;\nprovide your own root ID with --root-id or DLUT_PAN_ROOT_ID.")
	space.AddCommand(&cobra.Command{Use: "naosi", Short: "Show the NAOSI group and root IDs.", Long: "Print known NAOSI identifiers. This does not check or modify your membership.", Example: "  dlut-pan space naosi --json", Args: args(0), RunE: func(cmd *cobra.Command, _ []string) error {
		return o.render(cmd, map[string]any{"name": "NAOSI", "group_id": config.NAOSIGroup, "root_id": config.NAOSIRoot})
	}})
	root.AddCommand(auth, file, search, groups, space)
	return root
}

func remotePath(path string) error {
	if !strings.HasPrefix(path, "/") {
		return invalid("remote path must start with /")
	}
	return nil
}
func segment(id string) error {
	if strings.TrimSpace(id) == "" || id == "." || id == ".." || strings.ContainsAny(id, "/\\?#\r\n") {
		return invalid("IDs must be nonempty single path segments")
	}
	return nil
}
func (o *options) base() string {
	return strings.TrimRight(config.Value(o.baseURL, "DLUT_PAN_BASE_URL", config.BaseURL), "/")
}
func (o *options) auth() string  { return config.Value(o.authURL, "DLUT_PAN_AUTH_URL", config.AuthURL) }
func (o *options) space() string { return config.Value(o.rootID, "DLUT_PAN_ROOT_ID", config.NAOSIRoot) }
func (o *options) client() *api.Client {
	return &api.Client{BaseURL: o.base(), AuthURL: o.auth(), Timeout: o.timeout}
}

func stdinSecret(cmd *cobra.Command) (string, error) {
	const max = 64 << 10
	b, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), max+1))
	if err != nil {
		return "", errors.New("cannot read credential from stdin")
	}
	if len(b) > max {
		return "", invalid("stdin credential exceeds 64 KiB")
	}
	return strings.TrimSuffix(strings.TrimSuffix(string(b), "\n"), "\r"), nil
}

func (o *options) login(cmd *cobra.Command, client *api.Client) (any, error) {
	username := config.Value(o.username, "DLUT_PAN_USERNAME", "")
	if strings.TrimSpace(username) == "" {
		return nil, invalid("provide --username or DLUT_PAN_USERNAME for account authentication")
	}
	password := os.Getenv("DLUT_PAN_PASSWORD")
	if o.passwordStdin {
		var err error
		password, err = stdinSecret(cmd)
		if err != nil {
			return nil, err
		}
	} else if password == "" {
		f, ok := cmd.InOrStdin().(*os.File)
		if !ok || !term.IsTerminal(int(f.Fd())) {
			return nil, invalid("no terminal password prompt available; use DLUT_PAN_PASSWORD or --password-stdin")
		}
		fmt.Fprint(cmd.ErrOrStderr(), "Password: ")
		b, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(cmd.ErrOrStderr())
		if err != nil {
			return nil, errors.New("cannot read password from terminal")
		}
		password = string(b)
	}
	if password == "" {
		return nil, invalid("password must not be empty")
	}
	if strings.TrimSpace(o.device) == "" {
		return nil, invalid("--device must not be empty")
	}
	return client.Login(cmd.Context(), username, password, o.device)
}

func (o *options) authenticated(cmd *cobra.Command) (*api.Client, error) {
	c := o.client()
	// Explicit account input selects fresh sign-in even if an old token is in the environment.
	if o.username != "" || o.passwordStdin {
		if o.token != "" || o.tokenStdin {
			return nil, invalid("choose token authentication or username/password authentication")
		}
		_, err := o.login(cmd, c)
		return c, err
	}
	token := config.Value(o.token, "DLUT_PAN_TOKEN", "")
	if o.tokenStdin {
		var err error
		token, err = stdinSecret(cmd)
		if err != nil {
			return nil, err
		}
	}
	if token = strings.TrimSpace(token); token != "" {
		c.Token = token
		return c, nil
	}
	if o.tokenStdin {
		return nil, invalid("stdin token must not be empty")
	}
	_, err := o.login(cmd, c)
	return c, err
}

func (o *options) preview(cmd *cobra.Command, method, endpoint string, q url.Values, body any) error {
	if q == nil {
		q = url.Values{}
	}
	if q.Get("password") == "" {
		q.Set("token", "[REDACTED]")
	}
	return writeJSON(cmd.OutOrStdout(), map[string]any{"method": method, "url": endpoint + "?" + q.Encode(), "body": body, "dry_run": true})
}

func (o *options) request(cmd *cobra.Command, method, path string, q url.Values, body string) error {
	var decoded any
	if body != "" {
		_ = json.Unmarshal([]byte(body), &decoded)
	}
	if o.dryRun {
		return o.preview(cmd, method, o.base()+path, q, decoded)
	}
	c, err := o.authenticated(cmd)
	if err != nil {
		return err
	}
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	value, err := c.Call(cmd.Context(), method, path, q, reader)
	if err != nil {
		return err
	}
	return o.render(cmd, value)
}

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func (o *options) render(cmd *cobra.Command, value any) error {
	value = api.Redact(value)
	if o.json {
		return writeJSON(cmd.OutOrStdout(), value)
	}
	var entries []any
	switch v := value.(type) {
	case []any:
		entries = v
	case map[string]any:
		if list, ok := v["entries"].([]any); ok {
			entries = list
		} else if list, ok := v["items"].([]any); ok {
			entries = list
		} else {
			return writeJSON(cmd.OutOrStdout(), value)
		}
		if total, ok := v["total"]; ok {
			fmt.Fprintf(cmd.OutOrStdout(), "Total: %v\n", total)
		}
	default:
		return writeJSON(cmd.OutOrStdout(), value)
	}
	if len(entries) == 0 {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), "No entries.")
		return err
	}
	// Unknown schemas remain visible as JSON rather than guessing metadata fields.
	for _, entry := range entries {
		m, ok := entry.(map[string]any)
		if !ok || m["name"] == nil {
			return writeJSON(cmd.OutOrStdout(), value)
		}
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tSIZE\tPATH")
	for _, entry := range entries {
		m := entry.(map[string]any)
		typeValue := m["type"]
		if dir, ok := m["is_dir"].(bool); ok {
			typeValue = "file"
			if dir {
				typeValue = "directory"
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", cell(m["id"]), cell(m["name"]), cell(typeValue), cell(m["size"]), cell(m["path"]))
	}
	return w.Flush()
}

func cell(value any) string {
	if value == nil {
		return "-"
	}
	if m, ok := value.(map[string]any); ok {
		if display, exists := m["display_value"]; exists {
			value = display
		} else if bytes, exists := m["bytes"]; exists {
			value = bytes
		} else if name, exists := m["name"]; exists {
			value = name
		} else {
			b, _ := json.Marshal(m)
			value = string(b)
		}
	}
	return strings.Map(func(r rune) rune {
		if r < 32 || (r >= 127 && r <= 159) {
			return ' '
		}
		return r
	}, fmt.Sprint(value))
}
