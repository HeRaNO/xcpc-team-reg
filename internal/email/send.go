package email

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"github.com/HeRaNO/xcpc-team-reg/internal/contest"
	"github.com/HeRaNO/xcpc-team-reg/internal/dal/redis"
	"github.com/HeRaNO/xcpc-team-reg/internal/utils"
	"github.com/wneessen/go-mail"
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
	message := mail.NewMsg()
	if err := message.From(emailFrom); err != nil {
		return err
	}
	if err := message.To(*emailRecv); err != nil {
		return err
	}
	message.Subject(*subject)
	message.SetBodyString(mail.TypeTextHTML, string(content))

	return client.DialAndSend(message)
}

func SendEmailWithToken(ctx context.Context, email *string, emailType *string) berrors.Berror {
	if _, ok := emailActionMap[*emailType]; !ok {
		return errInvalidType
	}
	uid, err := redis.GetUserIDByEmail(ctx, email)
	if err != nil {
		return err
	}
	if *emailType == "register" {
		if uid != 0 {
			return errAlreadyRegistered
		}
	} else {
		if uid == 0 {
			return errNoRegisterRecord
		}
	}

	err = redis.GetEmailRequest(ctx, email)
	if err != nil {
		return err
	}

	token, err := utils.GenToken(contest.UserTokenLength)
	if err != nil {
		return err
	}

	err = redis.SetEmailToken(ctx, email, &token, emailTokenExpireTime)
	if err != nil {
		return err
	}
	err = redis.SetEmailAction(ctx, email, emailType, emailTokenExpireTime)
	if err != nil {
		return err
	}
	err = redis.SetEmailRequest(ctx, email, emailSendGapTime)
	if err != nil {
		return err
	}

	content := makeTokenEmail(emailTemplate, &token, emailType)
	subject := emailSubjectMap[*emailType]
	erro := sendEmail(email, &subject, content)
	if erro != nil {
		return berrors.New(berrors.ErrInternal, erro.Error())
	}
	return nil
}

func SendTeamAccountEmail(tmpl *template.Template, name *string, contestName *string, account *string, password *string, usrEmail *string, subject *string) error {
	content := makeTeamAccountEmail(tmpl, name, contestName, account, password)
	return sendEmail(usrEmail, subject, content)
}
