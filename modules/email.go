package modules

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	"github.com/jordan-wright/email"
	jsoniter "github.com/json-iterator/go"
)

func sendEmail(emailRecv *string, token *string, method *string) error {
	tmpl, err := template.ParseFiles("./template/email-verification.tmpl")

	if err != nil {
		return err
	}

	content := new(bytes.Buffer)
	tmpl.Execute(content, struct {
		Action string
		Time   string
		Token  string
		Sign   string
	}{
		Action: config.EmailActionMap[*method],
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Token:  *token,
		Sign:   config.EmailSign,
	})

	e := &email.Email{
		To:      []string{*emailRecv},
		From:    config.EmailFrom,
		Subject: config.EmailSubjectMap[*method],
		HTML:    content.Bytes(),
	}

	auth := smtp.PlainAuth("", config.EmailAddr, config.EmailPassword, config.EmailServer)
	err = e.SendWithTLS(config.SMTPHost, auth, &tls.Config{ServerName: config.EmailServer})

	return err
}

func SendValidationEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	info := model.EmailVerification{}
	err = jsoniter.Unmarshal(bd, &info)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	err = model.GetEmailRequest(r.Context(), &info.Email)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	token := util.GenToken(config.UserTokenLength)

	err = model.SetEmailToken(r.Context(), &info.Email, &token)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	err = model.SetEmailAction(r.Context(), &info.Email, &info.Type)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	model.SetEmailRequest(r.Context(), &info.Email)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	err = sendEmail(&info.Email, &token, &info.Type)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	util.SuccessResponse(w, r, "email token sended")
}
