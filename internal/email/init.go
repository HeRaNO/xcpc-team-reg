package email

import (
	"fmt"
	"html/template"

	"github.com/HeRaNO/xcpc-team-reg/internal/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Init(conf *config.EmailConfig) {
	if conf == nil {
		hlog.Fatal("Email config failed: conf is nil")
		panic("make static check happy")
	}
	emailSign, emailAddr = conf.EmailSign, conf.EmailAddr
	emailPassword, emailAlias = conf.EmailPassword, conf.EmailAlias
	emailFrom = fmt.Sprintf("%s <%s>", emailAlias, emailAddr)
	smtpAddr, smtpPort = conf.SMTPAddr, conf.SMTPPort
	smtpHost = fmt.Sprintf("%s:%d", smtpAddr, smtpPort)
	tmpl, err := template.ParseFiles("./configs/email-verification.tmpl")
	if err != nil {
		hlog.Fatalf("parse email template failed, err: %+v", err)
	}
	emailTemplate = tmpl
	hlog.Info("init email finished successfully")
}
