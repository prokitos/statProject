package clickhouseStat

import (
	"fmt"
	"time"
)

type (
	Key struct {
		Timestamp  int64
		Country    string
		Os         string
		Browser    string
		CampaginId uint32
	}

	Value struct {
		Requests    int64
		Impressions int64
	}

	Rows map[Key]Value
)

func (a Value) Assign(b Value) Value {
	res := a
	res.Requests += b.Requests
	res.Impressions += b.Impressions
	return res
}

func NewKey(k Key) Key {
	k.Timestamp = time.Now().Unix()
	k.Timestamp -= k.Timestamp % 60
	return k
}

func (m *Manager) StartInserting() {
	fmt.Println("insert")

	rows := m.Withdraw()
	if len(rows) == 0 {
		fmt.Println("nothing")
		return
	}

	err := m.writer.Insert(rows)
	if err != nil {
		fmt.Println("error")
		m.AppendAllRows(rows)
		return
	}

	fmt.Println("all goods writtern %d", len(rows))
}

func (m *Manager) Withdraw() Rows {
	m.mu.Lock()
	defer m.mu.Unlock()

	rows := m.rows
	m.rows = newRows()

	return rows
}

func (m *Manager) StartTimer() {
	fmt.Println("timer start")
	go m.LoopTimer()
}

func (m *Manager) LoopTimer() {
	for {
		select {

		case <-time.After(m.flushInterval):
			m.StartInserting()

		case <-m.ctx.Done():
			m.StartInserting()
			return

		}
	}
}
