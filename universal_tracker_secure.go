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

// Configuration structure
type Config struct {
	EtherscanAPIKey string `json:"etherscan_api_key"`
	BscscanAPIKey   string `json:"bscscan_api_key"`
	PolygonAPIKey   string `json:"polygon_api_key"`
}

var config Config

func main() {
	// Load configuration
	loadConfig()

	rootCmd := &cobra.Command{
		Use:   "tracker",
		Short: "Universal Blockchain Wallet Tracker",
		Long: `
🚀 Universal Blockchain Wallet Tracker v4.0
Supports: Bitcoin, Ethereum, BSC, Polygon

Examples:
  tracker 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa              # Auto-detect Bitcoin
  tracker 0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8      # Auto-detect Ethereum
  tracker 0x123...abc --network BSC                       # Force BSC network
  tracker 0x123...abc --flow                              # Show fund flow diagram
  tracker config                                           # Configure API keys`,
		Args: cobra.MaximumNArgs(1),
		Run:  runTracker,
	}

	// Subcommands
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure API keys securely",
		Run:   runConfig,
	}

	rootCmd.AddCommand(configCmd)

	// Flags
	rootCmd.Flags().StringP("network", "n", "auto", "Network: BTC, ETH, BSC, MATIC")
	rootCmd.Flags().IntP("limit", "l", 10, "Number of transactions")
	rootCmd.Flags().BoolP("flow", "f", false, "Show fund flow diagram")
	rootCmd.Flags().BoolP("export", "e", false, "Export to CSV")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() {
	// Priority order for API keys:
	// 1. Environment variables
	// 2. Config file (~/.wallet-tracker/config.json)
	// 3. Local config file (./tracker-config.json)

	// Check environment variables first
	config.EtherscanAPIKey = os.Getenv("ETHERSCAN_API_KEY")
	config.BscscanAPIKey = os.Getenv("BSCSCAN_API_KEY")
	config.PolygonAPIKey = os.Getenv("POLYGON_API_KEY")

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
				// Only override if not set by environment
				if config.EtherscanAPIKey == "" {
					config.EtherscanAPIKey = fileConfig.EtherscanAPIKey
				}
				if config.BscscanAPIKey == "" {
					config.BscscanAPIKey = fileConfig.BscscanAPIKey
				}
				if config.PolygonAPIKey == "" {
					config.PolygonAPIKey = fileConfig.PolygonAPIKey
				}
				break
			}
		}
	}
}

func runConfig(cmd *cobra.Command, args []string) {
	color.Cyan("🔧 Wallet Tracker Configuration")
	fmt.Println(strings.Repeat("─", 50))
	
	fmt.Println("\nCurrent configuration:")
	fmt.Printf("  Etherscan API: %s\n", maskAPIKey(config.EtherscanAPIKey))
	fmt.Printf("  BscScan API:   %s\n", maskAPIKey(config.BscscanAPIKey))
	fmt.Printf("  Polygon API:   %s\n", maskAPIKey(config.PolygonAPIKey))
	
	fmt.Println("\n📝 Configuration Methods (in order of priority):")
	
	color.Green("\n1. Environment Variables (Recommended):")
	fmt.Println("   export ETHERSCAN_API_KEY=your_key_here")
	fmt.Println("   export BSCSCAN_API_KEY=your_key_here")
	fmt.Println("   export POLYGON_API_KEY=your_key_here")
	
	color.Yellow("\n2. User Config File:")
	configDir := filepath.Join(os.Getenv("HOME"), ".wallet-tracker")
	fmt.Printf("   Create: %s/config.json\n", configDir)
	
	color.Blue("\n3. Local Config File:")
	fmt.Println("   Create: ./tracker-config.json")
	
	fmt.Println("\nExample config.json format:")
	fmt.Println(`{
  "etherscan_api_key": "YOUR_KEY_HERE",
  "bscscan_api_key": "YOUR_KEY_HERE",
  "polygon_api_key": "YOUR_KEY_HERE"
}`)
	
	fmt.Println("\n🔐 Security Tips:")
	fmt.Println("   • Never commit API keys to Git")
	fmt.Println("   • Add config files to .gitignore")
	fmt.Println("   • Use environment variables in production")
	
	// Offer to create config file
	fmt.Print("\nWould you like to create a config file now? (y/n): ")
	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) == "y" {
		createConfigFile()
	}
}

func createConfigFile() {
	configDir := filepath.Join(os.Getenv("HOME"), ".wallet-tracker")
	os.MkdirAll(configDir, 0700)
	
	configPath := filepath.Join(configDir, "config.json")
	
	var newConfig Config
	
	fmt.Print("\nEtherscan API Key (or press Enter to skip): ")
	fmt.Scanln(&newConfig.EtherscanAPIKey)
	
	fmt.Print("BscScan API Key (or press Enter to skip): ")
	fmt.Scanln(&newConfig.BscscanAPIKey)
	
	fmt.Print("Polygon API Key (or press Enter to skip): ")
	fmt.Scanln(&newConfig.PolygonAPIKey)
	
	data, _ := json.MarshalIndent(newConfig, "", "  ")
	
	if err := ioutil.WriteFile(configPath, data, 0600); err != nil {
		color.Red("❌ Failed to save config: %v", err)
	} else {
		color.Green("✅ Config saved to: %s", configPath)
		color.Yellow("   This file will be used automatically next time")
	}
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

func runTracker(cmd *cobra.Command, args []string) {
	// Show help if no wallet provided
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
		color.Green("✓ Auto-detected network: %s", network)
	}

	// Check API keys for ETH/BSC
	if (network == "ETH" || network == "BSC" || network == "MATIC") && !hasAPIKeys(network) {
		showAPIKeyWarning(network)
	}

	// Display header
	displayHeader(wallet, network)

	// Get transactions
	fmt.Println("🔄 Fetching blockchain data...")
	
	var transactions []Transaction
	var err error

	switch strings.ToUpper(network) {
	case "BTC":
		transactions, err = getBTCTransactions(wallet, limit)
	case "ETH":
		transactions, err = getETHTransactions(wallet, limit)
	case "BSC":
		transactions, err = getBSCTransactions(wallet, limit)
	case "MATIC":
		transactions, err = getPolygonTransactions(wallet, limit)
	default:
		color.Red("❌ Unsupported network: %s", network)
		return
	}

	if err != nil {
		color.Red("❌ Error: %v", err)
		return
	}

	if len(transactions) == 0 {
		color.Yellow("⚠️  No transactions found")
		return
	}

	// Display results
	displayTransactionTable(transactions, wallet)

	if showFlow {
		displayFundFlow(transactions, wallet)
	}

	displaySummary(transactions, wallet, network)
}

func showWelcome() {
	color.Cyan(`
╔══════════════════════════════════════════════════════════════════╗
║                 UNIVERSAL BLOCKCHAIN TRACKER                     ║
╚══════════════════════════════════════════════════════════════════╝`)
	
	fmt.Println("\n📊 Track any wallet across multiple blockchains!\n")
	
	fmt.Println("Supported Networks:")
	color.Yellow("  • Bitcoin (BTC) - Full support, no API key needed")
	color.Blue("  • Ethereum (ETH) - Requires Etherscan API key")
	color.Yellow("  • Binance Smart Chain (BSC) - Requires BscScan API key")
	color.Magenta("  • Polygon (MATIC) - Requires PolygonScan API key")
	
	fmt.Println("\nQuick Examples:")
}

func detectNetwork(wallet string) string {
	if strings.HasPrefix(wallet, "0x") && len(wallet) == 42 {
		return "ETH" // Default to ETH for 0x addresses
	} else if strings.HasPrefix(wallet, "1") || strings.HasPrefix(wallet, "bc1") || strings.HasPrefix(wallet, "3") {
		return "BTC"
	}
	return "UNKNOWN"
}

func hasAPIKeys(network string) bool {
	switch network {
	case "ETH":
		return config.EtherscanAPIKey != ""
	case "BSC":
		return config.BscscanAPIKey != ""
	case "MATIC":
		return config.PolygonAPIKey != ""
	}
	return false
}

func showAPIKeyWarning(network string) {
	color.Yellow("\n⚠️  No API key configured for %s", network)
	color.Yellow("   Using limited free tier (may have restrictions)")
	color.Yellow("   Run 'tracker config' to set up API keys\n")
	time.Sleep(2 * time.Second)
}

// Get Ethereum transactions
func getETHTransactions(wallet string, limit int) ([]Transaction, error) {
	if hasAPIKeys("ETH") {
		// Use Etherscan with API key
		url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=%d&sort=desc&apikey=%s",
			wallet, limit, config.EtherscanAPIKey)
		
		return getEtherscanTransactions(url, "ETH", wallet)
	}
	
	// Fallback to free tier
	return getBasicETHInfo(wallet)
}

// Get BSC transactions
func getBSCTransactions(wallet string, limit int) ([]Transaction, error) {
	if hasAPIKeys("BSC") {
		url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=%d&sort=desc&apikey=%s",
			wallet, limit, config.BscscanAPIKey)
		
		return getEtherscanTransactions(url, "BNB", wallet)
	}
	
	return getBasicInfo(wallet, "BSC", "BNB")
}

// Get Polygon transactions
func getPolygonTransactions(wallet string, limit int) ([]Transaction, error) {
	if hasAPIKeys("MATIC") {
		url := fmt.Sprintf("https://api.polygonscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=%d&sort=desc&apikey=%s",
			wallet, limit, config.PolygonAPIKey)
		
		return getEtherscanTransactions(url, "MATIC", wallet)
	}
	
	return getBasicInfo(wallet, "Polygon", "MATIC")
}

// Generic Etherscan-compatible API handler
func getEtherscanTransactions(url, symbol, wallet string) ([]Transaction, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if result.Status != "1" {
		if result.Message == "No transactions found" {
			return []Transaction{}, nil // Empty result, not an error
		}
		if strings.Contains(result.Message, "API Key") || strings.Contains(result.Message, "apikey") {
			return nil, fmt.Errorf("API key issue: %s\nPlease run './tracker config' to set up API keys", result.Message)
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

// Bitcoin transactions (existing implementation)
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
		
		// Determine transaction type and amount
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

// Display functions
func displayHeader(wallet, network string) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════════════════════╗")
	color.Cyan("║              UNIVERSAL BLOCKCHAIN TRACKER v4.0                   ║")
	color.Cyan("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	color.Yellow("📊 Tracking: %s", truncate(wallet))
	
	switch network {
	case "BTC":
		color.HiYellow("₿  Network: Bitcoin")
	case "ETH":
		color.Blue("⟠  Network: Ethereum")
	case "BSC":
		color.Yellow("🔶 Network: Binance Smart Chain")
	case "MATIC":
		color.Magenta("⬟  Network: Polygon")
	}
	
	color.White("🕒 Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("─", 70))
}

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
}

// Helper functions
func formatAddress(addr, myWallet string) string {
	if strings.EqualFold(addr, myWallet) {
		return "[TRACKED]"
	}
	return truncate(addr)
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

func getBasicETHInfo(wallet string) ([]Transaction, error) {
	return getBasicInfo(wallet, "Ethereum", "ETH")
}

func getBasicInfo(wallet, network, symbol string) ([]Transaction, error) {
	// price := getPriceForSymbol(symbol) // Not needed for basic info
	return []Transaction{{
		Hash:     "BALANCE_CHECK",
		From:     "Network",
		To:       wallet,
		Amount:   0,
		Symbol:   symbol,
		Time:     time.Now(),
		USDValue: 0,
		Type:     "INFO",
		Note:     fmt.Sprintf("%s - API key required for full history", network),
	}}, nil
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
	default:
		return 0
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
