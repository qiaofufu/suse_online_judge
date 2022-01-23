package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/models"
	"onlineJudge/response"
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

// SubmissionCode 提交代码
func SubmissionCode(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		Code string `json:"code"`
		Language string `json:"language"`
		ProblemID uint `json:"problem_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 创建提交
	submission, err := models.Submission{Code: dto.Code, UserID: user.ID, Language: dto.Language, ProblemID: dto.ProblemID, Author: user.UserName, Status: Pending}.Insert()
	if err != nil {
		response.Fail(ctx, nil, err.Error())
	} else {
		type responseDTO struct {
			SubmissionID uint `json:"submission_id"`
		}
		var res responseDTO
		res.SubmissionID = submission.ID
		models.JudgeMessage{}.Add(submission.ID)
		response.Success(ctx, res, "提交成功")
	}
}

// SubmissionExists 用户提交是否存在
func SubmissionExists(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		ProblemID uint `json:"problem_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 操作
	submission , err := models.Submission{}.GetSubmission("user_id = ? and problem_id = ? and status = ?", user.ID, dto.ProblemID, Accept)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
	} else {
		var res int
		if submission.ID == 0 {
			submission, err = models.Submission{}.GetSubmission("user_id = ? and problem_id", user.ID, dto.ProblemID)
			if submission.ID != 0 {
				res = -1
			}
		} else {
			res = 1
		}
		response.Success(ctx, res, "检验成功")
	}
}

// SubmissionList 提交列表
func SubmissionList(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		Myself int `json:"myself"`
		Result int `json:"result"`
		UserName string `json:"user_name"`
		ProblemID uint `json:"problem_id"`
		Page int `json:"page"`
		PageSize int `json:"page_size"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	var user models.User
	if dto.Myself != 0 {
		user, _ = getUserByJWT(ctx)
	}// 操作
	ret, err := models.Submission{}.GetsSubmissionRecord(user.ID, dto.Result, dto.UserName, dto.Page, dto.PageSize, 0, 0, dto.ProblemID)
	if err != nil {
		response.Fail(ctx, nil ,"查询失败")
	} else {
		response.Success(ctx, ret, "查询成功")
	}
}

// GetsContestSubmission 获取比赛提交
func GetsContestSubmission(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		Myself    int    `json:"myself"`
		Result    int    `json:"result"`
		UserName  string `json:"user_name"`
		ProblemID uint `json:"problem"`
		Page      int    `json:"page"`
		PageSize  int    `json:"page_size"`
		ContestID uint   `json:"contest_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	user, _ := getUserByJWT(ctx)
	if dto.Myself == 0 {
		user.ID = 0
	}
	// 操作
	ret, err := models.Submission{}.GetsSubmissionRecord(user.ID, dto.Result, dto.UserName, dto.Page, dto.PageSize, 1, dto.ContestID, dto.ProblemID)
	if err != nil {
		response.Fail(ctx, nil ,"查询失败")
	} else {
		response.Success(ctx, ret, "查询成功")
	}
}
