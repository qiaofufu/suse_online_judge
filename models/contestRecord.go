package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"onlineJudge/common"
)

type ContestRecord struct {
	gorm.Model
	UserID uint `json:"user_id"`
	ContestID uint `json:"contest_id"`
	AcceptedNumber int64 `json:"accepted_number" gorm:"default:0"`
	SubmissionNumber int64
	TotalTime int64 `json:"total_time"`
}

//// Update 更新记录
//func (t ContestRecord) Update(record ContestRecord) (ContestRecord, error) {
//	result := common.DB.Model(&ContestRecord{}).Where("id = ?", record.ID).Save(record)
//	if result.Error != nil {
//		return record, errors.New("更新失败")
//	} else {
//		return record, nil
//	}
//}

// GetContestRecordListByContestID 根据比赛ID获取比赛记录列表
func (t ContestRecord) GetContestRecordListByContestID(contestID uint) ([]ContestRecord, error) {
	var rankRecord []ContestRecord
	result := common.DB.Model(&ContestRecord{}).Where("contest_id = ?", contestID).Order("accepted_number desc").Order("total_time desc").Find(&rankRecord)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("比赛记录未找到")
		} else {
			return nil, errors.New("其他错误")
		}
	}
	return rankRecord, nil
}

// InsertContestRecord 插入比赛记录
func (t ContestRecord) InsertContestRecord() (ContestRecord,error) {

	result := common.DB.Create(&t)
	if result.Error != nil {
		return t, result.Error
	}
	return t, nil
}

// UpdateContestRecord 更新比赛记录
func (t ContestRecord) UpdateContestRecord() (ContestRecord, error) {
	result := common.DB.Save(&t)
	if result.Error != nil {
		return t, result.Error
	}
	return t, nil
}

// IsExist 检验是否存在记录
func (t ContestRecord) IsExist(query string, value ... interface{}) bool {
	var total int64
	result := common.DB.Model(&ContestRecord{}).Where(query, value...).Count(&total)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return total != 0
}

// GetContestRecord 获取比赛记录
func (t ContestRecord) GetContestRecord(query string, value ... interface{}) (ContestRecord, error) {
	result := common.DB.Model(&ContestRecord{}).Where(query, value...).First(&t)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return t, errors.New("记录未找到")
		} else {
			return t, errors.New("内部错误")
		}
	}
	return t, nil
}
