package entity

import "time"

type (
	CollectionName string

	WordData struct {
		WordTrans
		LastRepeat time.Time     `json:"last_repeat"`
		TimeDiff   time.Duration `json:"time_diff"`
	}

	UserWords struct {
		Words map[CollectionName][]WordData `json:"words"`
	}
)
