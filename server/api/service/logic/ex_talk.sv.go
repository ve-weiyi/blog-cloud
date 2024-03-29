package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

type TalkService struct {
	svcCtx *svc.ServiceContext
}

func NewTalkService(svcCtx *svc.ServiceContext) *TalkService {
	return &TalkService{
		svcCtx: svcCtx,
	}
}

// 创建Talk记录
func (s *TalkService) CreateTalk(reqCtx *request.Context, talk *entity.Talk) (data *entity.Talk, err error) {
	talk.UserID = reqCtx.UID
	return s.svcCtx.TalkRepository.Create(reqCtx, talk)
}

// 更新Talk记录
func (s *TalkService) UpdateTalk(reqCtx *request.Context, talk *entity.Talk) (data *entity.Talk, err error) {
	if talk.UserID != reqCtx.UID {
		return nil, apierr.ErrorUserNotPermission
	}
	return s.svcCtx.TalkRepository.Update(reqCtx, talk)
}

// 删除Talk记录
func (s *TalkService) DeleteTalk(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.TalkRepository.Delete(reqCtx, "id = ?", id)
}

// 查询Talk记录
func (s *TalkService) FindTalk(reqCtx *request.Context, id int) (data *entity.Talk, err error) {
	return s.svcCtx.TalkRepository.First(reqCtx, "id = ?", id)
}

// 批量删除Talk记录
func (s *TalkService) DeleteTalkByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.TalkRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取Talk记录
func (s *TalkService) FindTalkList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.Talk, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.TalkRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.TalkRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// 获取说说详情列表
func (s *TalkService) FindTalkDetailsList(reqCtx *request.Context, page *request.PageQuery) (list []*response.TalkDetailsDTO, total int64, err error) {
	talkList, total, err := s.FindTalkList(reqCtx, page)
	if err != nil {
		return nil, 0, err
	}

	// 查询分类下的文章数量

	for _, talk := range talkList {

		user, err := s.svcCtx.UserInformationRepository.First(reqCtx, "id = ?", talk.UserID)
		if err != nil {
			continue
		}

		var imgList []string
		jsonconv.JsonToObject(talk.Images, &imgList)
		data := &response.TalkDetailsDTO{
			ID:           talk.ID,
			UserID:       talk.UserID,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Content:      talk.Content,
			ImgList:      imgList,
			IsTop:        talk.IsTop,
			Status:       talk.Status,
			LikeCount:    10,
			CommentCount: 10,
			CreatedAt:    talk.CreatedAt,
			UpdatedAt:    talk.UpdatedAt,
		}

		list = append(list, data)
	}

	return list, total, err
}

// 获取说说详情
func (s *TalkService) FindTalkDetailsDTO(reqCtx *request.Context, id int) (data *response.TalkDetailsDTO, err error) {
	// 查询api信息
	talk, err := s.svcCtx.TalkRepository.First(reqCtx, "id = ?", id)
	if err != nil {
		return nil, err
	}

	user, err := s.svcCtx.UserInformationRepository.First(reqCtx, "id = ?", talk.UserID)
	if err != nil {
		return nil, err
	}

	var imgList []string
	jsonconv.JsonToObject(talk.Images, &imgList)
	data = &response.TalkDetailsDTO{
		ID:        talk.ID,
		UserID:    talk.UserID,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Content:   talk.Content,
		ImgList:   imgList,
		IsTop:     talk.IsTop,
		Status:    talk.Status,
		LikeCount: 10,
		CreatedAt: talk.CreatedAt,
		UpdatedAt: talk.UpdatedAt,
	}
	return data, nil
}

// 点赞说说
func (s *TalkService) LikeTalk(reqCtx *request.Context, id int) (data interface{}, err error) {
	return s.svcCtx.TalkRepository.LikeTalk(reqCtx, reqCtx.UID, id)
}
