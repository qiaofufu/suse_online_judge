package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"onlineJudge/common"
	"strconv"
	"time"
)

const (
	ContestInvisible = -1
	ContestUnderway = 0
	ContestNotStart = 1
	ContestEnded = 2
)

type Contest struct {
	gorm.Model
	Title string
	StartTime time.Time
	EndTime time.Time
	Password string	`gorm:"default = ''"`
	Status int `gorm:"default=-1"`
	Description string
	ContestCreatorID uint
	Answer string
}

// IsRunning 比赛是否进行中
func (c Contest) IsRunning(contestID uint) bool {
	var elem Contest
	result := common.DB.Model(&Contest{}).Where("id = ?", contestID).First(&elem)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return time.Now().Unix() < elem.EndTime.Unix() && time.Now().Unix() > elem.StartTime.Unix()
}

// GetsContestList 获取比赛列表
func (c Contest) GetsContestList(keyword string, status string, page int, pageSize int) (interface{}, error) {
	result := common.DB.Model(&Contest{}).Where("status != ?", ContestInvisible)
	if keyword != "" {
		result = result.Where("title like ?","%" + keyword + "%")
	}
	if status != "" {
		cur := time.Now().Format("2006-01-02 15:04:05")
		s , _ := strconv.Atoi(status)
		if s == ContestNotStart {
			result = result.Where("start_time > ?", cur)
		} else if s == ContestEnded {
			result = result.Where("end_time)< ?", cur)
		} else if s == ContestUnderway {
			result = result.Where("start_time < ? and end_time > ?", cur, cur)
		} else {
			return nil, errors.New("status 参数错误")
		}
	}
	var contestList []Contest
	result.Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc").Find(&contestList)
	type responseDTO struct {
		ID 				 uint 	   `json:"id"`
		Title            string    `json:"title"`
		StartTime        time.Time `json:"start_time"`
		EndTime          time.Time `json:"end_time"`
		HavePassword     bool      `json:"have_password"`
		Status           int       `json:"status"`
		Description      string    `json:"description"`
		ContestCreatorName string  `json:"contest_creator_name"`
	}
	var res []responseDTO
	for _, v := range contestList {
		name := Admin{}.getName(v.ContestCreatorID)
		elem := responseDTO{ID: v.ID, Title: v.Title, StartTime: v.StartTime, EndTime: v.EndTime, Status: v.Status,Description: v.Description,ContestCreatorName: name}
		if v.Password != "" {
			elem.HavePassword = true
		} else {
			elem.HavePassword = false
		}
		res = append(res, elem)
	}
	return res, nil
}

// GetContest 获取比赛
func (c Contest) GetContest(userID uint, contestId uint, password string) (Contest, error) {
	result := common.DB.Model(&Contest{}).Where("id = ?", contestId).First(&c)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return Contest{}, errors.New("记录未找到")
		} else {
			return Contest{}, errors.New("内部错误")
		}
	}

	if (Contest{}.AccessVerification(userID, contestId)) {
		return c, nil
	} else {
		if c.Password != "" {
			if password != c.Password {
				return c, errors.New("密码错误")
			}
		}
		_, err := ContestRecord{UserID: userID, ContestID: contestId}.InsertContestRecord()
		if err != nil {
			return c, errors.New("内部错误")
		}
		return c, nil
	}
}

// AccessVerification 比赛访问权限验证 true -- 通过 false -- 未通过
func (c Contest) AccessVerification(userID uint, contestID uint) bool {
	if (ContestRecord{}.IsExist("user_id = ? and contest_id = ?", userID, contestID)) {
		return true
	}
	return false
}

// InsertContest 新建比赛
func (c Contest) InsertContest() (Contest, error) {
	result := common.DB.Create(&c)
	if result.Error != nil {
		return c, result.Error
	}
	return c, nil
}

// UpdateContest 更新比赛
func (c Contest) UpdateContest() (Contest, error) {
	result := common.DB.Save(&c)
	if result.Error != nil {
		return c, result.Error
	}
	return c, nil
}

// AdminGetContest 管理员获取比赛内容
func (c Contest) AdminGetContest(contestID uint) (Contest, error) {
	result := common.DB.Model(Contest{}).Where("id= ?", contestID).First(&c)
	return c, result.Error
}

