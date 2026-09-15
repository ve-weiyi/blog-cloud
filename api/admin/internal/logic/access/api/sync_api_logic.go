package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/docs"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type SyncApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 同步接口列表
func NewSyncApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiLogic {
	return &SyncApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncApiLogic) SyncApi(req *types.EmptyReq) (resp *types.SyncApiResp, err error) {
	apis, err := l.getApisFromDocs()
	if err != nil {
		return nil, err
	}

	out, err := l.svcCtx.AccessService.SyncApis(l.ctx, &accessservice.SyncApisRequest{
		Apis: apis,
	})
	if err != nil {
		return nil, err
	}

	return &types.SyncApiResp{
		SuccessCount: out.SuccessCount,
	}, nil
}

type apiEntry struct {
	path   string
	method string
	name   string
	tags   []string
}

func (l *SyncApiLogic) getApisFromDocs() ([]*accessservice.CreateApiRequest, error) {
	doc, err := loads.Analyzed(json.RawMessage(docs.Docs), "")
	if err != nil {
		return nil, err
	}

	sp := doc.Spec()

	// 收集所有 API 操作及其标签
	var entries []apiEntry
	for path, item := range sp.Paths.Paths {
		if item.Get != nil {
			entries = append(entries, apiEntry{path, http.MethodGet, item.Get.Summary, item.Get.Tags})
		}
		if item.Post != nil {
			entries = append(entries, apiEntry{path, http.MethodPost, item.Post.Summary, item.Post.Tags})
		}
		if item.Put != nil {
			entries = append(entries, apiEntry{path, http.MethodPut, item.Put.Summary, item.Put.Tags})
		}
		if item.Delete != nil {
			entries = append(entries, apiEntry{path, http.MethodDelete, item.Delete.Summary, item.Delete.Tags})
		}
	}

	// 按第一个 Tag 分组，保持顺序
	tagGroups := make(map[string][]apiEntry)
	var tagOrder []string
	for _, e := range entries {
		tag := "未分类"
		if len(e.tags) > 0 {
			tag = e.tags[0]
		}
		if _, ok := tagGroups[tag]; !ok {
			tagOrder = append(tagOrder, tag)
		}
		tagGroups[tag] = append(tagGroups[tag], e)
	}

	// 构建树形结构：Tag 作为父节点，API 作为子节点
	var apis []*accessservice.CreateApiRequest
	for _, tag := range tagOrder {
		var children []*accessservice.CreateApiRequest
		for _, e := range tagGroups[tag] {
			traceable := int64(0)
			if e.method == http.MethodPut || e.method == http.MethodDelete {
				traceable = enums.APITraceableYes
			}
			children = append(children, &accessservice.CreateApiRequest{
				Path:      e.path,
				Name:      e.name,
				Method:    e.method,
				Traceable: traceable,
			})
		}
		apis = append(apis, &accessservice.CreateApiRequest{
			ParentId: 0,
			Name:     tag,
			Children: children,
		})
	}

	return apis, nil
}
