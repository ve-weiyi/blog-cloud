package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/logic"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type AppService struct {
	svcCtx              *svc.ServiceContext        //持有的repository层引用
	WebsiteService      *logic.WebsiteService      //博客
	AuthService         *logic.AuthService         //登录权限认证
	UserService         *logic.UserService         //用户
	UserAccountService  *logic.UserAccountService  //用户信息
	CaptchaService      *logic.CaptchaService      //验证码
	UploadService       *logic.UploadService       //文件上传
	ChatRecordService   *logic.ChatRecordService   //聊天记录
	RoleService         *logic.RoleService         //角色
	MenuService         *logic.MenuService         //菜单
	ApiService          *logic.ApiService          //接口
	ArticleService      *logic.ArticleService      //文章
	TagService          *logic.TagService          //文章标签
	CategoryService     *logic.CategoryService     //文章分类
	CommentService      *logic.CommentService      //评论
	PhotoService        *logic.PhotoService        //相片
	PhotoAlbumService   *logic.PhotoAlbumService   //相册
	PageService         *logic.PageService         //页面
	TalkService         *logic.TalkService         //说说
	FriendLinkService   *logic.FriendLinkService   //友链
	OperationLogService *logic.OperationLogService //操作记录
	RemarkService       *logic.RemarkService       //留言
	AIChatService       *logic.AIService           //AI聊天
}

func NewService(svcCtx *svc.ServiceContext) *AppService {
	return &AppService{
		svcCtx:              svcCtx,
		WebsiteService:      logic.NewWebsiteService(svcCtx),
		AuthService:         logic.NewAuthService(svcCtx),
		UserService:         logic.NewUserService(svcCtx),
		UserAccountService:  logic.NewUserAccountService(svcCtx),
		CaptchaService:      logic.NewCaptchaService(svcCtx),
		UploadService:       logic.NewUploadService(svcCtx),
		ChatRecordService:   logic.NewChatRecordService(svcCtx),
		RoleService:         logic.NewRoleService(svcCtx),
		MenuService:         logic.NewMenuService(svcCtx),
		ApiService:          logic.NewApiService(svcCtx),
		ArticleService:      logic.NewArticleService(svcCtx),
		TagService:          logic.NewTagService(svcCtx),
		CategoryService:     logic.NewCategoryService(svcCtx),
		CommentService:      logic.NewCommentService(svcCtx),
		PhotoService:        logic.NewPhotoService(svcCtx),
		PhotoAlbumService:   logic.NewPhotoAlbumService(svcCtx),
		PageService:         logic.NewPageService(svcCtx),
		TalkService:         logic.NewTalkService(svcCtx),
		FriendLinkService:   logic.NewFriendLinkService(svcCtx),
		OperationLogService: logic.NewOperationLogService(svcCtx),
		RemarkService:       logic.NewRemarkService(svcCtx),
		AIChatService:       logic.NewAIService(svcCtx),
	}
}
