package tiktok

/*
"privacy_level_options": ["PUBLIC_TO_EVERYONE", "MUTUAL_FOLLOW_FRIENDS", "SELF_ONLY"] 
*/

type PrivacyLevelOptions string

const (
    PUBLIC_TO_EVERYONE PrivacyLevelOptions = "PUBLIC_TO_EVERYONE"
    MUTUAL_FOLLOW_FRIENDS PrivacyLevelOptions = "MUTUAL_FOLLOW_FRIENDS"
    SELF_ONLY PrivacyLevelOptions = "SELF_ONLY"
)