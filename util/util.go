package util

import (
	"math/rand"
	"time"
)

// // RandomString 随机生成字符串的函数
//
//	func RandomString(n int) string {
//		var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
//		result := make([]byte, n)
//
//		rand.Seed(time.Now().Unix())
//		for i := range result {
//			result[i] = letters[rand.Intn(len(letters))]
//		}
//
//		return string(result)
//	}
//
// 字母集合
const letters = "asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP"

// RandomString 随机生成字符串的函数
func RandomString(n int) string {
	// 创建一个新的随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, n)

	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}

	return string(result)
}
