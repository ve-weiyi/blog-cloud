package mail

// 验证码邮件内容
type CaptchaEmail struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

type EmailContent struct {
	Title       string `json:"title"`        // 标题
	HeadBg      string `json:"head_bg"`      // 头部背景图
	HeadTips    string `json:"head_tips"`    // 头部提示
	DearUser    string `json:"dear"`         // 亲爱的用户
	Content     string `json:"content"`      // 内容
	ButtonTips  string `json:"button_tips"`  // 按钮提示
	ButtonLink  string `json:"button_url"`   // 按钮链接
	ContactUs   string `json:"contact_us"`   // 联系我们
	ContactLink string `json:"contact_link"` // 联系链接
}

func NewEmailContent() *EmailContent {
	return &EmailContent{
		Title:       "验证码邮件",
		HeadBg:      "",
		HeadTips:    "",
		DearUser:    "你好",
		Content:     "未填写的内容",
		ButtonTips:  "点击重置密码",
		ButtonLink:  "http://localhost:8888/blog",
		ContactUs:   "647166282@qq.com",
		ContactLink: "http://localhost:8888/blog",
	}
}

const TempSimpleCode = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        table,
        thead,
        tbody,
        tfoot,
        tr,
        td {
            padding: 0;
            margin: 0;
            border: 0;
            border-spacing: 0px;
        }
    </style>
</head>

<body>
    <div>
        <table style="width:700px;margin:0 auto;border:1px solid #f0f0f0;background-color: #fff;">
            <thead style='
                background-size:contain;
                background-image: url("{{.HeadBg}}");
                '
                >
                <tr>
                    <td>
                        <div style="box-sizing: border-box; height:130px;padding-top: 77px;padding-left:36px;">
                            <img src="" alt="" style="display: inline-block;width:40px;">
                        </div>
                        <div style="box-sizing: border-box; width:260px;height:130px;padding-left:36px;"> <span
                                style="font-family: 'PingFang SC'; font-weight: 500; font-size: 16px; line-height: 22px; color: #007DFF;">
                                {{.HeadTips}}
                            </span>
                        </div>
                    </td>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>
                        <div style="padding: 24px;">
                            <p style="
                            margin-top: 0;
                            margin-bottom: 16px;
                            font-style: normal;
                            font-weight: 600;
                            font-size: 14px;
                            line-height: 20px;
                            color: #222222;">
                                {{.DearUser}}
                            </p>
                            <p style="
                            margin-top: 0;
                            margin-bottom: 0;
                            font-style: normal;
                            font-weight: 600;
                            font-size: 14px;
                            line-height: 20px;
                            color: #222222;">
                                {{.Content}}
                            </p>
                        </div>
                    </td>
                </tr>
                <tr>
                    <td>
                        <div style="padding-top: 24px;padding-bottom:0;">
                            <a href=" {{.ButtonLink}}" style="text-decoration: none;">
                                <div style="
                                width: 160px;
                                height: 40px;
                                margin: 0 auto;
                                background: #128FFF;
                                border-radius: 8px;
                                cursor: pointer;
                                text-align: center;
                                line-height: 40px;
                                font-style: normal;
                                font-weight: 400;
                                font-size: 16px;
                                color:#fff;
                                "oncontextmenu="return false;">
                                    {{.ButtonTips}}
                                </div>
                            </a>
                        </div>
                    </td>
                </tr>
            </tbody>
            <tfoot>
                <tr>
                    <td>
                        <div style="
                        padding-top: 48px;
                        margin-bottom: 8px;
                        text-align: center;
                        font-style: normal;
                        font-weight: 400;
                        font-size: 12px;
                        line-height: 17px;
                        color: #BABABA;
                        ">
                            <span>若您有任何问题，可随时通过 <a href="{{.ContactLink}}">{{.ContactUs}}</a> 联系我们。</span>
                        </div>
                    </td>
                </tr>
            </tfoot>
        </table>
    </div>
</body>
</html>
`

const TempForgetPassword = `
您的账号{{.Username}}正在重置密码，验证码为 {{.Code}}，有效期15分钟！
`

const TempRegister = `
您好{{.Username}}，欢迎注册我的博客平台。您的验证码为 {{.Code}}，有效期15分钟！
`

const TempBind = `
您好{{.Username}}，您的账号正在尝试修改绑定邮箱。您的验证码为 {{.Code}}，有效期15分钟！
`
