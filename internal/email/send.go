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

func makeTokenEmail(tmpl *template.Template, token, method *string) ([]byte, error) {
	content := new(bytes.Buffer)
	err := tmpl.Execute(content, struct {
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
	if err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}

func makeTeamAccountEmail(tmpl *template.Template, name, contestName, account, password *string) ([]byte, error) {
	content := new(bytes.Buffer)
	err := tmpl.Execute(content, struct {
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
	if err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}

func sendEmail(emailRecv, subject *string, content []byte) error {
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

func SendEmailWithToken(ctx context.Context, email, emailType *string) berrors.Berror {
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

	content, erro := makeTokenEmail(emailTemplate, &token, emailType)
	if erro != nil {
		return berrors.New(berrors.ErrInternal, erro.Error())
	}
	subject := emailSubjectMap[*emailType]
	erro = sendEmail(email, &subject, content)
	if erro != nil {
		return berrors.New(berrors.ErrInternal, erro.Error())
	}
	return nil
}

func SendTeamAccountEmail(tmpl *template.Template, name, contestName, account, password, usrEmail, subject *string) error {
	content, err := makeTeamAccountEmail(tmpl, name, contestName, account, password)
	if err != nil {
		return err
	}
	return sendEmail(usrEmail, subject, content)
}
