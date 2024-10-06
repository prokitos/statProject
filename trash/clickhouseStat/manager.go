package clickhouseStat

import (
	"context"
	"sync"
	"time"
)

const (
	defaultCapacity = 1000
)

var GlobalManager *Manager

type (
	Writer interface {
		Insert(rows Rows) error
	}
	Manager struct {
		writer        Writer
		flushInterval time.Duration
		ctx           context.Context
		cancel        context.CancelFunc
		mu            sync.RWMutex
		rows          Rows
	}
)

func NewManager(w Writer, flushInterval time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	var temp *Manager = &Manager{
		writer:        w,
		flushInterval: flushInterval,
		ctx:           ctx,
		cancel:        cancel,
		rows:          newRows(),
	}

	GlobalManager = temp

	return temp

}

func newRows() Rows {
	return make(Rows, defaultCapacity)
}

func (m *Manager) Append(k Key, v Value) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.AddToMap(k, v)
}

func (m *Manager) AppendAllRows(row Rows) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range row {
		m.AddToMap(k, v)
	}
}

func (m *Manager) AddToMap(k Key, v Value) {
	current := m.rows[k]
	current = current.Assign(v)

	m.rows[k] = current
}
