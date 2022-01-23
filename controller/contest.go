package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/models"
	"onlineJudge/response"
)

// GetsContestList 获取比赛列表
func GetsContestList(ctx *gin.Context) {
	type requiredDTO struct {
		KeyWord  string `json:"key_word,omitempty"`
		Status   string `json:"status,omitempty"`
		Page     int    `json:"page,omitempty"`
		PageSize int    `json:"page_size,omitempty"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	contestList, err := models.Contest{}.GetsContestList(dto.KeyWord, dto.Status, dto.Page, dto.PageSize)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
	} else {
		response.Success(ctx, contestList, "")
	}
}

// GetContest 获取比赛信息
func GetContest(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		ID       uint   `json:"id,omitempty"`
		Password string `json:"password,omitempty"`
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
	// 处理
	contest, err := models.Contest{}.GetContest(user.ID, dto.ID, dto.Password)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
	} else {
		response.Success(ctx, contest, "")
	}
}

// GetContestProblemList 获取比赛题目列表
func GetContestProblemList(ctx *gin.Context) {
	// 数据绑定
	type requiredDTO struct {
		ContestID uint `json:"contest_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil ,"数据绑定失败")
		return
	}
	user, err := getUserByJWT(ctx)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	// 操作
	if (models.Contest{}.AccessVerification(user.ID, dto.ContestID)) {
		contestProblemList, err := models.Problem{}.GetContestProblemList(dto.ContestID)
		if err != nil {
			response.Fail(ctx, nil, err.Error())
		} else {
			response.Success(ctx, contestProblemList, "")
		}
	} else {
		response.Fail(ctx, nil, "不是比赛参与者，禁止访问")
	}

}



// GetContestRank 获取比赛排行
func GetContestRank(ctx *gin.Context) {
	type requiredDTO struct {
		ContestID uint `json:"contest_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	//
	type submissionInfo struct {
		ProblemId uint
		IsAC bool
		WrongNumber int
		AcTime int64
		IsFirstAC bool
	}
	type responseDTO struct {
		ContestRecord  models.ContestRecord
		SubmissionInfo []submissionInfo  `json:"submission_info"`
	}
	var ret []responseDTO
	contestRecordList, err := models.ContestRecord{}.GetContestRecordListByContestID(dto.ContestID)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	for _, value := range contestRecordList {
		var elem responseDTO
		elem.ContestRecord = value
		list, err := models.Submission{}.GetsSubmission("is_contest_submission = ? and contest_id = ? and user_id = ? ", true, value.ContestID, value.UserID)
		if err != nil {
			response.Fail(ctx, nil, err.Error())
			return
		}
		for _, v := range list {	// 遍历提交列表
			isAc := false
			wrongNumber := 1
			var acTime int64
			acTime = 0
			if v.Status == Accept {
				isAc = true
				wrongNumber = 0
				contest, _ := models.Contest{}.AdminGetContest(dto.ContestID)
				acTime = v.CreatedAt.Unix() - contest.StartTime.Unix()
			}
			isFirstAC := models.Submission{}.IsFirstAC(v.ID)
			flag := true
			for i, vv := range elem.SubmissionInfo {	// 遍历提交记录
				if v.ProblemID == vv.ProblemId {	// 存在记录
					flag = false
					if vv.IsAC == false  {
						elem.SubmissionInfo[i].WrongNumber += wrongNumber
						elem.SubmissionInfo[i].IsAC = isAc
						elem.SubmissionInfo[i].AcTime = acTime
						elem.SubmissionInfo[i].IsFirstAC = isFirstAC
					} else {
						continue
					}
				}
			}
			// 不存在记录
			if flag {
				elem.SubmissionInfo = append(elem.SubmissionInfo, submissionInfo{ProblemId: v.ProblemID, IsAC: isAc, WrongNumber: wrongNumber, AcTime: acTime, IsFirstAC: isFirstAC})
			}
		}
		ret = append(ret, elem)
	}
	response.Success(ctx, ret, "")
}