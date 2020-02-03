package mdsession

import (
	"github.com/gorilla/sessions"
)

func NewSessionCookieStoreStore(authKey, encryptionKey []byte) sessions.Store {
	return sessions.NewCookieStore(authKey, encryptionKey)
}
