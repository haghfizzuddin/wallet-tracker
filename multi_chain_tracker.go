package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "multi-chain-tracker",
		Short: "Multi-chain blockchain wallet tracker (BTC, ETH, BSC)",
	}

	trackCmd := &cobra.Command{
		Use:   "track",
		Short: "Track wallet transactions across multiple blockchains",
		RunE:  runTrack,
	}

	trackCmd.Flags().StringP("wallet", "w", "", "Wallet address to track")
	trackCmd.Flags().StringP("network", "n", "auto", "Network: BTC, ETH, BSC, or auto-detect")
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
		network = detectNetwork(wallet)
	}

	// Display header
	displayHeader(wallet, network)

	// Get transactions based on network
	fmt.Println("🔄 Fetching blockchain data...")
	
	var transactions []Transaction
	var err error

	switch strings.ToUpper(network) {
	case "BTC":
		transactions, err = getRealBTCTransactions(wallet, limit)
	case "ETH":
		transactions, err = getETHTransactions(wallet, limit)
	case "BSC":
		transactions, err = getBSCTransactions(wallet, limit)
	default:
		return fmt.Errorf("unsupported network: %s", network)
	}

	if err != nil {
		color.Red("❌ Error fetching transactions: %v", err)
		return err
	}

	if len(transactions) == 0 {
		color.Yellow("⚠️  No transactions found for this wallet")
		color.Yellow("\nNote: ETH/BSC require API keys for full data. Using limited free tier.")
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

func detectNetwork(wallet string) string {
	if strings.HasPrefix(wallet, "0x") && len(wallet) == 42 {
		// Could be ETH or BSC - default to ETH
		return "ETH"
	} else if strings.HasPrefix(wallet, "1") || strings.HasPrefix(wallet, "bc1") || strings.HasPrefix(wallet, "3") {
		return "BTC"
	}
	return "UNKNOWN"
}

// Ethereum transactions using Etherscan-compatible API
func getETHTransactions(wallet string, limit int) ([]Transaction, error) {
	// Using public Ethereum node (limited functionality)
	// For production, use Etherscan API with key
	
	// Try Ethplorer API (free tier, no key required)
	url := fmt.Sprintf("https://api.ethplorer.io/getAddressHistory/%s?apiKey=freekey&limit=%d", wallet, limit)
	
	resp, err := http.Get(url)
	if err != nil {
		// Fallback to basic Ethereum data
		return getBasicETHInfo(wallet)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return getBasicETHInfo(wallet)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getBasicETHInfo(wallet)
	}

	var ethResp struct {
		Operations []struct {
			Timestamp   int64   `json:"timestamp"`
			From        string  `json:"from"`
			To          string  `json:"to"`
			Value       float64 `json:"value"`
			Type        string  `json:"type"`
			Hash        string  `json:"transactionHash"`
		} `json:"operations"`
	}

	if err := json.Unmarshal(body, &ethResp); err != nil {
		return getBasicETHInfo(wallet)
	}

	ethPrice := getETHPrice()
	transactions := make([]Transaction, 0)

	for i, op := range ethResp.Operations {
		if i >= limit {
			break
		}

		txType := "IN"
		if strings.EqualFold(op.From, wallet) {
			txType = "OUT"
		}

		transactions = append(transactions, Transaction{
			Hash:     op.Hash,
			From:     op.From,
			To:       op.To,
			Amount:   op.Value,
			Symbol:   "ETH",
			Time:     time.Unix(op.Timestamp, 0),
			USDValue: op.Value * ethPrice,
			Type:     txType,
		})
	}

	return transactions, nil
}

// BSC transactions
func getBSCTransactions(wallet string, limit int) ([]Transaction, error) {
	// BSC uses similar structure to ETH
	// Using BSC public API (limited)
	
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=%d&sort=desc&apikey=YourAPIKeyToken", wallet, limit)
	
	// Without API key, return basic info
	return getBasicBSCInfo(wallet)
}

// Basic ETH info without full API
func getBasicETHInfo(wallet string) ([]Transaction, error) {
	// Get balance using public endpoint
	url := fmt.Sprintf("https://api.ethplorer.io/getAddressInfo/%s?apiKey=freekey", wallet)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info struct {
		ETH struct {
			Balance float64 `json:"balance"`
		} `json:"ETH"`
		CountTxs int `json:"countTxs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	ethPrice := getETHPrice()

	// Return balance info as a pseudo-transaction
	return []Transaction{
		{
			Hash:     "Balance Check",
			From:     "Network",
			To:       wallet,
			Amount:   info.ETH.Balance,
			Symbol:   "ETH",
			Time:     time.Now(),
			USDValue: info.ETH.Balance * ethPrice,
			Type:     "INFO",
			Note:     fmt.Sprintf("Current Balance - Total Txs: %d", info.CountTxs),
		},
	}, nil
}

// Basic BSC info
func getBasicBSCInfo(wallet string) ([]Transaction, error) {
	// Get BNB balance using public BSC RPC
	bnbPrice := getBNBPrice()
	
	// Return demo transaction for BSC
	return []Transaction{
		{
			Hash:     "BSC Balance Check",
			From:     "Network",
			To:       wallet,
			Amount:   0,
			Symbol:   "BNB",
			Time:     time.Now(),
			USDValue: 0,
			Type:     "INFO",
			Note:     "BSC data requires BscScan API key for full transaction history",
		},
	}, nil
}

// Get ETH price
func getETHPrice() float64 {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
	
	resp, err := http.Get(url)
	if err != nil {
		return 2350.0 // Fallback
	}
	defer resp.Body.Close()

	var price struct {
		Ethereum struct {
			USD float64 `json:"usd"`
		} `json:"ethereum"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 2350.0
	}

	return price.Ethereum.USD
}

// Get BNB price
func getBNBPrice() float64 {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=binancecoin&vs_currencies=usd"
	
	resp, err := http.Get(url)
	if err != nil {
		return 315.0 // Fallback
	}
	defer resp.Body.Close()

	var price struct {
		Binancecoin struct {
			USD float64 `json:"usd"`
		} `json:"binancecoin"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 315.0
	}

	return price.Binancecoin.USD
}

// Bitcoin implementation (existing)
func getRealBTCTransactions(wallet string, limit int) ([]Transaction, error) {
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

	btcPrice := getBTCPrice()
	transactions := make([]Transaction, 0)
	
	for _, tx := range btcResp.Txs {
		isIncoming := false
		var amount int64 = 0
		var from, to string
		var fee int64 = tx.Fee
		
		for _, out := range tx.Out {
			if out.Addr == wallet {
				isIncoming = true
				amount += out.Value
				to = wallet
				if len(tx.Inputs) > 0 && tx.Inputs[0].PrevOut.Addr != "" {
					from = tx.Inputs[0].PrevOut.Addr
				}
			}
		}
		
		if !isIncoming {
			for _, input := range tx.Inputs {
				if input.PrevOut.Addr == wallet {
					from = wallet
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
		
		if amount == 0 {
			continue
		}
		
		btcAmount := float64(amount) / 100000000
		feeAmount := float64(fee) / 100000000
		
		txType := "IN"
		if from == wallet {
			txType = "OUT"
		}
		
		transactions = append(transactions, Transaction{
			Hash:          tx.Hash,
			From:          from,
			To:            to,
			Amount:        btcAmount,
			Fee:           feeAmount,
			Symbol:        "BTC",
			Time:          time.Unix(tx.Time, 0),
			BlockHeight:   tx.BlockHeight,
			Confirmations: 0, // Would need current block height
			USDValue:      btcAmount * btcPrice,
			Type:          txType,
		})
	}
	
	return transactions, nil
}

func getBTCPrice() float64 {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"
	
	resp, err := http.Get(url)
	if err != nil {
		return 43000.0
	}
	defer resp.Body.Close()
	
	var price struct {
		Bitcoin struct {
			USD float64 `json:"usd"`
		} `json:"bitcoin"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 43000.0
	}
	
	return price.Bitcoin.USD
}

func displayHeader(wallet, network string) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════════════════════╗")
	color.Cyan("║           MULTI-CHAIN WALLET TRACKER V3.0                        ║")
	color.Cyan("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	color.Yellow("📊 Tracking Wallet: %s", wallet)
	color.Yellow("🌐 Network: %s", network)
	color.Yellow("🕒 Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	
	// Network-specific info
	switch network {
	case "ETH":
		color.Blue("⟠ Ethereum Network")
	case "BSC":
		color.Yellow("🔶 Binance Smart Chain")
	case "BTC":
		color.HiYellow("₿ Bitcoin Network")
	}
	
	fmt.Println(strings.Repeat("─", 70))
}

func displayTransactionTable(txs []Transaction, myWallet string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	
	headers := table.Row{"#", "Type", "Time", "From → To", "Amount", "USD Value"}
	
	// Add extra columns based on available data
	hasConfirmations := false
	hasFees := false
	hasNotes := false
	
	for _, tx := range txs {
		if tx.Confirmations > 0 {
			hasConfirmations = true
		}
		if tx.Fee > 0 {
			hasFees = true
		}
		if tx.Note != "" {
			hasNotes = true
		}
	}
	
	if hasConfirmations {
		headers = append(headers, "Confirms")
	}
	if hasFees {
		headers = append(headers, "Fee")
	}
	if hasNotes {
		headers = append(headers, "Note")
	}
	
	t.AppendHeader(headers)
	
	for i, tx := range txs {
		var txType string
		var typeColor text.Color
		
		switch tx.Type {
		case "IN":
			txType = "IN"
			typeColor = text.FgGreen
		case "OUT":
			txType = "OUT"
			typeColor = text.FgRed
		case "INFO":
			txType = "INFO"
			typeColor = text.FgCyan
		default:
			txType = "???"
			typeColor = text.FgWhite
		}
		
		timeStr := tx.Time.Format("01/02 15:04")
		if tx.Type == "INFO" {
			timeStr = "Current"
		}
		
		fromAddr := truncate(tx.From)
		toAddr := truncate(tx.To)
		if tx.From == myWallet {
			fromAddr = "[TRACKED]"
		}
		if tx.To == myWallet {
			toAddr = "[TRACKED]"
		}
		
		row := table.Row{
			i + 1,
			text.Colors{typeColor}.Sprint(txType),
			timeStr,
			fmt.Sprintf("%s → %s", fromAddr, toAddr),
			fmt.Sprintf("%.6f %s", tx.Amount, tx.Symbol),
			fmt.Sprintf("$%.2f", tx.USDValue),
		}
		
		if hasConfirmations {
			if tx.Confirmations > 0 {
				row = append(row, fmt.Sprintf("✅ %d", tx.Confirmations))
			} else {
				row = append(row, "-")
			}
		}
		
		if hasFees {
			if tx.Fee > 0 {
				row = append(row, fmt.Sprintf("%.8f", tx.Fee))
			} else {
				row = append(row, "-")
			}
		}
		
		if hasNotes {
			row = append(row, tx.Note)
		}
		
		t.AppendRow(row)
		
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
		if tx.Type == "INFO" {
			continue
		}
		
		if strings.EqualFold(tx.To, wallet) {
			inflows[tx.From] += tx.Amount
		} else {
			outflows[tx.To] += tx.Amount
		}
	}
	
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
			symbol := txs[0].Symbol
			price := getPriceForSymbol(symbol)
			fmt.Printf("  %-20s %s %.6f %s ($%.2f)\n", 
				truncate(addr), bar, amount, symbol, amount*price)
		}
	}
	
	fmt.Println()
	color.Yellow("                    ╔════════════════════════════╗")
	color.Yellow("                    ║      TRACKED WALLET        ║")
	color.Yellow("                    ║   %s   ║", centerText(truncate(wallet), 24))
	color.Yellow("                    ╚════════════════════════════════╝")
	fmt.Println()
	
	if len(outflows) > 0 {
		color.Red("📤 OUTFLOWS:")
		for addr, amount := range outflows {
			bar := generateBar(amount, maxOut)
			symbol := txs[0].Symbol
			price := getPriceForSymbol(symbol)
			fmt.Printf("  %-20s %s %.6f %s ($%.2f)\n", 
				truncate(addr), bar, amount, symbol, amount*price)
		}
	}
	
	if len(inflows) == 0 && len(outflows) == 0 {
		color.Yellow("No transaction flow data available")
	}
}

func displaySummary(txs []Transaction, wallet, network string) {
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	color.Cyan("📈 SUMMARY")
	fmt.Println(strings.Repeat("─", 70))
	
	var totalIn, totalOut, totalFees float64
	var inCount, outCount int
	symbol := network
	
	for _, tx := range txs {
		if tx.Type == "INFO" {
			continue
		}
		
		symbol = tx.Symbol
		
		if strings.EqualFold(tx.To, wallet) {
			totalIn += tx.Amount
			inCount++
		} else {
			totalOut += tx.Amount
			outCount++
		}
		totalFees += tx.Fee
	}
	
	price := getPriceForSymbol(symbol)
	netFlow := totalIn - totalOut
	
	fmt.Printf("🌐 Network: %s\n", network)
	fmt.Printf("📊 Total Transactions: %d\n", len(txs))
	
	if inCount > 0 || outCount > 0 {
		fmt.Printf("📥 Total Received: %.6f %s ($%.2f) in %d transactions\n", 
			totalIn, symbol, totalIn*price, inCount)
		fmt.Printf("📤 Total Sent: %.6f %s ($%.2f) in %d transactions\n", 
			totalOut, symbol, totalOut*price, outCount)
		
		if netFlow > 0 {
			color.Green("💰 Net Balance Change: +%.6f %s ($%.2f)\n", 
				netFlow, symbol, netFlow*price)
		} else if netFlow < 0 {
			color.Red("💰 Net Balance Change: %.6f %s ($%.2f)\n", 
				netFlow, symbol, netFlow*price)
		} else {
			fmt.Printf("💰 Net Balance Change: %.6f %s ($%.2f)\n", 
				netFlow, symbol, netFlow*price)
		}
		
		if totalFees > 0 {
			fmt.Printf("⛽ Total Fees: %.8f %s ($%.2f)\n", 
				totalFees, symbol, totalFees*price)
		}
	}
	
	// Add network-specific notes
	switch network {
	case "ETH":
		color.Yellow("\n📝 Note: Full ETH transaction history requires Etherscan API key")
		color.Yellow("   Current data from Ethplorer free tier (limited)")
	case "BSC":
		color.Yellow("\n📝 Note: BSC transaction history requires BscScan API key")
		color.Yellow("   Showing balance information only")
	}
}

func getPriceForSymbol(symbol string) float64 {
	switch symbol {
	case "BTC":
		return getBTCPrice()
	case "ETH":
		return getETHPrice()
	case "BNB":
		return getBNBPrice()
	default:
		return 0
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
	Type          string
	Note          string
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
