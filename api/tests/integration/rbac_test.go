package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/tests/integration/testhelpers"
)

func TestOwnerFullAccess(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.OwnerToken

	t.Run("owner_can_read_users_me", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/users/me", "", token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should not get 403 on users/me")
	})

	t.Run("owner_can_list_accounts", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/wallets", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("owner_can_create_account", func(t *testing.T) {
		body := `{"name":"Owner Savings","account_type":"asset","virtual_balance":0,"include_net_worth":true,"active":true}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/wallets", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should not get 403 when creating accounts")
	})

	t.Run("owner_can_list_transactions", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/transactions", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("owner_can_create_transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"withdrawal","amount":"10.00","description":"test tx","source_id":%d}`, testApp.Seed.WalletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/transactions", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should not get 403 when creating transactions")
	})

	t.Run("owner_can_list_budgets", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/budgets", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("owner_can_create_budget", func(t *testing.T) {
		body := `{"name":"Test Budget","amount":"1000.00"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/budgets", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should not get 403 when creating budgets")
	})

	t.Run("owner_can_list_groups", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/groups", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("owner_can_view_analytics", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/analytics/net-worth", "", token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"owner should not get 403 on analytics endpoints")
	})
}

func TestReadOnlyCanRead(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.ReadOnlyToken

	t.Run("read_only_can_list_accounts", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/wallets", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("read_only_can_list_transactions", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/transactions", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("read_only_can_list_budgets", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/budgets", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("read_only_can_list_categories", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/categories", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("read_only_can_read_users_me", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/users/me", "", token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"read_only should not get 403 on users/me")
	})
}

func TestReadOnlyCannotWrite(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.ReadOnlyToken

	t.Run("read_only_cannot_create_account", func(t *testing.T) {
		body := `{"name":"Blocked Account","account_type":"asset","virtual_balance":0}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/wallets", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only should get 403 when creating accounts")
	})

	t.Run("read_only_cannot_create_transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"withdrawal","amount":"10.00","description":"blocked","source_id":%d}`, testApp.Seed.WalletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/transactions", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only should get 403 when creating transactions")
	})

	t.Run("read_only_cannot_create_budget", func(t *testing.T) {
		body := `{"name":"Blocked Budget","amount":"500.00"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/budgets", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only should get 403 when creating budgets")
	})

	t.Run("read_only_cannot_create_tag", func(t *testing.T) {
		body := `{"tag":"blocked-tag"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/tags", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only should get 403 when creating tags")
	})

	t.Run("read_only_cannot_view_analytics", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/analytics/net-worth", "", token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"read_only should get 403 on analytics (requires view_reports)")
	})
}

func TestManageTransactionsCanCreateTx(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.TxUserToken

	t.Run("manage_transactions_can_list_transactions", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/transactions", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("manage_transactions_can_create_transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"withdrawal","amount":"25.00","description":"tx role test","source_id":%d}`, testApp.Seed.WalletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/transactions", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"manage_transactions should not get 403 when creating transactions")
	})

	t.Run("manage_transactions_can_list_bills", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/bills", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("manage_transactions_can_create_bill", func(t *testing.T) {
		body := `{"name":"Test Bill","amount_min":"10.00","amount_max":"10.00"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/bills", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"manage_transactions should not get 403 when creating bills")
	})
}

func TestManageTransactionsCannotCreateBudget(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.TxUserToken

	t.Run("manage_transactions_cannot_create_budget", func(t *testing.T) {
		body := `{"name":"Blocked Budget","amount":"500.00"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/budgets", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"manage_transactions should get 403 when creating budgets (requires manage_budgets)")
	})

	t.Run("manage_transactions_cannot_create_account", func(t *testing.T) {
		body := `{"name":"Blocked Account","account_type":"asset","virtual_balance":0}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/wallets", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"manage_transactions should get 403 when creating accounts (requires manage_meta)")
	})

	t.Run("manage_transactions_cannot_create_tag", func(t *testing.T) {
		body := `{"tag":"blocked-tag"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/tags", body, token)
		require.Equal(t, http.StatusForbidden, resp.StatusCode,
			"manage_transactions should get 403 when creating tags (requires manage_meta)")
	})
}

func TestFullRoleAccess(t *testing.T) {
	app := testApp.App
	token := testApp.Seed.FullToken

	t.Run("full_can_list_accounts", func(t *testing.T) {
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "GET", "/api/v1/wallets", "", token)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("full_can_create_account", func(t *testing.T) {
		body := `{"name":"Full User Account","account_type":"asset","virtual_balance":0,"include_net_worth":true,"active":true}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/wallets", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"full role should not get 403 on any write operation")
	})

	t.Run("full_can_create_transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"withdrawal","amount":"50.00","description":"full role tx","source_id":%d}`, testApp.Seed.WalletID)
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/transactions", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"full role should not get 403 when creating transactions")
	})

	t.Run("full_can_create_budget", func(t *testing.T) {
		body := `{"name":"Full Budget","amount":"2000.00"}`
		resp := testhelpers.MakeAuthenticatedRequest(t, app, "POST", "/api/v1/budgets", body, token)
		require.NotEqual(t, http.StatusForbidden, resp.StatusCode,
			"full role should not get 403 when creating budgets")
	})
}
