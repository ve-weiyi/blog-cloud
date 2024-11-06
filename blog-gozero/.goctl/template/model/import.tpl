import (
	"context"
    {{if .time}}"time"{{end}}

    {{if .withCache}}"github.com/redis/go-redis/v9"{{end}}
    "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)
