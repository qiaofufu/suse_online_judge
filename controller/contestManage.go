package controller

import (
	"github.com/gin-gonic/gin"
	"onlineJudge/common"
	"onlineJudge/models"
	"onlineJudge/response"
	"onlineJudge/util"
)

// AddContest 发布竞赛
func AddContest(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Password string `json:"password"`
		StartTime string `json:"start_time"`
		EndTime string `json:"end_time"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	startTime, err := util.GetTime(dto.StartTime)
	if err != nil {
		response.Fail(ctx, nil ,err.Error())
		return
	}
	endTime, err := util.GetTime(dto.EndTime)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	admin, err := getAdminJwt(ctx)
	if err != nil {
		response.Fail(ctx, nil,"获取管理员信息失败")
		return
	}
	// 数据库操作
	contest := models.Contest{Title: dto.Title, Description: dto.Description, Password: dto.Password, StartTime: startTime, EndTime: endTime,ContestCreatorID: admin.ID}
	result := common.DB.Model(&models.Contest{}).Create(&contest)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "创建比赛失败")
	} else {
		response.Success(ctx, contest, "创建比赛成功")
	}
}

// AdminGetsContestList 获取竞赛列表
func AdminGetsContestList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据库操作
	var contestList []models.Contest
	result := common.DB.Model(&models.Contest{}).Order("id desc").Find(&contestList)
	if result.Error != nil {
		response.Fail(ctx, nil, "获取数据错误")
	} else {
		response.Success(ctx, contestList, "获取成功")
	}
}

// UpdateContest 更新竞赛
func UpdateContest(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID uint `json:"alter_id"`
		Password string `json:"password"`
		StartTime string `json:"start_time"`
		EndTime string `json:"end_time"`
		Status int `json:"status"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	startTime, err := util.GetTime(dto.StartTime)
	if err != nil {
		response.Fail(ctx, nil, "数据绑定失败2")
		return
	}
	endTime, err := util.GetTime(dto.EndTime)
	if err != nil {
		response.Fail(ctx, nil, "数据绑定失败3")
		return
	}
	// 数据库操作
	result := common.DB.Model(&models.Contest{}).Where("id = ?", dto.AlterID).Updates(models.Contest{StartTime: startTime, EndTime: endTime, Status: dto.Status, Password: dto.Password})
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// DeleteContest 删除比赛
func DeleteContest(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		DeleteID uint `json:"delete_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Contest{})
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// ADDContestProblem 添加比赛问题
func ADDContestProblem(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		Title             string `json:"title"`
		Description       string `json:"description"`
		ContestID         uint   `json:"contest_id"`
		IsPublic          bool   `json:"is_public"`
		InputDescription  string `json:"input_description"`
		OutPutDescription string `json:"out_put_description"`
		Sample            string `json:"sample"`
		TimeLimit         int64  `json:"time_limit"`
		MemoryLimit       int64  `json:"memory_limit"`
		IsVisible         bool   `json:"is_visible"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	problem := models.Problem{Title: dto.Title, Description: dto.Description, ContestID: dto.ContestID, IsPublic: dto.IsPublic, InputDescription: dto.InputDescription,
		OutputDescription: dto.OutPutDescription, Samples: dto.Sample, TimeLimit: dto.TimeLimit, MemoryLimit: dto.MemoryLimit, IsVisible: dto.IsVisible}
	result := common.DB.Create(&problem)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "添加失败")
	} else {
		response.Success(ctx, nil, "添加成功")
	}
}

// DeleteContestProblem 删除比赛问题
func DeleteContestProblem(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		DeleteID uint `json:"delete_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Problem{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// UpdateContestProblem 更新比赛问题
func UpdateContestProblem(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID 		  uint   `json:"alter_id"`
		Title             string `json:"title"`
		Description       string `json:"description"`
		ContestID         uint   `json:"contest_id"`
		IsPublic          bool   `json:"is_public"`
		InputDescription  string `json:"input_description"`
		OutPutDescription string `json:"out_put_description"`
		Sample            string `json:"sample"`
		TimeLimit         int64  `json:"time_limit"`
		MemoryLimit       int64  `json:"memory_limit"`
		IsVisible         bool   `json:"is_visible"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	problem, _ := models.Problem{}.GetProblem(dto.AlterID)
	result := common.DB.Model(&models.Problem{}).Where("id = ?", problem.ID).Updates(models.Problem{Title: dto.Title, Description: dto.Description, ContestID: dto.ContestID, IsPublic: dto.IsPublic, InputDescription: dto.InputDescription,
		OutputDescription: dto.OutPutDescription, Samples: dto.Sample, TimeLimit: dto.TimeLimit, MemoryLimit: dto.MemoryLimit, IsVisible: dto.IsVisible})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "更新失败")
	} else {
		response.Success(ctx, nil, "更新成功")
	}
}

// GetsContestProblemList 获取比赛问题列表
func GetsContestProblemList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		ContestID uint `json:"contest_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	var problemList []models.Problem
	result := common.DB.Model(&models.Problem{}).Where("contest_id = ?", dto.ContestID).Find(&problemList)
	if result.Error != nil {
		response.Fail(ctx, nil, "获取失败")
	} else {
		response.Success(ctx, problemList, "获取成功")
	}
}

// AddExample 添加样例点
func AddExample(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		ProblemID uint `json:"problem_id"`
		Input string `json:"input"`
		Output string `json:"output"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	example := models.Example{Input: dto.Input, Output: dto.Output, ProblemID: dto.ProblemID}
	result := common.DB.Model(&models.Example{}).Create(&example)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "Add失败")
	} else {
		response.Success(ctx, nil, "Add成功")
	}
}

// UpdateExample 更新样例点
func UpdateExample(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		AlterID uint `json:"alter_id"`
		Input string `json:"input"`
		Output string `json:"output"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}

	// 数据库操作
	example := models.Example{Output: dto.Output, Input: dto.Input}
	result := common.DB.Where("id = ?", dto.AlterID).Updates(example)
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "Update失败")
	} else {
		response.Success(ctx, nil, "Update成功")
	}
}

// DeleteExample 删除样例点
func DeleteExample(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		DeleteID uint `json:"delete_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	result := common.DB.Where("id = ?", dto.DeleteID).Delete(&models.Example{})
	// 响应
	if result.Error != nil {
		response.Fail(ctx, nil, "删除失败")
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// GetsExampleList 获取样例点列表
func GetsExampleList(ctx *gin.Context) {
	// 权限验证
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		return
	}
	// 数据绑定
	type requiredDTO struct {
		ProblemID uint `json:"problem_id"`
	}
	var dto requiredDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
		return
	}
	// 数据库操作
	var exampleList []models.Example
	result := common.DB.Model(&models.Example{}).Where("problem_id = ?", dto.ProblemID).Order("id desc").Find(&exampleList)
	if result.Error != nil {
		response.Fail(ctx, nil, "获取数据错误")
	} else {
		response.Success(ctx, exampleList, "获取成功")
	}
}

// ReleaseAnswer 发布题解
func ReleaseAnswer(ctx *gin.Context) {
	if !roleCheck(ctx, ROLE_CONTESTMANAGE) {
		response.Fail(ctx, nil, "没有操作权限")
		return
	}
	type responseDTO struct {
		Content   string `json:"content"`
		ContestID uint   `json:"contest_id"`
	}
	var dto responseDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response.Fail(ctx, nil, "数据绑定失败")
	}

	contest, err := models.Contest{}.AdminGetContest(dto.ContestID)
	if err != nil {
		response.Fail(ctx, nil, "获取信息失败")
		return
	}
	contest.Answer = dto.Content
	contest.UpdateContest()
	response.Success(ctx, nil, "发布成功")
}