package photorpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/photorpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAlbumListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAlbumListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAlbumListLogic {
	return &FindAlbumListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询相册列表
func (l *FindAlbumListLogic) FindAlbumList(in *photorpc.FindAlbumListReq) (*photorpc.FindAlbumListResp, error) {
	page, size, sorts, conditions, params := convertAlbumQuery(in)

	result, err := l.svcCtx.TAlbumModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	cm, err := findPhotoCountGroupAlbum(l.ctx, l.svcCtx, result)
	if err != nil {
		return nil, err
	}

	var list []*photorpc.AlbumDetails
	for _, v := range result {
		list = append(list, convertAlbumOut(v, cm))
	}

	return &photorpc.FindAlbumListResp{
		List: list,
	}, nil
}

func convertAlbumQuery(in *photorpc.FindAlbumListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	if in.AlbumName != "" {
		conditions += " album_name like ?"
		params = append(params, "%"+in.AlbumName+"%")
	}

	return
}

func findPhotoCountGroupAlbum(ctx context.Context, svcCtx *svc.ServiceContext, list []*model.TAlbum) (acm map[int64]int, err error) {
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.Id)
	}

	// 查询每个 category_id 的文章数量
	var results []struct {
		AlbumID    int64 `gorm:"column:album_id"`
		PhotoCount int   `gorm:"column:photo_count"`
	}

	err = svcCtx.Gorm.Model(&model.TPhoto{}).
		Select("album_id, COUNT(*) as photo_count").
		Where("album_id IN ?", ids).
		Group("album_id").
		Order("album_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	acm = make(map[int64]int)
	for _, result := range results {
		acm[result.AlbumID] = result.PhotoCount
	}

	return acm, nil
}
