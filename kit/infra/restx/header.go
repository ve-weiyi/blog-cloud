package restx

/**
HTTP请求头部中常用的信息包括但不限于以下字段：

Host: 请求的目标主机和端口号。
User-Agent: 发出请求的客户端信息，如浏览器类型、版本等。
Accept: 客户端能够接收的内容类型。
Accept-Language: 客户端优先接受的语言。
Accept-Encoding: 客户端能够解码的编码方式，如gzip。
Authorization: 客户端提供的认证信息。
Connection: 控制不同网络连接选项。
Content-Type: 请求体的媒体类型。
Content-Length: 请求体的长度。
Cookie: 服务器发送到客户端的cookie。
Referer: 发起请求的原始URL。
X-Forwarded-For: 通过HTTP代理或负载均衡器发送请求的客户端的原始IP地址。
X-Real-IP: 类似于X-Forwarded-For，通常用于获取客户端的真实IP地址。

RemoteAddr和X-Real-IP在HTTP请求中的区别如下：
RemoteAddr: 是http.Request结构体的一个字段，它直接由Go的HTTP服务器提供，包含了发起请求的客户端的IP地址和端口号。这个地址可能是直接连接到服务器的客户端的地址，或者是最后一个代理服务器的地址，如果请求经过了代理。
X-Real-IP: 是一个非标准的HTTP请求头部字段，通常由反向代理服务器（如Nginx或HAProxy）设置，用于记录原始请求的客户端IP地址。当请求通过一个或多个代理时，X-Real-IP用于传递原始客户端的IP地址
在没有代理的情况下，RemoteAddr通常就是客户端的真实IP地址。但是在使用代理的情况下，RemoteAddr可能是代理服务器的地址，而X-Real-IP则是客户端的真实IP地址。因此，如果你的服务部署在使用了反向代理的环境中，通常需要检查X-Real-IP来获取客户端的真实IP地址。
*/

// grpc请求头部
const (
	// 自定义rpc请求头部，防止和grpc的头部冲突
	HeaderRPCRemoteAgent = "rpc-remote-agent"
	HeaderRPCRemoteIP    = "rpc-remote-ip"
)

// HTTP请求头部字段
const (
	// 通用请求头部
	HeaderRemoteAddr    = "remote-addr"
	HeaderUserAgent     = "user-agent"
	HeaderReferer       = "referer"
	HeaderXForwardedFor = "x-forwarded-for"
	HeaderXRealIP       = "x-real-ip"

	// 自定义请求头部
	HeaderAppName    = "app-name"
	HeaderAppVersion = "app-version"
	HeaderOsType     = "os-type"
	HeaderOsVersion  = "os-version"
	HeaderCountry    = "country"
	HeaderLanguage   = "language"
	HeaderTimezone   = "timezone"
	HeaderTimestamp  = "timestamp"

	// 用户授权信息（jwt）
	HeaderAuthorization = "authorization"

	// 用户授权信息
	HeaderUid   = "uid"
	HeaderToken = "token"

	// 管理平台授权信息（暂不使用）
	HeaderXAdminId    = "x-admin-id"
	HeaderXAdminToken = "x-admin-token"

	// 游客授权信息 token = md5(terminal-id,timestamp)
	HeaderXTerminalId    = "x-terminal-id"
	HeaderXTerminalToken = "x-terminal-token"

	// 防重放限制 sign=md5(id+ts+secret)
	HeaderXRequestId   = "x-request-id"
	HeaderXRequestTs   = "x-request-ts"
	HeaderXRequestSign = "x-request-sign"
)

var HeaderFields = []string{
	HeaderRemoteAddr,
	HeaderUserAgent,
	HeaderReferer,
	HeaderXForwardedFor,
	HeaderXRealIP,

	HeaderAppName,
	HeaderAppVersion,
	HeaderOsType,
	HeaderOsVersion,
	HeaderCountry,
	HeaderLanguage,
	HeaderTimezone,
	HeaderTimestamp,

	HeaderAuthorization,
	HeaderUid,
	HeaderToken,

	HeaderXAdminId,
	HeaderXAdminToken,

	HeaderXTerminalId,
	HeaderXTerminalToken,

	HeaderXRequestId,
	HeaderXRequestTs,
	HeaderXRequestSign,
}
