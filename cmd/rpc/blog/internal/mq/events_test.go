package mq

import "testing"

func TestSubscriptionSpecs(t *testing.T) {
	for _, spec := range Subscriptions() {
		t.Run(spec.Group, func(t *testing.T) {
			if _, err := spec.Prepare(); err != nil {
				t.Fatalf("Prepare: %v", err)
			}
			if spec.DeadLetter.Enabled && spec.DeadLetter.Destination == spec.Destination {
				t.Fatal("dead letter destination must differ from destination")
			}
			// 通知域三条死信队列共用一个 exchange：绑定必须收窄到本事件，
			// 留空（退化为 #）会让同一条死信复制进所有 dlq 队列。
			if spec.DeadLetter.Enabled && spec.DeadLetter.RoutingKey != spec.BindingKey {
				t.Fatalf("dead letter routing key = %q, want %q",
					spec.DeadLetter.RoutingKey, spec.BindingKey)
			}
		})
	}
}
