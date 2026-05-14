package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ajianaz/gofin-full/api/internal/domain"
	"github.com/ajianaz/gofin-full/api/internal/repository"
)

type ExportService struct {
	txRepo *repository.TransactionRepository
}

func NewExportService(txRepo *repository.TransactionRepository) *ExportService {
	return &ExportService{txRepo: txRepo}
}

// CSVRow represents a single row in a CSV export.
type CSVRow struct {
	Date               string
	Type               string
	Description        string
	Amount             string
	Currency           string
	Category           string
	SourceAccount      string
	DestinationAccount string
	Notes              string
	Tags               string
}

// ExportTransactionsCSV exports transactions as CSV.
func (s *ExportService) ExportTransactionsCSV(ctx context.Context, groupID uuid.UUID, w io.Writer, filter repository.TransactionFilter) error {
	groups, _, err := s.txRepo.ListGroups(ctx, groupID, filter)
	if err != nil {
		return fmt.Errorf("failed to list transactions: %w", err)
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	header := []string{"Date", "Type", "Description", "Amount", "Currency", "Category", "Source account", "Destination account", "Notes", "Tags"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, g := range groups {
		for _, j := range g.Journals {
			amount := getJournalAmount(j)
			sourceName := walletName(j.Source)
			destName := walletName(j.Destination)
			categoryName := ""
			if len(j.Categories) > 0 {
				categoryName = j.Categories[0].Name
			}

			record := []string{
				j.Date.Format("2006-01-02"),
				string(j.Type),
				sanitizeCSV(j.Description),
				amount.StringFixed(2),
				j.CurrencyID,
				sanitizeCSV(categoryName),
				sanitizeCSV(sourceName),
				sanitizeCSV(destName),
				sanitizeCSV(coalesceStr(j.Notes)),
				sanitizeCSV(joinTagNames(j.Tags)),
			}
			if err := csvWriter.Write(record); err != nil {
				return err
			}
		}
	}

	return nil
}

// ExportTransactionsOFX exports transactions in OFX (Open Financial Exchange) format.
func (s *ExportService) ExportTransactionsOFX(ctx context.Context, groupID uuid.UUID, w io.Writer, filter repository.TransactionFilter) error {
	groups, _, err := s.txRepo.ListGroups(ctx, groupID, filter)
	if err != nil {
		return fmt.Errorf("failed to list transactions: %w", err)
	}

	now := time.Now().UTC()
	fmt.Fprintf(w, "OFXHEADER:100\nDATA:OFXSGML\nVERSION:102\nSECURITY:NONE\nENCODING:USASCII\nCHARSET:1252\nCOMPRESSION:NONE\nOLDFILEUID:NONE\nNEWFILEUID:NONE\n\n")
	fmt.Fprintf(w, "<OFX>\n<SIGNONMSGSRSV1>\n<SONRS>\n<STATUS>\n<CODE>0\n<SEVERITY>INFO\n</STATUS>\n<DTSERVER>%s\n<LANGUAGE>ENG\n</SONRS>\n</SIGNONMSGSRSV1>\n", now.Format("20060102150405"))

	fmt.Fprintf(w, "<BANKMSGSRSV1>\n<STMTTRNRS>\n<TRNUID>1\n<STATUS>\n<CODE>0\n<SEVERITY>INFO\n</STATUS>\n<STMTRS>\n<CURDEF>USD\n<BANKACCTFROM>\n<ACCTID>Gofin\n<ACCTTYPE>CHECKING\n</BANKACCTFROM>\n<BANKTRANLIST>\n<DTSTART>%s\n<DTEND>%s\n",
		now.Format("20060102150405"), now.Format("20060102150405"))

	for _, g := range groups {
		for _, j := range g.Journals {
			amount := getJournalAmount(j)
			trnType := "CREDIT"
			if amount.IsNegative() {
				trnType = "DEBIT"
				amount = amount.Abs()
			}
			fmt.Fprintf(w, "<STMTTRN>\n<TRNTYPE>%s\n<DTPOSTED>%s\n<TRNAMT>%s\n<FITID>%d\n<NAME>%s\n</STMTTRN>\n",
				trnType, j.Date.Format("20060102"), amount.StringFixed(2), j.ID, escapeXML(j.Description))
		}
	}

	fmt.Fprintf(w, "</BANKTRANLIST>\n</STMTRS>\n</STMTTRNRS>\n</BANKMSGSRSV1>\n</OFX>\n")
	return nil
}

// ReconciliationResult represents the result of a reconciliation check.
type ReconciliationResult struct {
	Matched      int `json:"matched"`
	Unmatched    int `json:"unmatched"`
	TotalChecked int `json:"total_checked"`
}

// Reconcile matches imported transactions against existing records by date, amount, and description.
func (s *ExportService) Reconcile(ctx context.Context, groupID uuid.UUID, imports []CSVRow) (*ReconciliationResult, error) {
	existing, _, err := s.txRepo.ListGroups(ctx, groupID, repository.TransactionFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to list existing transactions: %w", err)
	}

	result := &ReconciliationResult{TotalChecked: len(imports)}

	for _, imp := range imports {
		amount, _ := decimal.NewFromString(imp.Amount)
		date, _ := time.Parse("2006-01-02", imp.Date)

		found := false
		for _, g := range existing {
			for _, j := range g.Journals {
				if j.Date.Format("2006-01-02") == date.Format("2006-01-02") {
					jAmount := getJournalAmount(j)
					if jAmount.Equal(amount) && j.Description == imp.Description {
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}

		if found {
			result.Matched++
		} else {
			result.Unmatched++
		}
	}

	return result, nil
}

func getJournalAmount(j domain.TransactionJournal) decimal.Decimal {
	// Source transactions are debits (negative for withdrawals)
	if len(j.SourceTransactions) > 0 {
		return j.SourceTransactions[0].Amount
	}
	return decimal.Zero
}

func walletName(w *domain.Wallet) string {
	if w == nil {
		return ""
	}
	return w.Name
}

func coalesceStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func joinTagNames(tags []domain.Tag) string {
	result := ""
	for i, t := range tags {
		if i > 0 {
			result += ","
		}
		result += t.Tag
	}
	return result
}

func sanitizeCSV(s string) string {
	if len(s) == 0 {
		return s
	}
	trimmed := strings.TrimSpace(s)
	if len(trimmed) == 0 {
		return s
	}
	c := trimmed[0]
	if c == '=' || c == '+' || c == '-' || c == '@' || c == '\t' || c == '\r' {
		return "'" + s + "'"
	}
	return s
}

func escapeXML(s string) string {
	result := make([]byte, 0, len(s))
	for _, c := range s {
		switch c {
		case '&':
			result = append(result, '&', 'a', 'm', 'p', ';')
		case '<':
			result = append(result, '&', 'l', 't', ';')
		case '>':
			result = append(result, '&', 'g', 't', ';')
		case '"':
			result = append(result, '&', 'q', 'u', 'o', 't', ';')
		case '\'':
			result = append(result, '&', 'a', 'p', 'o', 's', ';')
		default:
			result = append(result, byte(c))
		}
	}
	return string(result)
}
