package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"html/template"
	"net/smtp"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/jordan-wright/email"
)

func makeTokenEmail(tmpl *template.Template, token *string, method *string) []byte {
	content := new(bytes.Buffer)
	tmpl.Execute(content, struct {
		Action string
		Time   string
		Token  string
		Sign   string
	}{
		Action: emailActionMap[*method],
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Token:  *token,
		Sign:   emailSign,
	})
	return content.Bytes()
}

func makeTeamAccountEmail(tmpl *template.Template, name *string, contestName *string, account *string, password *string) []byte {
	content := new(bytes.Buffer)
	tmpl.Execute(content, struct {
		Name         string
		ContestName  string
		TeamAccount  string
		TeamPassword string
		Sign         string
	}{
		Name:         *name,
		ContestName:  *contestName,
		TeamAccount:  *account,
		TeamPassword: *password,
		Sign:         emailSign,
	})
	return content.Bytes()
}

func sendEmail(emailRecv *string, subject *string, content []byte) error {
	e := &email.Email{
		To:      []string{*emailRecv},
		From:    emailFrom,
		Subject: *subject,
		HTML:    content,
	}

	auth := smtp.PlainAuth("", emailAddr, emailPassword, smtpAddr)
	return e.SendWithTLS(smtpHost, auth, &tls.Config{ServerName: smtpAddr})
}

func SendEmailWithToken(ctx context.Context, email *string, emailType *string) (bool, error) {
	if _, ok := emailActionMap[*emailType]; !ok {
		return true, errors.New("unrecognized action")
	}
	uid, err := redis.GetUserIDByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if *emailType == "register" {
		if uid != 0 {
			return true, errors.New("email has already registered")
		}
	} else {
		if uid == 0 {
			return true, errors.New("no such user")
		}
	}

	err = redis.GetEmailRequest(ctx, email)
	if err != nil {
		return false, err
	}

	token, err := utils.GenToken(contest.UserTokenLength)
	if err != nil {
		return false, err
	}

	err = redis.SetEmailToken(ctx, email, &token, emailTokenExpireTime)
	if err != nil {
		return false, err
	}
	err = redis.SetEmailAction(ctx, email, emailType, emailTokenExpireTime)
	if err != nil {
		return false, err
	}
	err = redis.SetEmailRequest(ctx, email, emailSendGapTime)
	if err != nil {
		return false, err
	}

	content := makeTokenEmail(emailTemplate, &token, emailType)
	subject := emailSubjectMap[*emailType]
	err = sendEmail(email, &subject, content)
	return false, err
}

func SendTeamAccountEmail(tmpl *template.Template, name *string, contestName *string, account *string, password *string, usrEmail *string, subject *string) error {
	content := makeTeamAccountEmail(tmpl, name, contestName, account, password)
	return sendEmail(usrEmail, subject, content)
}
