package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type PhotoAlbumService struct {
	svcCtx *svc.ServiceContext
}

func NewPhotoAlbumService(svcCtx *svc.ServiceContext) *PhotoAlbumService {
	return &PhotoAlbumService{
		svcCtx: svcCtx,
	}
}

// 创建PhotoAlbum记录
func (s *PhotoAlbumService) CreatePhotoAlbum(reqCtx *request.Context, photoAlbum *entity.PhotoAlbum) (data *entity.PhotoAlbum, err error) {
	return s.svcCtx.PhotoAlbumRepository.Create(reqCtx, photoAlbum)
}

// 更新PhotoAlbum记录
func (s *PhotoAlbumService) UpdatePhotoAlbum(reqCtx *request.Context, photoAlbum *entity.PhotoAlbum) (data *entity.PhotoAlbum, err error) {
	return s.svcCtx.PhotoAlbumRepository.Update(reqCtx, photoAlbum)
}

// 删除PhotoAlbum记录
func (s *PhotoAlbumService) DeletePhotoAlbum(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.PhotoAlbumRepository.Delete(reqCtx, "id = ?", id)
}

// 查询PhotoAlbum记录
func (s *PhotoAlbumService) FindPhotoAlbum(reqCtx *request.Context, id int) (data *entity.PhotoAlbum, err error) {
	return s.svcCtx.PhotoAlbumRepository.First(reqCtx, "id = ?", id)
}

// 批量删除PhotoAlbum记录
func (s *PhotoAlbumService) DeletePhotoAlbumByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.PhotoAlbumRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取PhotoAlbum记录
func (s *PhotoAlbumService) FindPhotoAlbumList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.PhotoAlbum, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.PhotoAlbumRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.PhotoAlbumRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
