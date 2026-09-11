package tracelogx

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"
)

var _ Enforcer = &TraceEnforcer{}

type TraceEnforcer struct {
	mu           sync.RWMutex
	pr           permissionservice.PermissionService
	rds          *redis.Client
	rules        []traceRule
	policyLoaded bool
}

func NewTraceEnforcer(rds *redis.Client, pr permissionservice.PermissionService) *TraceEnforcer {
	h := &TraceEnforcer{pr: pr, rds: rds}
	h.startSubscribe()
	return h
}

func (s *TraceEnforcer) ReloadPolicy() error {
	return s.LoadPolicy()
}

func (s *TraceEnforcer) LoadPolicy() error {
	return s.loadRulesFromRPC()
}

func (s *TraceEnforcer) IsTraceLog(url string, method string) (bool, error) {
	if strings.TrimSpace(url) == "" {
		return false, errors.New("url is empty")
	}
	if strings.TrimSpace(method) == "" {
		return false, errors.New("method is empty")
	}

	s.mu.RLock()
	loaded := s.policyLoaded
	s.mu.RUnlock()
	if !loaded {
		if err := s.LoadPolicy(); err != nil {
			return false, err
		}
	}

	path := normalizePath(url)
	if path == "" {
		return false, errors.New("invalid url")
	}
	action := normalizeMethod(method)
	if action == "" {
		return false, errors.New("invalid method")
	}

	s.mu.RLock()
	rules := s.rules
	s.mu.RUnlock()

	segments := splitPathSegments(path)
	for _, rule := range rules {
		if rule.match(segments, action) {
			return true, nil
		}
	}
	return false, nil
}

func (s *TraceEnforcer) loadRulesFromRPC() error {
	if s.pr == nil {
		return errors.New("permission service is nil")
	}

	logx.Info("Loading trace rules from rpc...")

	resp, err := s.pr.ListApis(context.Background(), &permissionservice.ListApisRequest{})
	if err != nil {
		return err
	}

	flatApis := flattenApiTree(resp.List)
	seen := make(map[string]struct{})
	var rules []traceRule
	for _, api := range flatApis {
		if api == nil || api.Status == enums.APIStatusDisabled || api.Traceable != enums.APITraceableYes {
			continue
		}
		rule, ok := newTraceRule(api.Path, api.Method)
		if !ok {
			continue
		}
		key := rule.method + ":" + rule.path
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		rules = append(rules, rule)
	}

	s.mu.Lock()
	s.rules = rules
	s.policyLoaded = true
	s.mu.Unlock()

	return nil
}

func (s *TraceEnforcer) startSubscribe() {
	if s.rds == nil {
		return
	}
	go func() {
		for {
			if err := s.subscribe(); err != nil {
				logx.Errorf("pubsub subscription lost: %v, reconnecting...", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func (s *TraceEnforcer) subscribe() error {
	ctx := context.Background()
	sub := s.rds.Subscribe(ctx, cachekey.TraceInvalidateChannel)
	defer sub.Close()

	const debounce = 500 * time.Millisecond
	var timer *time.Timer

	for msg := range sub.Channel() {
		if msg.Channel != cachekey.TraceInvalidateChannel {
			continue
		}
		if timer == nil {
			timer = time.AfterFunc(debounce, func() {
				logx.Info("received trace policy invalidate, reloading...")
				if err := s.LoadPolicy(); err != nil {
					logx.Errorf("reload trace rules failed: %v", err)
				}
			})
		} else {
			timer.Reset(debounce)
		}
	}
	return errors.New("channel closed")
}

type traceRule struct {
	path     string
	method   string
	segments []string
}

func (r traceRule) match(segments []string, action string) bool {
	if r.method != "*" && r.method != normalizeMethod(action) {
		return false
	}
	return matchPathSegments(r.segments, segments)
}

func flattenApiTree(nodes []*permissionservice.Api) []*permissionservice.Api {
	var result []*permissionservice.Api
	for _, node := range nodes {
		result = append(result, node)
		if len(node.Children) > 0 {
			result = append(result, flattenApiTree(node.Children)...)
		}
	}
	return result
}

func newTraceRule(path, method string) (traceRule, bool) {
	path = normalizePath(path)
	if path == "" {
		return traceRule{}, false
	}
	method = normalizeMethod(method)
	if method == "" {
		return traceRule{}, false
	}
	return traceRule{path: path, method: method, segments: splitPathSegments(path)}, true
}
