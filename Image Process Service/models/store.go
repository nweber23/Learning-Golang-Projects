package models

import "sync"

type Store struct {
	users           map[string]*User
	images          map[string]*Image
	transformations map[string]*Transformation
	usersByUsername map[string]string
	mu              sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		users:           make(map[string]*User),
		images:          make(map[string]*Image),
		transformations: make(map[string]*Transformation),
		usersByUsername: make(map[string]string),
	}
}

func (s *Store) SaveUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil
}

func (s *Store) FindUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return nil, nil
}

func (s *Store) FindUserByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return nil, nil
}

func (s *Store) SaveImage(image *Image) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil
}

func (s *Store) FindImageByID(id string) (*Image, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return nil, nil
}

func (s *Store) ListImagesByUser(userID string, page, limit int) ([]*Image, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return nil, 0, nil
}

func (s *Store) DeleteImage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil
}

func (s *Store) SaveTransformation(transformation *Transformation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil
}

func (s *Store) FindTransformationByID(id string) (*Transformation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return nil, nil
}
