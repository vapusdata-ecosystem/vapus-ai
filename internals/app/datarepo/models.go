package datarepo

type LocalUserM struct {
	DisplayName       string
	FirstName         string
	LastName          string
	Email             string
	Organization      string
	PlatformRole      string
	OrganizationRoles []string
	AccountId         string
	ProfileImage      string
}
