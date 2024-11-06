package qq

type ErrResult struct {
	Error            int    `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// 临时票据结果
type OpenResult struct {
	OpenId  string `json:"openId"`
	Unionid string `json:"unionid"`
}

type TokenResult struct {
	AccessToken  string `json:"access_token" example:"30F378110D9E34CFE04EDF183165D0D0"`
	ExpiresIn    string `json:"expires_in" example:"7776000"`
	RefreshToken string `json:"refresh_token" example:"9AEA28A71B91AF087CB6D3986BA62D24"`
}

type RefreshResult struct {
	AccessToken  string `json:"access_token" example:"30F378110D9E34CFE04EDF183165D0D0"`
	ExpiresIn    string `json:"expires_in" example:"7776000"`
	RefreshToken string `json:"refresh_token" example:"9AEA28A71B91AF087CB6D3986BA62D24"`
}

type UserResult struct {
	Ret           int    `json:"ret" example:"0"`
	Msg           string `json:"msg" example:""`
	IsLost        int    `json:"is_lost" example:"0"`
	Nickname      string `json:"nickname" example:"静闻弦语"`
	FigureURL     string `json:"figureurl" example:"http://qzapp.qlogo.cn/qzapp/101993700/0661AA6B4844909CCF75C391A6DDB45A/30"`
	FigureURL1    string `json:"figureurl_1" example:"http://qzapp.qlogo.cn/qzapp/101993700/0661AA6B4844909CCF75C391A6DDB45A/50"`
	FigureURL2    string `json:"figureurl_2" example:"http://qzapp.qlogo.cn/qzapp/101993700/0661AA6B4844909CCF75C391A6DDB45A/100"`
	FigureURLQQ1  string `json:"figureurl_qq_1" example:"http://thirdqq.qlogo.cn/g?b=oidb&k=rgN0sF9KSaOsJhuJp0noZg&kti=ZXLrrAAAAAE&s=40&t=1644466423"`
	FigureURLQQ2  string `json:"figureurl_qq_2" example:"http://thirdqq.qlogo.cn/g?b=oidb&k=rgN0sF9KSaOsJhuJp0noZg&kti=ZXLrrAAAAAE&s=100&t=1644466423"`
	FigureURLQQ   string `json:"figureurl_qq" example:"http://thirdqq.qlogo.cn/g?b=oidb&k=rgN0sF9KSaOsJhuJp0noZg&kti=ZXLrrAAAAAE&s=640&t=1644466423"`
	FigureURLType string `json:"figureurl_type" example:"1"`

	// 2022年1月8日后，QQ互联接口将不再返回以下字段真实信息
	//Gender          string `json:"gender" example:"男"`
	//GenderType      int    `json:"gender_type" example:"2"`
	//Province        string `json:"province" example:"广东"`
	//City            string `json:"city" example:"深圳"`
	//Year            string `json:"year" example:"1990"`
	//Constellation   string `json:"constellation" example:"天蝎座"`
	//IsYellowVIP     string `json:"is_yellow_vip" example:"0"`
	//VIP             string `json:"vip" example:"0"`
	//YellowVIPLevel  string `json:"yellow_vip_level" example:"0"`
	//Level           string `json:"level" example:"0"`
	//IsYellowYearVIP string `json:"is_yellow_year_vip" example:"0"`
}
