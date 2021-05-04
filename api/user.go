package api

// User is a Unifi OS user.
type User struct {
	UniqueID   string `json:"unique_id"`
	Username   string `json:"username"`
	CreateTime Time   `json:"create_time"`
	UpdateTime Time   `json:"update_time"`

	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	SSO SSOUser `json:",inline"`

	LocalAccountExists bool `json:"local_account_exist"`
	CloudAccessGranted bool `json:"cloud_access_granted"`

	Owner      bool `json:"isOwner"`
	SuperAdmin bool `json:"isSuperAdmin"`

	Permissions map[string][]string `json:"permissions"`
	Scopes      []string            `json:"scopes"`
}

// SSOUser is user information from Ubiquti cloud single sign-on.
type SSOUser struct {
	UUID      string `json:"sso_uuid"`
	Username  string `json:"sso_username"`
	Email     string `json:"sso_account"`
	AvatarURL string `json:"sso_picture"`
}
