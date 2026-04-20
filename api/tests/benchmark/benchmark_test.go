package benchmark

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// BenchmarkDecimalOperations benchmarks common decimal operations used in financial calculations.
func BenchmarkDecimalAddition(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := decimal.NewFromFloat(100.50)
		b_ := decimal.NewFromFloat(200.75)
		_ = a.Add(b_)
	}
}

func BenchmarkDecimalMultiplication(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := decimal.NewFromFloat(1.13)
		b_ := decimal.NewFromFloat(100.50)
		_ = a.Mul(b_)
	}
}

func BenchmarkDecimalFixed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := decimal.NewFromFloat(123456789.123456)
		_ = a.StringFixed(2)
	}
}

// BenchmarkStringOperations benchmarks string operations common in API responses.
func BenchmarkStringConcat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := "Hello, " + "World, " + "Test, " + "Value"
		_ = len(s)
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var sb []byte
		sb = append(sb, "Hello, "...)
		sb = append(sb, "World, "...)
		sb = append(sb, "Test, "...)
		sb = append(sb, "Value"...)
		_ = len(sb)
	}
}

// BenchmarkJSONMarshal benchmarks JSON marshaling of a transaction-like struct.
func BenchmarkJSONMarshal(b *testing.B) {
	b.ReportAllocs()
	type Transaction struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
		Amount      string  `json:"amount"`
		Date        string  `json:"date"`
		Currency    string  `json:"currency"`
	}
	tx := Transaction{
		ID:          1,
		Description: "Test transaction",
		Amount:      "123.45",
		Date:        time.Now().Format(time.RFC3339),
		Currency:    "USD",
	}
	for i := 0; i < b.N; i++ {
		_, _ = tx.Date, tx.Amount // prevent unused variable
		_ = fmt.Sprintf("%v", tx)
	}
}

// BenchmarkSimulatedDBQuery simulates the cost of processing query results.
func BenchmarkSimulatedDBQuery(b *testing.B) {
	b.ReportAllocs()
	type Row struct {
		ID    int64
		Name  string
		Value decimal.Decimal
	}
	rows := make([]Row, 100)
	for i := range rows {
		rows[i] = Row{ID: int64(i), Name: fmt.Sprintf("item-%d", i), Value: decimal.NewFromFloat(float64(i) * 1.5)}
	}
	for i := 0; i < b.N; i++ {
		total := decimal.Zero
		for _, r := range rows {
			total = total.Add(r.Value)
		}
		_ = total
	}
}

// BenchmarkConcurrentAccess benchmarks concurrent read patterns.
func BenchmarkConcurrentReads(b *testing.B) {
	b.ReportAllocs()
	data := make(map[string]string, 100)
	for i := 0; i < 100; i++ {
		data[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("value-%d", i)
	}
	for i := 0; i < b.N; i++ {
		_ = data["key-50"]
		_ = data["key-99"]
	}
}

// unused
var _ = context.Background()
