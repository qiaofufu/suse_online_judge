package models

import (
	"errors"
	"gorm.io/gorm"
	"onlineJudge/common"
	"time"
)

const (
	Pending = 0
	Accept = 1
	TimeLimitExceeded = 2
	MemoryLimitExceeded = 3
	RuntimeError = 4
	SystemError = 5
	CompileError = 6
	wrongAnswer = 7
)

type Submission struct {
	gorm.Model
	Status              int    `json:"status,omitempty"`
	ProblemID           uint   `json:"problem_id,omitempty"`
	RunTime             int64  `json:"run_time,omitempty"`
	Memory              int64  `json:"memory,omitempty"`
	Language            string `json:"language,omitempty"`
	UserID              uint   `json:"user_id,omitempty"`
	Author              string `json:"author,omitempty"`
	Code                string `json:"code,omitempty"`
	IsContestSubmission bool   `json:"is_contest_submission,omitempty"`
	ContestID 			uint   `json:"contest_id"`
}

// GetSubmission 根据字段获取记录一条
func (s Submission) GetSubmission(query string, value ... interface{}) (Submission, error) {
	result := common.DB.Model(&Submission{}).Where(query, value ...).First(&s)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return s, nil
		}
		return s, result.Error
	}
	return s, nil
}

// GetsSubmission 根据字段获取记录多条
func (s Submission) GetsSubmission(query string, value ... interface{}) ([]Submission, error) {
	var ret []Submission
	result := common.DB.Model(&Submission{}).Where(query, value ... ).Find(&ret)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ret, nil
		}
		return ret, result.Error
	}
	return ret, nil
}

// GetSubmissionCount 获取提交数量
func (s Submission) GetSubmissionCount(query string, value ... interface{}) (total int64, err error) {
	result := common.DB.Model(&Submission{}).Where(query, value ...).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

// Insert 创建记录
func (s Submission) Insert()  (Submission, error) {
	problem, err := Problem{}.GetProblem(s.ProblemID)
	if err != nil {
		return s, err
	}
	if (problem.ContestID != 0 && Contest{}.IsRunning(problem.ContestID)) {
		s.IsContestSubmission = true
		s.ContestID = problem.ContestID
	}
	if s.ContestID != 0 && (Contest{}.AccessVerification(s.UserID, s.ContestID) == false){
		return s, errors.New("不是比赛的参与者，不能提交")
	}
	result := common.DB.Model(&Submission{}).Create(&s)
	if result.Error != nil {
		return s, errors.New("创建失败")
	}
	return s, nil
}

// Update 更新记录
func (s Submission) Update() (Submission, error) {
	result := common.DB.Model(&Submission{}).Where("id = ?",s.ID).Save(&s)
	if result.Error != nil {
		return s, errors.New("更新失败")
	} else {
		return s, nil
	}
}

// GetsSubmissionRecord 获取所有提交记录
func (s Submission) GetsSubmissionRecord(userId uint, status int, username string, page int, pageSize int, isContestSubmission int, contestID uint, problemID uint) (interface{}, error) {
	result := common.DB.Model(&Submission{})
	if isContestSubmission == 1 {
		result = result.Where("contest_id = ?", contestID)
	}
	if userId != 0 {
		result = result.Where("user_id = ?", userId)
	}
	if status != 0 {
		result = result.Where("status = ?", status)
	}
	if username != "" {
		result = result.Where("author = ?", username)
	}
	if problemID != 0 {
		result = result.Where("problem_id = ?", problemID)
	}
	var total int64
	result.Count(&total)
	var submissionList []Submission
	result = result.Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc").Find(&submissionList)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		}
	}

	type DTO struct {
		CreateTime time.Time `json:"create_time"`
		SubmissionID uint `json:"submission_id"`
		Status int `json:"status"`
		RunTime int64 `json:"run_time"`
		ProblemID uint `json:"problem_id"`
		Memory int64 `json:"memory"`
		Language string `json:"language"`
		AuthorName string `json:"author_name"`
		AuthorID uint `json:"author_id"`
	}
	type responseDTO struct {
		Result []DTO `json:"result"`
		Total int64 `json:"total"`
	}
	var dto []DTO
	for _, v := range submissionList {
		elem := DTO{ProblemID: v.ProblemID, Status: v.Status, SubmissionID: v.ID, CreateTime: v.CreatedAt, RunTime: v.RunTime, Memory: v.Memory, Language: v.Language, AuthorName: v.Author, AuthorID: v.UserID}
		dto = append(dto, elem)
	}
	return responseDTO{Total: total, Result: dto}, nil
}

// IsFirstAC 检查是否第一次AC
func (s Submission) IsFirstAC(id uint) bool {
	common.DB.Model(&Submission{}).Where("id = ?", id).First(&s)
	if s.Status != Accept {
		return false
	}
	var cnt int64
	common.DB.Model(&Submission{}).Where("problem_id = ? and status = ? and created_at < ?", s.ProblemID, Accept, s.CreatedAt).Count(&cnt)
 	return cnt == 0
}