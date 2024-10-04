package responseObject

import "clock/model"

type Rank struct {
	RankMsg []model.SimpleUser
	Name    string
	Date    string
}
