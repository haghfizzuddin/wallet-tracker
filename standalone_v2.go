package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// Simple standalone version of the enhanced tracker

func main() {
	rootCmd := &cobra.Command{
		Use:   "wallet-tracker-v2",
		Short: "Enhanced multi-chain wallet tracker",
	}

	trackCmd := &cobra.Command{
		Use:   "track",
		Short: "Track wallet transactions",
		RunE:  runTrack,
	}

	trackCmd.Flags().StringP("wallet", "w", "", "Wallet address to track")
	trackCmd.Flags().StringP("network", "n", "auto", "Network: BTC, ETH, BSC")
	trackCmd.Flags().BoolP("show-flow", "f", false, "Show fund flow")
	trackCmd.MarkFlagRequired("wallet")

	rootCmd.AddCommand(trackCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runTrack(cmd *cobra.Command, args []string) error {
	wallet, _ := cmd.Flags().GetString("wallet")
	network, _ := cmd.Flags().GetString("network")
	showFlow, _ := cmd.Flags().GetBool("show-flow")

	// Auto-detect network
	if network == "auto" {
		if strings.HasPrefix(wallet, "0x") {
			network = "ETH"
		} else if strings.HasPrefix(wallet, "1") || strings.HasPrefix(wallet, "bc1") {
			network = "BTC"
		}
	}

	// Display header
	displayHeader(wallet, network)

	// Mock transactions for demo
	transactions := getMockTransactions(wallet, network)

	// Display transaction table
	displayTransactionTable(transactions, wallet)

	// Show fund flow if requested
	if showFlow {
		displayFundFlow(transactions, wallet)
	}

	// Display summary
	displaySummary(transactions, wallet, network)

	return nil
}

func displayHeader(wallet, network string) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════════════════════╗")
	color.Cyan("║                    WALLET TRACKER V2.0                           ║")
	color.Cyan("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	color.Yellow("📊 Tracking Wallet: %s", wallet)
	color.Yellow("🌐 Network: %s", network)
	color.Yellow("🕒 Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("─", 70))
}

func displayTransactionTable(txs []Transaction, myWallet string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	
	t.AppendHeader(table.Row{"#", "Type", "Hash", "From → To", "Amount", "USD", "Time", "Status"})
	
	for i, tx := range txs {
		txType := "IN"
		typeColor := text.FgGreen
		if strings.EqualFold(tx.From, myWallet) {
			txType = "OUT"
			typeColor = text.FgRed
		}
		
		row := table.Row{
			i + 1,
			text.Colors{typeColor}.Sprint(txType),
			tx.Hash[:10] + "...",
			fmt.Sprintf("%s → %s", truncate(tx.From), truncate(tx.To)),
			fmt.Sprintf("%.4f %s", tx.Amount, tx.Symbol),
			fmt.Sprintf("$%.2f", tx.USDValue),
			tx.Time.Format("01/02 15:04"),
			"✅",
		}
		
		t.AppendRow(row)
	}
	
	fmt.Println(t.Render())
}

func displayFundFlow(txs []Transaction, wallet string) {
	fmt.Println()
	color.Cyan("💸 Fund Flow Visualization")
	fmt.Println(strings.Repeat("─", 70))
	
	inflows := make(map[string]float64)
	outflows := make(map[string]float64)
	
	for _, tx := range txs {
		if strings.EqualFold(tx.To, wallet) {
			inflows[tx.From] += tx.Amount
		} else {
			outflows[tx.To] += tx.Amount
		}
	}
	
	color.Green("📥 INFLOWS:")
	for addr, amount := range inflows {
		fmt.Printf("  %s %s %.4f\n", truncate(addr), generateBar(amount, 10), amount)
	}
	
	fmt.Println()
	color.Yellow("  ╔════════════════════╗")
	color.Yellow("  ║   YOUR WALLET      ║")
	color.Yellow("  ╚════════════════════╝")
	fmt.Println()
	
	color.Red("📤 OUTFLOWS:")
	for addr, amount := range outflows {
		fmt.Printf("  %s %s %.4f\n", truncate(addr), generateBar(amount, 10), amount)
	}
}

func displaySummary(txs []Transaction, wallet, network string) {
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	color.Cyan("📈 SUMMARY")
	
	var totalIn, totalOut float64
	for _, tx := range txs {
		if strings.EqualFold(tx.To, wallet) {
			totalIn += tx.Amount
		} else {
			totalOut += tx.Amount
		}
	}
	
	fmt.Printf("Total In:  %.4f %s\n", totalIn, network)
	fmt.Printf("Total Out: %.4f %s\n", totalOut, network)
	fmt.Printf("Net Flow:  %.4f %s\n", totalIn-totalOut, network)
}

func truncate(s string) string {
	if len(s) > 10 {
		return s[:6] + "..." + s[len(s)-4:]
	}
	return s
}

func generateBar(value, max float64) string {
	barLength := int((value / max) * 20)
	return strings.Repeat("█", barLength) + strings.Repeat("░", 20-barLength)
}

func getMockTransactions(wallet, network string) []Transaction {
	if network == "ETH" {
		return []Transaction{
			{
				Hash:     "0x1234567890abcdef1234567890abcdef",
				From:     "0x742d35Cc6634C0532925a3b844Bc9e7595f6b8e0",
				To:       wallet,
				Amount:   1.5,
				Symbol:   "ETH",
				Time:     time.Now().Add(-2 * time.Hour),
				USDValue: 3525.0,
			},
			{
				Hash:     "0xabcdef1234567890abcdef1234567890",
				From:     wallet,
				To:       "0x8894E0a0c962CB723c1976a4421c95949bE2D4E3",
				Amount:   0.5,
				Symbol:   "ETH",
				Time:     time.Now().Add(-24 * time.Hour),
				USDValue: 1175.0,
			},
		}
	}
	
	// Default BTC transactions
	return []Transaction{
		{
			Hash:     "a1b2c3d4e5f6789012345678901234567890",
			From:     "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			To:       wallet,
			Amount:   0.1,
			Symbol:   "BTC",
			Time:     time.Now().Add(-1 * time.Hour),
			USDValue: 4325.0,
		},
		{
			Hash:     "f1e2d3c4b5a6978012345678901234567890",
			From:     wallet,
			To:       "3FKj9W2FhYmB6KFdGeWihPpGfNeRHGU4Tz",
			Amount:   0.05,
			Symbol:   "BTC",
			Time:     time.Now().Add(-3 * time.Hour),
			USDValue: 2162.5,
		},
	}
}

type Transaction struct {
	Hash     string
	From     string
	To       string
	Amount   float64
	Symbol   string
	Time     time.Time
	USDValue float64
}
