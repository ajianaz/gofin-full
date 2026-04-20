package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ajianaz/gofin-full/api/internal/auth"
	"github.com/ajianaz/gofin-full/api/internal/repository"
	apperrors "github.com/ajianaz/gofin-full/api/pkg/errors"
)

// WalletRBAC checks wallet membership and enforces role-based access.
// requiredRole: "owner", "editor", "viewer"
//   - owner: full access (wallet creator)
//   - editor: can create/modify transactions
//   - viewer: read-only access
//
// Owner always has access. Editor has access if requiredRole is editor or viewer.
// Viewer only has access if requiredRole is viewer.
func WalletRBAC(memberRepo *repository.WalletMemberRepository, requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := auth.GetUser(c)
		if user == nil {
			return apperrors.New(401, "unauthorized")
		}

		walletID, err := c.ParamsInt("wallet_id")
		if err != nil {
			return apperrors.NewValidationError(map[string][]string{"wallet_id": {"invalid wallet id"}})
		}

		// Check if user is the wallet owner
		isOwner, err := memberRepo.IsWalletOwner(c.Context(), int64(walletID), user.ID)
		if err != nil {
			return apperrors.NewWithDetail(500, "failed to check wallet ownership", err.Error())
		}
		if isOwner {
			return c.Next()
		}

		// Check membership
		role, err := memberRepo.GetWalletRole(c.Context(), int64(walletID), user.ID)
		if err != nil {
			return apperrors.NewWithDetail(500, "failed to check wallet membership", err.Error())
		}
		if role == "" {
			return apperrors.New(403, "you do not have access to this wallet")
		}

		// Check role hierarchy: owner > editor > viewer
		switch requiredRole {
		case "owner":
			return apperrors.New(403, "only wallet owner can perform this action")
		case "editor":
			if role != "editor" {
				return apperrors.New(403, "editor or owner access required")
			}
		case "viewer":
			// both editor and viewer can access
		}

		return c.Next()
	}
}
