package domain

type OauthToken struct {
	OauthId         int    `gorm:"oauth_id;primaryKey;autoIncrement"`
	UserId          string `gorm:"column:user_id"`
	OauthProvider   string `gorm:"column:oauth_provider"`
	OauthProviderId string `gorm:"column:oauth_provider_id"`
}

func (OauthToken) TableName() string {
	return "oauth_token"
}
