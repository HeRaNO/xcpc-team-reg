package main

import (
	"context"
	"flag"
	"html/template"
	"sync"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/rdb"
	"github.com/HeRaNO/xcpc-team-reg/internal/email"
	"github.com/HeRaNO/xcpc-team-reg/pkg/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"golang.org/x/time/rate"
)

type TeamAcc struct {
	TeamAccount  string
	TeamPassword string
}

const maxEmailPerSecond = 15
const maxEmailWorkers = 100

func main() {
	initConfigFile := flag.String("c", "./configs/conf.yaml", "the path of configure file")
	subject := flag.String("s", "【通知】校赛初赛相关信息", "the subject of account email")
	flag.Parse()
	internal.InitConfig(initConfigFile)

	tmpl, err := template.ParseFiles("./configs/email-account.tmpl")
	if err != nil {
		hlog.Fatalf("cannot parse email template file, err: %s", err.Error())
	}

	contestName := contest.ContestInfo().Name

	usrs, err := rdb.GetAllUsersInTeam()
	if err != nil {
		hlog.Fatalf("cannot get all users, err: %s", err.Error())
	}
	teams, err := rdb.GetAllTeams()
	if err != nil {
		hlog.Fatalf("cannot get all teams, err: %s", err.Error())
	}

	teamMap := make(map[int64]TeamAcc)
	for _, team := range teams {
		teamMap[team.TeamID] = TeamAcc{
			TeamAccount:  team.TeamAccount,
			TeamPassword: team.TeamPassword,
		}
	}

	failed := sync.Map{}
	wg := sync.WaitGroup{}
	limiter := rate.NewLimiter(rate.Limit(maxEmailPerSecond), 1)
	jobs := make(chan model.User)
	ctx := context.Background()

	workerCnt := min(maxEmailWorkers, len(usrs))
	for range workerCnt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for usr := range jobs {
				if err := limiter.Wait(ctx); err != nil {
					hlog.Errorf("wait rate limiter failed, usr: %+v, err: %s", usr, err.Error())
					failed.Store(usr.UserID, true)
					continue
				}

				acc, ok := teamMap[usr.BelongTeam]
				if !ok {
					hlog.Errorf("send email error, usr: %+v, cannot find belong_team", usr)
					failed.Store(usr.UserID, true)
					continue
				}
				err := email.SendTeamAccountEmail(tmpl, &usr.Name, &contestName, &acc.TeamAccount, &acc.TeamPassword, &usr.Email, subject)
				if err != nil {
					hlog.Errorf("send email error, usr: %+v, err: %s", usr, err.Error())
					failed.Store(usr.UserID, true)
				} else {
					hlog.Infof("send email ok, user_id: %d", usr.UserID)
				}
			}
		}()
	}

	for _, usr := range usrs {
		jobs <- usr
	}
	close(jobs)
	wg.Wait()

	failedIDs := make([]int64, 0)
	failed.Range(func(key, value any) bool {
		failedIDs = append(failedIDs, key.(int64))
		return true
	})

	hlog.Infof("send email finished. failed id: %+v", failedIDs)
}
