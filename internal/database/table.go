package database

import "time"

type Statistic struct {
	Timestamp  time.Time
	Country    string
	Os         string
	Browser    string
	Request    int
	Impression int
}
