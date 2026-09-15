package discussionservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/jsonconv"
)

func convertMessageOut(in *model.TMessage) *discussionrpc.Message {
	return &discussionrpc.Message{
		Id:             in.Id,
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		MessageContent: in.MessageContent,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
	}
}

func convertCommentOut(in *model.TComment) *discussionrpc.Comment {
	return &discussionrpc.Comment{
		Id:             in.Id,
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		TopicId:        in.TopicId,
		ParentId:       in.ParentId,
		ReplyId:        in.ReplyId,
		ReplyUserId:    in.ReplyUserId,
		CommentContent: in.CommentContent,
		Type:           in.Type,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
		LikeCount:      in.LikeCount,
	}
}

func convertTalkOut(in *model.TTalk) *discussionrpc.Talk {
	var images []string
	jsonconv.JsonToAny(in.Images, &images)

	return &discussionrpc.Talk{
		Id:        in.Id,
		UserId:    in.UserId,
		Content:   in.Content,
		Images:    images,
		IsTop:     in.IsTop,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
		LikeCount: in.LikeCount,
	}
}
