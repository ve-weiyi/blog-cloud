package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/docs"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
)

type SyncApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 同步api列表
func NewSyncApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiListLogic {
	return &SyncApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncApiListLogic) SyncApiList(req *types.SyncApiReq) (resp *types.BatchResp, err error) {
	doc, err := loads.Analyzed(json.RawMessage(docs.Docs), "")
	if err != nil {
		panic(err)
	}

	sp := doc.Spec()
	routes := getRoutes(sp)

	// 分组
	groups := make(map[string][]*permissionrpc.AddApiReq)
	for k, v := range routes {
		for m, o := range v {
			if o != nil {
				// 自动判断是否需要记录日志
				var traceable int64 = 0
				switch m {
				case http.MethodPut:
					traceable = 1
				case http.MethodDelete:
					traceable = 1
				case http.MethodPost:
					if !strings.Contains(k, "list") && !strings.Contains(k, "get") {
						traceable = 1
					}
				default:
					break
				}
				child := &permissionrpc.AddApiReq{
					Id:        0,
					ParentId:  0,
					Path:      k,
					Name:      o.Summary,
					Method:    m,
					Traceable: traceable,
					Status:    0,
					Children:  nil,
				}

				var group = ""
				if len(o.Tags) > 0 {
					group = o.Tags[0]
				}
				groups[group] = append(groups[group], child)
			}
		}
	}

	var list []*permissionrpc.AddApiReq
	for g, children := range groups {
		root := &permissionrpc.AddApiReq{
			Id:        0,
			ParentId:  0,
			Path:      g,
			Name:      g,
			Method:    "",
			Traceable: 0,
			Status:    0,
			Children:  children,
		}
		list = append(list, root)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Path < list[j].Path
	})

	in := &permissionrpc.SyncApiListReq{
		Apis: list,
	}

	out, err := l.svcCtx.PermissionRpc.SyncApiList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}

func getRoutes(sp *spec.Swagger) map[string]map[string]*spec.Operation {
	// map[path][method] -> operation
	routes := make(map[string]map[string]*spec.Operation)

	for k, v := range sp.Paths.Paths {
		if routes[k] == nil {
			routes[k] = make(map[string]*spec.Operation)
		}

		if v.Get != nil {
			routes[k][http.MethodGet] = v.Get
		}

		if v.Put != nil {
			routes[k][http.MethodPut] = v.Put
		}

		if v.Post != nil {
			routes[k][http.MethodPost] = v.Post
		}

		if v.Delete != nil {
			routes[k][http.MethodDelete] = v.Delete
		}

		if v.Options != nil {
			routes[k][http.MethodOptions] = v.Options
		}

		if v.Head != nil {
			routes[k][http.MethodHead] = v.Head
		}
	}

	return routes
}
