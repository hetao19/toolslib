package smtp

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

type Smtp struct {
	Address    string
	Username   string
	Password   string
	TLS        bool
	Anonymous  bool
	SkipVerify bool
}

// Compatible with old API
func New(address, username, password string) *Smtp {
	return &Smtp{
		Address:  address,
		Username: username,
		Password: password,
	}
}

func NewSmtp(address, username, password string, tls, anonymous, skipVerify bool) *Smtp {
	return &Smtp{
		Address:    address,
		Username:   username,
		Password:   password,
		TLS:        tls,
		Anonymous:  anonymous,
		SkipVerify: skipVerify,
	}
}

func (s *Smtp) SendMail(from, receivers, subject, body string, contentType ...string) error {
	if s.Address == "" {
		return fmt.Errorf("address is necessary")
	}

	hp := strings.Split(s.Address, ":")
	if len(hp) != 2 {
		return fmt.Errorf("address format error")
	}

	arr := strings.Split(receivers, ";")
	count := len(arr)
	safeArr := make([]string, 0, count)
	for i := 0; i < count; i++ {
		if arr[i] == "" {
			continue
		}
		safeArr = append(safeArr, arr[i])
	}
	if len(safeArr) == 0 {
		return fmt.Errorf("receivers invalid")
	}

	receivers = strings.Join(safeArr, ";")

	b64 := base64.NewEncoding("")

	header := make(map[string]string)
	header["From"] = from
	header["To"] = receivers
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"

	ct := "text/plain; charset=UTF-8"
	if len(contentType) > 0 && contentType[0] == "html" {
		ct = "text/html; charset=UTF-8"
	}

	header["Content-Type"] = ct
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + b64.EncodeToString([]byte(body))

	var auth smtp.Auth
	if !s.Anonymous {
		auth = smtp.PlainAuth("", s.Username, s.Password, hp[0])

	}
	if s.TLS {
		return sendMailByTLS(s.Address, auth, from, strings.Split(receivers, ";"), []byte(message), s.SkipVerify)
	}
	return sendMail(s.Address, auth, from, strings.Split(receivers, ";"), []byte(message), s.SkipVerify)

}

func sendMail(addr string, a smtp.Auth, from string, to []string, msg []byte, skipVerify bool) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	hp := strings.Split(addr, ":")
	if len(hp) != 2 {
		return fmt.Errorf("address fromat error")
	}
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			InsecureSkipVerify: skipVerify,
			ServerName:         hp[0],
		}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if a != nil {
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func sendMailByTLS(addr string, a smtp.Auth, from string, to []string, msg []byte, skipVerify bool) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}
	hp := strings.Split(addr, ":")
	if len(hp) != 2 {
		return fmt.Errorf("address format error")
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		ServerName:         hp[0],
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)

	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, hp[0])
	if err != nil {
		return err
	}
	defer c.Close()
	if a != nil {
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return fmt.Errorf("smtp: A line must not contain CR or LF")
	}
	return nil
}
