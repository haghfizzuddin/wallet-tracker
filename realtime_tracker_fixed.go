package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "wallet-tracker-real",
		Short: "Real-time blockchain wallet tracker",
	}

	trackCmd := &cobra.Command{
		Use:   "track",
		Short: "Track wallet transactions from blockchain",
		RunE:  runTrack,
	}

	trackCmd.Flags().StringP("wallet", "w", "", "Wallet address to track")
	trackCmd.Flags().StringP("network", "n", "auto", "Network: BTC, ETH, BSC")
	trackCmd.Flags().IntP("limit", "l", 10, "Number of transactions to show")
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
	limit, _ := cmd.Flags().GetInt("limit")
	showFlow, _ := cmd.Flags().GetBool("show-flow")

	// Auto-detect network
	if network == "auto" {
		if strings.HasPrefix(wallet, "0x") {
			network = "ETH"
		} else if strings.HasPrefix(wallet, "1") || strings.HasPrefix(wallet, "bc1") || strings.HasPrefix(wallet, "3") {
			network = "BTC"
		}
	}

	// Display header
	displayHeader(wallet, network)

	// Get real transactions
	fmt.Println("🔄 Fetching real blockchain data...")
	transactions, err := getRealBTCTransactions(wallet, limit)
	if err != nil {
		color.Red("❌ Error fetching transactions: %v", err)
		color.Yellow("Using demo data for illustration...")
		transactions = getMockTransactions(wallet, network)
	}

	if len(transactions) == 0 {
		color.Yellow("⚠️  No transactions found for this wallet")
		return nil
	}

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

// Fetch real Bitcoin transactions
func getRealBTCTransactions(wallet string, limit int) ([]Transaction, error) {
	// Using blockchain.info API (free, no key required)
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s?limit=%d", wallet, limit)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var btcResp BTCAddressResponse
	if err := json.Unmarshal(body, &btcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Get current BTC price
	btcPrice := getBTCPrice()

	transactions := make([]Transaction, 0)
	
	for _, tx := range btcResp.Txs {
		// For each transaction, determine if it's incoming or outgoing
		isIncoming := false
		var amount int64 = 0
		var from, to string
		var fee int64 = tx.Fee
		
		// Check all outputs to find ones related to our wallet
		for _, out := range tx.Out {
			if out.Addr == wallet {
				// This is incoming to our wallet
				isIncoming = true
				amount += out.Value
				to = wallet
				// Get the sender (first input address that's not ours)
				if len(tx.Inputs) > 0 && tx.Inputs[0].PrevOut.Addr != "" {
					from = tx.Inputs[0].PrevOut.Addr
				}
			}
		}
		
		// If not incoming, check if it's outgoing
		if !isIncoming {
			for _, input := range tx.Inputs {
				if input.PrevOut.Addr == wallet {
					// This is outgoing from our wallet
					from = wallet
					// Find the main recipient (largest output that's not change)
					var maxOut int64 = 0
					for _, out := range tx.Out {
						if out.Addr != wallet && out.Value > maxOut {
							to = out.Addr
							amount = out.Value
							maxOut = out.Value
						}
					}
					break
				}
			}
		}
		
		// Skip if this transaction doesn't involve our wallet
		if amount == 0 {
			continue
		}
		
		// Convert satoshis to BTC
		btcAmount := float64(amount) / 100000000
		feeAmount := float64(fee) / 100000000
		
		transactions = append(transactions, Transaction{
			Hash:          tx.Hash,
			From:          from,
			To:            to,
			Amount:        btcAmount,
			Fee:           feeAmount,
			Symbol:        "BTC",
			Time:          time.Unix(tx.Time, 0),
			BlockHeight:   tx.BlockHeight,
			Confirmations: btcResp.NTx - tx.BlockIndex,
			USDValue:      btcAmount * btcPrice,
		})
	}
	
	return transactions, nil
}

// Get current BTC price
func getBTCPrice() float64 {
	url := "https://api.coinbase.com/v2/exchange-rates?currency=BTC"
	
	resp, err := http.Get(url)
	if err != nil {
		return 43000.0 // Fallback price
	}
	defer resp.Body.Close()
	
	var priceResp struct {
		Data struct {
			Rates map[string]string `json:"rates"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return 43000.0 // Fallback price
	}
	
	if usdRate, ok := priceResp.Data.Rates["USD"]; ok {
		var price float64
		fmt.Sscanf(usdRate, "%f", &price)
		return price
	}
	
	return 43000.0 // Fallback price
}

func displayHeader(wallet, network string) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════════════════════╗")
	color.Cyan("║              REAL-TIME WALLET TRACKER V2.0                       ║")
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
	
	t.AppendHeader(table.Row{"#", "Type", "Time", "From → To", "Amount", "USD Value", "Confirms", "Fee"})
	
	for i, tx := range txs {
		txType := "IN"
		typeColor := text.FgGreen
		if strings.EqualFold(tx.From, myWallet) {
			txType = "OUT"
			typeColor = text.FgRed
		}
		
		// Format time
		timeStr := tx.Time.Format("01/02 15:04")
		if time.Since(tx.Time) < 24*time.Hour {
			timeStr = tx.Time.Format("15:04")
		}
		
		// Format addresses - show "TRACKED" instead of "YOU"
		fromAddr := truncate(tx.From)
		toAddr := truncate(tx.To)
		if tx.From == myWallet {
			fromAddr = "[TRACKED]"
		}
		if tx.To == myWallet {
			toAddr = "[TRACKED]"
		}
		
		// Confirmation status
		confirmStatus := fmt.Sprintf("%d", tx.Confirmations)
		if tx.Confirmations == 0 {
			confirmStatus = "⏳ 0"
		} else if tx.Confirmations < 6 {
			confirmStatus = fmt.Sprintf("⚡ %d", tx.Confirmations)
		} else {
			confirmStatus = fmt.Sprintf("✅ %d", tx.Confirmations)
		}
		
		row := table.Row{
			i + 1,
			text.Colors{typeColor}.Sprint(txType),
			timeStr,
			fmt.Sprintf("%s → %s", fromAddr, toAddr),
			fmt.Sprintf("%.8f %s", tx.Amount, tx.Symbol),
			fmt.Sprintf("$%.2f", tx.USDValue),
			confirmStatus,
			fmt.Sprintf("%.8f", tx.Fee),
		}
		
		t.AppendRow(row)
		
		// Add separator for readability
		if i < len(txs)-1 {
			t.AppendSeparator()
		}
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
	
	// Find max for scaling
	maxIn := 0.0
	for _, v := range inflows {
		if v > maxIn {
			maxIn = v
		}
	}
	maxOut := 0.0
	for _, v := range outflows {
		if v > maxOut {
			maxOut = v
		}
	}
	
	if len(inflows) > 0 {
		color.Green("📥 INFLOWS:")
		for addr, amount := range inflows {
			bar := generateBar(amount, maxIn)
			fmt.Printf("  %-20s %s %.8f BTC ($%.2f)\n", 
				truncate(addr), bar, amount, amount*getBTCPrice())
		}
	}
	
	// Display wallet box properly centered
	fmt.Println()
	color.Yellow("                    ╔════════════════════════════╗")
	color.Yellow("                    ║      TRACKED WALLET        ║")
	color.Yellow("                    ║   %s   ║", centerText(truncate(wallet), 24))
	color.Yellow("                    ╚════════════════════════════╝")
	fmt.Println()
	
	if len(outflows) > 0 {
		color.Red("📤 OUTFLOWS:")
		for addr, amount := range outflows {
			bar := generateBar(amount, maxOut)
			fmt.Printf("  %-20s %s %.8f BTC ($%.2f)\n", 
				truncate(addr), bar, amount, amount*getBTCPrice())
		}
	}
}

func displaySummary(txs []Transaction, wallet, network string) {
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	color.Cyan("📈 SUMMARY")
	fmt.Println(strings.Repeat("─", 70))
	
	var totalIn, totalOut, totalFees float64
	var inCount, outCount int
	
	for _, tx := range txs {
		if strings.EqualFold(tx.To, wallet) {
			totalIn += tx.Amount
			inCount++
		} else {
			totalOut += tx.Amount
			outCount++
		}
		totalFees += tx.Fee
	}
	
	btcPrice := getBTCPrice()
	netFlow := totalIn - totalOut
	
	fmt.Printf("📊 Total Transactions: %d\n", len(txs))
	fmt.Printf("📥 Total Received: %.8f %s ($%.2f) in %d transactions\n", 
		totalIn, network, totalIn*btcPrice, inCount)
	fmt.Printf("📤 Total Sent: %.8f %s ($%.2f) in %d transactions\n", 
		totalOut, network, totalOut*btcPrice, outCount)
	
	if netFlow > 0 {
		color.Green("💰 Net Balance Change: +%.8f %s ($%.2f)\n", 
			netFlow, network, netFlow*btcPrice)
	} else if netFlow < 0 {
		color.Red("💰 Net Balance Change: %.8f %s ($%.2f)\n", 
			netFlow, network, netFlow*btcPrice)
	} else {
		fmt.Printf("💰 Net Balance Change: %.8f %s ($%.2f)\n", 
			netFlow, network, netFlow*btcPrice)
	}
	
	fmt.Printf("⛽ Total Fees Paid: %.8f %s ($%.2f)\n", 
		totalFees, network, totalFees*btcPrice)
	
	// Time range
	if len(txs) > 0 {
		oldest := txs[len(txs)-1].Time
		newest := txs[0].Time
		fmt.Printf("📅 Time Range: %s to %s\n", 
			oldest.Format("Jan 02, 2006"), newest.Format("Jan 02, 2006"))
	}
}

func truncate(s string) string {
	if len(s) > 12 {
		return s[:6] + "..." + s[len(s)-4:]
	}
	return s
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

func generateBar(value, max float64) string {
	if max == 0 {
		return ""
	}
	barLength := int((value / max) * 20)
	if barLength == 0 && value > 0 {
		barLength = 1
	}
	return strings.Repeat("█", barLength) + strings.Repeat("░", 20-barLength)
}

// Mock transactions for fallback
func getMockTransactions(wallet, network string) []Transaction {
	return []Transaction{
		{
			Hash:     "demo123...",
			From:     "1Demo...",
			To:       wallet,
			Amount:   0.1,
			Symbol:   "BTC",
			Time:     time.Now().Add(-1 * time.Hour),
			USDValue: 4300.0,
		},
	}
}

// Data structures
type Transaction struct {
	Hash          string
	From          string
	To            string
	Amount        float64
	Fee           float64
	Symbol        string
	Time          time.Time
	BlockHeight   int64
	Confirmations int64
	USDValue      float64
}

type BTCAddressResponse struct {
	Address string   `json:"address"`
	NTx     int64    `json:"n_tx"`
	Txs     []BTCTx  `json:"txs"`
}

type BTCTx struct {
	Hash        string     `json:"hash"`
	Time        int64      `json:"time"`
	BlockHeight int64      `json:"block_height"`
	BlockIndex  int64      `json:"block_index"`
	Fee         int64      `json:"fee"`
	Inputs      []BTCInput `json:"inputs"`
	Out         []BTCOut   `json:"out"`
}

type BTCInput struct {
	PrevOut BTCOut `json:"prev_out"`
}

type BTCOut struct {
	Value int64  `json:"value"`
	Addr  string `json:"addr"`
}
