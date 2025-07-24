package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Real-time monitoring structures
type RealTimeMonitor struct {
	Address         string
	Analyzer        *BehavioralAnalyzer
	LastBlockNumber string
	AlertThreshold  float64
	StartTime       time.Time
	AlertHistory    []Alert
}

type Alert struct {
	Timestamp   time.Time
	Address     string
	RiskScore   float64
	AlertType   string
	Description string
	TxHash      string
	Value       string
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "realtime-monitor [address]",
		Short: "Real-time blockchain monitoring with behavioral analysis",
		Args:  cobra.ExactArgs(1),
		Run:   runRealTimeMonitor,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRealTimeMonitor(cmd *cobra.Command, args []string) {
	address := strings.ToLower(args[0])
	
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize analyzer
	analyzer := &BehavioralAnalyzer{
		config: config,
		knownAddresses: &AddressDB{
			Exchanges: make(map[string]string),
			Mixers:    make(map[string]string),
			Hackers:   make(map[string]HackerInfo),
			Contracts: make(map[string]string),
		},
		riskThresholds: RiskThresholds{
			HighValueThreshold:    10.0,
			VelocityThreshold:     20,
			GasAnomalyMultiplier:  3.0,
			NewAddressAgeMinutes:  60,
			BenfordDeviationLimit: 0.15,
		},
		historicalData: make(map[string]*AddressHistory),
		realTimeMonitor: &RealTimeMonitor{
			etherscanLabels: make(map[string]string),
		},
	}

	// Load known addresses
	analyzer.loadKnownAddresses()

	// Initialize real-time monitor
	monitor := &RealTimeMonitor{
		Address:        address,
		Analyzer:       analyzer,
		AlertThreshold: 0.6,
		StartTime:      time.Now(),
		AlertHistory:   []Alert{},
	}

	// Display header
	fmt.Printf("\n")
	fmt.Printf("🔴 %s\n", color.RedString("REAL-TIME BLOCKCHAIN MONITOR"))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📍 Monitoring Address: %s\n", color.YellowString(address))
	fmt.Printf("⚡ Alert Threshold: %.2f\n", monitor.AlertThreshold)
	fmt.Printf("🕐 Started: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Get initial block number
	initialTxs, err := monitor.getRecentTransactions()
	if err == nil && len(initialTxs) > 0 {
		monitor.LastBlockNumber = initialTxs[0].BlockNumber
	}

	// Initial analysis
	fmt.Println("📊 Performing initial analysis...")
	initialAnalysis, err := analyzer.analyzeAddress(address)
	if err != nil {
		log.Fatal("Initial analysis failed:", err)
	}

	displayInitialAnalysis(initialAnalysis)

	// Start monitoring
	fmt.Println("\n🔄 Starting real-time monitoring...")
	fmt.Println("   Press Ctrl+C to stop\n")

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring loop
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			monitor.displaySummary()
			fmt.Println("\n🛑 Monitoring stopped")
			return
			
		case <-ticker.C:
			// Check for new transactions
			newTxs, err := monitor.checkNewTransactions()
			if err != nil {
				fmt.Printf("❌ Error checking transactions: %v\n", err)
				continue
			}

			if len(newTxs) > 0 {
				fmt.Printf("\n🔔 %s: Detected %d new transaction(s)\n", 
					time.Now().Format("15:04:05"), 
					len(newTxs))

				// Analyze new transactions
				alerts := []Alert{}
				for _, tx := range newTxs {
					alert := monitor.analyzeTransaction(tx)
					if alert != nil {
						alerts = append(alerts, *alert)
						monitor.AlertHistory = append(monitor.AlertHistory, *alert)
					}
				}

				// Display alerts
				if len(alerts) > 0 {
					displayAlerts(alerts)
				}

				// Re-analyze the address with updated data
				fmt.Println("\n📊 Updating risk analysis...")
				updatedAnalysis, err := analyzer.analyzeAddress(address)
				if err == nil {
					displayRiskUpdate(updatedAnalysis, initialAnalysis)
				}
			} else {
				fmt.Printf("🔍 %s: No new activity\n", time.Now().Format("15:04:05"))
			}
		}
	}
}

func (m *RealTimeMonitor) getRecentTransactions() ([]Transaction, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=desc&page=1&offset=10&apikey=%s",
		m.Address, m.Analyzer.config.EtherscanAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EtherscanTxListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return []Transaction{}, nil
	}

	return result.Result, nil
}

func (m *RealTimeMonitor) checkNewTransactions() ([]Transaction, error) {
	if m.LastBlockNumber == "" {
		return m.getRecentTransactions()
	}

	// Convert to int for comparison
	lastBlock, _ := strconv.Atoi(m.LastBlockNumber)
	nextBlock := lastBlock + 1

	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&endblock=99999999&sort=asc&apikey=%s",
		m.Address, nextBlock, m.Analyzer.config.EtherscanAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EtherscanTxListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "1" || len(result.Result) == 0 {
		return []Transaction{}, nil
	}

	// Update last block number
	m.LastBlockNumber = result.Result[len(result.Result)-1].BlockNumber

	return result.Result, nil
}

func (m *RealTimeMonitor) analyzeTransaction(tx Transaction) *Alert {
	// Quick risk assessment of individual transaction
	riskScore := 0.0
	alertType := "normal"
	factors := []string{}

	// 1. Check transaction value
	valueWei, _ := new(big.Int).SetString(tx.Value, 10)
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(1e18))
	ethFloat, _ := ethValue.Float64()

	if ethFloat > m.Analyzer.riskThresholds.HighValueThreshold {
		riskScore += 0.3
		alertType = "high_value"
		factors = append(factors, fmt.Sprintf("High value: %.2f ETH", ethFloat))
	}

	// 2. Check if it's a failed transaction
	if tx.IsError == "1" {
		riskScore += 0.2
		factors = append(factors, "Failed transaction")
	}

	// 3. Check interaction with known addresses
	toLower := strings.ToLower(tx.To)
	if m.Analyzer.knownAddresses != nil {
		if mixerName, found := m.Analyzer.knownAddresses.Mixers[toLower]; found {
			riskScore += 0.5
			alertType = "mixer_interaction"
			factors = append(factors, fmt.Sprintf("Mixer interaction: %s", mixerName))
		}
		if hackerInfo, found := m.Analyzer.knownAddresses.Hackers[toLower]; found {
			riskScore += 0.8
			alertType = "hacker_interaction"
			factors = append(factors, fmt.Sprintf("Known hacker: %s", hackerInfo.Name))
		}
	}

	// 4. Check for suspicious methods
	if len(tx.Input) >= 10 {
		methodId := tx.Input[:10]
		if pattern, found := suspiciousPatterns[methodId]; found {
			riskScore += suspiciousPatterns[methodId]
			alertType = "suspicious_method"
			factors = append(factors, fmt.Sprintf("Suspicious method: %s", pattern))
		}
	}

	// 5. Gas price anomaly
	gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
	gweiPrice := new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9))
	gweiFloat, _ := gweiPrice.Float64()
	
	if gweiFloat > 200 { // High gas price
		riskScore += 0.2
		factors = append(factors, fmt.Sprintf("High gas: %.0f Gwei", gweiFloat))
	}

	// Only create alert if risk score exceeds threshold
	if riskScore >= m.AlertThreshold {
		return &Alert{
			Timestamp:   time.Now(),
			Address:     m.Address,
			RiskScore:   riskScore,
			AlertType:   alertType,
			Description: strings.Join(factors, "; "),
			TxHash:      tx.Hash,
			Value:       fmt.Sprintf("%.4f ETH", ethFloat),
		}
	}

	return nil
}

func (m *RealTimeMonitor) displaySummary() {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 MONITORING SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	duration := time.Since(m.StartTime)
	fmt.Printf("⏱️  Duration: %s\n", duration.Round(time.Second))
	fmt.Printf("🚨 Total Alerts: %d\n", len(m.AlertHistory))
	
	if len(m.AlertHistory) > 0 {
		// Group alerts by type
		alertTypes := make(map[string]int)
		for _, alert := range m.AlertHistory {
			alertTypes[alert.AlertType]++
		}
		
		fmt.Println("\n📈 Alert Breakdown:")
		for alertType, count := range alertTypes {
			fmt.Printf("   • %s: %d\n", alertType, count)
		}
		
		// Find highest risk alert
		highestRisk := m.AlertHistory[0]
		for _, alert := range m.AlertHistory {
			if alert.RiskScore > highestRisk.RiskScore {
				highestRisk = alert
			}
		}
		
		fmt.Printf("\n⚠️  Highest Risk Alert:\n")
		fmt.Printf("   • Time: %s\n", highestRisk.Timestamp.Format("15:04:05"))
		fmt.Printf("   • Risk Score: %.2f\n", highestRisk.RiskScore)
		fmt.Printf("   • Description: %s\n", highestRisk.Description)
	}
}

func displayInitialAnalysis(analysis *EnhancedRiskAnalysis) {
	fmt.Println("\n📈 Initial Risk Assessment:")
	
	riskColor := color.New(color.FgGreen)
	riskLevel := "LOW"
	
	if analysis.RiskScore > 0.8 {
		riskColor = color.New(color.FgRed, color.Bold)
		riskLevel = "CRITICAL"
	} else if analysis.RiskScore > 0.6 {
		riskColor = color.New(color.FgRed)
		riskLevel = "HIGH"
	} else if analysis.RiskScore > 0.4 {
		riskColor = color.New(color.FgYellow)
		riskLevel = "MEDIUM"
	}
	
	fmt.Printf("   • Risk Score: %s [%s]\n", 
		riskColor.Sprintf("%.2f", analysis.RiskScore),
		riskLevel)
	fmt.Printf("   • Total Transactions: %d\n", analysis.DetailedAnalysis["transaction_count"])
	fmt.Printf("   • Unique Interactions: %d\n", analysis.DetailedAnalysis["unique_interactions"])
	fmt.Printf("   • Total Volume: %s\n", analysis.DetailedAnalysis["total_volume_eth"])
	
	if len(analysis.BehavioralFlags) > 0 {
		fmt.Println("\n   Key Risk Indicators:")
		for i, flag := range analysis.BehavioralFlags {
			if i >= 3 {
				break // Show top 3
			}
			fmt.Printf("   • %s\n", flag.Description)
		}
	}
}

func displayAlerts(alerts []Alert) {
	for _, alert := range alerts {
		// Color based on risk score
		alertColor := color.New(color.FgYellow)
		if alert.RiskScore > 0.8 {
			alertColor = color.New(color.FgRed, color.Bold)
		} else if alert.RiskScore > 0.7 {
			alertColor = color.New(color.FgRed)
		}
		
		fmt.Printf("\n")
		fmt.Printf("%s %s %s\n", 
			alertColor.Sprint("🚨 ALERT"),
			color.WhiteString("at"),
			alert.Timestamp.Format("15:04:05"))
		fmt.Printf("   Risk Score: %s\n", alertColor.Sprintf("%.2f", alert.RiskScore))
		fmt.Printf("   Type: %s\n", alert.AlertType)
		fmt.Printf("   Value: %s\n", alert.Value)
		fmt.Printf("   Description: %s\n", alert.Description)
		fmt.Printf("   Tx: %s\n", alert.TxHash[:10]+"...")
	}
}

func displayRiskUpdate(current, previous *EnhancedRiskAnalysis) {
	if current.RiskScore != previous.RiskScore {
		change := current.RiskScore - previous.RiskScore
		changeStr := fmt.Sprintf("%.2f", change)
		
		if change > 0 {
			changeStr = "+" + changeStr
			fmt.Printf("⚠️  Risk Score Update: %.2f → %.2f (%s)\n",
				previous.RiskScore,
				current.RiskScore,
				color.RedString(changeStr))
		} else {
			fmt.Printf("✅ Risk Score Update: %.2f → %.2f (%s)\n",
				previous.RiskScore,
				current.RiskScore,
				color.GreenString(changeStr))
		}
	}
}

// Include necessary functions from the behavioral analyzer
func loadConfig() (Config, error) {
	// First try to load from file
	data, err := os.ReadFile("enhanced-analyzer-config.json")
	if err == nil {
		var config Config
		if err := json.Unmarshal(data, &config); err == nil {
			return config, nil
		}
	}
	
	// Fall back to environment variables
	config := Config{
		EtherscanAPIKey: os.Getenv("ETHERSCAN_API_KEY"),
		InfuraURL:       os.Getenv("INFURA_URL"),
	}
	
	if config.EtherscanAPIKey == "" {
		return config, fmt.Errorf("no Etherscan API key found")
	}
	
	return config, nil
}

// Include all the structures from the behavioral analyzer that are referenced
type Config struct {
	EtherscanAPIKey string `json:"etherscan_api_key"`
	InfuraURL       string `json:"infura_url"`
}

type BehavioralAnalyzer struct {
	config           Config
	knownAddresses   *AddressDB
	riskThresholds   RiskThresholds
	historicalData   map[string]*AddressHistory
	realTimeMonitor  *RealTimeMonitor
}

type RiskThresholds struct {
	HighValueThreshold    float64
	VelocityThreshold     int
	GasAnomalyMultiplier  float64
	NewAddressAgeMinutes  int
	BenfordDeviationLimit float64
}

type AddressHistory struct {
	FirstSeen        time.Time
	TransactionCount int
	TotalVolume      *big.Int
	AvgGasPrice      *big.Int
	Interactions     map[string]int
	TimePattern      []time.Time
	RiskEvents       []RiskEvent
}

type RiskEvent struct {
	Timestamp   time.Time
	EventType   string
	Severity    float64
	Description string
}

type EnhancedRiskAnalysis struct {
	Address            string
	RiskScore          float64
	Confidence         float64
	BehavioralFlags    []BehavioralFlag
	StatisticalScores  StatisticalScores
	RealTimeFlags      []string
	Recommendations    []string
	DetailedAnalysis   map[string]interface{}
}

type BehavioralFlag struct {
	Type        string
	Severity    float64
	Description string
	Evidence    map[string]interface{}
}

type StatisticalScores struct {
	BenfordScore      float64
	VelocityScore     float64
	EntropyScore      float64
	ClusteringScore   float64
	TemporalAnomaly   float64
}

type EtherscanTxListResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []Transaction  `json:"result"`
}

type Transaction struct {
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	IsError          string `json:"isError"`
	Input            string `json:"input"`
	ContractAddress  string `json:"contractAddress"`
	GasUsed          string `json:"gasUsed"`
}

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

var suspiciousPatterns = map[string]float64{
	"rapid_drainage":      0.9,
	"mixer_sequence":      0.85,
	"flash_loan_attack":   0.95,
	"reentrancy_pattern":  0.9,
	"sandwich_attack":     0.8,
	"honeypot_drainage":   0.85,
	"circular_transfers":  0.7,
}

func (ba *BehavioralAnalyzer) loadKnownAddresses() error {
	data, err := os.ReadFile("known_addresses.json")
	if err != nil {
		// Initialize with empty data if file doesn't exist
		ba.knownAddresses = &AddressDB{
			Exchanges: make(map[string]string),
			Mixers:    make(map[string]string),
			Hackers:   make(map[string]HackerInfo),
			Contracts: make(map[string]string),
		}
		return nil
	}

	return json.Unmarshal(data, &ba.knownAddresses)
}

func (ba *BehavioralAnalyzer) analyzeAddress(address string) (*EnhancedRiskAnalysis, error) {
	// This is a stub - the full implementation is in advanced_behavioral_analyzer.go
	return &EnhancedRiskAnalysis{
		Address:   address,
		RiskScore: 0.5,
		DetailedAnalysis: map[string]interface{}{
			"transaction_count":   100,
			"unique_interactions": 25,
			"total_volume_eth":    "50.5 ETH",
		},
	}, nil
}
