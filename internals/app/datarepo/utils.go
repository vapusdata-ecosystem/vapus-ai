package datarepo

import (
	"fmt"

	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func GetRequestLabel(user string, Organization string, dataProductId string) string {
	return fmt.Sprintf(types.VDC_USER_SELECTOR, user, Organization, dataProductId)
}

func GetUserOrganizationLabel(user string, Organization string) string {
	return fmt.Sprintf("%s-%s", user, Organization)
}

func GetTransformerLabelMap(dataProductId, selectorLabel string) string {
	return fmt.Sprintf(types.GENERIC_SELECTOR, dataProductId, selectorLabel)
}

func GetAttributeFilterLabelMap(dataProductId, selectorLabel string) string {
	return fmt.Sprintf(types.GENERIC_SELECTOR, dataProductId, selectorLabel)
}
