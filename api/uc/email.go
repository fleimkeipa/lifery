package uc

import "github.com/fleimkeipa/lifery/repositories/interfaces"

type EmailUC struct {
	emailRepo interfaces.EmailInterfaces
}

func NewEmailUC(emailRepo interfaces.EmailInterfaces) *EmailUC {
	return &EmailUC{emailRepo: emailRepo}
}

func (uc *EmailUC) SendPasswordResetEmail(to, username, resetToken string) error {
	return uc.emailRepo.SendPasswordResetEmail(to, username, resetToken)
}
