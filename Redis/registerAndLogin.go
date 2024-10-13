package Redis

import (
	"clock/common"
	"clock/model"
	"context"
	"encoding/json"
	"time"
)

func SetUserInfo(user model.User) error {
	rdb := common.GetRDB()
	cbg := context.Background()

	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = rdb.Set(cbg, user.Telephone, userJson, 5*time.Minute).Err()

	return err
}

func GetUserInfo(telephone string) (model.User, error) {
	rdb := common.GetRDB()
	cbg := context.Background()

	userJson, err := rdb.Get(cbg, telephone).Result()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
