package responsex

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

// TestRulesCoverAllCodes 守住"新增标识必须同时确定 HTTP 映射"这条约束：
// 映射表与 bizcode.All 任一方向漏了都会在这里失败。
func TestRulesCoverAllCodes(t *testing.T) {
	declared := make(map[string]bool, len(bizcode.All))
	for _, code := range bizcode.All {
		if declared[code] {
			t.Errorf("bizcode.All 存在重复标识: %s", code)
		}
		declared[code] = true
	}

	for _, code := range bizcode.All {
		if _, ok := rules[code]; !ok {
			t.Errorf("标识 %s 缺少 HTTP 状态映射，请补入 responsex.rules", code)
		}
	}

	for code := range rules {
		if !declared[code] {
			t.Errorf("映射表存在未声明的标识 %s，请补入 bizcode.All 或删除该映射", code)
		}
	}
}

// TestChallengeByBizCode 守住 challenge 契约：
// challenge 由标识决定，不由状态码决定——同为 401/403，四种语义各有各的头。
func TestChallengeByBizCode(t *testing.T) {
	cases := []struct {
		name       string
		code       string
		wantStatus int
		wantHeader string // 空表示不应出现该头
	}{
		{"未登录按令牌失败给 invalid_token", bizcode.CodeUnauthenticated, http.StatusUnauthorized, `Bearer realm="blog", error="invalid_token"`},
		{"登录过期按令牌失败给 invalid_token", bizcode.CodeLoginExpired, http.StatusUnauthorized, `Bearer realm="blog", error="invalid_token"`},
		{"登录凭据失败不带错误参数", bizcode.CodeCredentialsInvalid, http.StatusUnauthorized, `Bearer realm="blog"`},
		{"请求签名失败用方案自身名", bizcode.CodeRequestSignInvalid, http.StatusUnauthorized, `Signature realm="blog"`},
		{"无权限给 insufficient_scope", bizcode.CodeNoPermission, http.StatusForbidden, `Bearer realm="blog", error="insufficient_scope"`},
		{"角色不匹配给 insufficient_scope", bizcode.CodeRoleNotMatch, http.StatusForbidden, `Bearer realm="blog", error="insufficient_scope"`},
		{"账号禁用不带挑战", bizcode.CodeAccountDisabled, http.StatusForbidden, ""},
		{"操作不允许不带挑战", bizcode.CodeOperationNotAllowed, http.StatusForbidden, ""},
		{"参数校验不带挑战", bizcode.CodeParamFormat, http.StatusUnprocessableEntity, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
			w := httptest.NewRecorder()

			Response(r, w, nil, bizerr.NewBizError(c.code, "boom"))

			if w.Code != c.wantStatus {
				t.Errorf("HTTP 状态 = %d, want %d", w.Code, c.wantStatus)
			}
			if got := w.Header().Get(headerWWWAuthenticate); got != c.wantHeader {
				t.Errorf("WWW-Authenticate = %q, want %q", got, c.wantHeader)
			}
		})
	}
}

// TestNoStoreOnNon2xxOnly：错误响应可能含请求相关细节，不应被缓存。
func TestNoStoreOnNon2xxOnly(t *testing.T) {
	t.Run("失败响应带 no-store", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
		w := httptest.NewRecorder()

		Response(r, w, nil, bizerr.NewBizError(bizcode.CodeParamFormat, "boom"))

		if got := w.Header().Get(headerCacheControl); got != cacheControlNoStore {
			t.Errorf("Cache-Control = %q, want %q", got, cacheControlNoStore)
		}
	})

	t.Run("成功响应不带 no-store", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
		w := httptest.NewRecorder()

		Response(r, w, map[string]any{"id": 1}, nil)

		if got := w.Header().Get(headerCacheControl); got != "" {
			t.Errorf("成功响应不应带 Cache-Control，得到 %q", got)
		}
	})
}

// TestRetryAfterNotFabricated：
// 给不出精确值时不臆造——宁可省略，也不写一个恒定的假值。
func TestRetryAfterNotFabricated(t *testing.T) {
	t.Run("未给出精确值时省略", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
		w := httptest.NewRecorder()

		Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "请求过于频繁"))

		if got := w.Header().Get(HeaderRetryAfter); got != "" {
			t.Errorf("限流未给出精确退避时长时不应臆造 Retry-After，得到 %q", got)
		}
	})

	t.Run("中间件已写精确值时保留", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
		w := httptest.NewRecorder()
		w.Header().Set(HeaderRetryAfter, "12")

		Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "请求过于频繁"))

		if got := w.Header().Get(HeaderRetryAfter); got != "12" {
			t.Errorf("Retry-After = %q, want %q（出口不得覆盖中间件给出的精确值）", got, "12")
		}
	})
}

// TestParamErrorsStatus 守住 400/422 的分界：
// 422 只用于请求**内容**语义不合法；请求头问题（缺失）归 400。
func TestParamErrorsStatus(t *testing.T) {
	cases := []struct {
		code       string
		wantStatus int
	}{
		{bizcode.CodeInvalidParam, http.StatusUnprocessableEntity},
		{bizcode.CodeParamFormat, http.StatusUnprocessableEntity},
		{bizcode.CodeParamValueInvalid, http.StatusUnprocessableEntity},
		{bizcode.CodeVerifyCodeError, http.StatusUnprocessableEntity},
		// 请求头缺失不是请求内容问题
		{bizcode.CodeParamMissing, http.StatusBadRequest},
	}

	for _, c := range cases {
		t.Run(c.code, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
			w := httptest.NewRecorder()

			Response(r, w, nil, bizerr.NewBizError(c.code, "boom"))

			if w.Code != c.wantStatus {
				t.Errorf("HTTP 状态 = %d, want %d", w.Code, c.wantStatus)
			}
		})
	}
}

// TestUnknownBizCodeFallsBackTo500 守住兜底：表外标识按 500 处理，但保留标识本身便于定位。
func TestUnknownBizCodeFallsBackTo500(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	Response(r, w, nil, bizerr.NewBizError("NOT_REGISTERED", "boom"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("表外标识的 HTTP 状态 = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	if got := decodeBody(t, w)["code"]; got != "NOT_REGISTERED" {
		t.Errorf("响应 code = %v, want %v（表外标识也应原样带出）", got, "NOT_REGISTERED")
	}
}

// TestViolationsGoIntoData：字段级明细放 data，不塞进 message。
func TestViolationsGoIntoData(t *testing.T) {
	bizErr := bizerr.NewBizError(bizcode.CodeParamFormat, "参数校验不通过")
	bizErr.Violations = []bizerr.FieldViolation{
		{Field: "title", Description: "不能为空"},
		{Field: "tags[2].name", Description: "长度超过 32"},
	}

	r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	Response(r, w, nil, bizErr)

	data, ok := decodeBody(t, w)["data"].(map[string]any)
	if !ok {
		t.Fatalf("data 不是对象: %v", decodeBody(t, w)["data"])
	}
	items, ok := data["field_violations"].([]any)
	if !ok {
		t.Fatalf("data.field_violations 不是数组: %v", data)
	}
	if len(items) != 2 {
		t.Fatalf("field_violations 条数 = %d, want 2", len(items))
	}
	first, _ := items[0].(map[string]any)
	if first["field"] != "title" {
		t.Errorf("field_violations[0].field = %v, want %v", first["field"], "title")
	}
}

// TestSuccessMessageIsEmpty：message 面向开发者，成功时为空串。
func TestSuccessMessageIsEmpty(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	Response(r, w, map[string]any{"id": 1}, nil)

	body := decodeBody(t, w)
	if body["message"] != "" {
		t.Errorf("成功响应 message = %v, want 空串", body["message"])
	}
	if body["code"] != bizcode.CodeSuccess {
		t.Errorf("成功响应 code = %v, want %v", body["code"], bizcode.CodeSuccess)
	}
}

// TestResponseStatusBadRequestHasNoCode：
// 请求解析失败写 400、**不带 code**、data.field_violations 必须存在。
func TestResponseStatusBadRequestHasNoCode(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	ResponseStatus(r, w, http.StatusBadRequest, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP 状态 = %d, want %d", w.Code, http.StatusBadRequest)
	}
	body := decodeBody(t, w)
	if _, ok := body["code"]; ok {
		t.Errorf("解析失败响应不应带 code，得到 %v", body["code"])
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data 不是对象: %v", body["data"])
	}
	items, ok := data["field_violations"].([]any)
	if !ok {
		t.Fatalf("data.field_violations 必须是数组（定位不到就给空数组）: %v", data)
	}
	if len(items) != 0 {
		t.Errorf("field_violations 条数 = %d, want 0", len(items))
	}
	if body["trace_id"] == nil {
		t.Error("trace_id 字段必须存在")
	}
}

// TestRequestInvalidKeepsCause：
// 请求读取/解析失败写 400、不带 code，且**原始错误必须保留**（进日志与 message，供开发者排错）。
func TestRequestInvalidKeepsCause(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	RequestInvalid(r, w, errors.New("field title is not set"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP 状态 = %d, want %d", w.Code, http.StatusBadRequest)
	}
	body := decodeBody(t, w)
	if _, ok := body["code"]; ok {
		t.Errorf("解析失败响应不应带 code，得到 %v", body["code"])
	}
	if body["message"] != "field title is not set" {
		t.Errorf("message = %v, want 原始错误文本（排错信息不得丢失）", body["message"])
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data 不是对象: %v", body["data"])
	}
	items, ok := data["field_violations"].([]any)
	if !ok {
		t.Fatalf("data.field_violations 必须是数组: %v", data)
	}
	if len(items) != 0 {
		t.Errorf("field_violations 条数 = %d, want 0", len(items))
	}
}

// TestResponseStatusCreatedKeepsEnvelope：
// 除本入口与 RequestInvalid 两条路径外，所有响应都是信封，201 也不例外。
func TestResponseStatusCreatedKeepsEnvelope(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	ResponseStatus(r, w, http.StatusCreated, map[string]any{"id": 7})

	if w.Code != http.StatusCreated {
		t.Errorf("HTTP 状态 = %d, want %d", w.Code, http.StatusCreated)
	}
	body := decodeBody(t, w)
	if body["code"] != bizcode.CodeSuccess {
		t.Errorf("201 的 code = %v, want %v", body["code"], bizcode.CodeSuccess)
	}
	data, ok := body["data"].(map[string]any)
	if !ok || data["id"] != float64(7) {
		t.Errorf("201 的 data = %v, want 含 id=7", body["data"])
	}
}

// TestResponseStatusRejectsNonAllowlisted 守住白名单：显式状态出口不得被当作业务状态的后门。
func TestResponseStatusRejectsNonAllowlisted(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
	w := httptest.NewRecorder()

	ResponseStatus(r, w, http.StatusNotFound, nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("白名单外状态应兜底为 500，得到 %d", w.Code)
	}
	if got := decodeBody(t, w)["code"]; got != bizcode.CodeInternalServerError {
		t.Errorf("兜底响应的 code = %v, want %v", got, bizcode.CodeInternalServerError)
	}
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体不是 JSON: %v, body=%q", err, w.Body.String())
	}

	return body
}
