package tiktok

import "errors"

var (
	PrivacyLevelWrong = errors.New("Privacy Level is not correct!")
	PhotoModeWrong    = errors.New("Photo mode is not correct!")
)

// Input errors returned before any HTTP call is made.
var (
	// ErrClientKeyRequired and ErrClientSecretRequired are returned by
	// NewTikTok: every TikTok endpoint refuses a request from an app it cannot
	// identify, so failing at construction beats failing later with an opaque
	// API error.
	ErrClientKeyRequired    = errors.New("tiktok: client key is required")
	ErrClientSecretRequired = errors.New("tiktok: client secret is required")
)
