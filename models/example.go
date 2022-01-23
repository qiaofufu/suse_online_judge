package models

import (
	"errors"
	"gorm.io/gorm"
	"onlineJudge/common"
)

type Example struct {
	gorm.Model
	Input      string `json:"input"`
	Output     string `json:"output"`
	ProblemID  uint   `json:"problem_id"`
}

func (e Example) GetsExampleList(problemID uint) ([]Example, error) {
	var list []Example
	result := common.DB.Model(Example{}).Where("problem_id = ?", problemID).Find(&list)
	if result.Error != nil {
		return nil, errors.New("内部错误")
	}
	return list, nil
}
