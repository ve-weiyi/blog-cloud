package task

import (
	"github.com/robfig/cron/v3"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/timer"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
)

func Run(svcCtx *svc.ServiceContext) {

	tm := timer.NewTimerTask()

	taskOnline := NewTaskClearChatOnline(svcCtx.Redis)
	_, err := tm.AddTaskByFunc("clear_chat_online",
		"0 */5 * * * *",
		taskOnline.Run,
		cron.WithSeconds(),
	)
	if err != nil {
		panic(err)
	}

}
