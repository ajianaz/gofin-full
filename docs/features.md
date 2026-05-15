# Features

Comprehensive overview of Gofin's capabilities.

## 💰 Financial Tracking

### Wallets (Accounts)
Multi-type financial containers that serve as the core of your finance tracking.

| Wallet Type | Use Case |
|------------|----------|
| Bank Account | Checking, savings, deposit accounts |
| Cash | Physical cash tracking |
| Credit Card | Credit card with liability tracking |
| Loan | Loans with balance tracking |
| Investment | Investment portfolios |
| Other | Custom account types |

Each wallet supports:
- Multi-currency with exchange rate conversion
- Opening balance and virtual balance
- IBAN and account number fields
- Notes and attachments
- Include/exclude from net worth calculation
- Grouping via object groups

### Transactions
Full double-entry bookkeeping system. Every transaction has a **source** and **destination** wallet.

- **Transaction types:** Deposit, withdrawal, transfer, opening balance
- **Split transactions** — Divide a single transaction across multiple categories
- **Tags** — Flexible labeling system
- **Categories** — Hierarchical categorization (parent/child)
- **Attachments** — Upload receipts, invoices (multi-file support)
- **Notes** — Markdown notes on any transaction
- **Foreign currency** — Automatic conversion to primary currency
- **Search & filter** — By type, date range, account, category, budget, tag

### Budgets
Set spending limits per period to stay on track.

- **Budget limits** per period (daily, weekly, monthly, yearly, custom)
- **Auto-budget modes:** None, Reset, Rollover, Adjusted
- **Spent/remaining indicators** with visual progress
- **Linked transactions** — See which transactions count toward a budget

### Piggy Banks
Savings goals linked to a wallet account.

- **Target amount** with progress tracking
- **Add/remove money** events with history
- **Start date and target date**
- Progress visualization

### Bills
Recurring bills with expected dates and amount ranges.

- **Min/max amount range** for variable bills
- **Repeat frequency** (weekly, monthly, quarterly, yearly)
- **Next expected date** calculation
- **Warning status** for overdue bills
- **Skip** and end date support

### Recurring Transactions
Automate repetitive transactions.

- **Repetition types:** Daily, weekly, monthly, quarterly, yearly
- **Custom repeat moment** (e.g., "every 3rd Friday")
- **Skip** certain occurrences
- **End date** support
- **Transaction templates** — Pre-fill amount, description, and accounts

## 📊 Analytics & Reports

| Report | Description |
|--------|-------------|
| **Spending by Category** | Pie chart with category breakdown over a date range |
| **Spending by Period** | Bar/line chart with period selector (day/week/month/year) |
| **Net Worth** | Line chart tracking assets vs liabilities over time |
| **Dashboard** | Summary cards (income, expenses, net worth) with quick-add |

## 🔄 Automation

### Rules Engine
Powerful automation based on transaction triggers.

- **Rule groups** — Organize rules with ordering and active/inactive status
- **Triggers** — Match transactions by field, operator, and value (e.g., "description contains 'STARBUCKS'")
- **Actions** — Auto-categorize, set tags, modify values
- **Strict mode** — Stop processing if rule doesn't match
- **Test button** — Preview rule results before activating

### Export
Download your financial data in standard formats.

- **CSV** — Comma-separated values for spreadsheet import
- **OFX** — Open Financial Exchange for bank software
- **Filters** — Date range and account selection

## 👥 Collaboration

### Groups
Multi-user support with group-level access control.

- Switch between groups via the group switcher
- Each group has its own set of wallets, transactions, budgets
- Group-level roles (see [RBAC & Permissions](/rbac))

### Wallet Sharing
Share specific wallets with other users in your group.

- **Owner** — Full access (wallet creator)
- **Editor** — Can create and modify transactions
- **Viewer** — Read-only access

## 🔔 Notifications
Real-time updates via Server-Sent Events (SSE).

- Unread notification panel
- Mark read / mark all as read
- Live updates without page refresh

## 🌐 Internationalization

| Feature | Details |
|---------|---------|
| **Languages** | Indonesian (Bahasa Indonesia), English |
| **Dark mode** | Toggle or follow system preference |
| **Currency** | 150+ currencies with exchange rate support |
| **Responsive** | Works on desktop, tablet, and mobile |

## 🔑 Authentication & Security

See the [Security](/security) page for full details.

- JWT with refresh tokens
- OAuth2 (Google, GitHub)
- Optional Keycloak OIDC integration
- Password strength policy (min 8 chars, 3 of 4 types)
- Login lockout (email + IP keyed)
- Rate limiting (configurable)
- HSTS in production
- 2FA support (TOTP)
