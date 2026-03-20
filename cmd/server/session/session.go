package session

import (
	"gosecureskeleton/cmd/server/objects"
	"gosecureskeleton/cmd/server/util"
	"sync"
)

type SessionData struct {
	User objects.User
}

type SessionStore struct {
	Tokens map[string]SessionData
	mutex  sync.Mutex
}

func newSessionStore() *SessionStore {
	return &SessionStore{
		Tokens: make(map[string]SessionData),
	}
}

func (s *SessionStore) Create(user objects.User) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	token, err := util.GenerateToken()
	if err != nil {
		return "", err
	}

	s.Tokens[token] = SessionData{user}
	return token, nil
}

func (s *SessionStore) Lookup(token string) (SessionData, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, ok := s.Tokens[token]
	return data, ok
}

func (s *SessionStore) Delete(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Tokens, token)
}

var Sessions SessionStore = *newSessionStore()
