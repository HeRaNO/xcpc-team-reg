package email

import (
	"html/template"
	"time"
)

const (
	emailSendGapTime     = 2 * time.Minute
	emailTokenExpireTime = 10 * time.Minute
	stuEmailSuffix       = "@std.uestc.edu.cn"
)

var smtpAddr, smtpHost string
var smtpPort int
var emailSign, emailAddr, emailPassword, emailAlias, emailFrom string
var emailTemplate *template.Template

var emailActionMap = map[string]string{
	"register": "注册账户",
	"reset":    "重置密码",
}

var emailSubjectMap = map[string]string{
	"register": "邮箱验证邮件 - 注册账户",
	"reset":    "邮箱验证邮件 - 重置密码",
}
