package poloniex

import (
	"fmt"
	"sync"
)

type OrderObserver interface {
	Observe(side, symbol, orderID string) error
	Items(orderID string) (ServerableObject, error)
	// TODO: Delete after order completely fill.
	Delete(orderID string) error
	Lock() error
	Unlock()
	IsObservable(orderID string) bool
}

type ServerableObject struct {
	side    string
	symbol  string
	orderID string
}

// WebsocketObserver реализация OrderObserver для Websocket.
// Используется для синхронизации WS & REST
type WebsocketObserver struct {
	items   map[string]ServerableObject
	itemsMu sync.RWMutex
	mu      sync.Mutex
}

func NewWebsocketObserver() *WebsocketObserver {
	return &WebsocketObserver{
		items: make(map[string]ServerableObject),
	}
}

func (w *WebsocketObserver) IsObservable(orderID string) bool {
	w.itemsMu.RLock()
	defer w.itemsMu.RUnlock()

	if _, ok := w.items[orderID]; ok {
		return true
	}

	return false
}

func (w *WebsocketObserver) Observe(side, symbol, orderID string) error {
	w.itemsMu.RLock()

	if _, ok := w.items[orderID]; ok {
		w.itemsMu.RUnlock()
		return fmt.Errorf("already exists: %v", orderID)
	}

	w.itemsMu.RUnlock()
	w.itemsMu.Lock()
	w.items[orderID] = ServerableObject{
		side:    side,
		symbol:  symbol,
		orderID: orderID,
	}
	w.itemsMu.Unlock()

	return nil
}

func (w *WebsocketObserver) Items(orderID string) (ServerableObject, error) {
	w.itemsMu.RLock()

	if value, ok := w.items[orderID]; ok {
		w.itemsMu.RUnlock()
		return value, nil
	}

	w.itemsMu.RUnlock()
	return ServerableObject{}, fmt.Errorf("orderID %v not registered", orderID)
}

func (w *WebsocketObserver) Delete(orderID string) error {
	w.itemsMu.RLock()
	if _, ok := w.items[orderID]; !ok {
		w.itemsMu.RUnlock()
		return fmt.Errorf("not found: %v", orderID)
	}
	w.itemsMu.RUnlock()

	w.itemsMu.Lock()
	delete(w.items, orderID)
	w.itemsMu.Unlock()

	return nil
}

// Lock TODO: Сделать кастомный Locker, чтобы возвращать ошибку, что блокировка длится дольше T
func (w *WebsocketObserver) Lock() error {
	w.mu.Lock()
	return nil
}

func (w *WebsocketObserver) Unlock() {
	w.mu.Unlock()
}

// NilObserver пустая реализации без синхронизаций. Используется, если получение трейдов из WebSocket не нужен.
type NilObserver struct{}

func NewNilObserver() *NilObserver {
	return &NilObserver{}
}

func (n *NilObserver) Observe(_, _ string, _ int64) error {
	return nil
}

func (n *NilObserver) Delete(_ int64) error {
	return nil
}

func (n *NilObserver) Lock() error {
	return nil
}

func (n *NilObserver) Unlock() {
}

func (n *NilObserver) IsObservable(orderID int64) bool {
	return false
}
