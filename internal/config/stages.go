package config

type Stage string

const (
	PRE                   Stage = "PRE"
	POST_SUCCESS          Stage = "POST_SUCCESS"
	POST_SUCCESS_TRANSFER Stage = "POST_SUCCESS_TRANSFER"
	POST_FAILURE          Stage = "POST_FAILURE"
)
