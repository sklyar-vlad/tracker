package resend

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

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
