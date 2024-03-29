package oauth

// 注册来源
type RegisterSource string

const (
	RegisterSourceHand       RegisterSource = "hand"        //手动添加
	RegisterSourceMobile     RegisterSource = "mobile"      //手机一键登录
	RegisterSourceSms        RegisterSource = "sms"         //手机短信
	RegisterSourceWxMini     RegisterSource = "wx_mini"     //微信小程序
	RegisterSourceBaiduMini  RegisterSource = "baidu_mini"  //百度小程序
	RegisterSourceAlipayMini RegisterSource = "alipay_mini" //支付宝小程序
	RegisterSourceDouYinMini RegisterSource = "douyin_mini" //抖音小程序
	RegisterSourceWechat     RegisterSource = "weixin"      //微信登录（APP通过微信登录）
	RegisterSourceQQ         RegisterSource = "qq"          //QQ登录
	RegisterSourceAlipay     RegisterSource = "alipay"      //支付宝登录
	RegisterSourceDouYin     RegisterSource = "douyin"      //抖音登录
	RegisterSourceWeibo      RegisterSource = "weibo"       //微博登录
)
