package tmpl

const Api = `import http from "@/utils/request"

interface {{.UpperStartCamelName}} {
	{{- range .Fields}}
    {{.FieldValueName}}: any
	{{- end}}
}

/** 增 */
export function create{{.UpperStartCamelName}}Api(data?: object): Promise<IApiResponseData<{{.UpperStartCamelName}}>> {
  return http.request<IApiResponseData<{{.UpperStartCamelName}}>>({
    url: "/api/v1/{{.SnakeName}}",
    method: "post",
    data,
  })
}

/** 改 */
export function update{{.UpperStartCamelName}}Api(data?: object): Promise<IApiResponseData<{{.UpperStartCamelName}}>> {
  return http.request<IApiResponseData<{{.UpperStartCamelName}}>>({
    url: "/api/v1/{{.SnakeName}}",
    method: "put",
    data,
  })
}

/** 删 */
export function delete{{.UpperStartCamelName}}Api(id: number): Promise<IApiResponseData<{{.UpperStartCamelName}}>> {
  return http.request<IApiResponseData<{{.UpperStartCamelName}}>>({
    url: ` + "`/api/v1/{{.SnakeName}}/${id}`" + `,
    method: "delete",
  })
}

/** 查 */
export function find{{.UpperStartCamelName}}Api(id: number): Promise<IApiResponseData<{{.UpperStartCamelName}}>> {
  return http.request<IApiResponseData<{{.UpperStartCamelName}}>>({
    url: ` + "`/api/v1/{{.SnakeName}}/${id}`" + `,
    method: "get",
  })
}

/** 删除 批量操作 */
export function delete{{.UpperStartCamelName}}ByIdsApi(ids: number[]): Promise<IApiResponseData<{{.UpperStartCamelName}}>> {
  return http.request<IApiResponseData<{{.UpperStartCamelName}}>>({
    url: "/api/v1/{{.SnakeName}}/batch_delete",
    method: "delete",
    data: ids,
  })
}

/** 查询 分页列表 */
export function find{{.UpperStartCamelName}}ListApi(page?: Page): Promise<IApiResponseData<PageResult<{{.UpperStartCamelName}}>>> {
  return http.request<IApiResponseData<PageResult<{{.UpperStartCamelName}}>>>({
    url: "/api/v1/{{.SnakeName}}/list",
    method: "post",
    data: page,
  })
}
`
