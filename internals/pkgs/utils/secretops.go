package dmutils

import (
	guuid "github.com/google/uuid"
)

func GetSecretName(resource, resourceId, attribute string) string {
	if resource == "" {
		return guuid.NewString() + "_" + attribute
	}
	return resourceId + "_" + attribute
}
