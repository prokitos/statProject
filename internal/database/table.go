package database

type Statistic struct {
	Timestamp  string
	Country    string
	Os         string
	Browser    string
	Request    int
	Impression int
}
