package service

import (
	"context"

	"github.com/ajianaz/gofin-full/api/internal/repository"
	"github.com/ajianaz/gofin-full/api/internal/sse"
)

type NotificationService struct {
	repo *repository.NotificationRepository
	hub  *sse.Hub
}

func NewNotificationService(repo *repository.NotificationRepository, hub *sse.Hub) *NotificationService {
	return &NotificationService{repo: repo, hub: hub}
}

// Notify creates a notification for a user and pushes it via SSE.
func (s *NotificationService) Notify(ctx context.Context, userID int64, channel, notifType, title, message string) error {
	n, err := s.repo.Create(ctx, userID, channel, notifType, title, message)
	if err != nil {
		return err
	}

	// Push to connected SSE clients
	if s.hub != nil {
		s.hub.SendToUser(userID, sse.Event{
			Type: "notification",
			Data: map[string]interface{}{
				"id":      n.ID,
				"channel": n.Channel,
				"type":    n.Type,
				"title":   n.Title,
				"message": n.Message,
				"read":    n.Read,
			},
		})
	}

	return nil
}

// NotifySecurity sends a security notification (always enabled).
func (s *NotificationService) NotifySecurity(ctx context.Context, userID int64, title, message string) error {
	return s.Notify(ctx, userID, "email", "security", title, message)
}

// NotifyLogin sends a login notification.
func (s *NotificationService) NotifyLogin(ctx context.Context, userID int64) {
	_ = s.NotifySecurity(ctx, userID, "New login detected", "A new login to your account was detected.")
}

// NotifyPasswordChange sends a password change notification.
func (s *NotificationService) NotifyPasswordChange(ctx context.Context, userID int64) {
	_ = s.NotifySecurity(ctx, userID, "Password changed", "Your account password was changed.")
}

// NotifyWalletShared sends a notification when a wallet is shared.
func (s *NotificationService) NotifyWalletShared(ctx context.Context, userID int64, walletName string) {
	_ = s.Notify(ctx, userID, "email", "wallet_shared", "Wallet shared with you",
		"You have been granted access to wallet: "+walletName)
}

// NotifyRoleChanged sends a notification when a member's role is changed.
func (s *NotificationService) NotifyRoleChanged(ctx context.Context, userID int64, walletName, newRole string) {
	_ = s.Notify(ctx, userID, "email", "role_changed", "Wallet role updated",
		"Your role on wallet "+walletName+" has been changed to "+newRole)
}

// NotifyAccessRevoked sends a notification when access is revoked.
func (s *NotificationService) NotifyAccessRevoked(ctx context.Context, userID int64, walletName string) {
	_ = s.Notify(ctx, userID, "email", "access_revoked", "Wallet access revoked",
		"Your access to wallet "+walletName+" has been removed.")
}
