package models

import (
	"errors"
	"gorm.io/gorm"
	"onlineJudge/common"
)

type ProblemStatusRecord struct {
	gorm.Model
	ProblemID uint `json:"problem_id"`
	UserID uint `json:"user_id"`
	Status int `json:"status"`
}

// Update 更新问题状态记录
func (p ProblemStatusRecord) Update(problemID uint, userID uint, status int)  (ProblemStatusRecord, error) {
	result := common.DB.Model(&ProblemStatusRecord{}).Where("problem_id = ? and user_id = ?", problemID, userID).First(&p)
	if result.Error == nil {
		// 有记录
		if status == Accept && p.Status != Accept {
			p.Status = Accept
		}
		result := common.DB.Save(&p)
		if result.Error != nil {
			return p, errors.New("更新失败")
		}
	} else if result.Error == gorm.ErrRecordNotFound {
		// 没有记录
		result := common.DB.Create(&ProblemStatusRecord{ProblemID: problemID, UserID: userID, Status: status})
		if result.Error != nil {
			return p, errors.New("创建问题提交记录失败")
		}
	} else {
		// 其他错误
		return p, errors.New("内部错误1")
	}
	return p, nil
}

// GetBy 根据字段获取一条记录
func (p ProblemStatusRecord) GetBy(query string, value ... interface{}) (ProblemStatusRecord, error) {
	result := common.DB.Model(&ProblemStatusRecord{}).Where(query, value).First(&p)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return p, errors.New("记录未找到")
		}
		return p, result.Error
	}
	return p, nil
}
