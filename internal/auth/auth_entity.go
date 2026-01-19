package auth

// SignUpData represents user registration data for domain layer
type SignUpData struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Photo           string
}

// SignInData represents user login data for domain layer
type SignInData struct {
	Email    string
	Password string
}
