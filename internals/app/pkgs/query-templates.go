package apppkgs

import (
	"fmt"

	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
)

func ListResourceWithGovernance(ctxClaim map[string]string) string {
	if ctxClaim != nil {
		return fmt.Sprintf(`status != 'DELETED' AND deleted_at IS NULL AND 
		(
	(scope='PLATFORM_SCOPE') OR (scope='ORGANIZATION_SCOPE' AND organization='%s') OR (scope='USER_SCOPE' AND organization='%s' AND created_by='%s')
		)
		 ORDER BY created_at DESC`, ctxClaim[encryption.ClaimOrganizationKey], ctxClaim[encryption.ClaimOrganizationKey], ctxClaim[encryption.ClaimUserIdKey])
	}
	return ""
}
