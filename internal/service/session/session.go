package session

import (
	"2024_2_FIGHT-CLUB/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
)

type ServiceSession struct {
	store *sessions.CookieStore
}

func NewSessionService(store *sessions.CookieStore) *ServiceSession {
	return &ServiceSession{store: store}
}

func (s *ServiceSession) LogoutSession(r *http.Request, w http.ResponseWriter) error {
	session, _ := s.store.Get(r, "session_id")
	if session.IsNew {
		return errors.New("no such session")
	}
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *ServiceSession) GetUserID(r *http.Request, w http.ResponseWriter) (string, error) {
	session, _ := s.store.Get(r, "session_id")
	if session.IsNew {
		return "", errors.New("no active session")
	}
	return session.Values["id"].(string), nil
}

func (s *ServiceSession) GetSessionData(r *http.Request) (*map[string]interface{}, error) {
	session, _ := s.store.Get(r, "session_id")

	if session.IsNew {
		return nil, errors.New("no active session")
	}

	userID := session.Values["id"].(string)
	Avatar, okAvatar := session.Values["avatar"].(string)

	sessionData := map[string]interface{}{}
	if okAvatar {
		sessionData = map[string]interface{}{
			"id":     userID,
			"avatar": Avatar,
		}
	} else {
		sessionData = map[string]interface{}{
			"id":     userID,
			"avatar": "",
		}
	}

	return &sessionData, nil
}

func (s *ServiceSession) CreateSession(r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
	session, _ := s.store.Get(r, "session_id")

	if !session.IsNew {
		return "", errors.New("session already exists")
	}

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
		return "", errors.New("failed to generate session id")
	}

	session.Values["session_id"] = sessionID

	// Сохраняем сессию
	err = session.Save(r, w)
	if err != nil {
		return "", errors.New("failed to save sessions")
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
