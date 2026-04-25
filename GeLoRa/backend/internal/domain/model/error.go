package domainmodel

import "errors"

var (
	ErrQueryBuilder = errors.New("Query building failed")

	ErrNodeMidExists = errors.New("Node Manufacture ID already exists")
	ErrNodeNotFound  = errors.New("Node can not be found")
	ErrNodeNoChange  = errors.New("Node has no change")

	ErrUserUsernameExists     = errors.New("User username already exists")
	ErrUserNotFound           = errors.New("User can not be found")
	ErrUserInvalidCredentials = errors.New("User credentials are invalid")

	ErrSessionNotFound     = errors.New("Session can not be found")
	ErrSessionNoIdentifier = errors.New("Session requires at least one identifier")
	ErrSessionAlreadyActive = errors.New("Session is already active")
)
