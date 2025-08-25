package session

import (
	"os"

	"github.com/gorilla/sessions"
)

const (
	AdminSessionName = "admin"
)

var (
	AdminSession = sessions.NewCookieStore([]byte(os.Getenv("ADMIN_SESSION_KEY")))
)