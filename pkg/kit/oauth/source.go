package oauth

// 注册来源
type Platform string

const (
	PlatformHand       = "hand"        //手动添加
	PlatformMobile     = "mobile"      //手机一键登录
	PlatformSms        = "sms"         //手机短信
	PlatformWxMini     = "wx_mini"     //微信小程序
	PlatformBaiduMini  = "baidu_mini"  //百度小程序
	PlatformAlipayMini = "alipay_mini" //支付宝小程序
	PlatformDouYinMini = "douyin_mini" //抖音小程序
	PlatformWechat     = "weixin"      //微信登录（APP通过微信登录）
	PlatformQQ         = "qq"          //QQ登录
	PlatformAlipay     = "alipay"      //支付宝登录
	PlatformDouYin     = "douyin"      //抖音登录
	PlatformWeibo      = "weibo"       //微博登录
)
