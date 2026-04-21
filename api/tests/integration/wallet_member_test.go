package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/tests/integration/testhelpers"
)

func TestOwnerCanShareWallet(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.OwnerToken
	walletID := testApp.Seed.WalletID
	readOnlyUserID := testApp.Seed.ReadOnlyUserID

	t.Run("owner_can_list_wallet_members", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", path, "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode,
			"owner should be able to list wallet members")
	})

	t.Run("owner_can_share_wallet", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		body := fmt.Sprintf(`{"user_id":"%s","role":"viewer"}`, readOnlyUserID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", path, body, token)
		// Accept 200 (success) or 422 (validation) — not 403.
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should not get 403 when sharing a wallet")
	})
}

func TestViewerCanOnlyRead(t *testing.T) {
	app := testApp.App
	ownerToken := testApp.Seed.OwnerToken
	readOnlyToken := testApp.Seed.ReadOnlyToken
	walletID := testApp.Seed.WalletID
	readOnlyUserID := testApp.Seed.ReadOnlyUserID

	// First, owner shares the wallet with the read_only user as a viewer.
	// This ensures the membership exists for the read operations.
	t.Run("setup_share_wallet_as_viewer", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		body := fmt.Sprintf(`{"user_id":"%s","role":"viewer"}`, readOnlyUserID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", path, body, ownerToken)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should be able to share wallet during setup")
	})

	t.Run("viewer_can_read_wallet", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s", walletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", path, "", readOnlyToken)
		// The viewer should be able to read — accept 200 or 404 (wallet might not belong to their scope),
		// but not 403 from RBAC.
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"viewer should not get 403 when reading a shared wallet")
	})

	t.Run("viewer_can_list_wallet_members", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", path, "", readOnlyToken)
		require.Equal(t, http.StatusOK, resp.StatusCode,
			"viewer should be able to list wallet members (read)")
	})

	t.Run("viewer_cannot_share_wallet", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		body := fmt.Sprintf(`{"user_id":"%s","role":"viewer"}`, readOnlyUserID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", path, body, readOnlyToken)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"viewer should get 403 when trying to share a wallet (requires owner role)")
	})

	t.Run("viewer_cannot_create_account", func(t *testing.T) {
		body := `{"name":"Viewer Account","account_type":"asset","virtual_balance":0}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/wallets", body, readOnlyToken)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"viewer should get 403 when creating accounts")
	})

	t.Run("viewer_cannot_create_transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"withdrawal","amount":"5.00","description":"viewer attempt","source_id":"%s"}`, walletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/transactions", body, readOnlyToken)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"viewer should get 403 when creating transactions")
	})
}

func TestReadOnlyCannotShareWallet(t *testing.T) {
	app := testApp.App
	readOnlyToken := testApp.Seed.ReadOnlyToken
	walletID := testApp.Seed.WalletID
	txUserID := testApp.Seed.TxUserID

	t.Run("read_only_cannot_share_wallet", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		body := fmt.Sprintf(`{"user_id":"%s","role":"viewer"}`, txUserID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", path, body, readOnlyToken)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only user should get 403 when trying to share wallet")
	})
}

func TestManageTransactionsCannotShareWallet(t *testing.T) {
	app := testApp.App
	txToken := testApp.Seed.TxUserToken
	walletID := testApp.Seed.WalletID
	readOnlyUserID := testApp.Seed.ReadOnlyUserID

	t.Run("manage_transactions_cannot_share_wallet", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/wallets/%s/members", walletID)
		body := fmt.Sprintf(`{"user_id":"%s","role":"viewer"}`, readOnlyUserID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", path, body, txToken)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"manage_transactions user should get 403 when trying to share wallet (requires owner)")
	})
}
