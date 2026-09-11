package discussionservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

func convertChatOut(in *model.TChat) *discussionrpc.Chat {
	return &discussionrpc.Chat{
		Id:        in.Id,
		UserId:    in.UserId,
		DeviceId:  in.DeviceId,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		Nickname:  in.Nickname,
		Avatar:    in.Avatar,
		Type:      in.Type,
		Content:   in.Content,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}

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
