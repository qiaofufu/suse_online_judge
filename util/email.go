package util

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"onlineJudge/common"
)

// VerifyEmailCode 验证email验证码
func VerifyEmailCode(emailCode string, email string) bool {
	redis := common.RedisClient
	n, err := redis.Exists(CodeKey(email)).Result()
	if err != nil {
		panic(err)
	}
	if n > 0 {
		res, err := redis.Get(CodeKey(email)).Result()
		redis.Del(CodeKey(email))
		if err != nil {
			panic(err)
		} else {
			return res == emailCode
		}
	}
	return false
}

func SendEmail(email string, code string) bool {
	m := gomail.NewMessage()
	m.SetHeader("From",viper.GetString("email.address"))
	m.SetHeader("To",email)
	m.SetAddressHeader("Aa",viper.GetString("email.address"), "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello you register code is <b>"+code+"</b>!")
	d := gomail.NewDialer(viper.GetString("email.host"),viper.GetInt("email.port"),viper.GetString("email.address"),viper.GetString("email.password"))
	if err := d.DialAndSend(m); err != nil {
		return false
	}
	return true
}