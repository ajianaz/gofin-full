package domain

import "time"

// WalletMemberRole represents the role of a wallet member.
type WalletMemberRole string

const (
	WalletMemberRoleOwner  WalletMemberRole = "owner"
	WalletMemberRoleEditor WalletMemberRole = "editor"
	WalletMemberRoleViewer WalletMemberRole = "viewer"
)

// WalletMember represents per-wallet access control for sharing.
type WalletMember struct {
	ID        int64      `json:"id" db:"id"`
	WalletID  int64      `json:"wallet_id" db:"wallet_id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	Role      string     `json:"role" db:"role"` // owner, editor, viewer
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`

	// Joined
	User   *User   `json:"user,omitempty" db:"-"`
	Wallet *Wallet `json:"wallet,omitempty" db:"-"`
}

// IsOwner returns true if the member is the wallet owner.
func (m *WalletMember) IsOwner() bool {
	return m.Role == string(WalletMemberRoleOwner)
}

// CanWrite returns true if the member can create/modify transactions.
func (m *WalletMember) CanWrite() bool {
	return m.Role == string(WalletMemberRoleOwner) || m.Role == string(WalletMemberRoleEditor)
}

// CanShare returns true if the member can manage sharing (owner only).
func (m *WalletMember) CanShare() bool {
	return m.Role == string(WalletMemberRoleOwner)
}
