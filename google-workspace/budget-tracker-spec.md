# Personal Budget & Expense Tracker - Specification

## Project Info
- Spreadsheet located in Google Drive > workspace folder
- Built with `gws` CLI (googleworkspace/cli v0.5.0)

## Sheet Structure (19 sheets)

### 1. DASHBOARD (first tab)
- Annual overview: total income, expenses, money left
- Spending by category breakdown
- Debt progress summary (per-debt balance vs original, % paid off)
- Month-over-month trends (income, expenses, savings, debt, leftover, cumulative)
- Daily average metrics

### 2. JAN-DEC (12 monthly tabs)
- **Income Summary**: Description, Date, Planned, Actual, Difference
- **Bills Summary**: Description, Due Date, Planned, Actual, Difference
- **Expenses Summary**: Description, Planned, Actual, Difference
- **Subscriptions Summary**: Description, Due Date, Planned, Actual, Difference
- **Debt Summary**: Description, Planned, Actual, Difference
- **Savings Summary**: Description, Planned, Actual, Difference
- **Financial Summary**: Rollover + Income - Bills - Expenses - Debt - Savings - Subscriptions = Total Leftover
- **Financial Insights**: Average Daily Income, Average Daily Spending, Average Daily Bills
- Rollover carries between months

### 3. TRANSACTIONS (single log)
- Fields: Date, Section, Subcategory, Amount, Description, Account
- Section dropdown: Income, Bills, Expenses, Subscriptions, Debt, Savings
- Subcategory dropdown: two-level style ("Bills > Mortgage", "Expenses > Groceries", etc.)
- Account dropdown: configurable list of bank accounts / credit cards
- Primary data entry sheet

### 4. BILL CALENDAR
- Monthly calendar grid (6 weeks)
- Recurring bills checklist with checkboxes for paid/unpaid
- Bill name, due date, amount

### 5. ANNUAL TOTALS
- Top summary: Total Monthly Income, Expenses, Savings, Daily Averages across all months
- Per-section breakdowns: Annual Income, Bills, Expenses, Subscriptions summaries
- Each item shown across JAN-DEC with Planned, Actual, Difference columns

### 6. DEBT TRACKER (snowball method)
- Overview: Current Debt Load, Total Monthly Payments, Total Monthly Interest, Debt Free Date
- Configuration: Starting Month, Extra Monthly Payment
- Per-debt cards: Balance, Min Payment, Interest Rate, Total Interest Owed, Debt Free Date
- 60-month payment projection table per debt (Payment, Interest, Balance)
- Snowball logic: extra payment applies to first debt, cascades when paid off

### 7. SAVINGS GOALS
- Overview: Total Goal Amount, Total Saved, Total Remaining, Overall %
- Goal table: Goal Name, Target Amount, Saved So Far, Remaining, % Complete, Monthly Contribution, Months Until Goal
- Contribution History log: Date, Goal Name, Amount, Notes

### 8. NET WORTH
- Overview: Set Date, Net Worth, Total Assets, Total Liabilities, YTD Progress, % YTD
- Net Worth Summary: monthly grid (START, JAN-DEC) for Net Worth, Total Assets, Total Liabilities
- Total Assets section: Cash & Bank Accounts, Investments, Properties, Other Assets
- Total Liabilities section: per-debt entries matching debt tracker
- YTD Progress and % tracking

## Styling & Design

### Color Theme: Cool & Modern, Light & Clean
- **Background**: White/light gray
- **Income**: Bright blue (#4285F4)
- **Bills**: Slate gray (#5F6368)
- **Expenses**: Indigo (#3F51B5)
- **Debt**: Rose (#E91E63)
- **Savings**: Mint green (#00BFA5)
- **Subscriptions**: Lavender (#7C4DFF)

### Typography
- **Font**: Inter
- **Sheet titles**: 18px bold
- **Section headers**: 14px bold, white text on colored banner
- **Column headers**: 11px bold
- **Data cells**: 10px regular
- **Summary/totals**: 11px bold

### Row Styling
- Alternating row stripes (subtle tint) for readability
- Banding on Transactions and Debt Tracker projection tables

### Alerts
- Over-budget / negative values: Red text (#D32F2F) AND light red background (#FFEBEE)

### General
- Frozen header rows on all sheets
- Currency formatting on money columns
- Date formatting on date columns
- Percentage formatting where applicable

## Technical Notes
- CLI: `gws` (googleworkspace/cli v0.5.0)
- Auth: `gws auth login -s drive,sheets`, then `gws auth export --unmasked > credentials.json`
- Environment variable for API calls:
  ```
  export GOOGLE_WORKSPACE_CLI_CREDENTIALS_FILE=credentials.json
  ```
- Use Python subprocess for complex JSON payloads (avoids shell escaping issues)
- Sheet IDs:
  - DASHBOARD: 1
  - JAN-DEC: 10-21
  - TRANSACTIONS: 100
  - BILL CALENDAR: 101
  - ANNUAL TOTALS: 102
  - DEBT TRACKER: 103
  - SAVINGS GOALS: 104
  - NET WORTH: 105

## Known Items to Refine
- Monthly sheet "Actual" columns not yet wired to TRANSACTIONS via SUMIFS (currently manual entry)
- Bill Calendar is a static template (not auto-generated from dates)
- Debt Tracker Total Interest and Debt Free Date summary cells need formulas
- Annual Totals month columns need formula links to monthly sheets (currently blank)
- Rollover between months not yet linked (each month's rollover cell is manual)
