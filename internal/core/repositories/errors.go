package repositories

import "errors"

var (
	NoChangesHappenedError = errors.New("no changes happened")
	RegisterAlreadyExists  = errors.New("register already exists")
)
