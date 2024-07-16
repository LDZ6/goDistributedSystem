package models

/**
使用redis需实现Store中的三个方法
type Store interface {
    // Set sets the digits for the captcha id.
    Set(id string, value string)
    // Get returns stored digits for the captcha id. Clear indicates
    // whether the captcha must be deleted from the store.
    Get(id string, clear bool) string
    //Verify captcha's answer directly
    Verify(id, answer string, clear bool) bool
}
*/

import (
	"context"
	"fmt"
	"time"
)

var ctx = context.Background()

const CAPTCHA = "captcha:"

type RedisStore struct {
}

// 实现设置 captcha 的方法
func (r RedisStore) Set(id string, value string) error {
	key := CAPTCHA + id
	err := RedisDb.Set(ctx, key, value, time.Minute*2).Err()
	return err
}

// 实现获取 captcha 的方法
func (r RedisStore) Get(id string, clear bool) string {
	key := CAPTCHA + id
	//获取 captcha
	val, err := RedisDb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//如果clear == true, 则删除
	if clear {
		err := RedisDb.Del(ctx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// 实现验证 captcha 的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := RedisStore{}.Get(id, clear)
	return v == answer
}
