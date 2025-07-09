package client

import "github.com/go-stomp/stomp/v3/frame"

type txStore struct {
	transactions map[string][]*frame.Frame
}

func newTxStore() *txStore {
	return &txStore{
		transactions: make(map[string][]*frame.Frame),
	}
}

func (txs *txStore) Begin(tx string) error {
	if _, exists := txs.transactions[tx]; exists {
		return errTransactionExists
	}
	txs.transactions[tx] = make([]*frame.Frame, 0)
	return nil
}

func (txs *txStore) Add(tx string, f *frame.Frame) error {
	frames, exists := txs.transactions[tx]
	if !exists {
		return errTransactionNotFound
	}
	// 移除 transaction 头部
	f.Header.Del(frame.Transaction)
	txs.transactions[tx] = append(frames, f)
	return nil
}

func (txs *txStore) Commit(tx string, commitFunc func(f *frame.Frame) error) error {
	frames, exists := txs.transactions[tx]
	if !exists {
		return errTransactionNotFound
	}

	for _, f := range frames {
		if err := commitFunc(f); err != nil {
			return err
		}
	}

	delete(txs.transactions, tx)
	return nil
}

func (txs *txStore) Abort(tx string) error {
	if _, exists := txs.transactions[tx]; !exists {
		return errTransactionNotFound
	}
	delete(txs.transactions, tx)
	return nil
}

func (txs *txStore) Clear() {
	txs.transactions = make(map[string][]*frame.Frame)
}
