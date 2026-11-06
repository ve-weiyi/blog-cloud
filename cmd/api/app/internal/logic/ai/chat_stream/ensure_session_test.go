package chat_stream

import (
	"context"
	"errors"
	"strings"
	"testing"

	"google.golang.org/grpc"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/aiservice"
	"github.com/ve-weiyi/blog-cloud/infra/ai"
)

type fakeAiService struct {
	created       bool
	nextSessionID int64
	getSessionErr error
}

func (f *fakeAiService) CreateSession(ctx context.Context, in *aiservice.CreateSessionRequest, opts ...grpc.CallOption) (*aiservice.CreateSessionResponse, error) {
	f.created = true
	return &aiservice.CreateSessionResponse{Session: &aiservice.Session{Id: f.nextSessionID}}, nil
}

func (f *fakeAiService) ListSessions(ctx context.Context, in *aiservice.ListSessionsRequest, opts ...grpc.CallOption) (*aiservice.ListSessionsResponse, error) {
	return &aiservice.ListSessionsResponse{}, nil
}

func (f *fakeAiService) GetSession(ctx context.Context, in *aiservice.GetSessionRequest, opts ...grpc.CallOption) (*aiservice.GetSessionResponse, error) {
	if f.getSessionErr != nil {
		return nil, f.getSessionErr
	}
	return &aiservice.GetSessionResponse{Session: &aiservice.Session{Id: in.Id}}, nil
}

func (f *fakeAiService) ListMessages(ctx context.Context, in *aiservice.ListMessagesRequest, opts ...grpc.CallOption) (*aiservice.ListMessagesResponse, error) {
	return &aiservice.ListMessagesResponse{}, nil
}

func (f *fakeAiService) AppendMessage(ctx context.Context, in *aiservice.AppendMessageRequest, opts ...grpc.CallOption) (*aiservice.AppendMessageResponse, error) {
	return &aiservice.AppendMessageResponse{}, nil
}

func (f *fakeAiService) DeleteSession(ctx context.Context, in *aiservice.DeleteSessionRequest, opts ...grpc.CallOption) (*aiservice.DeleteSessionResponse, error) {
	return &aiservice.DeleteSessionResponse{}, nil
}

func newTestLogic(aiSvc aiservice.AiService) *AiChatStreamLogic {
	return NewAiChatStreamLogic(context.Background(), &svc.ServiceContext{
		AiService: aiSvc,
		AiClient:  ai.NewClient(ai.Config{Model: "test-model"}),
	})
}

// 游客不建会话。
func TestEnsureSessionGuest(t *testing.T) {
	fake := &fakeAiService{nextSessionID: 42}
	logic := newTestLogic(fake)

	sessionID, err := logic.ensureSession(aictx.Identity{Guest: true, DeviceId: "d1"}, 5, "问题")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sessionID != 0 || fake.created {
		t.Fatalf("guest should not create session, got id=%d created=%v", sessionID, fake.created)
	}
}

// 登录用户没有会话 id 时新建。
func TestEnsureSessionCreates(t *testing.T) {
	fake := &fakeAiService{nextSessionID: 42}
	logic := newTestLogic(fake)

	sessionID, err := logic.ensureSession(aictx.Identity{UserId: "u1"}, 0, "第一个问题")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sessionID != 42 || !fake.created {
		t.Fatalf("want created session 42, got id=%d created=%v", sessionID, fake.created)
	}
}

// 传入的会话 id 必须校验归属，失败时不采信、不新建。
func TestEnsureSessionRejectsForeignSession(t *testing.T) {
	fake := &fakeAiService{nextSessionID: 42, getSessionErr: errors.New("not found")}
	logic := newTestLogic(fake)

	_, err := logic.ensureSession(aictx.Identity{UserId: "u1"}, 7, "问题")
	if err == nil {
		t.Fatal("expected ownership error")
	}
	if fake.created {
		t.Fatal("should not create a new session when requested id is invalid")
	}
}

// 游客历史最多 50 条，单条截到 2000 字，非法角色丢弃。
func TestCapGuestHistory(t *testing.T) {
	items := make([]types.AiChatMessage, 0, 60)
	for i := 0; i < 60; i++ {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}
		content := strings.Repeat("字", 3000)
		if i == 55 {
			role = "system"
		}
		items = append(items, types.AiChatMessage{Role: role, Content: content})
	}

	history := capGuestHistory(items)

	if len(history) != 49 {
		t.Fatalf("want 49 valid messages, got %d", len(history))
	}
	for _, message := range history {
		if len([]rune(message.Content)) > maxGuestMessageRunes {
			t.Fatalf("content should be truncated, got %d runes", len([]rune(message.Content)))
		}
	}
}
