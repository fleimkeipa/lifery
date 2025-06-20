package interfaces

type EmailInterfaces interface {
	SendPasswordResetEmail(to, username, resetToken string) error
}
