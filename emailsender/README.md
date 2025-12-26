# PocketBase Email Sender Module

A lightweight Go module that enables server-side triggering of PocketBase authentication emails (verification magic links and OTPs) when using PocketBase as a backend framework.

PocketBase normally only allows requesting verification emails or OTPs from client-side SDKs. This module provides a server-side `EmailSender` interface to call the same built-in endpoints via HTTP, making it possible to send these emails directly from your Go services or custom authentication logic.

## Features

- Send verification emails (magic link) to registered users.
- Request and retrieve OTP ID for registered users.
- Context-aware requests with optional custom HTTP client.
- Simple, dependency-free implementation.

## Installation

```bash
go get github.com/jim-ww/pocketbase-modules/emailsender
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/yourusername/emailsender"
)

func main() {
    // Initialize sender pointing to your PocketBase instance
    emailSender := emailsender.New("http://localhost:8090")
    // Optional: emailsender.New("http://localhost:8090", emailsender.WithCustomHTTPClient(myClient))

    // Send verification email (magic link)
    if err := emailSender.RequestVerification(context.Background(), "user@example.com"); err != nil {
        log.Fatal(err)
    }

    // Request OTP and get the otpId (useful for custom login flows)
    otpId, err := emailSender.RequestOTP(context.Background(), "user@example.com")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("OTP sent, otpId: %s", otpId)
}
```

**Important**: The user with the given email must already exist in the users collection.

## API

```go
type EmailSender interface {
    RequestVerification(ctx context.Context, email string) error
    RequestOTP(ctx context.Context, email string) (otpId string, err error)
}
```
- New(pocketbaseAddress string, opts ...Option) EmailSender
- WithCustomHTTPClient(client *http.Client) Option
