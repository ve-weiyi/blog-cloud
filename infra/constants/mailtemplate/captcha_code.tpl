<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮箱验证 - Trae Blog</title>
    <style>
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
            color: #2c3e50;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 2px solid #eee;
            padding-bottom: 20px;
        }
        .header h1 {
            color: #16a085;
            font-size: 28px;
            margin: 0;
        }
        .header p {
            color: #7f8c8d;
            margin-top: 10px;
        }
        .content {
            background-color: #f9f9f9;
            padding: 30px;
            border-radius: 12px;
            line-height: 1.8;
            box-shadow: 0 2px 8px rgba(0,0,0,0.05);
        }
        .code {
            font-size: 28px;
            font-weight: bold;
            color: #16a085;
            text-align: center;
            padding: 20px;
            margin: 20px 0;
            background-color: #e8f6f3;
            border-radius: 8px;
            letter-spacing: 2px;
        }
        .footer {
            margin-top: 30px;
            font-size: 13px;
            color: #95a5a6;
            text-align: center;
            line-height: 1.6;
        }
        .social-links {
            margin-top: 20px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Blog</h1>
            <p>你能做的，岂止如此</p>
        </div>
        <div class="content">
            <p>Hi，<strong>{{.Username}}</strong></p>
            <p>{{.Content}}</p>
            <p>您的验证码为：</p>
            <div class="code">{{.Code}}</div>
            <p>温馨提示：</p>
            <ul>
                <li>验证码有效期为15分钟</li>
                <li>如非本人操作，请忽略此邮件</li>
            </ul>
        </div>
        <div class="footer">
            若您有任何问题，可随时通过 <a href="https://blog.veweiyi.cn">647166282@qq.com</a> 联系我们。
            <div class="social-links">
                <p>© 2024 Blog. All rights reserved.</p>
            </div>
        </div>
    </div>
</body>
</html>
