package database

import (
	"context"
	"sync"
	"time"
)

type Manager struct {
	writer        *ClickDatabase
	flushInterval time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.RWMutex
	rows          []Statistic
}

var GlobalManager *Manager

func NewManager(db *ClickDatabase, flushInterval time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	var temp *Manager = &Manager{
		writer:        db,
		flushInterval: flushInterval,
		ctx:           ctx,
		cancel:        cancel,
		rows:          newRows(),
	}

	GlobalManager = temp

	return temp

}

func newRows() []Statistic {
	return make([]Statistic, 0)
}

type Rows Statistic

func (m *Manager) Withdraw() []Statistic {
	m.mu.Lock()
	defer m.mu.Unlock()

	rows := m.rows
	m.rows = newRows()

	return rows
}

func (m *Manager) Summary(row []Statistic) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rows = append(m.rows, row...)
}

// добавление данных при запросе в менеджер.
func (m *Manager) Adding(data Statistic) {
	m.rows = append(m.rows, data)
}
