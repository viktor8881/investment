package email

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strconv"
)

type settings struct {
	host     string
	port     int64
	username string
	password string
	to       []string
	from     string
	subject  string
}

var config settings

func SetSetting(host string, port int64, username, password string, to []string, from, subject string) {
	config.host = host
	config.port = port
	config.username = username
	config.password = password
	config.to = to
	config.from = from
	config.subject = subject
}

func SendAnalytics(templateFileName string, data interface{}) (bool, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return false, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return false, err
	}
	body := buf.String()

	auth := smtp.PlainAuth("", config.username, config.password, config.host)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + config.subject + "!\n"
	msg := []byte(subject + mime + "\n" + body)
	addr := config.host + ":" + strconv.Itoa(int(config.port))

	if err := smtp.SendMail(addr, auth, config.from, config.to, msg); err != nil {
		return false, err
	}
	return true, nil
}
