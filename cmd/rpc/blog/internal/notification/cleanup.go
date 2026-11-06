package notification

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// TerminateObject 终止对象关联事件并删除条目与投递，必须在调用方的事务里执行。
//
// 顺序不能反：扇出每批写入前会复查事件状态，先终止就能保证它不再补写新条目；
// 反过来先删条目，扇出随后又会把死链写回来。
func (s *Service) TerminateObject(ctx context.Context, tx *gorm.DB, objectType string, objectIDs []int64) error {
	if len(objectIDs) == 0 {
		return nil
	}
	return s.terminateObject(ctx, tx, objectType, objectIDs)
}

// TerminateTopic 回收主题（文章/说说）及其下评论产生的通知。
//
// 主题删了评论还在库里，但读者点进去只会看到 404，所以挂在上面的条目要一并回收。
func (s *Service) TerminateTopic(ctx context.Context, tx *gorm.DB, topicType string, topicIDs []int64) error {
	if len(topicIDs) == 0 {
		return nil
	}
	if err := s.terminateObject(ctx, tx, topicType, topicIDs); err != nil {
		return err
	}

	commentType, ok := commentTypeOf(topicType)
	if !ok {
		return nil
	}
	// 评论 id 用子查询下推：主题下评论可能上千，没必要先取回再拼 IN；
	// 真正要删的是条目，条目数量远小于评论数
	conditions := fmt.Sprintf(
		"object_type = ? AND object_id IN (SELECT id FROM %s WHERE type = ? AND topic_id IN ?)",
		s.comments.TableName())
	return s.terminateWhere(ctx, tx, conditions, notify.ObjectComment, commentType, topicIDs)
}

// RetractInteraction 撤销「某个触发者在某个对象上的一次互动」产生的条目与投递。
//
// 对象本身还在（文章还在，只是那条评论被驳回），因此不能按对象回收；
// 可分组场景一个触发者对同一对象只落一条条目，按去重键即可精确定位。
//
// 条件里不带收件人：「按对象分组」的场景中一个对象只有一个收件人（文章、说说的作者），
// 所以不会波及他人。多个收件人共享同一对象的场景一旦出现，这里必须补上收件人条件。
func (s *Service) RetractInteraction(ctx context.Context, tx *gorm.DB, sceneKey, objectType string, objectID int64, actorID string) error {
	scene, ok := notify.Lookup(sceneKey)
	if !ok {
		return fmt.Errorf("notification: unknown scene %q", sceneKey)
	}
	if scene.Grouping != notify.GroupingObject {
		return nil
	}
	key := InteractionKey(scene, objectType, objectID, actorID)
	if _, err := s.deleteEntries(ctx, tx, "dedupe_key = ?", key); err != nil {
		return fmt.Errorf("notification: retract interaction: %w", err)
	}
	return nil
}

// terminateObject 是 TerminateObject 的无判空版本，供同一事务内已确认非空的调用方使用。
func (s *Service) terminateObject(ctx context.Context, tx *gorm.DB, objectType string, objectIDs []int64) error {
	return s.terminateWhere(ctx, tx, "object_type = ? AND object_id IN ?", objectType, objectIDs)
}

// terminateWhere 按条件终止关联事件，并删除对应的条目与投递。
//
// 条件里同时约束 object_type 与 object_id，调用方可以用固定列表，也可以用子查询。
func (s *Service) terminateWhere(ctx context.Context, tx *gorm.DB, conditions string, args ...interface{}) error {
	statuses := []string{StatusDone, StatusFailed, StatusTerminated}
	if _, err := s.events.WithTx(tx).UpdateFields(ctx,
		map[string]interface{}{"status": StatusTerminated},
		conditions+" AND status NOT IN ?", append(args, statuses)...); err != nil {
		return fmt.Errorf("notification: terminate events: %w", err)
	}

	_, err := s.deleteEntries(ctx, tx, conditions, args...)
	return err
}

// deleteEntries 删除满足条件的条目及其投递记录。
func (s *Service) deleteEntries(ctx context.Context, tx *gorm.DB, conditions string, args ...interface{}) (int64, error) {
	rows, err := s.notifs.WithTx(tx).FindAll(ctx, conditions, args...)
	if err != nil {
		return 0, fmt.Errorf("notification: query entries: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.Id)
	}

	if _, err := s.deliveries.WithTx(tx).DeleteBatch(ctx, "notification_id IN ?", ids); err != nil {
		return 0, fmt.Errorf("notification: delete deliveries of entries: %w", err)
	}
	deleted, err := s.notifs.WithTx(tx).DeleteBatch(ctx, "id IN ?", ids)
	if err != nil {
		return 0, fmt.Errorf("notification: delete entries: %w", err)
	}
	return deleted, nil
}

// commentTypeOf 把主题类型映射成 t_comment.type 的取值；不可评论的主题返回 false。
func commentTypeOf(topicType string) (int64, bool) {
	switch topicType {
	case notify.ObjectArticle:
		return enums.CommentTypeArticle, true
	case notify.ObjectTalk:
		return enums.CommentTypeTalk, true
	default:
		return 0, false
	}
}
