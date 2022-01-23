package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"onlineJudge/common"
)

type Problem struct {
	gorm.Model
	Title      string `json:"title"`
	IsPublic   bool   `json:"is_public"`

	ContestID uint `json:"contest_id"`
	// HTML
	Description       string `json:"description"`
	InputDescription  string `json:"input_description"`
	OutputDescription string `json:"output_description"`
	Samples           string `json:"samples"`

	TimeLimit   int64 `json:"time_limit"`
	MemoryLimit int64 `json:"memory_limit"`

	IsVisible bool `json:"is_visible"`

	SubmissionNumber int64 `json:"submission_number"`
	AcceptedNumber   int64 `json:"accepted_number"`
}

// GetProblem 获取题目
func (p Problem) GetProblem(problemID uint) (Problem, error) {
	result := common.DB.Model(&Problem{}).Where("id = ?", problemID).First(&p)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return p, errors.New("内部错误1")
		} else {
			return p, errors.New("题目不存在")
		}
	}
	return p, nil
}

// GetsProblemList 获取题目列表
func (p Problem) GetsProblemList(page int, pageSize int, keyWord string) (interface{}, error) {
	type responseDTO struct {
		Result []Problem `json:"result"`
		Total  int64     `json:"total"`
	}
	var ret responseDTO
	result := common.DB.Model(&Problem{})

	if keyWord != "" {
		result.Where("title like ?", "%" + keyWord + "%")
	}
	result = result.Where("is_visible  = ?", true).Limit(pageSize).Offset((page - 1) * pageSize).Find(&ret.Result)
	ret.Total = result.RowsAffected
	fmt.Println(len(ret.Result))
	if result.Error != nil {
		return ret, result.Error
	} else {
		return ret, nil
	}
}

// GetContestProblemList 获取比赛问题列表
func (p Problem) GetContestProblemList(contestID uint) (interface{}, error)  {
	var contestProblemList []Problem
	result := common.DB.Model(Problem{}).Where("contest_id = ?", contestID).Find(&contestProblemList)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return contestProblemList, errors.New("记录未找到")
		} else {
			return contestProblemList, errors.New("内部错误")
		}
	}
	return contestProblemList, nil
}

// Update 更新问题
func (p Problem) Update() (Problem, error) {
	result := common.DB.Save(&p)
	if result.Error != nil {
		return p, result.Error
	}
	return p, nil
}
