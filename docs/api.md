# API notes

Protocol source: [dlut-pan API reference](https://github.com/NAOSI-DLUT/dlut-skills/blob/main/skills/dlut-pan/references/api-reference.md), reviewed 2026-10-12.

| CLI operation | Method and endpoint | Parameters/body |
| --- | --- | --- |
| Account authentication | `POST http://202.118.66.144:8088/v1/auth/sign_in` | Query: username, password, link_device, link_name, ldap_name; empty body |
| Directory listing | `POST /v1/fileops/list_folder` | Query: root_id, path, token, locale=zh_CN; empty body |
| Name search | `GET /v1/search/files` | Query: query, root_id, token, path (optional), offset, limit |
| Full-text search | `GET /v1/search/full_text_files` | Same search parameters |
| Group search | `GET /v1/search/groups` | Query: query, token, offset, limit; no root_id |
| Group join | `POST /v1/users/{user_id}/groups` | Query: group_id, token; JSON body: remarks |
| File download | `GET /v1/roots/{root_id}/files/{file_id}` | Query: token; binary response or storage redirect |

The API base is `http://pan.dlut.edu.cn/v1`. Query parameters are URL-encoded.
`auth login` performs authentication and saves the token and API endpoint in the
system credential store. Remote operations require and automatically reuse that
session; they never implicitly sign in. Dry runs require no session and make no
HTTP requests. A 401 prompts re-login; a 403 reports permissions without deleting
the saved session. No token-refresh or server-side logout endpoint is documented.

The client accepts JSON values without forcing undocumented fields into a
model. Known `entries`/`items` arrays with names are rendered as tables; unknown
schemas remain JSON. `--json` preserves server structure while replacing secret
fields (keys containing token, password, secret or cookie). Login only returns
a local success object, never the server's identity/session payload.

HTTP failures report status and a local hint, without response-body text or
credential-bearing request URLs. A successful login response must contain a
nonempty top-level token. The maximum JSON response is 16 MiB. API/business
errors returned with a success HTTP status and an undocumented schema may still
appear as JSON; inspect the response when integrating new server versions.

Only file downloads follow redirects. Cross-origin redirects remove query keys
`token`, `password` and `username` while retaining storage signatures. HTTPS to
HTTP download redirects are refused. Authentication and metadata redirects are
reported as errors, preventing credentials from being forwarded.

Downloads use an adjacent temporary file, verify transfer completion and content
length when provided, and clean up incomplete files. Existing paths are refused
atomically using a hard link unless `--overwrite` was explicitly selected. This
requires a local filesystem with hard-link support; unsupported filesystems
return an error rather than falling back to a race-prone overwrite. Overwrite
uses the operating system's rename semantics. HTML responses are rejected as
possible login pages, so intentionally downloading an HTML file is currently
unsupported. ZIP, PDF and other binary responses stream without memory buffering.

No live-account verification was performed during initial implementation.
Fixtures validate the published protocol, not upstream service availability,
token lifetime, access rights or index coverage.
