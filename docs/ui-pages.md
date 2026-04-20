# Gofin UI Pages / Screens

Complete list of pages needed for the Gofin frontend, derived from the API endpoints (125+), business flows (14), and research documents (22).

---

## 1. Authentication (4 screens)

| # | Screen | Description |
|---|--------|-------------|
| 1 | Login | Email/password form, OAuth buttons (Google, GitHub), "Forgot password" link |
| 2 | Register | Email, password, confirm password, optional invite code |
| 3 | Forgot Password | Email input to request reset link |
| 4 | Reset Password | New password form with token |

## 2. Onboarding (2 screens)

| # | Screen | Description |
|---|--------|-------------|
| 5 | Setup Wizard | Step 1: Language, Step 2: Currency + bank account creation, Step 3: Optional savings account |
| 6 | 2FA Setup | QR code display, backup codes download, verification input |

## 3. Main Shell (1 template)

| # | Screen | Description |
|---|--------|-------------|
| 7 | App Shell | Sidebar navigation, top bar (user avatar, group switcher, notifications bell), main content area |

## 4. Dashboard (1 screen)

| # | Screen | Description |
|---|--------|-------------|
| 8 | Dashboard | Account balance cards (from frontpageAccounts preference), recent transactions, spending summary, bill reminders, quick-add transaction |

## 5. Transactions (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 9 | Transaction List | Paginated table with filters (type, date range, account, category, budget), sort, search |
| 10 | Transaction Detail | Full view: source/destination, amounts (original + native), tags, categories, budget, bill, notes, attachments, audit history |
| 11 | Transaction Create/Edit | Form: type selector, date, description, amount, source/destination account, category, budget, bill, piggy bank, tags, notes, foreign currency toggle; split transaction support |

## 6. Wallets / Accounts (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 12 | Wallet List | Grid/list of all wallets with balances, type filters (asset, cash, credit card, loan, etc.) |
| 13 | Wallet Detail | Balance overview, linked transactions, piggy banks, notes, attachments, member sharing |
| 14 | Wallet Create/Edit | Form: name, type, currency, IBAN, opening balance, virtual balance, liability fields, credit card fields, include in net worth toggle |

## 7. Budgets (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 15 | Budget List | All budgets with spent/remaining indicators, period selector, auto-budget badge |
| 16 | Budget Detail | Budget limits per period, available budget calculation, linked transactions |
| 17 | Budget Create/Edit | Form: name, active toggle, auto-budget type (none/reset/rollover/adjusted), amount, period |

## 8. Piggy Banks (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 18 | Piggy Bank List | Per-account list with progress bars, target amounts, current saved amounts |
| 19 | Piggy Bank Detail | Progress chart (savings over time), add/remove money form, event history |
| 20 | Piggy Bank Create/Edit | Form: name, target amount, linked account, start date, target date |

## 9. Bills (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 21 | Bill List | All bills with next expected date, amount range, overdue indicator |
| 22 | Bill Detail | Pay date calendar, linked transactions, warning status |
| 23 | Bill Create/Edit | Form: name, amount min/max, currency, start date, repeat frequency, skip, end date |

## 10. Recurring Transactions (2 screens)

| # | Screen | Description |
|---|--------|-------------|
| 24 | Recurrence List | All recurring transactions with next occurrence, repeat type, active status |
| 25 | Recurrence Create/Edit | Form: title, repetition type, repeat moment, skip, end date, transaction templates |

## 11. Categories (2 screens)

| # | Screen | Description |
|---|--------|-------------|
| 26 | Category List | All categories with transaction counts |
| 27 | Category Create/Edit | Form: name, parent category (optional) |

## 12. Tags (2 screens)

| # | Screen | Description |
|---|--------|-------------|
| 28 | Tag List | All tags with usage counts |
| 29 | Tag Create/Edit | Form: tag text, date, description |

## 13. Rules Engine (4 screens)

| # | Screen | Description |
|---|--------|-------------|
| 30 | Rule Group List | All rule groups with order, active status, rule count |
| 31 | Rule Group Create/Edit | Form: title, active toggle, order, stop processing flag |
| 32 | Rule List | All rules within a group with trigger/action summary |
| 33 | Rule Create/Edit | Form: title, strict mode, stop processing, trigger builder (operator + value), action builder (type + value), test button |

## 14. Analytics / Reports (4 screens)

| # | Screen | Description |
|---|--------|-------------|
| 34 | Spending by Category | Pie chart with category breakdown, date range picker |
| 35 | Spending by Period | Bar/line chart with period selector (day/week/month/year) |
| 36 | Net Worth | Line chart over time, assets vs liabilities breakdown |
| 37 | Reports Dashboard | Summary cards (income, expenses, net worth), chart navigation |

## 15. Groups & Wallet Sharing (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 38 | Group List | All groups user belongs to, current group indicator |
| 39 | Group Switcher | Dropdown/modal to switch active group context |
| 40 | Wallet Members | Member list with roles, invite form (user selector + role), role update, remove |

## 16. Notifications (1 screen)

| # | Screen | Description |
|---|--------|-------------|
| 41 | Notification Panel | Dropdown/panel with unread list, mark read/mark all, real-time SSE updates |

## 17. User Settings (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 42 | Profile Settings | Email, name, password change |
| 43 | Preferences | List of user preferences (currency, view settings, front page accounts, notification toggles, 2FA) |
| 44 | API Keys | List of personal access tokens, create new key, delete |

## 18. Currencies (2 screens)

| # | Screen | Description |
|---|--------|-------------|
| 45 | Currency List | All available currencies with codes, symbols |
| 46 | Exchange Rates | Exchange rate table, add/delete rates |

## 19. Admin (2 screens)

| # | Screen | Description |
|---|--------|-------------|
| 47 | Admin User Management | System-wide user list, create user, admin role assignment |
| 48 | Admin Feature Flags | Toggle features (export, webhooks, debts, expression engine, running balance) |

## 20. Export & Audit (3 screens)

| # | Screen | Description |
|---|--------|-------------|
| 49 | Export | Format selection (CSV/OFX), date range, account filter, download |
| 50 | Reconciliation | Account selector, date range, cleared transaction selection, difference calculation |
| 51 | Audit Log | Table of transaction changes (who, what, when, before/after), filter by action/user |

## 21. Inline Components (3 components)

| # | Component | Description |
|---|----------|-------------|
| 52 | Attachment Manager | Upload area, file list with download/delete, 2-step upload, preview for images/PDFs |
| 53 | Notes Editor | Markdown/text editor for notes on any entity (polymorphic) |
| 54 | Object Groups | Grouping/ordering containers for accounts, bills, piggy banks |

---

## Summary

- **Total screens: 54**
- **Design status**: Screen #1 (Login) has a design in Pencil
- **Priority order**: Login > Register > Dashboard > Transaction List/Create > Wallet List/Create > Budget > Piggy Bank > Analytics > Settings > Admin
- **RBAC consideration**: All screens must respect the 21 group-level roles and 3 wallet-level roles for visibility and action permissions
- **Multi-currency**: Every monetary value needs original + primary currency conversion display
- **Real-time**: Notification panel uses SSE for live updates
- **Polymorphic**: Attachment, Note, Location components are reusable across entity types
