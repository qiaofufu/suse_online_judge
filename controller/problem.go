package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/models"
	"onlineJudge/response"
)

// GetProblemContent 获取普通题目内容
func GetProblemContent(ctx *gin.Context) {
	type requiredDTO struct {
		ProblemID uint `json:"problem_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil ,"数据绑定失败")
		return
	}
	ret, err := models.Problem{}.GetProblem(dto.ProblemID)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
	} else {
		response.Success(ctx, ret, "")
	}
}

// GetContestProblemContent 获取比赛题目内容
func GetContestProblemContent(ctx *gin.Context) {
	type requiredDTO struct {
		ProblemID uint `json:"problem_id,omitempty"`
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
	ret, err := models.Problem{}.GetProblem(dto.ProblemID)
	if (models.Contest{}.AccessVerification(user.ID, ret.ContestID)) {
		if err != nil {
			response.Fail(ctx, nil, err.Error())
		} else {
			response.Success(ctx, ret, "")
		}
	}

}

// GetsProblemList 获取问题列表
func GetsProblemList(ctx *gin.Context) {
	type requiredDTO struct {
		KeyWord  string `json:"key_word"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	ret, err := models.Problem{}.GetsProblemList(dto.Page, dto.PageSize, dto.KeyWord)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
	} else {
		response.Success(ctx, ret, "")
	}
}


