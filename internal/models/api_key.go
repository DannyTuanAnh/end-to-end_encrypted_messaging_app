package models

type GenerateAPIKeyArgs struct {
	ClientType string
}

type RevokeAPIKeyArgs struct {
	KeyID     string
	RevokeAll bool
}
