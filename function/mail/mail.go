package mail

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
)

type Sender interface {
	Send(ctx context.Context, subject, body, recipient string) (string, string, error)
}

func New(domain, apiKey string) Sender {
	return &senderImpl{
		domain: domain,
		apiKey: apiKey,
	}
}

type senderImpl struct {
	domain string
	apiKey string
}

func (s *senderImpl) Send(ctx context.Context, subject, body, recipient string) (string, string, error) {
	mg := mailgun.NewMailgun(s.domain, s.apiKey)
	sender := fmt.Sprintf("postmaster@%s", s.domain)
	m := mg.NewMessage(sender, subject, body, recipient)
	status, id, err := mg.Send(ctx, m)
	if err != nil {
		return "", "", fmt.Errorf("failed to send mail: %w", err)
	}
	return status, id, nil
}
