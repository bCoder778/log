package log

import (
	"github.com/go-gomail/gomail"
	"strconv"
)

type EMail struct {
	opt *EMailOption
}

type EMailOption struct {
	User   string
	Name   string
	Pass   string
	Host   string
	Port   string
	Target []string
}

func NewEmail(opt *EMailOption) *EMail {
	return &EMail{opt: opt}
}

func (e *EMail) SendEmail(subject string, body string) error {
	m := gomail.NewMessage()

	port, err := strconv.Atoi(e.opt.Port)
	if err != nil {
		return err
	}
	m.SetHeader("From", m.FormatAddress(e.opt.User, e.opt.Name)) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", e.opt.Target...) //发送给多个用户
	m.SetHeader("Subject", subject)    //设置邮件主题
	m.SetBody("text/html", body)       //设置邮件正文

	d := gomail.NewDialer(e.opt.Host, port, e.opt.User, e.opt.Pass)

	return d.DialAndSend(m)
}
