package contest

import (
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	name, note               string
	startTime, endTime       time.Time
	startTimeStr, endTimeStr string
	idSchoolMap              map[int]string
	validStuIDLength         map[int]bool
	validTshirtSize          map[string]bool
	tshirtSize               []string
)

func Init(conf *config.ContestConfig) {
	if conf == nil {
		hlog.Fatal("Contest config failed: conf is nil")
		panic("make static check happy")
	}
	name, note = conf.Name, conf.Note
	startTimeStr, endTimeStr = conf.StartTime, conf.EndTime
	loc, err := time.LoadLocation("Local")
	if err != nil {
		hlog.Fatalf("Contest config failed: cannot load local time config, err: %+v", err)
	}
	st, err := time.ParseInLocation("2006-01-02 15:04:05", conf.StartTime, loc)
	if err != nil {
		hlog.Fatalf("Contest config failed: cannot parse start time, err: %+v", err)
	}
	ed, err := time.ParseInLocation("2006-01-02 15:04:05", conf.EndTime, loc)
	if err != nil {
		hlog.Fatalf("Contest config failed: cannot parse end time, err: %+v", err)
	}
	if !st.Before(ed) {
		hlog.Fatal("Contest config failed: contest start time is not before contest end time")
	}
	startTime, endTime = st, ed

	idSchoolMap = make(map[int]string)
	validStuIDLength = make(map[int]bool)
	validTshirtSize = make(map[string]bool)
	tshirtSize = conf.ValidTshirtSize
	for i, schoolName := range conf.SchoolName {
		idSchoolMap[i+1] = schoolName
	}
	for _, stuIDLength := range ValidStuIDLength {
		validStuIDLength[stuIDLength] = true
	}
	for _, size := range conf.ValidTshirtSize {
		validTshirtSize[size] = true
	}
	hlog.Info("init contest finished successfully")
}
