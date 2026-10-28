package outbox

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

type fakeStore struct {
	rows           []Record
	stats          Stats
	published      []int64
	publishedToken string
	failed         []failedCall
	dead           []int64
	deadToken      string
}

type failedCall struct {
	id      int64
	token   string
	retries int
	next    time.Time
	cause   error
}

func (f *fakeStore) Enqueue(context.Context, *gorm.DB, mqx.Message) error { return nil }

func (f *fakeStore) Claim(context.Context, int) ([]Record, error) { return f.rows, nil }

func (f *fakeStore) MarkPublished(_ context.Context, ids []int64, token string) error {
	f.published = append(f.published, ids...)
	f.publishedToken = token
	return nil
}

func (f *fakeStore) MarkFailed(_ context.Context, id int64, token string, retries int, next time.Time, cause error) error {
	f.failed = append(f.failed, failedCall{id: id, token: token, retries: retries, next: next, cause: cause})
	return nil
}

func (f *fakeStore) MarkDead(_ context.Context, id int64, token string, _ error) error {
	f.dead = append(f.dead, id)
	f.deadToken = token
	return nil
}

func (f *fakeStore) ResetStuck(context.Context, time.Time) (int64, error) { return 0, nil }

func (f *fakeStore) DeletePublished(context.Context, time.Time) (int64, error) { return 0, nil }

func (f *fakeStore) Stats(context.Context) (Stats, error) { return f.stats, nil }

type fakeMetrics struct {
	calls   int
	pending int64
	dead    int64
	oldest  time.Duration
}

func (f *fakeMetrics) OutboxStats(pending, dead int64, oldestAge time.Duration) {
	f.calls++
	f.pending = pending
	f.dead = dead
	f.oldest = oldestAge
}

type fakePublisher struct {
	err       error
	published []mqx.Message
}

func (f *fakePublisher) Publish(_ context.Context, msg *mqx.Message) (*mqx.PublishResult, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.published = append(f.published, *msg)
	return &mqx.PublishResult{ID: msg.ID}, nil
}

func newTestRecord() Record {
	now := time.Now()
	return Record{
		ID:            1,
		EventID:       "event-1",
		Destination:   "test.exchange",
		RoutingKey:    "test.event",
		Headers:       `{"mqx.event_id":"event-1"}`,
		Payload:       []byte(`{"id":1}`),
		Status:        StatusPublishing,
		LockToken:     "token-1",
		NextAttemptAt: now,
		OccurredAt:    now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func TestRunOncePublishes(t *testing.T) {
	store := &fakeStore{rows: []Record{newTestRecord()}}
	pub := &fakePublisher{}
	relay := NewRelay(store, pub)

	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if len(pub.published) != 1 {
		t.Fatalf("published = %d, want 1", len(pub.published))
	}
	if len(store.published) != 1 || store.published[0] != 1 {
		t.Fatalf("mark published = %v, want [1]", store.published)
	}
	if store.publishedToken != "token-1" {
		t.Fatalf("published token = %q, want token-1", store.publishedToken)
	}
	if len(store.failed) != 0 || len(store.dead) != 0 {
		t.Fatalf("unexpected failed=%v dead=%v", store.failed, store.dead)
	}
}

func TestRunOnceRetriesOnFailure(t *testing.T) {
	store := &fakeStore{rows: []Record{newTestRecord()}}
	pub := &fakePublisher{err: errors.New("broker down")}
	relay := NewRelay(store, pub)
	before := time.Now()

	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if len(store.published) != 0 {
		t.Fatalf("should not mark published: %v", store.published)
	}
	if len(store.failed) != 1 {
		t.Fatalf("failed = %d, want 1", len(store.failed))
	}
	if store.failed[0].retries != 1 {
		t.Fatalf("retries = %d, want 1", store.failed[0].retries)
	}
	if store.failed[0].token != "token-1" {
		t.Fatalf("failed token = %q, want token-1", store.failed[0].token)
	}
	if !store.failed[0].next.After(before) {
		t.Fatalf("next attempt = %v, want after %v", store.failed[0].next, before)
	}
}

func TestRunOnceMarksDeadAfterMaxRetries(t *testing.T) {
	row := newTestRecord()
	row.RetryCount = defaultMaxRetries - 1
	store := &fakeStore{rows: []Record{row}}
	pub := &fakePublisher{err: errors.New("permanent")}
	relay := NewRelay(store, pub)

	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if len(store.dead) != 1 || store.dead[0] != 1 {
		t.Fatalf("dead = %v, want [1]", store.dead)
	}
	if store.deadToken != "token-1" {
		t.Fatalf("dead token = %q, want token-1", store.deadToken)
	}
	if len(store.failed) != 0 {
		t.Fatalf("failed = %v, want none", store.failed)
	}
}

func TestRunOnceCorruptHeadersMarksDead(t *testing.T) {
	row := newTestRecord()
	row.Headers = "not-json"
	store := &fakeStore{rows: []Record{row}}
	pub := &fakePublisher{}
	relay := NewRelay(store, pub)

	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if len(pub.published) != 0 {
		t.Fatalf("should not publish: %d", len(pub.published))
	}
	if len(store.dead) != 1 {
		t.Fatalf("dead = %v, want [1]", store.dead)
	}
}

func TestRunOnceMarksDeadOnPermanentPublishError(t *testing.T) {
	store := &fakeStore{rows: []Record{newTestRecord()}}
	pub := &fakePublisher{err: mqx.Permanent(errors.New("permission denied"))}
	relay := NewRelay(store, pub)

	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if len(store.dead) != 1 {
		t.Fatalf("dead = %v, want [1]", store.dead)
	}
	if len(store.failed) != 0 {
		t.Fatalf("failed = %v, want none", store.failed)
	}
}

func TestRunOnceKeepsRetryingWhenBrokerUnavailable(t *testing.T) {
	row := newTestRecord()
	row.RetryCount = defaultMaxRetries - 1
	store := &fakeStore{rows: []Record{row}}
	pub := &fakePublisher{err: mqx.ErrUnavailable}
	relay := NewRelay(store, pub)

	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if len(store.dead) != 0 {
		t.Fatalf("broker unavailable must not mark dead: %v", store.dead)
	}
	if len(store.failed) != 1 {
		t.Fatalf("failed = %d, want 1", len(store.failed))
	}
}

func TestRunOnceNoRows(t *testing.T) {
	store := &fakeStore{}
	relay := NewRelay(store, &fakePublisher{})
	if err := relay.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
}

func TestRecordMessage(t *testing.T) {
	row := newTestRecord()
	msg, err := row.Message()
	if err != nil {
		t.Fatalf("Message: %v", err)
	}
	if msg.ID != row.EventID || msg.Destination != row.Destination || msg.RoutingKey != row.RoutingKey {
		t.Fatalf("message = %+v", msg)
	}
	if string(msg.Body) != string(row.Payload) || msg.Headers["mqx.event_id"] != row.EventID {
		t.Fatalf("message payload = %+v", msg)
	}
}

func TestRelayReportsOutboxStats(t *testing.T) {
	oldest := time.Now().Add(-3 * time.Minute)
	store := &fakeStore{stats: Stats{Pending: 5, Dead: 2, Oldest: &oldest}}
	recorder := &fakeMetrics{}
	relay := NewRelay(store, &fakePublisher{}, WithMetrics(recorder))

	relay.reportStats(context.Background())

	if recorder.calls != 1 {
		t.Fatalf("OutboxStats calls = %d, want 1", recorder.calls)
	}
	if recorder.pending != 5 || recorder.dead != 2 {
		t.Fatalf("stats = pending %d dead %d, want 5/2", recorder.pending, recorder.dead)
	}
	if recorder.oldest < 3*time.Minute {
		t.Fatalf("oldest age = %v, want >= 3m", recorder.oldest)
	}
}

func TestRelayStartStop(t *testing.T) {
	relay := NewRelay(&fakeStore{}, &fakePublisher{})
	relay.Start()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	relay.Stop(ctx)
}
