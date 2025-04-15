package email

import (
	"fmt"
	"html/template"

	"github.com/HeRaNO/xcpc-team-reg/internal/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/wneessen/go-mail"
)

func Init(conf *config.EmailConfig) {
	if conf == nil {
		hlog.Fatal("Email config failed: conf is nil")
		panic("make static check happy")
	}
	emailSign = conf.EmailSign
	emailFrom = fmt.Sprintf("%s <%s>", conf.EmailAlias, conf.EmailAddr)

	opts := []mail.Option{
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(conf.EmailAddr),
		mail.WithPassword(conf.EmailPassword),
	}

	switch conf.SMTPEncMethod {
	case "", "plain", "ssl", "tls":
	default:
		hlog.Warnf("unrecognised encrypt method: %s, fallback to plain", conf.SMTPEncMethod)
	}

	if conf.SMTPPort == 465 || conf.SMTPEncMethod == "ssl" {
		opts = append(opts, mail.WithSSLPort(false))
	} else if conf.SMTPPort == 587 || conf.SMTPEncMethod == "tls" {
		opts = append(opts, mail.WithTLSPortPolicy(mail.TLSMandatory))
	} else if conf.SMTPPort != 25 {
		opts = append(opts, mail.WithPort(conf.SMTPPort))
	}

	c, err := mail.NewClient(conf.SMTPAddr, opts...)
	if err != nil {
		hlog.Fatalf("create email client failed, err: %+v", err)
	}
	client = c

	tmpl, err := template.ParseFiles("./configs/email-verification.tmpl")
	if err != nil {
		hlog.Fatalf("parse email template failed, err: %+v", err)
	}
	emailTemplate = tmpl
	hlog.Info("init email finished successfully")
}
