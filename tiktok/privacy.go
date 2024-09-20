package tiktok

/*
"privacy_level_options": ["PUBLIC_TO_EVERYONE", "MUTUAL_FOLLOW_FRIENDS", "SELF_ONLY"] 
*/

type PrivacyLevelOptions string

const (
    PUBLIC_TO_EVERYONE PrivacyLevelOptions 		= "PUBLIC_TO_EVERYONE"
    MUTUAL_FOLLOW_FRIENDS PrivacyLevelOptions 	= "MUTUAL_FOLLOW_FRIENDS"
	FOLLOWER_OF_CREATOR PrivacyLevelOptions 	= "FOLLOWER_OF_CREATOR"
    SELF_ONLY PrivacyLevelOptions 				= "SELF_ONLY"
)

// CheckPrivacyLevel checks if the given string matches any PrivacyLevelOptions
func CheckPrivacyLevel(level string) bool {
	switch PrivacyLevelOptions(level) {
	case PUBLIC_TO_EVERYONE, MUTUAL_FOLLOW_FRIENDS, FOLLOWER_OF_CREATOR, SELF_ONLY:
		return true
	default:
		return false
	}
}