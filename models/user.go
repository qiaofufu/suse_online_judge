package models

import (
	"errors"
	"gorm.io/gorm"
	"onlineJudge/common"
)

type User struct {
	gorm.Model
	UserName 		string `json:"user_name"`
	UserEmail 		string `json:"user_email"`
	UserPassword 	string `json:"user_password"`
}

func (t User) GetByID(ID uint) (User, error) {
	var user User
	result := common.DB.Model(&User{}).Where("id = ?", ID).First(&user)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return user, errors.New("内部错误1")
		} else {
			return user, errors.New("题目不存在")
		}
	}
	return user, nil
}

func (t User) Rank() (interface{}, error) {
	var userList []User
	type responseDTO struct {
		UserName string
		AcNumber int64
		Integral uint
	}
	var rank []responseDTO
	common.DB.Model(&User{}).Find(&userList)

	for _, v := range userList {
		var integral UserIntegral
		common.DB.Model(&UserIntegral{}).Where("user_id = ?", v.ID).First(&integral)
		var acNumber int64
		common.DB.Model(&ProblemStatusRecord{}).Where("user_id = ? and status = ?", v.ID, Accept).Count(&acNumber)
		rank = append(rank, responseDTO{UserName: v.UserName, AcNumber: acNumber, Integral: integral.Value})
	}
	return rank, nil
}