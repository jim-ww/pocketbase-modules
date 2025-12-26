package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

func LoadAuthContextFromCookie(authCookieName string) func(re *core.RequestEvent) error {
	return func(re *core.RequestEvent) error {
		cookie, err := re.Request.Cookie(authCookieName)
		if err != nil || cookie.Value == "" {
			return re.Next()
		}

		record, err := re.App.FindAuthRecordByToken(cookie.Value)
		if err != nil {
			// token is invalid
			ClearAuthCookie(re.Response, authCookieName)
			return re.Next()
		}
		re.Auth = record

		return re.Next()
	}
}

func SetAuthToken(authCookieName string) func(re *core.RequestEvent, user *core.Record) error {
	return func(re *core.RequestEvent, user *core.Record) error {
		s, err := user.NewAuthToken()
		if err != nil {
			return fmt.Errorf("failed to generate new auth token: %w", err)
		}

		re.SetCookie(&http.Cookie{
			Name:     authCookieName,
			Value:    s,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		})

		return nil
	}
}

func AuthorizeAccessToRecord(re *core.RequestEvent, record *core.Record, rule *string) error {
	requestInfo, err := re.RequestInfo()
	if err != nil {
		return err
	}

	canAccess, err := re.App.CanAccessRecord(record, requestInfo, rule)
	if err != nil {
		return err
	}

	if !canAccess {
		return re.ForbiddenError("Forbidden", nil)
	}

	return nil
}

func ClearAuthCookie(w http.ResponseWriter, authCookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
