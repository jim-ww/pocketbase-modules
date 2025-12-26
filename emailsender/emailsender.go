package emailsender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmailSender interface {
	// RequestVerification will send verification email containing link with token to sign in (magic link) to specified email.
	//
	// note: user with specified email must already be registered (exist in database)
	RequestVerification(ctx context.Context, email string) error
	// RequestOTP will send email containing One-Time-Password to sign in, to specified email.
	//
	// note: user with specified email must already be registered (exist in database)
	RequestOTP(ctx context.Context, email string) (otpId string, err error)
}

type EmailSenderImpl struct {
	// Address of pocketbase instance, where to send requests to. (e.g. http://localhost:8090)
	pocketbaseAddress string
	client            *http.Client
}

type Option func(*EmailSenderImpl)

func WithCustomHTTPClient(client *http.Client) Option {
	return func(es *EmailSenderImpl) {
		es.client = client
	}
}

func New(pocketbaseAddress string, opts ...Option) EmailSender {
	emailSender := &EmailSenderImpl{pocketbaseAddress: pocketbaseAddress, client: http.DefaultClient}
	for _, opt := range opts {
		opt(emailSender)
	}
	return emailSender
}

type emailDTO struct {
	Email string `json:"email"`
}

func (es *EmailSenderImpl) RequestVerification(ctx context.Context, email string) error {
	dtoJSON, err := json.Marshal(emailDTO{
		Email: email,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, es.pocketbaseAddress+"/api/collections/users/request-verification", bytes.NewReader(dtoJSON))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	if _, err := es.client.Do(req); err != nil {
		return fmt.Errorf("failed to send http POST request for verification: %w", err)
	}

	return nil
}

func (es *EmailSenderImpl) RequestOTP(ctx context.Context, email string) (otpId string, err error) {
	dtoJSON, err := json.Marshal(emailDTO{
		Email: email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, es.pocketbaseAddress+"/api/collections/users/request-otp", bytes.NewReader(dtoJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send http POST request to get OTP: %w", err)
	}
	defer resp.Body.Close()

	respBody := struct {
		OTPID string `json:"otpId"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	return respBody.OTPID, nil
}
