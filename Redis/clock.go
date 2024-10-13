package Redis

import (
	"clock/common"
	"context"
	"time"
)

func SetIsClock(uid string) error {
	rdb := common.GetRDB()
	cbg := context.Background()

	key := isClock + uid
	err := rdb.Set(cbg, key, "1", time.Hour*24).Err()
	return err
}

func GetIsClock(uid string) (bool, error) {
	rdb := common.GetRDB()
	cbg := context.Background()
	key := isClock + uid

	res, err := rdb.Get(cbg, key).Result()
	if res == "1" && err == nil {
		return true, nil
	}

	return false, err
}

func DeleteIsClock(uid string) error {
	rdb := common.GetRDB()
	cbg := context.Background()
	key := isClock + uid

	err := rdb.Del(cbg, key).Err()
	return err
}
