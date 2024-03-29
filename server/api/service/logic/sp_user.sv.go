package logic

import (
	"mime/multipart"
	"path"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr/httperr"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/cache"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/mail"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/sqlx"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/crypto"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/temputil"
)

type UserService struct {
	svcCtx *svc.ServiceContext
}

func NewUserService(svcCtx *svc.ServiceContext) *UserService {
	return &UserService{
		svcCtx: svcCtx,
	}
}

// 分页获取UserAccount记录
func (s *UserService) FindUserList(reqCtx *request.Context, page *request.PageQuery) (list []*response.UserDTO, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()
	// 查询账号信息
	userAccounts, err := s.svcCtx.UserAccountRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.svcCtx.UserAccountRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	var ids []int
	for _, ua := range userAccounts {
		ids = append(ids, ua.ID)
	}

	//获取用户信息
	infos, err := s.svcCtx.UserInformationRepository.FindALL(reqCtx, "id in (?)", ids)
	if err != nil {
		return nil, 0, err
	}

	var infoMap = make(map[int]*entity.UserInformation)
	for _, info := range infos {
		infoMap[info.ID] = info
	}

	for _, account := range userAccounts {
		info := infoMap[account.ID]
		// 查询账号角色信息
		roles, _ := s.svcCtx.RoleRepository.FindUserRoles(reqCtx, account.ID)

		item := &response.UserDTO{
			ID:           account.ID,
			Username:     account.Username,
			Nickname:     info.Nickname,
			Status:       account.Status,
			Avatar:       info.Avatar,
			Intro:        info.Intro,
			Website:      info.Website,
			Email:        info.Email,
			RegisterType: account.RegisterType,
			IpAddress:    account.IpAddress,
			IpSource:     account.IpSource,
			CreatedAt:    account.CreatedAt,
			UpdatedAt:    account.UpdatedAt,
			Roles:        convertRoleList(roles),
		}

		list = append(list, item)
	}

	return list, total, nil
}

// 获取在线用户列表
func (s *UserService) FindOnlineUserList(reqCtx *request.Context, page *request.PageQuery) (list []*response.UserDTO, total int64, err error) {
	keys, err := s.svcCtx.UserAccountRepository.Online(reqCtx, page.Page, page.PageSize)
	if err != nil {
		return nil, 0, err
	}

	s.svcCtx.Log.JsonIndent("names", keys)
	page.Page = 0
	page.PageSize = 0
	page.Conditions = append(page.Conditions, sqlx.NewCondition("id in (?)", keys))
	return s.FindUserList(reqCtx, page)
}

func (s *UserService) FindUserAreaList(reqCtx *request.Context, page *request.PageQuery) (result []*response.UserAreaDTO, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()
	list, err := s.svcCtx.UserAccountRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.svcCtx.UserAccountRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	// 分类
	AreaMap := make(map[string]int)
	for _, item := range list {
		key := item.IpSource
		if _, ok := AreaMap[key]; ok {
			AreaMap[key]++
		} else {
			AreaMap[key] = 1
		}
	}

	for k, v := range AreaMap {
		result = append(result, &response.UserAreaDTO{
			Name:  k,
			Value: v,
		})
	}
	return result, int64(len(result)), nil
}

func (s *UserService) FindUserLoginHistoryList(reqCtx *request.Context, page *request.PageQuery) (result []*response.LoginHistory, total int64, err error) {
	//获取用户
	account, err := s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", reqCtx.UID)
	if err != nil {
		return nil, 0, apierr.NewApiError(httperr.CodeForbidden, "用户不存在！")
	}

	// 添加用户id条件
	c := &sqlx.Condition{Field: "user_id", Value: account.ID, Rule: "=", Flag: "AND"}
	page.Conditions = append(page.Conditions, c)
	cond, args := page.ConditionClause()
	order := page.OrderClause()
	histories, err := s.svcCtx.UserLoginHistoryRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.svcCtx.UserLoginHistoryRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range histories {
		his := convertLoginHistory(item)
		result = append(result, his)
	}

	return result, total, nil
}

func (s *UserService) DeleteUserLoginHistoryByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	//获取用户
	account, err := s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", reqCtx.UID)
	if err != nil {
		return 0, apierr.NewApiError(httperr.CodeForbidden, "用户不存在！")
	}

	// 添加用户id条件
	return s.svcCtx.UserLoginHistoryRepository.Delete(reqCtx, "id in (?) and user_id = ?", ids, account.ID)
}

func (s *UserService) SendForgetPwdEmail(reqCtx *request.Context, req *request.UserEmail) (resp interface{}, err error) {
	// 验证用户是否存在
	account, err := s.svcCtx.UserAccountRepository.LoadUserByUsername(reqCtx, req.Username)
	if account == nil {
		return nil, apierr.ErrorUserNotExist
	}

	// 获取code
	key := cache.WrapCacheKey(constant.ForgetPassword, req.Username)
	code := s.svcCtx.Captcha.GetCodeCaptcha(key)
	data := mail.CaptchaEmail{
		Username: req.Username,
		Code:     code,
	}

	// 组装邮件内容
	content, err := temputil.TempParseString(mail.TempForgetPassword, data)
	if err != nil {
		return nil, err
	}

	msg := &mail.EmailMessage{
		To:      []string{req.Username},
		Subject: "忘记密码",
		Content: content,
		Type:    0,
	}
	// 发送邮件
	err = s.svcCtx.EmailPublisher.PublishMessage([]byte(jsonconv.ObjectToJson(msg)))
	if err != nil {
		return nil, err
	}
	return true, nil
}

func (s *UserService) ResetPassword(reqCtx *request.Context, req *request.ResetPasswordReq) (resp interface{}, err error) {
	// 验证code是否正确
	key := cache.WrapCacheKey(constant.ForgetPassword, req.Username)
	if !s.svcCtx.Captcha.VerifyCaptcha(key, req.Code) {
		return nil, apierr.ErrorCaptchaVerify
	}

	// 验证用户是否存在
	account, err := s.svcCtx.UserAccountRepository.LoadUserByUsername(reqCtx, req.Username)
	if account == nil {
		return nil, apierr.ErrorUserNotExist
	}

	// 更新密码
	account.Password = crypto.BcryptHash(req.Password)
	_, err = s.svcCtx.UserAccountRepository.Update(reqCtx, account)
	if err != nil {
		return nil, err
	}

	return true, nil
}

// 修改用户角色
func (s *UserService) UpdateUserAvatar(reqCtx *request.Context, file *multipart.FileHeader) (data *entity.UserInformation, err error) {
	label := "avatar"
	url, err := s.svcCtx.Uploader.UploadFile(path.Join(cast.ToString(reqCtx.UID), label), file)
	if err != nil {
		return nil, err
	}

	// 保存上传记录
	up := &entity.UploadRecord{
		UserID:   reqCtx.UID,
		Label:    label,
		FileName: file.Filename,
		FileSize: int(file.Size),
		FileMd5:  crypto.MD5V([]byte(file.Filename)),
		FileURL:  url,
	}

	_, err = s.svcCtx.UploadRecordRepository.Create(reqCtx, up)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	information, err := s.svcCtx.UserInformationRepository.First(reqCtx, "id = ?", reqCtx.UID)
	if err != nil {
		return nil, err
	}
	information.Avatar = url

	return s.svcCtx.UserInformationRepository.Update(reqCtx, information)
}

// 修改用户角色
func (s *UserService) UpdateUserRoles(reqCtx *request.Context, req *request.UpdateUserRolesReq) (data interface{}, err error) {

	return s.svcCtx.RoleRepository.UpdateUserRoles(reqCtx, req.UserId, req.RoleIds)
}

// 修改用户状态
func (s *UserService) UpdateUserStatus(reqCtx *request.Context, req *entity.UserAccount) (data *entity.UserAccount, err error) {
	// 创建db
	account, err := s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", req.ID)
	if err != nil {
		return nil, err
	}

	account.Status = req.Status
	_, err = s.svcCtx.UserAccountRepository.Update(reqCtx, account)
	if err != nil {
		return nil, err
	}

	return account, err
}

// 修改用户信息
func (s *UserService) UpdateUserInfo(reqCtx *request.Context, req *request.UserInfoReq) (data *entity.UserInformation, err error) {
	info, err := s.svcCtx.UserAccountRepository.FindUserInfo(reqCtx, reqCtx.UID)
	if err != nil {
		return nil, err
	}

	info.Nickname = req.Nickname
	info.Intro = req.Intro
	info.Website = req.Website
	info.Avatar = req.Avatar
	_, err = s.svcCtx.UserInformationRepository.Update(reqCtx, info)
	if err != nil {
		return nil, err
	}

	return info, err
}

func (s *UserService) GetUserInfo(reqCtx *request.Context, userId int) (data *response.UserInfo, err error) {
	account, err := s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", userId)
	if err != nil {
		return nil, apierr.NewApiError(httperr.CodeForbidden, "用户不存在！")
	}

	//获取用户信息
	info, err := s.svcCtx.UserAccountRepository.FindUserInfo(reqCtx, account.ID)
	if err != nil {
		return nil, err
	}

	accountLikeSet, _ := s.svcCtx.ArticleRepository.FindUserLikeArticle(reqCtx, account.ID)
	commentLikeSet, _ := s.svcCtx.CommentRepository.FindUserLikeComment(reqCtx, account.ID)
	talkLikeSet, _ := s.svcCtx.TalkRepository.FindUserLikeTalk(reqCtx, account.ID)

	roles, err := s.svcCtx.RoleRepository.FindUserRoles(reqCtx, account.ID)
	data = &response.UserInfo{
		ID:             account.ID,
		Username:       account.Username,
		Nickname:       info.Nickname,
		Avatar:         info.Avatar,
		Intro:          info.Intro,
		Website:        info.Website,
		Email:          info.Email,
		ArticleLikeSet: accountLikeSet,
		CommentLikeSet: commentLikeSet,
		TalkLikeSet:    talkLikeSet,
		Roles:          convertRoleList(roles),
	}

	return data, nil
}

func (s *UserService) GetUserMenus(reqCtx *request.Context, req interface{}) (data []*response.MenuDetailsDTO, err error) {
	//查询用户信息
	account, err := s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", reqCtx.UID)
	if err != nil {
		return nil, err
	}

	//查询用户角色
	roles, err := s.svcCtx.RoleRepository.FindUserRoles(reqCtx, account.ID)
	if err != nil {
		return nil, err
	}

	//查询角色权限,取交集
	menuMaps := make(map[int]*entity.Menu)
	for _, item := range roles {
		menus, err := s.svcCtx.RoleRepository.FindRoleMenus(reqCtx, item.ID)
		if err != nil {
			return nil, err
		}
		// 去重
		for _, m := range menus {
			if _, ok := menuMaps[m.ID]; !ok {
				menuMaps[m.ID] = m
			}
		}
	}

	var list []*entity.Menu
	for _, v := range menuMaps {
		list = append(list, v)
	}

	var out response.MenuDetailsDTO
	out.Children = getMenuChildren(out, list)

	return out.Children, err
}

func (s *UserService) GetUserApis(reqCtx *request.Context, req interface{}) (data []*response.ApiDetailsDTO, err error) {
	//查询用户信息
	account, err := s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", reqCtx.UID)
	if err != nil {
		return nil, err
	}

	//查询用户角色
	roles, err := s.svcCtx.RoleRepository.FindUserRoles(reqCtx, account.ID)
	if err != nil {
		return nil, err
	}

	//查询角色权限,取交集
	menuMaps := make(map[int]*entity.Api)
	for _, item := range roles {
		menus, err := s.svcCtx.RoleRepository.FindRoleApis(reqCtx, item.ID)
		if err != nil {
			return nil, err
		}
		// 去重
		for _, m := range menus {
			if _, ok := menuMaps[m.ID]; !ok {
				menuMaps[m.ID] = m
			}
		}
	}

	var list []*entity.Api
	for _, v := range menuMaps {
		list = append(list, v)
	}

	var out response.ApiDetailsDTO
	out.Children = getApiChildren(out, list)

	return out.Children, err
}
