package tiktok

import "errors"

// Input errors, returned before any HTTP call is made.
var (
	// ErrPrivacyLevel is returned when the privacy level is not one of the
	// values TikTok accepts (see CheckPrivacyLevel).
	ErrPrivacyLevel = errors.New("tiktok: privacy level is not valid")
	// ErrPostMode is returned when the post mode is neither DIRECT_POST nor
	// MEDIA_UPLOAD (see CheckPostMode).
	ErrPostMode = errors.New("tiktok: post mode is not valid")
	// ErrClientKeyRequired and ErrClientSecretRequired are returned by
	// NewTikTok: every TikTok endpoint refuses a request from an app it cannot
	// identify, so failing at construction beats failing later with an opaque
	// API error.
	ErrClientKeyRequired    = errors.New("tiktok: client key is required")
	ErrClientSecretRequired = errors.New("tiktok: client secret is required")
)

// Deprecated aliases. They are the same error values, so errors.Is keeps
// working for existing callers; only the message changed.
var (
	// Deprecated: use ErrPrivacyLevel.
	PrivacyLevelWrong = ErrPrivacyLevel
	// Deprecated: use ErrPostMode.
	PhotoModeWrong = ErrPostMode
)
