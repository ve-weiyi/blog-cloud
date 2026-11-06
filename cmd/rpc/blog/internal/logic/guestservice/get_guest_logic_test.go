package guestservicelogic

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mssola/useragent"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/guestrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	testGuestDeviceID = "d41d8cd98f00b204e9800998ecf8427e"
	testUserAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type fakeGuestModel struct {
	model.TGuestModel

	guest        *model.TGuest
	findErr      error
	insertErr    error
	inserted     *model.TGuest
	updateCalls  int
	updateFields map[string]interface{}
}

func (f *fakeGuestModel) FindOneByDeviceId(_ context.Context, _ string) (*model.TGuest, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	if f.guest == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.guest, nil
}

func (f *fakeGuestModel) Insert(_ context.Context, in *model.TGuest) (int64, error) {
	if f.insertErr != nil {
		if errors.Is(f.insertErr, gorm.ErrDuplicatedKey) {
			now := time.Now()
			f.guest = &model.TGuest{
				Id:        9,
				DeviceId:  in.DeviceId,
				Os:        "existing",
				Browser:   "existing",
				IpAddress: "10.0.0.1",
				IpSource:  "existing",
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
		return 0, f.insertErr
	}

	in.Id = 1
	f.inserted = in
	return 1, nil
}

func (f *fakeGuestModel) UpdateFields(_ context.Context, fields map[string]interface{}, _ string, _ ...interface{}) (int64, error) {
	f.updateCalls++
	f.updateFields = fields
	return 1, nil
}

func testGuestContext(deviceID string) context.Context {
	md := metadata.Pairs(
		metax.HeaderRPCRemoteIP, "127.0.0.1",
		metax.HeaderRPCRemoteAgent, testUserAgent,
	)
	if deviceID != "" {
		md.Set(headerx.HeaderXDeviceId, deviceID)
	}
	return metadata.NewIncomingContext(context.Background(), md)
}

func stubIPSource(t *testing.T) {
	t.Helper()
	old := ipSourceLookup
	ipSourceLookup = func(string) string { return "本地" }
	t.Cleanup(func() { ipSourceLookup = old })
}

func expectedEnv() (osName, browser string) {
	ua := useragent.New(testUserAgent)
	browser, _ = ua.Browser()
	return ua.OS(), browser
}

func TestGetGuest_CreatesNewIdentity(t *testing.T) {
	stubIPSource(t)
	fake := &fakeGuestModel{}

	logic := NewGetGuestLogic(testGuestContext(""), &svc.ServiceContext{TGuestModel: fake})
	resp, err := logic.GetGuest(&guestrpc.GetGuestRequest{})
	if err != nil {
		t.Fatalf("GetGuest() error: %v", err)
	}

	if fake.inserted == nil {
		t.Fatal("expected a new guest row")
	}
	if len(resp.Guest.DeviceId) != 32 {
		t.Fatalf("DeviceId = %q, want 32 hex chars", resp.Guest.DeviceId)
	}
	if resp.Guest.DeviceId != fake.inserted.DeviceId {
		t.Fatalf("response DeviceId = %q, inserted = %q", resp.Guest.DeviceId, fake.inserted.DeviceId)
	}
}

func TestGetGuest_ExistingUnchanged(t *testing.T) {
	stubIPSource(t)
	osName, browser := expectedEnv()

	now := time.Now()
	fake := &fakeGuestModel{guest: &model.TGuest{
		Id:        1,
		DeviceId:  testGuestDeviceID,
		Os:        osName,
		Browser:   browser,
		IpAddress: "127.0.0.1",
		IpSource:  "本地",
		CreatedAt: now,
		UpdatedAt: now,
	}}

	logic := NewGetGuestLogic(testGuestContext(testGuestDeviceID), &svc.ServiceContext{TGuestModel: fake})
	resp, err := logic.GetGuest(&guestrpc.GetGuestRequest{})
	if err != nil {
		t.Fatalf("GetGuest() error: %v", err)
	}

	if resp.Guest.DeviceId != testGuestDeviceID {
		t.Fatalf("DeviceId = %q, want %q", resp.Guest.DeviceId, testGuestDeviceID)
	}
	if fake.updateCalls != 0 {
		t.Fatalf("UpdateFields calls = %d, want 0", fake.updateCalls)
	}
}

func TestGetGuest_ExistingChanged(t *testing.T) {
	stubIPSource(t)
	osName, browser := expectedEnv()

	now := time.Now()
	fake := &fakeGuestModel{guest: &model.TGuest{
		Id:        1,
		DeviceId:  testGuestDeviceID,
		Os:        "old-os",
		Browser:   "old-browser",
		IpAddress: "10.0.0.1",
		IpSource:  "old-source",
		CreatedAt: now,
		UpdatedAt: now,
	}}

	logic := NewGetGuestLogic(testGuestContext(testGuestDeviceID), &svc.ServiceContext{TGuestModel: fake})
	if _, err := logic.GetGuest(&guestrpc.GetGuestRequest{}); err != nil {
		t.Fatalf("GetGuest() error: %v", err)
	}

	if fake.updateCalls != 1 {
		t.Fatalf("UpdateFields calls = %d, want 1", fake.updateCalls)
	}
	if fake.updateFields["os"] != osName || fake.updateFields["browser"] != browser {
		t.Fatalf("update fields = %v, want os=%q browser=%q", fake.updateFields, osName, browser)
	}
}

func TestGetGuest_InvalidCandidateGetsNewIdentity(t *testing.T) {
	stubIPSource(t)
	fake := &fakeGuestModel{}

	// 36 位 UUID 是用户 ID 的形态，不能作为游客标识
	logic := NewGetGuestLogic(testGuestContext("0190f3a1-7c2e-7c3a-9f4e-1a2b3c4d5e6f"), &svc.ServiceContext{TGuestModel: fake})
	resp, err := logic.GetGuest(&guestrpc.GetGuestRequest{})
	if err != nil {
		t.Fatalf("GetGuest() error: %v", err)
	}

	if len(resp.Guest.DeviceId) != 32 {
		t.Fatalf("DeviceId = %q, want new 32 hex id", resp.Guest.DeviceId)
	}
	if fake.inserted == nil {
		t.Fatal("expected a new guest row")
	}
}

func TestGetGuest_InsertConflictReusesRow(t *testing.T) {
	stubIPSource(t)
	fake := &fakeGuestModel{insertErr: gorm.ErrDuplicatedKey}

	logic := NewGetGuestLogic(testGuestContext(""), &svc.ServiceContext{TGuestModel: fake})
	resp, err := logic.GetGuest(&guestrpc.GetGuestRequest{})
	if err != nil {
		t.Fatalf("GetGuest() error: %v", err)
	}

	if resp.Guest.Id != 9 {
		t.Fatalf("Guest.Id = %d, want existing row 9", resp.Guest.Id)
	}
	if fake.updateCalls != 0 {
		t.Fatalf("UpdateFields calls = %d, want 0", fake.updateCalls)
	}
}
