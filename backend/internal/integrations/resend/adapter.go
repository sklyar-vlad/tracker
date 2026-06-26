package resend

import (
	"go.uber.org/zap"

	"github.com/sklyar-vlad/selfDev/internal/config"
)

type adapter struct {
	client *client
	config config.ConfigEmailSender
	logger *zap.Logger
}

func NewAdapter(cfg config.ConfigEmailSender, logger *zap.Logger) *adapter {
	return &adapter{
		client: newClient(cfg.ApiKey),
		config: cfg,
		logger: logger,
	}
}

func (a *adapter) SendEmailVerification(email string) error {
	msg, err := newParams(a.config.From, a.config.Html, a.config.Subject, []string{email})
	if err != nil {
		a.logger.Error("failed create message request", zap.Error(err))
		return err
	}

	err = a.client.sendEmail(msg)
	a.logger.Info("msg", zap.String("msg", string(msg.Html)))
	if err != nil {
		a.logger.Error("failed send email message", zap.Error(err))
		return err
	}

	return nil
}
