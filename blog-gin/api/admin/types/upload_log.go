package types

type UploadLogQuery struct {
	PageQuery
	FilePath string `json:"file_path,optional"` // 文件路径
	FileName string `json:"file_name,optional"` // 文件名称
	FileType string `json:"file_type,optional"` // 文件类型
}
