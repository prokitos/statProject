package database

import (
	"fmt"
	"time"
)

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

// каждые 10 секунд заходит сюда
func (m *Manager) StartInserting() {
	fmt.Println("insert")

	// безопасно получаем данные из менеджера
	rows := m.Withdraw()
	if len(rows) == 0 {
		fmt.Println("nothing")
		return
	}

	// если данные есть. то записываем их. иначе добавить текущие в мэнеджер, и пустить на новый круг.
	err := m.writer.ClickHouseInsert(rows)
	if err != nil {
		fmt.Println("error")
		m.Summary(rows)
		return
	}

	fmt.Println("all goods writtern")
}
