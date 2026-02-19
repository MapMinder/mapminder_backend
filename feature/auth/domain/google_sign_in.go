package domain

const (
	GoogleProvider = "google"
)

// GoogleSignIn
type GoogleSignInCreds struct {
	IdToken  string
	ClientId string
}

type GoogleIdToken struct {
	Issuer    string
	Audience  string
	Subject   string
	IssuedAt  int64
	ExpiresAt int64
	Claims    GoogleIdTokenClaims
}

type GoogleIdTokenClaims struct {
	Issuer          string `json:"iss"`
	Audience        string `json:"aud"`
	AuthorizedParty string `json:"azp,omitempty"`
	Subject         string `json:"sub"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`

	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`

	Nonce  string `json:"nonce,omitempty"`
	AtHash string `json:"at_hash,omitempty"`

	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}
