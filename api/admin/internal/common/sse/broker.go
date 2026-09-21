// Package sse 承载通知 SSE 推送的进程内管道：泛型广播器 Broker[T]（每频道一实例）、
// Redis 订阅器（Subscriber）、到推送通道的桥接（Bridge）。
//
// 三者都不导入 svc.ServiceContext：导入 svc 会把整套 rpc client（含 protobuf 注册）
// 拉进测试二进制，本仓库里 proto 文件名重复注册会直接 panic，因此这条管道保持
// 无 svc 依赖以便单测（同样的约束见 permissionx / tracelogx 的测试）。
//
// 与 scan 项目 core/sse.Broker 的差异：不缓存"最近一次的值"。本项目的推送是"信号"
// 而非"快照"，客户端连上后会回查权威数据，回放旧信号既多余又可能已失效。
package sse

import (
	"sync"
)

// SubscriberBuffer 单订阅者投递缓冲容量。写满即丢弃当前帧——慢消费者不阻塞
// 广播方，也不断开连接（丢帧由客户端重连时的回查补齐）。
const SubscriberBuffer = 64

// Broker 单频道广播器：多个订阅者共享同一事件流。
type Broker[T any] struct {
	mu     sync.RWMutex
	subs   map[chan T]struct{}
	closed bool
}

// NewBroker 创建广播器
func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{subs: make(map[chan T]struct{})}
}

// Subscribe 注册订阅者，返回接收通道与注销函数。
//
// 注销函数幂等，且会关闭通道（消费方以 ok=false 退出）。关闭是安全的：
// 注销持写锁、Publish 持读锁，二者互斥，不会向已关闭的通道发送。
func (b *Broker[T]) Subscribe() (<-chan T, func()) {
	ch := make(chan T, SubscriberBuffer)

	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		close(ch)
		return ch, func() {}
	}
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			b.mu.Lock()
			defer b.mu.Unlock()
			if _, ok := b.subs[ch]; !ok {
				// 已被 Close 收走：通道在那时已关闭，这里不能再关一次
				return
			}
			delete(b.subs, ch)
			close(ch)
		})
	}

	return ch, cancel
}

// Publish 向全部订阅者投递；订阅者缓冲满时丢弃该帧。
func (b *Broker[T]) Publish(v T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subs {
		select {
		case ch <- v:
		default: // 慢消费者：丢弃本帧
		}
	}
}

// Len 当前订阅者数。可用于"无人观看时跳过取数"这类优化。
func (b *Broker[T]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs)
}

// Close 关闭所有订阅者通道并停止扇出。进程退出时调用。
func (b *Broker[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}
	b.closed = true
	for ch := range b.subs {
		delete(b.subs, ch)
		close(ch)
	}
}
