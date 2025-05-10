package model

type UserRole uint

const (
	AdminRole  UserRole = 7
	EditorRole UserRole = 5
	ViewerRole UserRole = 1
)
