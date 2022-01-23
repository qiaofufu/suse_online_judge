package judge

import (
	"onlineJudge/models"
	"time"
)

// Judge 开启judge服务
func Judge() {
	for ; ; {
		if (models.JudgeMessage{}.Empty()) {
			time.Sleep(time.Duration(1) * time.Second)
		} else {
			judge, err := models.JudgeMessage{}.GetNextJudgeMessage()
			if err == nil {
				RunJudge(judge.SubmissionID)
			}
		}
	}
}

func RunJudge(sid uint) {
	submission, err := models.Submission{}.GetSubmission("id = ?", sid)
	if err != nil {
		return
	}
	status, memory, runTime := judgeProblem(submission.ProblemID, submission.Code, submission.Language)
	submission.Status = status
	submission.Memory = memory
	submission.RunTime = runTime


	problem, err := models.Problem{}.GetProblem(submission.ProblemID)
	problem.SubmissionNumber ++
	if status == models.Accept {
		problem.AcceptedNumber ++
	}
	problem.Update()

	models.ProblemStatusRecord{}.Update(submission.ProblemID, submission.UserID, status)
	if submission.IsContestSubmission == true {
		cnt, _ := models.Submission{}.GetSubmissionCount("user_id  = ? and problem_id = ? and status = ? and contest_id = ?", submission.UserID, submission.ProblemID, models.Accept, submission.ContestID)
		if cnt == 0 {
			record, _ := models.ContestRecord{}.GetContestRecord("user_id = ? and contest_id = ?", submission.UserID, submission.ContestID)
			record.SubmissionNumber++
			if status == models.Accept {
				record.AcceptedNumber++
				record.TotalTime += runTime
				models.UserIntegral{}.AddIntegral(submission.UserID, 1)
			} else {
				record.TotalTime += int64(0.25 * float64(runTime))
			}
			record.UpdateContestRecord()
		}
	}
	submission.Update()
}
