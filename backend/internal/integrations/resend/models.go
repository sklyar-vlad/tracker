package resend

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/resend/resend-go/v3"
)

type emailData struct {
	VerifyURL string
}

func newHtml(token, url string) (string, error) {
	verifyURL := fmt.Sprintf(url, token)
	tmpl, err := template.ParseFiles("/app/templates/email.html")
	if err != nil {
		return "", err
	}

	data := emailData{
		VerifyURL: verifyURL,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func newParams(from, html, subject string, to []string) (*resend.SendEmailRequest, error) {
	if to == nil {
		return nil, fmt.Errorf("invalid email")
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      to,
		Html:    html,
		Subject: subject,
	}
	return params, nil
}
