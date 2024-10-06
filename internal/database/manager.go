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
	rows          Statistic
}

var GlobalManager *Manager
var defaultCapacity = 1000

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

func newRows() Statistic {
	return Statistic{}
}

type Rows Statistic

func (m *Manager) Withdraw() Statistic {
	m.mu.Lock()
	defer m.mu.Unlock()

	rows := m.rows
	m.rows = newRows()

	return rows
}

func (m *Manager) Summary(row Statistic) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rows.Request += row.Request
	m.rows.Impression += row.Impression
	m.rows.Browser = row.Browser
	m.rows.Os = row.Os
	m.rows.Country = row.Country
	m.rows.Timestamp = row.Timestamp

}

// добавление данных при запросе в менеджер.
func (m *Manager) Adding(data Statistic) {
	m.rows = data
}
