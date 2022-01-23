package judge

import (
	"onlineJudge/util"
	"strconv"
)

// judgeProblem 判题
func judgeProblem(pid uint, code string, language string) (status int, memory int64, runTime int64){
	status, _ = strconv.Atoi(util.RandomCode(1))
	return status, memory, runTime
}