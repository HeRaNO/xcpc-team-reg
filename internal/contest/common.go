package contest

import (
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
)

func BeforeRegTime(now time.Time) bool {
	return now.Before(startTime)
}

func AfterRegTime(now time.Time) bool {
	return now.After(endTime)
}

func IsValidSchool(id int) bool {
	_, ok := idSchoolMap[id]
	return ok
}

func IsValidStuID(id *string) bool {
	return utils.IsNumber(id) && validStuIDLength[len(*id)]
}

func IsValidTshirtSize(size *string) bool {
	return validTshirtSize[*size]
}

func ContestInfo() model.ContestInfo {
	return model.ContestInfo{
		Name:       name,
		StartTime:  startTimeStr,
		EndTime:    endTimeStr,
		Note:       note,
		TshirtSize: tshirtSize,
	}
}

func GetIDSchoolMap() map[int]string {
	return idSchoolMap
}
