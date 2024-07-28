package tiktok

import "errors"

var (
	PrivacyLevelWrong = errors.New("Privacy Level is not correct!")
	PhotoModeWrong = errors.New("Photo mode is not correct!")
)
