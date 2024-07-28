package tiktok

type PostModeOptions string

const (
    DIRECT_POST PostModeOptions = "DIRECT_POST"
    MEDIA_UPLOAD PostModeOptions = "MEDIA_UPLOAD"
)

// CheckPostMode checks if the given string matches any PostModeOptions
func CheckPostMode(mode string) bool {
	switch PostModeOptions(mode) {
	case DIRECT_POST, MEDIA_UPLOAD:
		return true
	default:
		return false
	}
}