package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type EmailCredentials struct {
	Username string
	Password string
	Smtp     string
}

func GetEmailCredentials() (*EmailCredentials, error) {
	EMAIL_USERNAME, err := GetEnv("EMAIL_USERNAME")
	if err != nil {
		return nil, err
	}

	EMAIL_PASSWORD, err := GetEnv("EMAIL_PASSWORD")
	if err != nil {
		return nil, err
	}

	EMAIL_SMTP := "smtp.gmail.com"

	return &EmailCredentials{
		Username: *EMAIL_USERNAME,
		Password: *EMAIL_PASSWORD,
		Smtp:     EMAIL_SMTP,
	}, nil
}

func SendCreateUserEmail(
	templateName string,
	to string,
	subject string,
	password string,
) error {
	credentials, err := GetEmailCredentials()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", credentials.Username, credentials.Password, credentials.Smtp)

	t, err := template.ParseFiles(fmt.Sprintf("./src/privates/%s", templateName))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, headers)))

	if err := t.Execute(&body, struct {
		Email    string
		Password string
	}{
		Email:    to,
		Password: password,
	}); err != nil {
		return err
	}

	if err := smtp.SendMail(credentials.Smtp+":587", auth, credentials.Username, []string{to}, body.Bytes()); err != nil {
		return err
	}

	return nil
}

func SendEmail(
	templateName string,
	to []string,
	subject string,
	data any,
) error {
	credentials, err := GetEmailCredentials()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", credentials.Username, credentials.Password, credentials.Smtp)

	t, err := template.ParseFiles(fmt.Sprintf("./src/privates/%s", templateName))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	if _, err := body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, headers))); err != nil {
		return nil
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	if err := smtp.SendMail(
		credentials.Smtp+":587",
		auth,
		credentials.Username,
		to,
		body.Bytes(),
	); err != nil {
		return err
	}

	return nil
}
