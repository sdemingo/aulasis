package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"
)

const BOUNDARY = "9p23e030u5x30oz6nmhiu2302x54aesrda"

type EmailMsg struct {
	From            string
	To              []string
	Subject         string
	Body            string
	BodyContentType string
	Attachments     map[string][]byte
}


func (m *EmailMsg) Attach(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	_, fileName := filepath.Split(file)
	m.Attachments[fileName] = b
	return nil
}


func NewEmailMsg(subject string, body string) *EmailMsg {
	m := &EmailMsg{Subject: subject, Body: body}
	m.BodyContentType= "text/plain"
	m.Attachments = make(map[string][]byte)
	return m
}


func (m *EmailMsg) GetBytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.From + "\n")
	buf.WriteString("To: " + strings.Join(m.To, ",") + "\n")
	buf.WriteString("Subject: " + m.Subject + "\n")
	buf.WriteString("MIME-Version: 1.0\n")

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + BOUNDARY + "\n")
		buf.WriteString("--" + BOUNDARY + "\n")
	}
	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\n", m.BodyContentType))
	buf.WriteString(m.Body)

	if len(m.Attachments) > 0 {
		for k, v := range m.Attachments {
			buf.WriteString("\n\n--" + BOUNDARY + "\n")
			buf.WriteString("Content-Type: application/octet-stream\n")
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString("Content-Disposition: attachment; filename=\"" + k + "\"\n\n")

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString("\n--" + BOUNDARY)
		}
		buf.WriteString("--")
	}
	return buf.Bytes()
}


func (m *EmailMsg) SendEmailMsg(smtpData SmtpConfig) error {
	auth:=smtp.PlainAuth("", smtpData.User, smtpData.Password, smtpData.Addr)
	addr:=smtpData.Addr+":"+smtpData.Port
	return smtp.SendMail(addr, auth, m.From, m.To, m.GetBytes())
}


