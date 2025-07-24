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
	"github.com/spf13/cobra"
)

// Dynamic address database
type AddressDB struct {
	Exchanges map[string]string            `json:"exchanges"`
	Mixers    map[string]string            `json:"mixers"`
	Hackers   map[string]HackerInfo        `json:"hackers"`
	Contracts map[string]string            `json:"contracts"`
}

type HackerInfo struct {
	Name         string `json:"name"`
	AmountStolen string `json:"amount_stolen"`
	Date         string `json:"date"`
	HackType     string `json:"hack_type"`
}

// Pattern detection
type RiskIndicators struct {
	IsLargeValue      bool
	IsKnownHacker     bool
	IsMixerInteraction bool
	IsSuspiciousMethod bool
	IsFailedTx        bool
	GasAnomalies      bool
}

// Transaction structures - Fixed to match Etherscan response
type TxData struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Result  *TxResult `json:"result"`
}

type TxResult struct {
	BlockHash   string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	Gas         string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	Hash        string `json:"hash"`
	Input       string `json:"input"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Nonce       string `json:"nonce"`
}

type ReceiptData struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      int            `json:"id"`
	Result  *ReceiptResult `json:"result"`
}

type ReceiptResult struct {
	Status  string     `json:"status"`
	Logs    []LogEntry `json:"logs"`
	GasUsed string     `json:"gasUsed"`
	From    string     `json:"from"`
	To      string     `json:"to"`
}

type LogEntry struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

// Known suspicious methods
var suspiciousMethods = map[string]string{
	"0x3ccfd60b": "withdraw() - Potential Reentrancy",
	"0x2e1a7d4d": "withdraw(uint256) - Potential Reentrancy", 
	"0x853828b6": "withdrawAll() - Mass Withdrawal",
	"0x5cffe9de": "flashLoan() - Flash Loan",
	"0xab9c4b5d": "flashLoanSimple() - Flash Loan",
	"0x095ea7b3": "approve() - Token Approval",
	"0x23b872dd": "transferFrom() - Token Transfer",
	"0xa9059cbb": "transfer() - Token Transfer",
	"0x7ff36ab5": "swapExactETHForTokens() - DEX Swap",
	"0x38ed1739": "swapExactTokensForTokens() - DEX Swap",
}

// Load addresses from file
func loadKnownAddresses() (*AddressDB, error) {
	data, err := ioutil.ReadFile("known_addresses.json")
	if err != nil {
		// Return empty DB if file doesn't exist
		return &AddressDB{
			Exchanges: make(map[string]string),
			Mixers:    make(map[string]string),
			Hackers:   make(map[string]HackerInfo),
			Contracts: make(map[string]string),
		}, nil
	}
	
	var db AddressDB
	json.Unmarshal(data, &db)
	return &db, nil
}

// Algorithmic risk detection
func calculateRiskScore(indicators RiskIndicators) (int, []string) {
	score := 0
	var reasons []string
	
	// Pattern-based scoring
	if indicators.IsLargeValue {
		score += 25
		reasons = append(reasons, "Large value transfer detected")
	}
	
	if indicators.IsKnownHacker {
		score += 40
		reasons = append(reasons, "Known malicious address")
	}
	
	if indicators.IsMixerInteraction {
		score += 35
		reasons = append(reasons, "Privacy mixer interaction")
	}
	
	if indicators.IsSuspiciousMethod {
		score += 20
		reasons = append(reasons, "Suspicious contract method")
	}
	
	if indicators.IsFailedTx {
		score += 10
		reasons = append(reasons, "Failed transaction attempt")
	}
	
	if indicators.GasAnomalies {
		score += 15
		reasons = append(reasons, "Abnormal gas usage pattern")
	}
	
	return score, reasons
}

// Detect patterns algorithmically
func detectPatterns(tx *TxResult, receipt *ReceiptResult) RiskIndicators {
	var indicators RiskIndicators
	
	// Large value detection (>100 ETH)
	if tx.Value != "" && tx.Value != "0x0" && tx.Value != "0x" {
		valueBig := new(big.Int)
		valueBig.SetString(strings.TrimPrefix(tx.Value, "0x"), 16)
		threshold := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
		if valueBig.Cmp(threshold) > 0 {
			indicators.IsLargeValue = true
		}
	}
	
	// Gas anomaly detection
	if tx.Gas != "" {
		gasBig := new(big.Int)
		gasBig.SetString(strings.TrimPrefix(tx.Gas, "0x"), 16)
		// Unusually high gas (>1M)
		if gasBig.Cmp(big.NewInt(1000000)) > 0 {
			indicators.GasAnomalies = true
		}
	}
	
	// Failed transaction
	if receipt != nil && receipt.Status == "0x0" {
		indicators.IsFailedTx = true
	}
	
	// Suspicious methods
	if len(tx.Input) >= 10 {
		method := tx.Input[:10]
		if _, ok := suspiciousMethods[method]; ok {
			indicators.IsSuspiciousMethod = true
		}
	}
	
	return indicators
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "analyzer",
		Short: "Dynamic Blockchain Analyzer",
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze [target]",
		Short: "Analyze transaction or address",
		Args:  cobra.ExactArgs(1),
		Run:   analyze,
	}

	rootCmd.AddCommand(analyzeCmd)
	rootCmd.Execute()
}

func analyze(cmd *cobra.Command, args []string) {
	target := args[0]
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	
	if apiKey == "" {
		if data, err := ioutil.ReadFile("enhanced-analyzer-config.json"); err == nil {
			var config map[string]string
			json.Unmarshal(data, &config)
			apiKey = config["etherscan_api_key"]
		}
	}

	if apiKey == "" {
		color.Red("❌ No API key found")
		return
	}

	// Load known addresses
	addressDB, _ := loadKnownAddresses()

	fmt.Println("\n" + strings.Repeat("=", 60))
	color.Cyan("🔍 DYNAMIC BLOCKCHAIN ANALYSIS")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Target: %s\n", target)
	fmt.Printf("Time: %s\n", time.Now().Format("15:04:05"))
	fmt.Println(strings.Repeat("=", 60) + "\n")

	if len(target) == 66 {
		analyzeTxDynamic(target, apiKey, addressDB)
	} else if len(target) == 42 {
		analyzeAddressDynamic(target, apiKey, addressDB)
	}
}

func analyzeTxDynamic(txHash, apiKey string, db *AddressDB) {
	// Fetch transaction
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s", txHash, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		color.Red("❌ Failed to fetch data")
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	
	var txData TxData
	json.Unmarshal(body, &txData)

	if txData.Result == nil || txData.Result.Hash == "" {
		color.Red("❌ Transaction not found")
		return
	}

	tx := txData.Result

	// Fetch receipt
	receiptURL := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%s&apikey=%s", txHash, apiKey)
	receiptResp, _ := http.Get(receiptURL)
	defer receiptResp.Body.Close()
	
	receiptBody, _ := ioutil.ReadAll(receiptResp.Body)
	var receiptData ReceiptData
	json.Unmarshal(receiptBody, &receiptData)

	// Display enhanced info
	color.Yellow("📋 TRANSACTION DETAILS:")
	fmt.Printf("Block: %s\n", tx.BlockNumber)
	fmt.Printf("From: %s", tx.From)
	
	// Check against dynamic database
	fromLabel := checkAddress(tx.From, db)
	if fromLabel != "" {
		color.Red(" [%s]", fromLabel)
	}
	fmt.Println()
	
	fmt.Printf("To:   %s", tx.To)
	toLabel := checkAddress(tx.To, db)
	if toLabel != "" {
		color.Red(" [%s]", toLabel)
	}
	fmt.Println()

	// Parse and display value
	if tx.Value != "" && tx.Value != "0x0" && tx.Value != "0x" {
		valueBig := new(big.Int)
		valueBig.SetString(strings.TrimPrefix(tx.Value, "0x"), 16)
		valueEth := new(big.Float).SetInt(valueBig)
		valueEth.Quo(valueEth, big.NewFloat(1e18))
		
		ethValue, _ := valueEth.Float64()
		fmt.Printf("Value: %.6f ETH ($%.2f)\n", ethValue, ethValue*2350)
		
		if ethValue > 100 {
			color.Red("⚠️  HIGH VALUE TRANSFER!")
		}
	} else {
		fmt.Println("Value: 0 ETH (Contract Interaction)")
	}

	// Method analysis
	if len(tx.Input) >= 10 {
		methodId := tx.Input[:10]
		fmt.Printf("Method: %s", methodId)
		
		if methodName, ok := suspiciousMethods[methodId]; ok {
			color.Red(" - %s", methodName)
		} else if tx.Input == "0x" {
			fmt.Print(" - Simple Transfer")
		} else {
			fmt.Print(" - Contract Interaction")
		}
		fmt.Println()
	}

	// Gas analysis
	if tx.Gas != "" {
		gasBig := new(big.Int)
		gasBig.SetString(strings.TrimPrefix(tx.Gas, "0x"), 16)
		gasLimit := gasBig.Int64()
		
		fmt.Printf("Gas Limit: %d", gasLimit)
		if gasLimit > 1000000 {
			color.Yellow(" (High Gas Usage)")
		}
		fmt.Println()
	}

	// Status from receipt
	if receiptData.Result != nil {
		if receiptData.Result.Status == "0x1" {
			color.Green("Status: ✅ Success")
		} else {
			color.Red("Status: ❌ Failed")
		}
		
		if len(receiptData.Result.Logs) > 0 {
			fmt.Printf("Events: %d logs emitted\n", len(receiptData.Result.Logs))
		}
	}

	// Algorithmic pattern detection
	indicators := detectPatterns(tx, receiptData.Result)
	
	// Check if addresses are known
	if hackerInfo, ok := db.Hackers[strings.ToLower(tx.From)]; ok {
		indicators.IsKnownHacker = true
		color.Red("\n🚨 ALERT: Known hacker - %s", hackerInfo.Name)
		fmt.Printf("   Stolen: %s on %s\n", hackerInfo.AmountStolen, hackerInfo.Date)
		fmt.Printf("   Type: %s\n", hackerInfo.HackType)
	}
	
	if _, ok := db.Mixers[strings.ToLower(tx.To)]; ok {
		indicators.IsMixerInteraction = true
	}

	// Calculate risk score algorithmically
	score, reasons := calculateRiskScore(indicators)
	
	fmt.Println("\n" + strings.Repeat("-", 60))
	color.Cyan("🎯 ALGORITHMIC ANALYSIS:")
	
	if len(reasons) > 0 {
		for _, reason := range reasons {
			fmt.Printf("• %s\n", reason)
		}
	} else {
		fmt.Println("• No risk indicators detected")
	}
	
	// Risk level
	var riskLevel string
	var riskColor color.Attribute
	
	if score >= 80 {
		riskLevel = "CRITICAL"
		riskColor = color.FgRed
	} else if score >= 60 {
		riskLevel = "HIGH"
		riskColor = color.FgHiRed
	} else if score >= 40 {
		riskLevel = "MEDIUM"
		riskColor = color.FgYellow
	} else if score >= 20 {
		riskLevel = "LOW"
		riskColor = color.FgHiYellow
	} else {
		riskLevel = "MINIMAL"
		riskColor = color.FgGreen
	}
	
	fmt.Println()
	color.New(riskColor).Printf("Risk Score: %d/100 (%s)\n", score, riskLevel)
	
	// Smart recommendations based on patterns
	if score > 0 {
		fmt.Println("\n💡 AUTOMATED RECOMMENDATIONS:")
		
		if indicators.IsKnownHacker {
			fmt.Println("• Alert all connected addresses")
			fmt.Println("• Trace fund movement patterns")
			fmt.Println("• Report to exchange security teams")
		}
		
		if indicators.IsMixerInteraction {
			fmt.Println("• Implement enhanced monitoring")
			fmt.Println("• Check for clustering patterns")
			fmt.Println("• Analyze timing correlations")
		}
		
		if indicators.IsLargeValue {
			fmt.Println("• Verify transaction legitimacy")
			fmt.Println("• Check for related transactions")
		}
		
		if indicators.IsSuspiciousMethod {
			fmt.Println("• Review contract interaction")
			fmt.Println("• Check for known exploits using this method")
		}
		
		if indicators.GasAnomalies {
			fmt.Println("• Investigate contract complexity")
			fmt.Println("• Check for infinite loops or DOS attempts")
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 60))
}

func checkAddress(address string, db *AddressDB) string {
	lower := strings.ToLower(address)
	
	if info, ok := db.Hackers[lower]; ok {
		return "🚨 " + info.Name
	}
	
	if name, ok := db.Mixers[lower]; ok {
		return "🌀 " + name
	}
	
	if name, ok := db.Exchanges[lower]; ok {
		return "💱 " + name
	}
	
	if name, ok := db.Contracts[lower]; ok {
		return "📜 " + name
	}
	
	return ""
}

func analyzeAddressDynamic(address, apiKey string, db *AddressDB) {
	color.Yellow("📊 ADDRESS ANALYSIS: %s", address)
	fmt.Println()
	
	// Check if address is in database
	label := checkAddress(address, db)
	if label != "" {
		color.Red("⚠️  IDENTIFIED: %s", label)
		
		if hackerInfo, ok := db.Hackers[strings.ToLower(address)]; ok {
			fmt.Printf("\nDetails:\n")
			fmt.Printf("  Amount Stolen: %s\n", hackerInfo.AmountStolen)
			fmt.Printf("  Date: %s\n", hackerInfo.Date)
			fmt.Printf("  Type: %s\n", hackerInfo.HackType)
		}
	}
	
	// Get balance first
	balanceURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=%s", address, apiKey)
	balanceResp, _ := http.Get(balanceURL)
	defer balanceResp.Body.Close()
	
	balanceBody, _ := ioutil.ReadAll(balanceResp.Body)
	var balanceResult struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}
	json.Unmarshal(balanceBody, &balanceResult)
	
	if balanceResult.Status == "1" && balanceResult.Result != "" {
		balanceBig := new(big.Int)
		balanceBig.SetString(balanceResult.Result, 10)
		balanceEth := new(big.Float).SetInt(balanceBig)
		balanceEth.Quo(balanceEth, big.NewFloat(1e18))
		balance, _ := balanceEth.Float64()
		fmt.Printf("💰 Current Balance: %.6f ETH ($%.2f)\n", balance, balance*2350)
	}
	
	// Fetch recent transactions for pattern analysis
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=desc&apikey=%s", 
		address, apiKey)
	
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Status string `json:"status"`
		Result []struct {
			Hash      string `json:"hash"`
			From      string `json:"from"`
			To        string `json:"to"`
			Value     string `json:"value"`
			TimeStamp string `json:"timeStamp"`
			IsError   string `json:"isError"`
			Input     string `json:"input"`
		} `json:"result"`
	}
	json.Unmarshal(body, &result)
	
	if len(result.Result) == 0 {
		fmt.Println("\n🔇 No recent transactions found")
		return
	}
	
	// Show transaction count
	fmt.Printf("\n📤 Recent Transactions: %d\n", len(result.Result))
	
	// Pattern analysis
	fmt.Println("\n🔬 BEHAVIORAL ANALYSIS:")
	
	// Look for patterns
	mixerInteractions := 0
	largeTransfers := 0
	failedTxs := 0
	suspiciousMethodCount := 0
	totalVolume := big.NewInt(0)
	uniqueAddresses := make(map[string]bool)
	
	for _, tx := range result.Result {
		// Track unique addresses
		if strings.ToLower(tx.From) == strings.ToLower(address) {
			uniqueAddresses[tx.To] = true
			
			// Check for mixer
			if _, ok := db.Mixers[strings.ToLower(tx.To)]; ok {
				mixerInteractions++
			}
		} else {
			uniqueAddresses[tx.From] = true
		}
		
		// Failed transactions
		if tx.IsError == "1" {
			failedTxs++
		}
		
		// Check methods
		if len(tx.Input) >= 10 {
			if _, ok := suspiciousMethods[tx.Input[:10]]; ok {
				suspiciousMethodCount++
			}
		}
		
		// Check for large values
		if tx.Value != "" && tx.Value != "0" {
			valueBig := new(big.Int)
			valueBig.SetString(tx.Value, 10)
			totalVolume.Add(totalVolume, valueBig)
			
			if valueBig.Cmp(new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18))) > 0 {
				largeTransfers++
			}
		}
	}
	
	// Calculate total volume in ETH
	totalEth := new(big.Float).SetInt(totalVolume)
	totalEth.Quo(totalEth, big.NewFloat(1e18))
	volume, _ := totalEth.Float64()
	
	fmt.Printf("• Total volume: %.6f ETH ($%.2f)\n", volume, volume*2350)
	fmt.Printf("• Unique addresses interacted: %d\n", len(uniqueAddresses))
	
	if mixerInteractions > 0 {
		color.Red("• %d mixer interactions detected 🌀", mixerInteractions)
	}
	
	if largeTransfers > 0 {
		color.Yellow("• %d large value transfers (>10 ETH) 💰", largeTransfers)
	}
	
	if failedTxs > 0 {
		color.Yellow("• %d failed transactions ❌", failedTxs)
	}
	
	if suspiciousMethodCount > 0 {
		color.Red("• %d suspicious method calls detected 🚨", suspiciousMethodCount)
	}
	
	// Activity pattern
	if len(result.Result) > 0 {
		// Parse timestamps correctly
		firstTime, _ := strconv.ParseInt(result.Result[len(result.Result)-1].TimeStamp, 10, 64)
		lastTime, _ := strconv.ParseInt(result.Result[0].TimeStamp, 10, 64)
		
		duration := time.Duration(lastTime-firstTime) * time.Second
		
		if duration.Hours() < 24 && len(result.Result) > 5 {
			color.Yellow("• Burst activity pattern detected ⚠️")
		}
		
		fmt.Printf("• Activity span: %.1f days\n", duration.Hours()/24)
	}
	
	// Risk assessment for address
	addressRisk := 0
	if mixerInteractions > 0 {
		addressRisk += 40
	}
	if largeTransfers > 3 {
		addressRisk += 20
	}
	if failedTxs > 2 {
		addressRisk += 15
	}
	if suspiciousMethodCount > 0 {
		addressRisk += 25
	}
	
	if addressRisk > 0 {
		fmt.Println("\n⚠️  ADDRESS RISK ASSESSMENT:")
		
		var riskLevel string
		if addressRisk >= 60 {
			riskLevel = "HIGH RISK"
			color.Red("Risk Score: %d/100 (%s)", addressRisk, riskLevel)
		} else if addressRisk >= 30 {
			riskLevel = "MEDIUM RISK"
			color.Yellow("Risk Score: %d/100 (%s)", addressRisk, riskLevel)
		} else {
			riskLevel = "LOW RISK"
			fmt.Printf("Risk Score: %d/100 (%s)\n", addressRisk, riskLevel)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 60))
}
