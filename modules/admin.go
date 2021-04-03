package modules

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
)

func ImportTeamInfo(w http.ResponseWriter, r *http.Request) {
	// read the file
	// team.team_account <- account, team.team_password <- password

	if err := r.ParseMultipartForm(config.MaxUploadSize); err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	file, fileHeader, err := r.FormFile("team_table")
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	if fileHeader.Size > config.MaxUploadSize {
		util.ErrorResponse(w, r, "file too large", config.ERR_INTERNAL)
		return
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	fileReader := bytes.NewReader(fileBytes)
	csvReader := csv.NewReader(fileReader)
	rowCnt := 0

	teamAccPwd := map[int64][]string{}
	failedID := make([]int64, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		rowCnt++
		if len(row) != 3 {
			errmsg := fmt.Sprintf("On line #%d: item count less than 3", rowCnt)
			util.ErrorResponse(w, r, errmsg, config.ERR_WRONGINFO)
			return
		}

		teamID, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			errmsg := fmt.Sprintf("On line #%d: the first item is not team_id", rowCnt)
			util.ErrorResponse(w, r, errmsg, config.ERR_WRONGINFO)
			return
		}

		if _, ok := teamAccPwd[teamID]; ok {
			errmsg := fmt.Sprintf("On line #%d: there's a team_id same as this row", rowCnt)
			util.ErrorResponse(w, r, errmsg, config.ERR_WRONGINFO)
			return
		}

		accPwd := make([]string, 0, 2)
		accPwd = append(accPwd, row[1])
		accPwd = append(accPwd, row[2])
		teamAccPwd[teamID] = accPwd
	}

	for teamID, accPwd := range teamAccPwd {
		err := model.SetTeamAccPwdByID(r.Context(), teamID, &accPwd[0], &accPwd[1])
		if err != nil {
			failedID = append(failedID, teamID)
		}
	}

	util.SuccessResponseWithTotal(w, r, failedID, len(failedID))
}

func ExportTeamInfo(w http.ResponseWriter, r *http.Request) {
	// read database
	// data -> csv, download

	csvContent := make([]byte, 0)
	buf := bytes.NewBuffer(csvContent)
	csvWriter := csv.NewWriter(buf)

	allTeamID, err := model.GetAllTeamIDs(r.Context())

	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	var teamTable [][]string

	for _, id := range allTeamID {
		allTeamInfo, err := model.GetTeamInfoByTeamID(r.Context(), id)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
			return
		}
		teamInfo := make([]string, 0)
		teamInfo = append(teamInfo, fmt.Sprintf("%d", allTeamInfo.TeamID))
		teamInfo = append(teamInfo, allTeamInfo.TeamName)
		for _, member := range allTeamInfo.TeamMember {
			teamInfo = append(teamInfo, member.Name)
			teamInfo = append(teamInfo, member.School)
			teamInfo = append(teamInfo, member.StuID)
		}

		teamTable = append(teamTable, teamInfo)
	}

	err = csvWriter.WriteAll(teamTable)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	csvWriter.Flush()
	util.File(w, r, buf.Bytes())
}
