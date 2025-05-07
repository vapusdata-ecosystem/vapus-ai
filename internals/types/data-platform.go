package types

type FileUserType string

const (
	UserType  FileUserType = "user"
	GroupType FileUserType = "group"
)

func (ut FileUserType) String() string {
	return string(ut)
}

const (
	SECRETENGINE = "secreteengine"
)
