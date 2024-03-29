package properties

type Upload struct {
	Mode   string       `mapstructure:"mode"`
	Local  UploadConfig `mapstructure:"local"`
	Aliyun UploadConfig `mapstructure:"aliyun"`
	Qiniu  UploadConfig `mapstructure:"qiniu"`
}

type UploadConfig struct {
	Zone            string `mapstructure:"zone" json:"zone" yaml:"zone"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" json:"access-key-secret" yaml:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name" json:"bucket-name" yaml:"bucket-name"`
	BucketUrl       string `mapstructure:"bucket-url" json:"bucket-url" yaml:"bucket-url"` // 项目访问地址
	BasePath        string `mapstructure:"base-path" json:"base-path" yaml:"base-path"`    // 文件访问目录
}
