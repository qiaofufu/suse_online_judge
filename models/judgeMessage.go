package models

import (
	"gorm.io/gorm"
	"onlineJudge/common"
)

type JudgeMessage struct {
	gorm.Model
	SubmissionID uint
}

// Add 添加判题消息
func (m JudgeMessage) Add(submissionID uint) error {
	m.SubmissionID = submissionID
	result := common.DB.Create(&m)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Empty 判断判题消息是否为空
func (m JudgeMessage) Empty() bool {
	var cnt int64
	common.DB.Model(&JudgeMessage{}).Count(&cnt)
	if cnt == 0 {
		return true
	} else {
		return false
	}
}

// GetNextJudgeMessage 获取下一条judge消息
func (m JudgeMessage) GetNextJudgeMessage() (JudgeMessage, error) {
	result := common.DB.Model(&JudgeMessage{}).First(&m)
	if result.Error != nil {
		return m, result.Error
	} else {
		common.DB.Model(&JudgeMessage{}).Delete(&m)
		return m, nil
	}
}