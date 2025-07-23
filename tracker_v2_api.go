package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// Chain configurations
var chains = map[string]struct {
	ID     int
	Symbol string
	Name   string
}{
	"ETH":      {1, "ETH", "Ethereum"},
	"BSC":      {56, "BNB", "Binance Smart Chain"},
	"MATIC":    {137, "MATIC", "Polygon"},
	"ARB":      {42161, "ETH", "Arbitrum"},
	"OP":       {10, "ETH", "Optimism"},
	"BASE":     {8453, "ETH", "Base"},
	"AVAX":     {43114, "AVAX", "Avalanche"},
	"FTM":      {250, "FTM", "Fantom"},
	"BLAST":    {81457, "ETH", "Blast"},
	"SCROLL":   {534352, "ETH", "Scroll"},
}

// Configuration structure
type Config struct {
	EtherscanAPIKey string `json:"etherscan_api_key"` // Works for ALL chains now!
}

var config Config

func main() {
	// Load configuration
	loadConfig()

	rootCmd := &cobra.Command{
		Use:   "tracker",
		Short: "Universal Blockchain Wallet Tracker (V2 API)",
		Long: `
🚀 Universal Blockchain Wallet Tracker v5.0
Using Etherscan V2 API - One key for 50+ chains!

Supported Networks: ETH, BSC, MATIC, ARB, OP, BASE, AVAX, FTM, BLAST, SCROLL

Examples:
  tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045              # Auto-detect
  tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --network ARB # Arbitrum
  tracker 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa                      # Bitcoin
  tracker config                                                    # Setup`,
		Args: cobra.MaximumNArgs(1),
		Run:  runTracker,
	}

	// Subcommands
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure API key (one key for all chains!)",
		Run:   runConfig,
	}

	rootCmd.AddCommand(configCmd)

	// Flags
	rootCmd.Flags().StringP("network", "n", "auto", "Network: ETH, BSC, MATIC, ARB, OP, BASE, etc.")
	rootCmd.Flags().IntP("limit", "l", 10, "Number of transactions")
	rootCmd.Flags().BoolP("flow", "f", false, "Show fund flow diagram")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() {
	// Check environment variable first
	config.EtherscanAPIKey = os.Getenv("ETHERSCAN_API_KEY")

	// Try to load from config file
	configPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".wallet-tracker", "config.json"),
		"tracker-config.json",
		".tracker-config.json",
	}

	for _, path := range configPaths {
		if data, err := ioutil.ReadFile(path); err == nil {
			var fileConfig Config
			if err := json.Unmarshal(data, &fileConfig); err == nil {
				if config.EtherscanAPIKey == "" {
					config.EtherscanAPIKey = fileConfig.EtherscanAPIKey
				}
				break
			}
		}
	}
}

func runConfig(cmd *cobra.Command, args []string) {
	color.Cyan("🔧 Wallet Tracker Configuration (V2 API)")
	fmt.Println(strings.Repeat("─", 50))
	
	fmt.Println("\n✨ Great news! You now only need ONE API key for all chains!")
	
	fmt.Printf("\nCurrent Etherscan API Key: %s\n", maskAPIKey(config.EtherscanAPIKey))
	
	fmt.Println("\n📝 Configuration Methods:")
	
	color.Green("\n1. Environment Variable (Recommended):")
	fmt.Println("   export ETHERSCAN_API_KEY=your_key_here")
	
	color.Yellow("\n2. Config File:")
	fmt.Println("   Create tracker-config.json with:")
	fmt.Println(`   {
     "etherscan_api_key": "YOUR_KEY_HERE"
   }`)
	
	fmt.Println("\n🔑 Get your API key from:")
	fmt.Println("   https://etherscan.io/apis")
	
	fmt.Println("\n📊 This single key now works for:")
	for code, chain := range chains {
		fmt.Printf("   • %s - %s (Chain ID: %d)\n", code, chain.Name, chain.ID)
	}
	
	// Offer to create config
	fmt.Print("\nWould you like to set up your API key now? (y/n): ")
	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) == "y" {
		createConfigFile()
	}
}

func createConfigFile() {
	var apiKey string
	fmt.Print("\nEnter your Etherscan API Key: ")
	fmt.Scanln(&apiKey)
	
	if apiKey == "" {
		color.Red("No key entered. Exiting.")
		return
	}
	
	configData := Config{EtherscanAPIKey: apiKey}
	data, _ := json.MarshalIndent(configData, "", "  ")
	
	if err := ioutil.WriteFile("tracker-config.json", data, 0600); err != nil {
		color.Red("❌ Failed to save config: %v", err)
	} else {
		color.Green("✅ Config saved! You can now track wallets on all supported chains!")
	}
}

func runTracker(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showWelcome()
		cmd.Help()
		return
	}

	wallet := args[0]
	network, _ := cmd.Flags().GetString("network")
	limit, _ := cmd.Flags().GetInt("limit")
	showFlow, _ := cmd.Flags().GetBool("flow")

	// Auto-detect network
	if network == "auto" {
		network = detectNetwork(wallet)
		if network != "BTC" {
			network = "ETH" // Default EVM chain
		}
		color.Green("✓ Auto-detected network: %s", network)
	}

	// Convert network to uppercase
	network = strings.ToUpper(network)

	// Display header
	displayHeader(wallet, network)

	// Get transactions
	fmt.Println("🔄 Fetching blockchain data...")
	
	var transactions []Transaction
	var err error

	if network == "BTC" {
		transactions, err = getBTCTransactions(wallet, limit)
	} else if chain, ok := chains[network]; ok {
		// Use V2 API for all EVM chains
		transactions, err = getEVMTransactionsV2(wallet, chain.ID, chain.Symbol, limit)
	} else {
		color.Red("❌ Unsupported network: %s", network)
		fmt.Println("\nSupported networks:")
		for code, chain := range chains {
			fmt.Printf("  %s - %s\n", code, chain.Name)
		}
		return
	}

	if err != nil {
		color.Red("❌ Error: %v", err)
		return
	}

	if len(transactions) == 0 {
		color.Yellow("⚠️  No transactions found for this wallet")
		return
	}

	// Display results
	displayTransactionTable(transactions, wallet)

	if showFlow {
		displayFundFlow(transactions, wallet)
	}

	displaySummary(transactions, wallet, network)
}

// V2 API implementation for all EVM chains
func getEVMTransactionsV2(wallet string, chainID int, symbol string, limit int) ([]Transaction, error) {
	if config.EtherscanAPIKey == "" {
		return nil, fmt.Errorf("no API key configured. Run './tracker config' to set up")
	}

	// V2 API endpoint - single endpoint for all chains!
	url := fmt.Sprintf("https://api.etherscan.io/v2/api?chainid=%d&module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=%d&sort=desc&apikey=%s",
		chainID, wallet, limit, config.EtherscanAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			BlockNumber       string `json:"blockNumber"`
			TimeStamp         string `json:"timeStamp"`
			Hash              string `json:"hash"`
			From              string `json:"from"`
			To                string `json:"to"`
			Value             string `json:"value"`
			Gas               string `json:"gas"`
			GasPrice          string `json:"gasPrice"`
			GasUsed           string `json:"gasUsed"`
			IsError           string `json:"isError"`
			Confirmations     string `json:"confirmations"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "1" {
		if result.Message == "No transactions found" {
			return []Transaction{}, nil
		}
		if strings.Contains(result.Message, "API") || strings.Contains(result.Message, "key") {
			return nil, fmt.Errorf("API key issue: %s\nPlease run './tracker config' to set up", result.Message)
		}
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	price := getPriceForSymbol(symbol)
	transactions := make([]Transaction, 0)

	for _, tx := range result.Result {
		// Convert values
		value := parseWei(tx.Value)
		gasUsed := parseWei(tx.GasUsed)
		gasPrice := parseWei(tx.GasPrice)
		fee := gasUsed * gasPrice
		
		timestamp, _ := parseInt64(tx.TimeStamp)
		
		// Determine transaction type
		txType := "OUT"
		if strings.EqualFold(tx.To, wallet) && strings.EqualFold(tx.From, wallet) {
			txType = "SELF"
		} else if strings.EqualFold(tx.To, wallet) {
			txType = "IN"
		}

		transactions = append(transactions, Transaction{
			Hash:     tx.Hash,
			From:     tx.From,
			To:       tx.To,
			Amount:   value,
			Fee:      fee,
			Symbol:   symbol,
			Time:     time.Unix(timestamp, 0),
			USDValue: value * price,
			Type:     txType,
			Status:   tx.IsError == "0",
		})
	}

	return transactions, nil
}

// Display functions
func displayHeader(wallet, network string) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════════════════════╗")
	color.Cyan("║         UNIVERSAL BLOCKCHAIN TRACKER v5.0 (V2 API)              ║")
	color.Cyan("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	color.Yellow("📊 Tracking: %s", truncate(wallet))
	
	if chain, ok := chains[network]; ok {
		fmt.Printf("🌐 Network: %s (Chain ID: %d)\n", chain.Name, chain.ID)
	} else if network == "BTC" {
		color.HiYellow("₿  Network: Bitcoin")
	} else {
		fmt.Printf("🌐 Network: %s\n", network)
	}
	
	color.White("🕒 Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("─", 70))
}

func showWelcome() {
	color.Cyan(`
╔══════════════════════════════════════════════════════════════════╗
║           UNIVERSAL BLOCKCHAIN TRACKER (V2 API)                  ║
╚══════════════════════════════════════════════════════════════════╝`)
	
	fmt.Println("\n🚀 Track any wallet across 50+ blockchains with ONE API key!\n")
	
	fmt.Println("Supported Networks:")
	color.Yellow("  • Bitcoin (BTC) - No API key needed")
	
	fmt.Println("\n  With Etherscan API key (one key for all!):")
	for code, chain := range chains {
		fmt.Printf("  • %s - %s\n", code, chain.Name)
	}
	
	fmt.Println("\nExamples:")
	fmt.Println("  tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045       # Ethereum")
	fmt.Println("  tracker 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045 --network ARB")
	fmt.Println("  tracker config                                           # Setup API key")
}

// Helper functions remain the same...
func detectNetwork(wallet string) string {
	if strings.HasPrefix(wallet, "0x") && len(wallet) == 42 {
		return "ETH"
	} else if strings.HasPrefix(wallet, "1") || strings.HasPrefix(wallet, "bc1") || strings.HasPrefix(wallet, "3") {
		return "BTC"
	}
	return "UNKNOWN"
}

func maskAPIKey(key string) string {
	if key == "" {
		return "Not configured"
	}
	if len(key) < 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func truncate(s string) string {
	if len(s) > 15 {
		return s[:6] + "..." + s[len(s)-4:]
	}
	return s
}

func parseWei(s string) float64 {
	if s == "" {
		return 0
	}
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f / 1e18
}

func parseInt64(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func formatAddress(addr, myWallet string) string {
	if strings.EqualFold(addr, myWallet) {
		return "[TRACKED]"
	}
	return truncate(addr)
}

// Display functions (same as before)
func displayTransactionTable(txs []Transaction, myWallet string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	
	t.AppendHeader(table.Row{"#", "Type", "Time", "From → To", "Amount", "USD Value", "Status"})
	
	for i, tx := range txs {
		var typeColor text.Color
		switch tx.Type {
		case "IN":
			typeColor = text.FgGreen
		case "OUT":
			typeColor = text.FgRed
		case "SELF":
			typeColor = text.FgYellow
		default:
			typeColor = text.FgWhite
		}
		
		timeStr := tx.Time.Format("01/02 15:04")
		
		fromAddr := formatAddress(tx.From, myWallet)
		toAddr := formatAddress(tx.To, myWallet)
		
		status := "✅"
		if !tx.Status {
			status = "❌"
		}
		
		row := table.Row{
			i + 1,
			text.Colors{typeColor}.Sprint(tx.Type),
			timeStr,
			fmt.Sprintf("%s → %s", fromAddr, toAddr),
			fmt.Sprintf("%.6f %s", tx.Amount, tx.Symbol),
			fmt.Sprintf("$%.2f", tx.USDValue),
			status,
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
		if tx.Type == "IN" {
			inflows[tx.From] += tx.Amount
		} else if tx.Type == "OUT" {
			outflows[tx.To] += tx.Amount
		}
	}
	
	if len(inflows) > 0 {
		color.Green("📥 INFLOWS:")
		for addr, amount := range inflows {
			fmt.Printf("  %-20s %.6f %s\n", truncate(addr), amount, txs[0].Symbol)
		}
	}
	
	fmt.Println()
	color.Yellow("  [%s]", truncate(wallet))
	fmt.Println()
	
	if len(outflows) > 0 {
		color.Red("📤 OUTFLOWS:")
		for addr, amount := range outflows {
			fmt.Printf("  %-20s %.6f %s\n", truncate(addr), amount, txs[0].Symbol)
		}
	}
}

func displaySummary(txs []Transaction, wallet, network string) {
	fmt.Println()
	fmt.Println(strings.Repeat("═", 70))
	color.Cyan("📈 SUMMARY")
	
	var totalIn, totalOut, totalFees float64
	var inCount, outCount int
	
	for _, tx := range txs {
		if tx.Type == "IN" {
			totalIn += tx.Amount
			inCount++
		} else if tx.Type == "OUT" {
			totalOut += tx.Amount
			outCount++
		}
		totalFees += tx.Fee
	}
	
	symbol := txs[0].Symbol
	price := getPriceForSymbol(symbol)
	
	fmt.Printf("\n📥 Received: %.6f %s ($%.2f) - %d txs\n", totalIn, symbol, totalIn*price, inCount)
	fmt.Printf("📤 Sent: %.6f %s ($%.2f) - %d txs\n", totalOut, symbol, totalOut*price, outCount)
	fmt.Printf("💰 Net: %.6f %s ($%.2f)\n", totalIn-totalOut, symbol, (totalIn-totalOut)*price)
	if totalFees > 0 {
		fmt.Printf("⛽ Fees: %.6f %s ($%.2f)\n", totalFees, symbol, totalFees*price)
	}
}

// Bitcoin implementation (unchanged)
func getBTCTransactions(wallet string, limit int) ([]Transaction, error) {
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s?limit=%d", wallet, limit)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var btcResp BTCAddressResponse
	if err := json.Unmarshal(body, &btcResp); err != nil {
		return nil, err
	}

	btcPrice := getBTCPrice()
	transactions := make([]Transaction, 0)
	
	for _, tx := range btcResp.Txs {
		var amount int64 = 0
		var from, to string
		var txType string = "UNKNOWN"
		
		for _, out := range tx.Out {
			if out.Addr == wallet {
				amount += out.Value
				to = wallet
				if len(tx.Inputs) > 0 && tx.Inputs[0].PrevOut.Addr != "" {
					from = tx.Inputs[0].PrevOut.Addr
				}
				txType = "IN"
			}
		}
		
		if txType == "UNKNOWN" {
			for _, input := range tx.Inputs {
				if input.PrevOut.Addr == wallet {
					from = wallet
					for _, out := range tx.Out {
						if out.Addr != wallet && out.Value > amount {
							to = out.Addr
							amount = out.Value
						}
					}
					txType = "OUT"
					break
				}
			}
		}
		
		if amount == 0 {
			continue
		}
		
		btcAmount := float64(amount) / 100000000
		feeAmount := float64(tx.Fee) / 100000000
		
		transactions = append(transactions, Transaction{
			Hash:     tx.Hash,
			From:     from,
			To:       to,
			Amount:   btcAmount,
			Fee:      feeAmount,
			Symbol:   "BTC",
			Time:     time.Unix(tx.Time, 0),
			USDValue: btcAmount * btcPrice,
			Type:     txType,
			Status:   true,
		})
	}
	
	return transactions, nil
}

// Price functions
func getPriceForSymbol(symbol string) float64 {
	switch symbol {
	case "BTC":
		return getBTCPrice()
	case "ETH":
		return getETHPrice()
	case "BNB":
		return getBNBPrice()
	case "MATIC":
		return getMaticPrice()
	case "AVAX":
		return getPrice("avalanche-2", 35)
	case "FTM":
		return getPrice("fantom", 0.35)
	default:
		return getETHPrice() // Default for ETH-based L2s
	}
}

func getBTCPrice() float64 {
	return getPrice("bitcoin", 43000)
}

func getETHPrice() float64 {
	return getPrice("ethereum", 2350)
}

func getBNBPrice() float64 {
	return getPrice("binancecoin", 315)
}

func getMaticPrice() float64 {
	return getPrice("matic-network", 0.85)
}

func getPrice(coin string, fallback float64) float64 {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coin)
	
	resp, err := http.Get(url)
	if err != nil {
		return fallback
	}
	defer resp.Body.Close()
	
	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fallback
	}
	
	if price, ok := result[coin]["usd"]; ok {
		return price
	}
	
	return fallback
}

// Data structures
type Transaction struct {
	Hash     string
	From     string
	To       string
	Amount   float64
	Fee      float64
	Symbol   string
	Time     time.Time
	USDValue float64
	Type     string
	Status   bool
	Note     string
}

type BTCAddressResponse struct {
	Txs []BTCTx `json:"txs"`
}

type BTCTx struct {
	Hash   string     `json:"hash"`
	Time   int64      `json:"time"`
	Fee    int64      `json:"fee"`
	Inputs []BTCInput `json:"inputs"`
	Out    []BTCOut   `json:"out"`
}

type BTCInput struct {
	PrevOut BTCOut `json:"prev_out"`
}

type BTCOut struct {
	Value int64  `json:"value"`
	Addr  string `json:"addr"`
}
