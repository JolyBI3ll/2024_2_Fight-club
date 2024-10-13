package service

import (
	"2024_2_FIGHT-CLUB/domain"
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/sessions"
	"net/http"
)

type SessionService struct {
	store *sessions.CookieStore
}

func NewSessionService(store *sessions.CookieStore) *SessionService {
	return &SessionService{store: store}
}

func (s *SessionService) CreateSession(r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
	session, _ := s.store.Get(r, "session_id")

	session.Values["id"] = user.UUID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email
	if user.Name != "" {
		session.Values["name"] = user.Name
	}
	if user.Avatar != "" {
		session.Values["avatar"] = user.Avatar
	}

	sessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}

	session.Values["session_id"] = sessionID

	// Сохраняем сессию
	err = session.Save(r, w)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
