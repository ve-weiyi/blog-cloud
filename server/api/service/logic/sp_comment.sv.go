package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/sqlx"
)

// 分页获取Comment记录
func (s *CommentService) FindCommentDetailsList(reqCtx *request.Context, page *request.PageQuery) (list []*response.CommentDTO, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	commentList, err := s.svcCtx.CommentRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.svcCtx.CommentRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	var userIds []int
	var commentIds []int
	for _, item := range commentList {
		userIds = append(userIds, item.UserID)
		commentIds = append(commentIds, item.ID)
	}

	// 查询用户
	users, _ := s.svcCtx.UserInformationRepository.FindALL(reqCtx, "id in (?)", userIds)

	var userMap = make(map[int]*entity.UserInformation)
	for _, item := range users {
		userMap[item.ID] = item
	}

	for _, item := range commentList {
		// 查询评论下所有回复列表,只显示五条
		replyList, count, _ := s.FindCommentReplyList(reqCtx, item.ID, &request.PageQuery{
			PageLimit: sqlx.PageLimit{
				Page:     1,
				PageSize: 5,
			},
		})
		// 查询当前评论下所有回复列表
		data := &response.CommentDTO{
			ID:             item.ID,
			UserID:         item.UserID,
			CommentContent: item.CommentContent,
			LikeCount:      100,
			CreatedAt:      item.CreatedAt,
			ReplyCount:     count,
			ReplyDTOList:   replyList,
		}

		// 用户信息
		info, _ := userMap[item.UserID]
		if info != nil {
			data.Nickname = info.Nickname
			data.Avatar = info.Avatar
			data.Website = info.Website
		}

		// 回复的用户信息
		//rinfo, _ := userMap[item.ReplyUserID]
		//if rinfo != nil {
		//	data.ReplyUserID = rinfo.ID
		//	data.ReplyNickname = rinfo.Nickname
		//	data.ReplyWebsite = rinfo.Website
		//}

		list = append(list, data)
	}

	return
}

// 查询Comment记录
func (s *CommentService) FindCommentReplyList(reqCtx *request.Context, commentId int, page *request.PageQuery) (list []*response.ReplyDTO, total int64, err error) {
	page.Conditions = append(page.Conditions, &sqlx.Condition{Field: "parent_id", Rule: "=", Value: commentId})

	cond, args := page.ConditionClause()
	order := page.OrderClause()
	// 查询评论下所有回复列表
	replyList, err := s.svcCtx.CommentRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.svcCtx.CommentRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	// 收集需要查询的用户id
	var userIds []int
	for _, item := range replyList {
		userIds = append(userIds, item.UserID)
		userIds = append(userIds, item.ReplyUserID)
	}

	// 查询用户
	users, _ := s.svcCtx.UserInformationRepository.FindALL(reqCtx, "id in (?)", userIds)

	var userMap = make(map[int]*entity.UserInformation)
	for _, item := range users {
		userMap[item.ID] = item
	}

	// 组装返回数据
	for _, item := range replyList {

		data := &response.ReplyDTO{
			ID:             item.ID,
			ParentID:       item.ParentID,
			UserID:         item.UserID,
			ReplyUserID:    item.ReplyUserID,
			CommentContent: item.CommentContent,
			LikeCount:      5,
			CreatedAt:      item.CreatedAt,
		}

		// 用户信息
		info, _ := userMap[item.UserID]
		if info != nil {
			data.Nickname = info.Nickname
			data.Avatar = info.Avatar
			data.Website = info.Website
		}

		// 回复的用户信息
		rinfo, _ := userMap[item.ReplyUserID]
		if rinfo != nil {
			data.ReplyUserID = rinfo.ID
			data.ReplyNickname = rinfo.Nickname
			data.ReplyWebsite = rinfo.Website
		}

		list = append(list, data)
	}
	return list, total, nil
}

// 查询Comment后台记录
func (s *CommentService) FindCommentBackList(reqCtx *request.Context, page *request.PageQuery) (list []*response.CommentBackDTO, total int64, err error) {
	// 使用用户昵称查询
	cd := sqlx.FindCondition(page.Conditions, "username")
	if cd != nil {
		accounts, err := s.svcCtx.UserAccountRepository.FindALL(reqCtx, "username like ?")
		if err != nil {
			return nil, 0, err
		}

		var userIds []int
		for _, item := range accounts {
			userIds = append(userIds, item.ID)
		}
		// 替换查询条件
		cd.Field = "user_id"
		cd.Value = userIds
		cd.Rule = "in"
	}

	cond, args := page.ConditionClause()
	order := page.OrderClause()
	// 查询评论下所有回复列表
	commentList, err := s.svcCtx.CommentRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	total, err = s.svcCtx.CommentRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}

	// 收集需要查询的用户id
	var userIds []int
	var articleIds []int
	for _, item := range commentList {
		userIds = append(userIds, item.UserID)
		userIds = append(userIds, item.ReplyUserID)
		articleIds = append(articleIds, item.TopicID)
	}

	// 查询用户
	users, _ := s.svcCtx.UserInformationRepository.FindALL(reqCtx, "id in (?)", userIds)
	var userMap = make(map[int]*entity.UserInformation)
	for _, item := range users {
		userMap[item.ID] = item
	}
	// 查询文章
	articles, _ := s.svcCtx.ArticleRepository.FindALL(reqCtx, "id in (?)", articleIds)
	var articleMap = make(map[int]*entity.Article)
	for _, item := range articles {
		articleMap[item.ID] = item
	}

	// 组装返回数据
	for _, item := range commentList {

		data := &response.CommentBackDTO{
			ID:             item.ID,
			Avatar:         "",
			Nickname:       "",
			ReplyNickname:  "",
			ArticleTitle:   "",
			CommentContent: item.CommentContent,
			Type:           item.Type,
			IsReview:       item.IsReview,
			CreatedAt:      item.CreatedAt,
		}

		// 用户信息
		info, _ := userMap[item.UserID]
		if info != nil {
			data.Avatar = info.Avatar
			data.Nickname = info.Nickname
		}

		// 回复的用户信息
		rinfo, _ := userMap[item.ReplyUserID]
		if rinfo != nil {
			data.ReplyNickname = rinfo.Nickname
		}

		// 回复的文章信息
		aInfo, _ := articleMap[item.TopicID]
		if aInfo != nil {
			data.ArticleTitle = aInfo.ArticleTitle
		}

		list = append(list, data)
	}
	return list, total, nil
}

// 点赞Comment
func (s *CommentService) LikeComment(reqCtx *request.Context, commentId int) (data interface{}, err error) {

	return s.svcCtx.CommentRepository.LikeComment(reqCtx, reqCtx.UID, commentId)
}
