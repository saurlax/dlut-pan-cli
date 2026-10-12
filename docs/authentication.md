# Authentication

先登录一次，CLI 自动保存登录态；后续命令直接使用已保存的会话，不需要再提供
学号、密码或 token，也不会隐式登录或再次提示输入密码。

```powershell
dlut-pan auth login
dlut-pan file list /
dlut-pan search files '操作系统'
dlut-pan auth status
dlut-pan auth logout
```

## 登录与保存

`auth login` 在终端询问账号，并隐藏输入密码。也可指定账号：

```powershell
dlut-pan auth login --username YOUR_ACCOUNT
```

登录成功后，只把 token 和它对应的 API 地址保存到系统凭证库：Windows
Credential Manager、macOS Keychain 或 Linux Secret Service。学号和密码不保存，
token 不打印、不写入项目文件。凭证库不可用或被锁定时，命令报错，不会谎报
登录成功或回退到明文文件。Linux 需要可用的 D-Bus 会话和 Secret Service。

再次登录成功会替换旧会话；登录失败时保留原会话。默认只管理一个当前会话。
后续命令自动使用登录时保存的 API 地址；如果显式选择另一个 `--base-url`，
需要先在该地址登录，防止旧 token 被发送给另一个服务。

无人值守登录可使用 `DLUT_PAN_USERNAME` 和 `DLUT_PAN_PASSWORD`，或者由密码
管理器通过管道传入密码：

```text
YOUR_SECRET_PROVIDER | dlut-pan auth login --username YOUR_ACCOUNT --password-stdin
```

`YOUR_SECRET_PROVIDER` 是占位符。stdin 读取至 EOF，最多 64 KiB，只移除末尾
的一组 CR/LF；密码中的空格和特殊字符保持不变。`--password-stdin` 优先于
环境密码，否则环境密码优先于终端提示。这些输入只在 `auth login` 中使用。

## 状态、过期与退出

- `auth status` 查询本地是否保存会话，不输出 token，也不向服务端发请求。
  已保存不代表会话尚未过期；上游没有已确认的 token 状态或刷新接口。
- 未登录执行远程命令时，立即提示运行 `dlut-pan auth login`，不向服务器发送请求。
- 服务端返回 401 时，提示重新登录。403 表示权限或群组问题，不擅自删除会话。
- `auth logout` 删除本地会话；之后必须重新登录。没有会话时退出仍然成功。
  这是本地退出，上游没有已确认的服务端注销接口。

帮助、补全、`space naosi` 和 `--dry-run` 不需要登录。文件/搜索/群组命令不再
接受 `--username`、`--password-stdin` 或 token 参数；`DLUT_PAN_TOKEN` 不再使用。

## 导入浏览器 CAS 会话

在浏览器打开学校云盘 `/cas` 页面，手动完成统一认证、验证码或二次认证。
回到云盘主页后，按上游说明在自己的开发者工具中读取嵌套会话结构：

```javascript
const storage = JSON.parse(localStorage.getItem('jStorage') || '{}');
const session = JSON.parse(storage.data || '{}');
// session.token: access token
// session.userid: internal user ID
// session.groupIDs: group membership / root IDs
```

仅在自己的浏览器中使用这些值，不贴到聊天、issue 或日志。如果页面改版，
不要假设该结构仍有效。CLI 不会自动打开浏览器或截取会话。

用隐藏输入把浏览器 token 导入并保存一次：

```powershell
$panSecureToken = Read-Host 'Pan token' -AsSecureString
$panCredential = [System.Management.Automation.PSCredential]::new('token', $panSecureToken)
$panCredential.GetNetworkCredential().Password | dlut-pan auth login --token-stdin
Remove-Variable panSecureToken, panCredential
dlut-pan file list /
```

导入仅保存 token，没有服务端验证；第一次远程调用会检验是否被接受。

## 传输与空间

上游说明给出的 API 和登录地址使用 HTTP，密码和 token 位于 query string。
CLI 隐藏错误和预览中的凭证，但 HTTP 网络传输仍为明文。只在可信网络使用；
如果学校提供经过确认的 HTTPS 地址，可在登录时设置 `--base-url` 和 `--auth-url`。
CLI 不禁用证书校验，也不自动回退到 HTTP。

默认空间是 NAOSI。未加入时可能返回 403，可以显式加入：

```text
dlut-pan group join ej7f8s176d8f --user-id YOUR_INTERNAL_USER_ID --yes
```

这里的内部 `userid` 不等于学号。`relation.is_activated=true` 表示加入，
`false` 表示等待审批。CLI 不自动加入群组；个人空间用 `--root-id YOUR_ROOT` 指定。
