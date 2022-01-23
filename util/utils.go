package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func RandomCode(n int) string {
	var letters = []byte("1234567890")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}


func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn((len(letters)))]
	}
	return string(result)
}

// GetTime 获取time
func GetTime(timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	fmt.Println(timeStr)
	ret, err := time.Parse(layout, timeStr)

	if err != nil {
		return time.Now(), err
	} else {
		return ret, nil
	}
}


func CheckEmailMatch(email string) bool {
	emailMatch, _ := regexp.MatchString("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])",email)
	return emailMatch
}

func CheckUsernameMatch(username string) bool {
	usernameMatch, _ := regexp.MatchString("^[a-zA-z\\d]{8,32}$", username)
	return usernameMatch
}

func CheckPasswordMatch(password string) bool {
	passwordMatch, _ := regexp.MatchString("^[a-zA-Z\\d]{8,32}$", password)
	return passwordMatch
}