package user_service

import "fmt"

type User struct {
	ID   string
	Name string
}

type UserService struct {
	users map[string]User
}

func NewUserService() *UserService {
	return &UserService{users: make(map[string]User)}
}

func (s *UserService) CreateUser(id, name string) (User, error) {
	if id == "" {
		return User{}, fmt.Errorf("user ID cannot be empty")
	}
	if _, exists := s.users[id]; exists {
		return User{}, fmt.Errorf("user %s already exists", id)
	}
	nu := User{ID: id, Name: name}
	s.users[id] = nu
	fmt.Printf("[UserService] User created: %+v\n", nu)
	return nu, nil
}

func (s *UserService) GetUser(id string) (User, error) {
	user, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %s not found", id)
	}
	fmt.Printf("[UserService] User retrieved: %+v\n", user)
	return user, nil
}

func (s *UserService) UpdateUser(id, newName string) error {
	if id == "" || newName == "" {
		return fmt.Errorf("id and newName cannot be empty for UpdateUser")
	}
	user, ok := s.users[id]
	if !ok {
		return fmt.Errorf("user %s not found", id)
	}
	user.Name = newName
	s.users[id] = user
	fmt.Printf("[UserService] UpdateUser called for ID: %s, NewName: %s\n", id, newName)
	return nil
}
