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

	s.users[user.ID] = user
	s.usersByUsername[user.Username] = user.ID
	return nil
}

func (s *Store) FindUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, ok := s.usersByUsername[username]
	if !ok {
		return nil, nil
	}
	return s.users[userID], nil
}

func (s *Store) FindUserByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (s *Store) SaveImage(image *Image) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.images[image.ID] = image
	return nil
}

func (s *Store) FindImageByID(id string) (*Image, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	image, ok := s.images[id]
	if !ok {
		return nil, nil
	}
	return image, nil
}

func (s *Store) ListImagesByUser(userID string, page, limit int) ([]*Image, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userImages []*Image
	for _, img := range s.images {
		if img.UserID == userID {
			userImages = append(userImages, img)
		}
	}

	total := len(userImages)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		return []*Image{}, total, nil
	}
	if end > total {
		end = total
	}

	return userImages[start:end], total, nil
}

func (s *Store) DeleteImage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.images, id)
	return nil
}

func (s *Store) SaveTransformation(transformation *Transformation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transformations[transformation.ID] = transformation
	return nil
}

func (s *Store) FindTransformationByID(id string) (*Transformation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	transformation, ok := s.transformations[id]
	if !ok {
		return nil, nil
	}
	return transformation, nil
}
