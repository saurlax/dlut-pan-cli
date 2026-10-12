# Authentication

凭证只存在于当前进程中，不保存到文件、系统凭证库或日志。`auth login`
只验证账号；后续命令必须再次提供认证输入。CLI 不会输出登录返回的 token。

## 账号登录

```powershell
dlut-pan auth login --username YOUR_ACCOUNT
dlut-pan file list / --username YOUR_ACCOUNT
dlut-pan search files '操作系统' --username YOUR_ACCOUNT
```

密码在终端中隐藏输入。同一次 `file list` 调用会先登录，然后使用返回的 token
列目录。用户名是学号/工号；群组加入所需的内部 `userid` 是另一个标识。

无人值守调用可使用 `DLUT_PAN_USERNAME` 与 `DLUT_PAN_PASSWORD` 环境变量，或由
密码管理器把密码通过管道传入：

```text
YOUR_SECRET_PROVIDER | dlut-pan file list / --username YOUR_ACCOUNT --password-stdin
```

这里的 `YOUR_SECRET_PROVIDER` 是占位符。stdin 会读取至 EOF，最多 64 KiB，
只移除末尾的一组 CR/LF；密码中的空格和特殊字符保持不变。

## 浏览器 CAS token

在浏览器打开学校云盘的 `/cas` 页面并手动完成统一认证、验证码或二次认证。
登录回云盘主页后，按上游说明在开发者工具中读取嵌套的本地会话结构：

```javascript
const storage = JSON.parse(localStorage.getItem('jStorage') || '{}');
const session = JSON.parse(storage.data || '{}');
// session.token: access token
// session.userid: internal user ID
// session.groupIDs: group membership / root IDs
```

仅在自己的浏览器本地使用这些值，不把它们贴到聊天、issue、日志或项目文件。
如果页面改版，不要假设此结构仍然有效。CLI 不会自动打开浏览器或截取会话。

PowerShell 可从隐藏提示读入 token，再传给 CLI：

```powershell
$panSecureToken = Read-Host 'Pan token' -AsSecureString
$panCredential = [System.Management.Automation.PSCredential]::new('token', $panSecureToken)
$panCredential.GetNetworkCredential().Password | dlut-pan file list / --token-stdin
Remove-Variable panSecureToken, panCredential
```

脚本也可从受控的父进程传入 `DLUT_PAN_TOKEN`。`--token` 支持显式覆盖，但会让
token 进入命令行参数和可能的 shell 历史，优先使用 stdin 或环境变量。

## 优先级

1. 显式 `--username` 或 `--password-stdin` 选择本次账号登录。
2. 否则使用 `--token-stdin`、`--token` 或 `DLUT_PAN_TOKEN`。
3. 没有 token 时，使用 `DLUT_PAN_USERNAME` 和密码输入登录。

显式账号参数与显式 token 参数冲突时返回用法错误。`--password-stdin`
覆盖 `DLUT_PAN_PASSWORD`，否则环境密码优先于隐藏终端提示。
无终端、无 stdin 密码、无环境密码时，立即失败，不尝试空密码登录。

## 传输与空间

上游说明给出的 API 和登录地址使用 HTTP，密码和 token 放在 query string 中。
CLI 隐藏预览、错误中的凭证并禁止登录请求重定向，但 HTTP 网络传输仍为明文。
仅在可信网络使用；如果学校提供经过确认的 HTTPS 接口，可用 `--base-url`
和 `--auth-url` 覆盖。CLI 不禁用证书验证，不自动回退到 HTTP。

默认空间是 NAOSI；未加入时可能返回 403。可以显式加入：

```text
dlut-pan group join ej7f8s176d8f --user-id YOUR_INTERNAL_USER_ID --yes
```

`relation.is_activated` 为 `true` 表示加入，为 `false` 表示等待审批。
CLI 从不自动加入群组或猜测个人 root。用 `--root-id YOUR_ROOT` 指定个人空间。
