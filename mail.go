package log

import (
	"github.com/go-gomail/gomail"
	"strconv"
	"time"
)

type EMail struct {
	opt *EMailOption
	msg map[string]int64
}

type EMailOption struct {
	Title             string
	User              string
	Name              string
	Pass              string
	Host              string
	Port              string
	Target            []string
	DuplicateRemoval  bool
	DuplicateInterval int64
}

func NewEmail(opt *EMailOption) *EMail {
	if opt.DuplicateInterval == 0 {
		opt.DuplicateInterval = 60 * 60
	}
	e := &EMail{opt: opt}
	if e.opt.DuplicateRemoval {
		e.msg = make(map[string]int64)
		go e.expired()
	}
	return e
}

func (e *EMail) expired() {
	t := time.NewTicker(time.Second * time.Duration(e.opt.DuplicateInterval))
	for {
		select {
		case <-t.C:
			now := time.Now()
			for key, time := range e.msg {
				if now.Unix()-time > e.opt.DuplicateInterval {
					delete(e.msg, key)
				}
			}
		}
	}
}

func (e *EMail) SendEmail(subject string, body string) error {
	if e.opt.DuplicateRemoval {
		key := subject + body
		_, exist := e.msg[key]
		if exist {
			return nil
		} else {
			e.msg[key] = time.Now().Unix()
		}
	}
	m := gomail.NewMessage()
	if subject == "" {
		subject = e.opt.Title
	}
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
