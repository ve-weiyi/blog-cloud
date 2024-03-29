package repository

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/logic"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/svc"
)

// 对应go-zero rpc层服务
type AppRepository struct {
	svcCtx                     *svc.RepositoryContext            //持有的repository层引用
	ApiRepository              *logic.ApiRepository              //api路由
	ArticleRepository          *logic.ArticleRepository          //文章
	CasbinRuleRepository       *logic.CasbinRuleRepository       //casbin_rule
	MenuRepository             *logic.MenuRepository             //菜单
	RoleRepository             *logic.RoleRepository             //角色
	RoleApiRepository          *logic.RoleApiRepository          //角色-api关联
	RoleMenuRepository         *logic.RoleMenuRepository         //角色-菜单关联
	UserAccountRepository      *logic.UserAccountRepository      //用户登录信息
	UserInformationRepository  *logic.UserInformationRepository  //用户信息
	UserLoginHistoryRepository *logic.UserLoginHistoryRepository //用户登录历史
	UserOauthRepository        *logic.UserOauthRepository        //第三方登录信息
	UserRoleRepository         *logic.UserRoleRepository         //用户-角色关联
	CategoryRepository         *logic.CategoryRepository         //文章分类
	FriendLinkRepository       *logic.FriendLinkRepository       //文章分类
	TagRepository              *logic.TagRepository              //文章标签
	PageRepository             *logic.PageRepository             //页面
	CommentRepository          *logic.CommentRepository          //评论
	PhotoRepository            *logic.PhotoRepository            //照片
	PhotoAlbumRepository       *logic.PhotoAlbumRepository       //相册
	TalkRepository             *logic.TalkRepository             //说说
	ArticleTagRepository       *logic.ArticleTagRepository       //文章标签映射
	ChatRecordRepository       *logic.ChatRecordRepository       //聊天记录
	UniqueViewRepository       *logic.UniqueViewRepository       //页面访问数量
	OperationLogRepository     *logic.OperationLogRepository     //操作记录
	RemarkRepository           *logic.RemarkRepository           //留言
	WebsiteConfigRepository    *logic.WebsiteConfigRepository    //网站设置
	UploadRecordRepository     *logic.UploadRecordRepository     //文件上传
	ChatSessionRepository      *logic.ChatSessionRepository      //聊天会话
	ChatMessageRepository      *logic.ChatMessageRepository      //聊天消息
}

func NewRepository(svcCtx *svc.RepositoryContext) *AppRepository {
	return &AppRepository{
		svcCtx:                     svcCtx,
		ApiRepository:              logic.NewApiRepository(svcCtx),
		ArticleRepository:          logic.NewArticleRepository(svcCtx),
		CasbinRuleRepository:       logic.NewCasbinRuleRepository(svcCtx),
		MenuRepository:             logic.NewMenuRepository(svcCtx),
		RoleRepository:             logic.NewRoleRepository(svcCtx),
		RoleApiRepository:          logic.NewRoleApiRepository(svcCtx),
		RoleMenuRepository:         logic.NewRoleMenuRepository(svcCtx),
		UserAccountRepository:      logic.NewUserAccountRepository(svcCtx),
		UserInformationRepository:  logic.NewUserInformationRepository(svcCtx),
		UserLoginHistoryRepository: logic.NewUserLoginHistoryRepository(svcCtx),
		UserOauthRepository:        logic.NewUserOauthRepository(svcCtx),
		UserRoleRepository:         logic.NewUserRoleRepository(svcCtx),
		CategoryRepository:         logic.NewCategoryRepository(svcCtx),
		FriendLinkRepository:       logic.NewFriendLinkRepository(svcCtx),
		TagRepository:              logic.NewTagRepository(svcCtx),
		PageRepository:             logic.NewPageRepository(svcCtx),
		CommentRepository:          logic.NewCommentRepository(svcCtx),
		PhotoRepository:            logic.NewPhotoRepository(svcCtx),
		PhotoAlbumRepository:       logic.NewPhotoAlbumRepository(svcCtx),
		TalkRepository:             logic.NewTalkRepository(svcCtx),
		ArticleTagRepository:       logic.NewArticleTagRepository(svcCtx),
		ChatRecordRepository:       logic.NewChatRecordRepository(svcCtx),
		UniqueViewRepository:       logic.NewUniqueViewRepository(svcCtx),
		OperationLogRepository:     logic.NewOperationLogRepository(svcCtx),
		RemarkRepository:           logic.NewRemarkRepository(svcCtx),
		WebsiteConfigRepository:    logic.NewWebsiteConfigRepository(svcCtx),
		UploadRecordRepository:     logic.NewUploadRecordRepository(svcCtx),
		ChatSessionRepository:      logic.NewChatSessionRepository(svcCtx),
		ChatMessageRepository:      logic.NewChatMessageRepository(svcCtx),
	}
}
