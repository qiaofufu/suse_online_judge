package models

//type ContestProblemRecord struct {
//	gorm.Model
//	ProblemID uint `json:"problem_id"`
//	ContestID uint `json:"contest_id"`
//	CreateID uint
//
//}
//
//// Update 更新比赛问题提交记录
//func (p ContestProblemRecord) Update(userID uint, problemID uint, contestID uint,isAC bool, isFirstAC bool, addTime int64) (ContestProblemRecord, error) {
//	var record ContestProblemRecord
//	var contestRecord ContestRecord
//	contestRecord, err := contestRecord.GetRecord(userID, contestID)
//	if err != nil {
//		return record, err
//	}
//	if !isAC {	// 没AC 就不是第一个AC的
//		isFirstAC = false
//	}
//	result := common.DB.Model(&ContestProblemRecord{}).Where("user_id = ? and problem_id = ?", userID, problemID).First(&record)
//	if result.Error == nil {
//		// 有记录
//		record.IsFirst = isFirstAC
//		if record.IsAC == false { 	// 本题未AC
//			if isAC == true {
//				record.TotalTime += addTime
//				contestRecord.TotalTime += addTime
//				record.IsAC = true
//				contestRecord.AcceptedNumber++
//			} else {
//				record.TotalTime += int64(0.5 * float64(addTime))
//				contestRecord.TotalTime += int64(0.5 * float64(addTime))
//				record.ErrorNumber ++
//				contestRecord.WrongNumber ++
//			}
//		}
//		result := common.DB.Save(&record)
//		if result.Error != nil {
//			return record, errors.New("更新失败")
//		}
//	} else if result.Error == gorm.ErrRecordNotFound {
//		// 没有记录
//		record = ContestProblemRecord{UserID: userID, ProblemID: problemID, ContestID: contestID, IsAC: isAC, IsFirst: isFirstAC}
//		if !isAC {
//			record.ErrorNumber ++
//			contestRecord.WrongNumber ++
//			record.TotalTime += int64(0.5 * float64(addTime))
//			contestRecord.TotalTime += int64(0.5 * float64(addTime))
//		} else {
//			contestRecord.AcceptedNumber ++
//			contestRecord.TotalTime += addTime
//		}
//		result := common.DB.Create(&record)
//		if result.Error != nil {
//			return record, errors.New("创建失败")
//		}
//
//	} else {
//		// 其他错误
//		return record, errors.New("内部错误")
//	}
//	contestRecord.Update(contestRecord)
//	return record, nil
//}
//
//// GetIsAC 获取是否AC
//func (t ContestProblemRecord) GetIsAC(userID uint, problemID uint) (bool, error) {
//	var record ContestProblemRecord
//	result := common.DB.Model(&ContestProblemRecord{}).Where("user_id = ? and problem_id = ?", userID, problemID).First(&record)
//	if result.Error != nil {
//		if result.Error != gorm.ErrRecordNotFound {
//			return false, errors.New("内部错误1")
//		} else {
//			return false, errors.New("题目不存在")
//		}
//	}
//	return record.IsAC, nil
//}
//
//// IsExistsACSubmissionByProblemID 根据问题id检验是否存在ac提交
//func (p ContestProblemRecord) IsExistsACSubmissionByProblemID(problemID uint) bool {
//	var cnt int64
//	common.DB.Model(&ContestProblemRecord{}).Where("problem_id = ? and is_ac = ?", problemID, Accept).Count(&cnt)
//	fmt.Println(cnt)
//	if cnt == 0 {
//		return false
//	} else {
//		return true
//	}
//}
