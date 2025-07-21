package doubles

import "fmt"

type User struct {
	ID    string
	Email string
	Name  string
}

type UserNotifier interface {
	NotifyUserCreated(user User) error
}

type UserCreationService struct {
	notifier UserNotifier
	// ... other dependencies like a UserRepository
}

func NewUserCreationService(notifier UserNotifier) *UserCreationService {
	return &UserCreationService{notifier: notifier}
}

func (s *UserCreationService) CreateUserAndNotify(id, email, name string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("email cannot be empty")
	}
	user := User{ID: id, Email: email, Name: name}
	// ... logic to save user to a repository ...
	fmt.Printf("Service: User %s created.\n", user.ID)

	err := s.notifier.NotifyUserCreated(user)
	if err != nil {
		// Log the notification error but don't fail user creation for this demo
		fmt.Printf("Service: Failed to notify user %s creation (non-fatal): %v\n", user.ID, err)
	}
	return user, nil
}
