package resend

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

type client struct {
	apiClient *resend.Client
}

func newClient(key string) *client {
	return &client{
		apiClient: resend.NewClient(key),
	}
}

func (c *client) sendEmail(msg *resend.SendEmailRequest) error {
	_, err := c.apiClient.Emails.Send(msg)
	if err != nil {
		return fmt.Errorf("failed sent email for verification: %v", err)
	}

	return nil
}
