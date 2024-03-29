package weibo

// 微博token返回结果
type TokenResult struct {
	AccessToken string `json:"access_token" example:"2.00OYpWYGPTpttCaf929b916cL6FMXD"` //访问令牌
	RemindIn    string `json:"remind_in" example:"157679999"`                           //过期时间
	ExpiresIn   int    `json:"expires_in" example:"157679999"`                          //过期时间
	Uid         string `json:"uid" example:"6007017078"`                                //用户ID
	IsRealName  string `json:"isRealName" example:"true"`                               //是否实名认证
}

type UserResult struct {
	ID               int64  `json:"id" example:"1404376560"`
	ScreenName       string `json:"screen_name" example:"zaku"`
	Name             string `json:"name" example:"zaku"`
	Province         string `json:"province" example:"11"`
	City             string `json:"city" example:"5"`
	Location         string `json:"location" example:"北京 朝阳区"`
	Description      string `json:"description" example:"人生五十年，乃如梦如幻；有生斯有死，壮士复何憾。"`
	URL              string `json:"url" example:"http://blog.sina.com.cn/zaku"`
	ProfileImageURL  string `json:"profile_image_url" example:"http://tp1.sinaimg.cn/1404376560/50/0/1"`
	Domain           string `json:"domain" example:"zaku"`
	Gender           string `json:"gender" example:"m"`
	FollowersCount   int    `json:"followers_count" example:"1204"`
	FriendsCount     int    `json:"friends_count" example:"447"`
	StatusesCount    int    `json:"statuses_count" example:"2908"`
	FavouritesCount  int    `json:"favourites_count" example:"0"`
	CreatedAt        string `json:"created_at" example:"Fri Aug 28 00:00:00 +0800 2009"`
	Following        bool   `json:"following" example:"false"`
	AllowAllActMsg   bool   `json:"allow_all_act_msg" example:"false"`
	GeoEnabled       bool   `json:"geo_enabled" example:"true"`
	Verified         bool   `json:"verified" example:"false"`
	Status           Status `json:"status"`
	AllowAllComment  bool   `json:"allow_all_comment" example:"true"`
	AvatarLarge      string `json:"avatar_large" example:"http://tp1.sinaimg.cn/1404376560/180/0/1"`
	VerifiedReason   string `json:"verified_reason" example:""`
	FollowMe         bool   `json:"follow_me" example:"false"`
	OnlineStatus     int    `json:"online_status" example:"0"`
	BiFollowersCount int    `json:"bi_followers_count" example:"215"`
}

type Status struct {
	CreatedAt           string        `json:"created_at" example:"Tue May 24 18:04:53 +0800 2011"`
	ID                  int64         `json:"id" example:"11142488790"`
	Text                string        `json:"text" example:"我的相机到了。"`
	Source              string        `json:"source" example:"<a href='http://weibo.com' rel='nofollow'>新浪微博</a>"`
	Favorited           bool          `json:"favorited" example:"false"`
	Truncated           bool          `json:"truncated" example:"false"`
	InReplyToStatusID   string        `json:"in_reply_to_status_id" example:""`
	InReplyToUserID     string        `json:"in_reply_to_user_id" example:""`
	InReplyToScreenName string        `json:"in_reply_to_screen_name" example:""`
	Geo                 interface{}   `json:"geo" example:"null"`
	MID                 string        `json:"mid" example:"5610221544300749636"`
	Annotations         []interface{} `json:"annotations" example:"[]"`
	RepostsCount        int           `json:"reposts_count" example:"5"`
	CommentsCount       int           `json:"comments_count" example:"8"`
}
