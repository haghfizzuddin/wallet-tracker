package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Configuration
type Config struct {
	EtherscanAPIKey string `json:"etherscan_api_key"`
	InfuraURL       string `json:"infura_url"`
}

// Enhanced structures for behavioral analysis
type BehavioralAnalyzer struct {
	config           Config
	knownAddresses   *AddressDB
	riskThresholds   RiskThresholds
	historicalData   map[string]*AddressHistory
	realTimeMonitor  *RealTimeMonitor
}

type RiskThresholds struct {
	HighValueThreshold    float64 // ETH
	VelocityThreshold     int     // transactions per hour
	GasAnomalyMultiplier  float64
	NewAddressAgeMinutes  int
	BenfordDeviationLimit float64
}

type AddressHistory struct {
	FirstSeen        time.Time
	TransactionCount int
	TotalVolume      *big.Int
	AvgGasPrice      *big.Int
	Interactions     map[string]int // address -> count
	TimePattern      []time.Time
	RiskEvents       []RiskEvent
}

type RiskEvent struct {
	Timestamp   time.Time
	EventType   string
	Severity    float64
	Description string
}

type RealTimeMonitor struct {
	etherscanLabels map[string]string
	lastUpdate      time.Time
}

// Enhanced risk analysis
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

// API response structures
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

type EtherscanLabelResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Labels  string `json:"labels"`
	} `json:"result"`
}

// Known suspicious patterns
var suspiciousPatterns = map[string]float64{
	"rapid_drainage":      0.9,  // Multiple large withdrawals in short time
	"mixer_sequence":      0.85, // Typical mixer interaction pattern
	"flash_loan_attack":   0.95, // Flash loan patterns
	"reentrancy_pattern":  0.9,  // Reentrancy attack signature
	"sandwich_attack":     0.8,  // MEV sandwich pattern
	"honeypot_drainage":   0.85, // Draining honeypot contracts
	"circular_transfers":  0.7,  // Money laundering pattern
}

var knownAddresses *AddressDB

func init() {
	// Initialize with some known addresses if file doesn't exist
	knownAddresses = &AddressDB{
		Exchanges: make(map[string]string),
		Mixers:    make(map[string]string),
		Hackers:   make(map[string]HackerInfo),
		Contracts: make(map[string]string),
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "behavioral-analyzer [address]",
		Short: "Advanced blockchain security analyzer with behavioral patterns",
		Args:  cobra.ExactArgs(1),
		Run:   runAnalysis,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runAnalysis(cmd *cobra.Command, args []string) {
	address := strings.ToLower(args[0])
	
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize analyzer
	analyzer := &BehavioralAnalyzer{
		config: config,
		knownAddresses: knownAddresses,
		riskThresholds: RiskThresholds{
			HighValueThreshold:    10.0,  // 10 ETH
			VelocityThreshold:     20,    // 20 tx/hour
			GasAnomalyMultiplier:  3.0,   // 3x average
			NewAddressAgeMinutes:  60,    // 1 hour
			BenfordDeviationLimit: 0.15,  // 15% deviation
		},
		historicalData: make(map[string]*AddressHistory),
		realTimeMonitor: &RealTimeMonitor{
			etherscanLabels: make(map[string]string),
			lastUpdate:      time.Now(),
		},
	}

	// Load known addresses
	if err := analyzer.loadKnownAddresses(); err != nil {
		fmt.Printf("Warning: Could not load known addresses: %v\n", err)
	}

	// Perform comprehensive analysis
	fmt.Printf("\n🔍 Analyzing address: %s\n", color.YellowString(address))
	fmt.Println(strings.Repeat("=", 80))

	analysis, err := analyzer.analyzeAddress(address)
	if err != nil {
		log.Fatal("Analysis failed:", err)
	}

	// Display results
	displayResults(analysis)
}

func (ba *BehavioralAnalyzer) analyzeAddress(address string) (*EnhancedRiskAnalysis, error) {
	analysis := &EnhancedRiskAnalysis{
		Address:          address,
		BehavioralFlags:  []BehavioralFlag{},
		RealTimeFlags:    []string{},
		Recommendations:  []string{},
		DetailedAnalysis: make(map[string]interface{}),
	}

	// 1. Fetch transaction history
	fmt.Println("📊 Fetching transaction history...")
	transactions, err := ba.fetchTransactionHistory(address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	
	if len(transactions) == 0 {
		analysis.RealTimeFlags = append(analysis.RealTimeFlags, "No transaction history found")
		analysis.RiskScore = 0.1
		analysis.Confidence = 0.1
		return analysis, nil
	}

	fmt.Printf("   Found %d transactions\n", len(transactions))

	// 2. Check real-time Etherscan labels
	fmt.Println("🏷️  Checking real-time labels...")
	if labels := ba.checkEtherscanLabels(address); labels != "" {
		analysis.RealTimeFlags = append(analysis.RealTimeFlags, fmt.Sprintf("Etherscan Label: %s", labels))
	}

	// 3. Behavioral pattern analysis
	fmt.Println("🧠 Analyzing behavioral patterns...")
	behavioralFlags := ba.analyzeBehavioralPatterns(address, transactions)
	analysis.BehavioralFlags = append(analysis.BehavioralFlags, behavioralFlags...)

	// 4. Statistical analysis
	fmt.Println("📈 Performing statistical analysis...")
	analysis.StatisticalScores = ba.performStatisticalAnalysis(transactions)

	// 5. Real-time risk indicators
	fmt.Println("⚡ Checking real-time risk indicators...")
	realTimeRisks := ba.checkRealTimeRisks(address, transactions)
	analysis.RealTimeFlags = append(analysis.RealTimeFlags, realTimeRisks...)

	// 6. Calculate final risk score
	analysis.RiskScore, analysis.Confidence = ba.calculateFinalRiskScore(analysis)

	// 7. Generate recommendations
	analysis.Recommendations = ba.generateRecommendations(analysis)

	// 8. Add detailed analysis data
	analysis.DetailedAnalysis["transaction_count"] = len(transactions)
	analysis.DetailedAnalysis["first_tx_time"] = ba.getFirstTransactionTime(transactions)
	analysis.DetailedAnalysis["last_tx_time"] = ba.getLastTransactionTime(transactions)
	analysis.DetailedAnalysis["unique_interactions"] = ba.countUniqueInteractions(transactions)
	analysis.DetailedAnalysis["total_volume_eth"] = ba.calculateTotalVolume(transactions)

	return analysis, nil
}

func (ba *BehavioralAnalyzer) fetchTransactionHistory(address string) ([]Transaction, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=desc&apikey=%s",
		address, ba.config.EtherscanAPIKey)

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
		return []Transaction{}, nil // Return empty list if no transactions
	}

	return result.Result, nil
}

func (ba *BehavioralAnalyzer) analyzeBehavioralPatterns(address string, txs []Transaction) []BehavioralFlag {
	flags := []BehavioralFlag{}

	// 1. Velocity analysis - transactions per time window
	velocityFlag := ba.analyzeVelocity(txs)
	if velocityFlag != nil {
		flags = append(flags, *velocityFlag)
	}

	// 2. Value concentration - large value movements
	concentrationFlag := ba.analyzeValueConcentration(txs)
	if concentrationFlag != nil {
		flags = append(flags, *concentrationFlag)
	}

	// 3. New address behavior
	newAddressFlag := ba.analyzeNewAddressBehavior(txs)
	if newAddressFlag != nil {
		flags = append(flags, *newAddressFlag)
	}

	// 4. Gas anomalies
	gasFlags := ba.analyzeGasAnomalies(txs)
	flags = append(flags, gasFlags...)

	// 5. Interaction patterns
	interactionFlags := ba.analyzeInteractionPatterns(txs)
	flags = append(flags, interactionFlags...)

	// 6. Time-based patterns
	timeFlags := ba.analyzeTimePatterns(txs)
	flags = append(flags, timeFlags...)

	return flags
}

func (ba *BehavioralAnalyzer) analyzeVelocity(txs []Transaction) *BehavioralFlag {
	if len(txs) < 2 {
		return nil
	}

	// Group transactions by hour
	hourlyTxs := make(map[int64]int)
	for _, tx := range txs {
		timestamp, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		hourKey := timestamp / 3600
		hourlyTxs[hourKey]++
	}

	// Find maximum transactions in any hour
	maxTxPerHour := 0
	for _, count := range hourlyTxs {
		if count > maxTxPerHour {
			maxTxPerHour = count
		}
	}

	if maxTxPerHour > ba.riskThresholds.VelocityThreshold {
		return &BehavioralFlag{
			Type:        "high_velocity",
			Severity:    float64(maxTxPerHour) / float64(ba.riskThresholds.VelocityThreshold),
			Description: fmt.Sprintf("Detected %d transactions in one hour (threshold: %d)", maxTxPerHour, ba.riskThresholds.VelocityThreshold),
			Evidence: map[string]interface{}{
				"max_tx_per_hour": maxTxPerHour,
				"threshold":       ba.riskThresholds.VelocityThreshold,
			},
		}
	}

	return nil
}

func (ba *BehavioralAnalyzer) analyzeValueConcentration(txs []Transaction) *BehavioralFlag {
	// Look for rapid accumulation or distribution of funds
	incomingValue := big.NewInt(0)
	outgoingValue := big.NewInt(0)
	
	for _, tx := range txs {
		value, _ := new(big.Int).SetString(tx.Value, 10)
		if strings.EqualFold(tx.To, txs[0].From) {
			incomingValue.Add(incomingValue, value)
		} else {
			outgoingValue.Add(outgoingValue, value)
		}
	}

	// Check for rapid drainage pattern
	if len(txs) >= 5 {
		recentTxs := txs[:5]
		recentOutgoing := big.NewInt(0)
		for _, tx := range recentTxs {
			if !strings.EqualFold(tx.To, txs[0].From) {
				value, _ := new(big.Int).SetString(tx.Value, 10)
				recentOutgoing.Add(recentOutgoing, value)
			}
		}
		
		ethValue := new(big.Float).Quo(new(big.Float).SetInt(recentOutgoing), big.NewFloat(1e18))
		ethFloat, _ := ethValue.Float64()
		
		if ethFloat > ba.riskThresholds.HighValueThreshold {
			return &BehavioralFlag{
				Type:        "rapid_drainage",
				Severity:    0.9,
				Description: fmt.Sprintf("Rapid outgoing transfers detected: %.2f ETH in recent transactions", ethFloat),
				Evidence: map[string]interface{}{
					"recent_outgoing_eth": ethFloat,
					"transaction_count":   5,
				},
			}
		}
	}

	return nil
}

func (ba *BehavioralAnalyzer) analyzeNewAddressBehavior(txs []Transaction) *BehavioralFlag {
	if len(txs) == 0 {
		return nil
	}

	// Get first transaction time
	firstTxTime, _ := strconv.ParseInt(txs[len(txs)-1].TimeStamp, 10, 64)
	addressAge := time.Now().Unix() - firstTxTime
	ageMinutes := addressAge / 60

	// Check if address is new and transacting large amounts
	if ageMinutes < int64(ba.riskThresholds.NewAddressAgeMinutes) {
		totalValue := big.NewInt(0)
		for _, tx := range txs {
			value, _ := new(big.Int).SetString(tx.Value, 10)
			totalValue.Add(totalValue, value)
		}
		
		ethValue := new(big.Float).Quo(new(big.Float).SetInt(totalValue), big.NewFloat(1e18))
		ethFloat, _ := ethValue.Float64()
		
		if ethFloat > ba.riskThresholds.HighValueThreshold {
			return &BehavioralFlag{
				Type:        "new_address_high_value",
				Severity:    0.85,
				Description: fmt.Sprintf("New address (age: %d minutes) transacting %.2f ETH", ageMinutes, ethFloat),
				Evidence: map[string]interface{}{
					"address_age_minutes": ageMinutes,
					"total_value_eth":     ethFloat,
				},
			}
		}
	}

	return nil
}

func (ba *BehavioralAnalyzer) analyzeGasAnomalies(txs []Transaction) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	if len(txs) < 3 {
		return flags
	}

	// Calculate average gas price
	totalGasPrice := big.NewInt(0)
	validTxCount := 0
	
	for _, tx := range txs {
		gasPrice, ok := new(big.Int).SetString(tx.GasPrice, 10)
		if ok && gasPrice.Cmp(big.NewInt(0)) > 0 {
			totalGasPrice.Add(totalGasPrice, gasPrice)
			validTxCount++
		}
	}

	if validTxCount == 0 {
		return flags
	}

	avgGasPrice := new(big.Int).Div(totalGasPrice, big.NewInt(int64(validTxCount)))
	threshold := new(big.Int).Mul(avgGasPrice, big.NewInt(int64(ba.riskThresholds.GasAnomalyMultiplier)))

	// Check for anomalous gas prices
	for _, tx := range txs {
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		if gasPrice.Cmp(threshold) > 0 {
			gweiPrice := new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9))
			gweiFloat, _ := gweiPrice.Float64()
			
			flags = append(flags, BehavioralFlag{
				Type:        "gas_anomaly",
				Severity:    0.7,
				Description: fmt.Sprintf("Transaction %s used abnormally high gas: %.2f Gwei", tx.Hash[:10]+"...", gweiFloat),
				Evidence: map[string]interface{}{
					"tx_hash":     tx.Hash,
					"gas_price":   gweiFloat,
					"avg_gas":     new(big.Float).Quo(new(big.Float).SetInt(avgGasPrice), big.NewFloat(1e9)),
				},
			})
		}
	}

	return flags
}

func (ba *BehavioralAnalyzer) analyzeInteractionPatterns(txs []Transaction) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	// Count interactions with different addresses
	interactions := make(map[string]int)
	methodCalls := make(map[string]int)
	
	for _, tx := range txs {
		if tx.To != "" {
			interactions[strings.ToLower(tx.To)]++
		}
		
		// Analyze method calls
		if len(tx.Input) >= 10 {
			methodId := tx.Input[:10]
			methodCalls[methodId]++
			
			// Check for suspicious methods
			if desc, found := suspiciousPatterns[methodId]; found {
				flags = append(flags, BehavioralFlag{
					Type:        "suspicious_method",
					Severity:    0.8,
					Description: fmt.Sprintf("Called suspicious method: %s", desc),
					Evidence: map[string]interface{}{
						"method_id": methodId,
						"tx_hash":   tx.Hash,
					},
				})
			}
		}
	}

	// Check for mixer interactions
	for addr := range interactions {
		if ba.knownAddresses != nil && ba.knownAddresses.Mixers != nil {
			if mixerName, found := ba.knownAddresses.Mixers[addr]; found {
				flags = append(flags, BehavioralFlag{
					Type:        "mixer_interaction",
					Severity:    0.85,
					Description: fmt.Sprintf("Interacted with known mixer: %s", mixerName),
					Evidence: map[string]interface{}{
						"mixer_address": addr,
						"mixer_name":    mixerName,
						"interactions":  interactions[addr],
					},
				})
			}
		}
	}

	// Check for circular patterns
	if ba.detectCircularPattern(txs) {
		flags = append(flags, BehavioralFlag{
			Type:        "circular_pattern",
			Severity:    0.75,
			Description: "Detected circular transaction pattern (possible money laundering)",
			Evidence: map[string]interface{}{
				"pattern": "circular_transfers",
			},
		})
	}

	return flags
}

func (ba *BehavioralAnalyzer) analyzeTimePatterns(txs []Transaction) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	if len(txs) < 3 {
		return flags
	}

	// Analyze time gaps between transactions
	timeGaps := []int64{}
	for i := 1; i < len(txs); i++ {
		time1, _ := strconv.ParseInt(txs[i-1].TimeStamp, 10, 64)
		time2, _ := strconv.ParseInt(txs[i].TimeStamp, 10, 64)
		gap := time1 - time2 // Note: transactions are sorted desc
		if gap > 0 {
			timeGaps = append(timeGaps, gap)
		}
	}

	// Check for automated/bot behavior (very consistent timing)
	if len(timeGaps) >= 5 {
		variance := calculateVariance(timeGaps)
		if variance < 10.0 { // Very low variance suggests automation
			flags = append(flags, BehavioralFlag{
				Type:        "automated_behavior",
				Severity:    0.6,
				Description: "Transaction timing suggests automated/bot behavior",
				Evidence: map[string]interface{}{
					"timing_variance": variance,
					"sample_size":     len(timeGaps),
				},
			})
		}
	}

	// Check for burst patterns
	burstCount := 0
	for _, gap := range timeGaps {
		if gap < 60 { // Less than 1 minute between transactions
			burstCount++
		}
	}
	
	if float64(burstCount)/float64(len(timeGaps)) > 0.5 {
		flags = append(flags, BehavioralFlag{
			Type:        "burst_pattern",
			Severity:    0.7,
			Description: "Detected burst transaction pattern",
			Evidence: map[string]interface{}{
				"burst_ratio": float64(burstCount) / float64(len(timeGaps)),
				"burst_count": burstCount,
			},
		})
	}

	return flags
}

func (ba *BehavioralAnalyzer) performStatisticalAnalysis(txs []Transaction) StatisticalScores {
	scores := StatisticalScores{}
	
	// 1. Benford's Law analysis on transaction values
	scores.BenfordScore = ba.calculateBenfordScore(txs)
	
	// 2. Velocity score
	scores.VelocityScore = ba.calculateVelocityScore(txs)
	
	// 3. Entropy score for address generation
	scores.EntropyScore = ba.calculateEntropyScore(txs)
	
	// 4. Clustering coefficient
	scores.ClusteringScore = ba.calculateClusteringScore(txs)
	
	// 5. Temporal anomaly score
	scores.TemporalAnomaly = ba.calculateTemporalAnomalyScore(txs)
	
	return scores
}

func (ba *BehavioralAnalyzer) calculateBenfordScore(txs []Transaction) float64 {
	if len(txs) < 10 {
		return 0.0
	}

	// Expected Benford's Law distribution
	benfordExpected := []float64{0.301, 0.176, 0.125, 0.097, 0.079, 0.067, 0.058, 0.051, 0.046}
	
	// Count first digits of transaction values
	digitCounts := make([]int, 9)
	totalCount := 0
	
	for _, tx := range txs {
		value, ok := new(big.Int).SetString(tx.Value, 10)
		if !ok || value.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		
		// Get first digit
		valueStr := value.String()
		if len(valueStr) > 0 && valueStr[0] >= '1' && valueStr[0] <= '9' {
			digit := int(valueStr[0] - '1')
			digitCounts[digit]++
			totalCount++
		}
	}
	
	if totalCount < 10 {
		return 0.0
	}
	
	// Calculate chi-square statistic
	chiSquare := 0.0
	for i := 0; i < 9; i++ {
		observed := float64(digitCounts[i]) / float64(totalCount)
		expected := benfordExpected[i]
		if expected > 0 {
			chiSquare += math.Pow(observed-expected, 2) / expected
		}
	}
	
	// Normalize to 0-1 scale (higher = more suspicious)
	// Chi-square > 0.15 is considered suspicious
	return math.Min(chiSquare/ba.riskThresholds.BenfordDeviationLimit, 1.0)
}

func (ba *BehavioralAnalyzer) calculateVelocityScore(txs []Transaction) float64 {
	if len(txs) < 2 {
		return 0.0
	}

	// Calculate transactions per hour over different time windows
	velocities := []float64{}
	
	// 1-hour windows
	for i := 0; i < len(txs)-1; i++ {
		time1, _ := strconv.ParseInt(txs[i].TimeStamp, 10, 64)
		time2, _ := strconv.ParseInt(txs[i+1].TimeStamp, 10, 64)
		
		timeDiff := float64(time1 - time2)
		if timeDiff > 0 && timeDiff < 3600 { // Within 1 hour
			velocity := 3600.0 / timeDiff
			velocities = append(velocities, velocity)
		}
	}
	
	if len(velocities) == 0 {
		return 0.0
	}
	
	// Get maximum velocity
	maxVelocity := 0.0
	for _, v := range velocities {
		if v > maxVelocity {
			maxVelocity = v
		}
	}
	
	// Normalize (20 tx/hour = 1.0)
	return math.Min(maxVelocity/20.0, 1.0)
}

func (ba *BehavioralAnalyzer) calculateEntropyScore(txs []Transaction) float64 {
	// Calculate entropy of interacted addresses
	interactions := make(map[string]int)
	total := 0
	
	for _, tx := range txs {
		if tx.To != "" {
			interactions[strings.ToLower(tx.To)]++
			total++
		}
	}
	
	if total == 0 {
		return 0.0
	}
	
	entropy := 0.0
	for _, count := range interactions {
		p := float64(count) / float64(total)
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	
	// Normalize by log2(n) where n is number of unique addresses
	if len(interactions) > 1 {
		maxEntropy := math.Log2(float64(len(interactions)))
		return entropy / maxEntropy
	}
	
	return 0.0
}

func (ba *BehavioralAnalyzer) calculateClusteringScore(txs []Transaction) float64 {
	// Build interaction graph
	neighbors := make(map[string]map[string]bool)
	
	for _, tx := range txs {
		from := strings.ToLower(tx.From)
		to := strings.ToLower(tx.To)
		
		if from != "" && to != "" {
			if neighbors[from] == nil {
				neighbors[from] = make(map[string]bool)
			}
			if neighbors[to] == nil {
				neighbors[to] = make(map[string]bool)
			}
			neighbors[from][to] = true
			neighbors[to][from] = true
		}
	}
	
	// Calculate clustering coefficient
	totalCoeff := 0.0
	nodeCount := 0
	
	for node, nodeNeighbors := range neighbors {
		if len(nodeNeighbors) < 2 {
			continue
		}
		
		// Count triangles
		triangles := 0
		neighborList := []string{}
		for n := range nodeNeighbors {
			neighborList = append(neighborList, n)
		}
		
		for i := 0; i < len(neighborList); i++ {
			for j := i + 1; j < len(neighborList); j++ {
				// Check if neighbors are connected
				if neighbors[neighborList[i]] != nil && neighbors[neighborList[i]][neighborList[j]] {
					triangles++
				}
			}
		}
		
		possibleTriangles := len(neighborList) * (len(neighborList) - 1) / 2
		if possibleTriangles > 0 {
			totalCoeff += float64(triangles) / float64(possibleTriangles)
			nodeCount++
		}
		_ = node // Use node variable
	}
	
	if nodeCount > 0 {
		return totalCoeff / float64(nodeCount)
	}
	
	return 0.0
}

func (ba *BehavioralAnalyzer) calculateTemporalAnomalyScore(txs []Transaction) float64 {
	if len(txs) < 3 {
		return 0.0
	}
	
	// Analyze time gaps and look for anomalies
	timeGaps := []float64{}
	for i := 1; i < len(txs); i++ {
		time1, _ := strconv.ParseInt(txs[i-1].TimeStamp, 10, 64)
		time2, _ := strconv.ParseInt(txs[i].TimeStamp, 10, 64)
		gap := float64(time1 - time2)
		if gap > 0 {
			timeGaps = append(timeGaps, gap)
		}
	}
	
	if len(timeGaps) < 2 {
		return 0.0
	}
	
	// Calculate mean and standard deviation
	mean := 0.0
	for _, gap := range timeGaps {
		mean += gap
	}
	mean /= float64(len(timeGaps))
	
	variance := 0.0
	for _, gap := range timeGaps {
		variance += math.Pow(gap-mean, 2)
	}
	variance /= float64(len(timeGaps))
	stdDev := math.Sqrt(variance)
	
	// Count anomalies (gaps > 2 standard deviations from mean)
	anomalyCount := 0
	for _, gap := range timeGaps {
		if math.Abs(gap-mean) > 2*stdDev {
			anomalyCount++
		}
	}
	
	// Return ratio of anomalies
	return float64(anomalyCount) / float64(len(timeGaps))
}

func (ba *BehavioralAnalyzer) checkEtherscanLabels(address string) string {
	// In a real implementation, this would query Etherscan's label API
	// For now, we'll check our known addresses
	if ba.knownAddresses != nil {
		if ba.knownAddresses.Exchanges != nil {
			if name, found := ba.knownAddresses.Exchanges[address]; found {
				return fmt.Sprintf("Exchange: %s", name)
			}
		}
		if ba.knownAddresses.Mixers != nil {
			if name, found := ba.knownAddresses.Mixers[address]; found {
				return fmt.Sprintf("Mixer: %s", name)
			}
		}
		if ba.knownAddresses.Hackers != nil {
			if info, found := ba.knownAddresses.Hackers[address]; found {
				return fmt.Sprintf("Known Hacker: %s", info.Name)
			}
		}
		if ba.knownAddresses.Contracts != nil {
			if name, found := ba.knownAddresses.Contracts[address]; found {
				return fmt.Sprintf("Contract: %s", name)
			}
		}
	}
	return ""
}

func (ba *BehavioralAnalyzer) checkRealTimeRisks(address string, txs []Transaction) []string {
	risks := []string{}
	
	// Check for failed transactions
	failedCount := 0
	for _, tx := range txs {
		if tx.IsError == "1" {
			failedCount++
		}
	}
	
	if failedCount > len(txs)/4 {
		risks = append(risks, fmt.Sprintf("High failure rate: %d/%d transactions failed", failedCount, len(txs)))
	}
	
	// Check for contract creation
	contractsCreated := 0
	for _, tx := range txs {
		if tx.ContractAddress != "" {
			contractsCreated++
		}
	}
	
	if contractsCreated > 0 {
		risks = append(risks, fmt.Sprintf("Created %d contracts", contractsCreated))
	}
	
	// Check for interaction with known suspicious addresses
	suspiciousInteractions := ba.checkSuspiciousInteractions(txs)
	if len(suspiciousInteractions) > 0 {
		risks = append(risks, suspiciousInteractions...)
	}
	
	return risks
}

func (ba *BehavioralAnalyzer) checkSuspiciousInteractions(txs []Transaction) []string {
	interactions := []string{}
	checked := make(map[string]bool)
	
	for _, tx := range txs {
		to := strings.ToLower(tx.To)
		if to != "" && !checked[to] {
			checked[to] = true
			
			// Check if it's a known malicious address
			if ba.knownAddresses != nil && ba.knownAddresses.Hackers != nil {
				if info, found := ba.knownAddresses.Hackers[to]; found {
					interactions = append(interactions, fmt.Sprintf("Interacted with known hacker: %s (%s)", info.Name, to[:10]+"..."))
				}
			}
		}
	}
	
	return interactions
}

func (ba *BehavioralAnalyzer) calculateFinalRiskScore(analysis *EnhancedRiskAnalysis) (float64, float64) {
	totalScore := 0.0
	totalWeight := 0.0
	
	// Behavioral flags contribution
	for _, flag := range analysis.BehavioralFlags {
		weight := flag.Severity
		totalScore += flag.Severity * weight
		totalWeight += weight
	}
	
	// Statistical scores contribution
	statWeights := map[string]float64{
		"benford":   0.7,
		"velocity":  0.8,
		"entropy":   0.5,
		"clustering": 0.6,
		"temporal":  0.7,
	}
	
	totalScore += analysis.StatisticalScores.BenfordScore * statWeights["benford"]
	totalScore += analysis.StatisticalScores.VelocityScore * statWeights["velocity"]
	totalScore += analysis.StatisticalScores.EntropyScore * statWeights["entropy"]
	totalScore += analysis.StatisticalScores.ClusteringScore * statWeights["clustering"]
	totalScore += analysis.StatisticalScores.TemporalAnomaly * statWeights["temporal"]
	
	for _, weight := range statWeights {
		totalWeight += weight
	}
	
	// Real-time flags contribution
	realTimeWeight := 0.9
	if len(analysis.RealTimeFlags) > 0 {
		totalScore += float64(len(analysis.RealTimeFlags)) * 0.2 * realTimeWeight
		totalWeight += realTimeWeight
	}
	
	// Calculate weighted average
	finalScore := 0.0
	if totalWeight > 0 {
		finalScore = totalScore / totalWeight
	}
	
	// Calculate confidence based on amount of data
	txCount := 0
	if val, ok := analysis.DetailedAnalysis["transaction_count"]; ok {
		if count, ok := val.(int); ok {
			txCount = count
		}
	}
	
	confidence := math.Min(float64(len(analysis.BehavioralFlags))/5.0, 1.0) * 0.5 +
		math.Min(float64(txCount)/50.0, 1.0) * 0.5
	
	return math.Min(finalScore, 1.0), confidence
}

func (ba *BehavioralAnalyzer) generateRecommendations(analysis *EnhancedRiskAnalysis) []string {
	recommendations := []string{}
	
	if analysis.RiskScore > 0.8 {
		recommendations = append(recommendations, "⚠️  CRITICAL: This address shows multiple high-risk indicators. Avoid interaction.")
		recommendations = append(recommendations, "📞 Report to relevant authorities if you've been affected.")
	} else if analysis.RiskScore > 0.6 {
		recommendations = append(recommendations, "⚡ HIGH RISK: Exercise extreme caution with this address.")
		recommendations = append(recommendations, "🔍 Perform additional due diligence before any interaction.")
	} else if analysis.RiskScore > 0.4 {
		recommendations = append(recommendations, "⚠️  MEDIUM RISK: Some suspicious patterns detected.")
		recommendations = append(recommendations, "💡 Monitor this address closely for further activity.")
	} else if analysis.RiskScore > 0.2 {
		recommendations = append(recommendations, "ℹ️  LOW RISK: Minor anomalies detected.")
		recommendations = append(recommendations, "✅ Generally safe but maintain standard precautions.")
	} else {
		recommendations = append(recommendations, "✅ MINIMAL RISK: No significant suspicious patterns detected.")
	}
	
	// Specific recommendations based on flags
	for _, flag := range analysis.BehavioralFlags {
		switch flag.Type {
		case "mixer_interaction":
			recommendations = append(recommendations, "🌀 Consider blockchain analysis tools to trace fund origins.")
		case "high_velocity":
			recommendations = append(recommendations, "⏱️  Monitor for potential automated attack patterns.")
		case "new_address_high_value":
			recommendations = append(recommendations, "🆕 Verify the legitimacy of this newly active address.")
		}
	}
	
	return recommendations
}

// Helper functions
func (ba *BehavioralAnalyzer) detectCircularPattern(txs []Transaction) bool {
	if len(txs) < 4 {
		return false
	}
	
	// Simple circular detection: A -> B -> C -> A
	fromTo := make(map[string][]string)
	
	for _, tx := range txs {
		from := strings.ToLower(tx.From)
		to := strings.ToLower(tx.To)
		if from != "" && to != "" {
			fromTo[from] = append(fromTo[from], to)
		}
	}
	
	// Look for cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var hasCycle func(string) bool
	hasCycle = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		
		for _, neighbor := range fromTo[node] {
			if !visited[neighbor] {
				if hasCycle(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}
		
		recStack[node] = false
		return false
	}
	
	for node := range fromTo {
		if !visited[node] {
			if hasCycle(node) {
				return true
			}
		}
	}
	
	return false
}

func calculateVariance(values []int64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	
	// Calculate mean
	sum := int64(0)
	for _, v := range values {
		sum += v
	}
	mean := float64(sum) / float64(len(values))
	
	// Calculate variance
	variance := 0.0
	for _, v := range values {
		variance += math.Pow(float64(v)-mean, 2)
	}
	
	return variance / float64(len(values))
}

func (ba *BehavioralAnalyzer) getFirstTransactionTime(txs []Transaction) string {
	if len(txs) == 0 {
		return "N/A"
	}
	
	timestamp, _ := strconv.ParseInt(txs[len(txs)-1].TimeStamp, 10, 64)
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func (ba *BehavioralAnalyzer) getLastTransactionTime(txs []Transaction) string {
	if len(txs) == 0 {
		return "N/A"
	}
	
	timestamp, _ := strconv.ParseInt(txs[0].TimeStamp, 10, 64)
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func (ba *BehavioralAnalyzer) countUniqueInteractions(txs []Transaction) int {
	unique := make(map[string]bool)
	
	for _, tx := range txs {
		if tx.To != "" {
			unique[strings.ToLower(tx.To)] = true
		}
	}
	
	return len(unique)
}

func (ba *BehavioralAnalyzer) calculateTotalVolume(txs []Transaction) string {
	total := big.NewInt(0)
	
	for _, tx := range txs {
		value, _ := new(big.Int).SetString(tx.Value, 10)
		total.Add(total, value)
	}
	
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(total), big.NewFloat(1e18))
	return fmt.Sprintf("%.4f ETH", ethValue)
}

func (ba *BehavioralAnalyzer) loadKnownAddresses() error {
	data, err := ioutil.ReadFile("known_addresses.json")
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

func loadConfig() (Config, error) {
	// First try to load from file
	data, err := ioutil.ReadFile("enhanced-analyzer-config.json")
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
		config.EtherscanAPIKey = "YourEtherscanAPIKey" // Replace with actual key
	}
	
	return config, nil
}

func displayResults(analysis *EnhancedRiskAnalysis) {
	// Header
	fmt.Println("\n📊 ANALYSIS RESULTS")
	fmt.Println(strings.Repeat("=", 80))
	
	// Risk Score with color
	riskColor := color.New(color.FgGreen)
	if analysis.RiskScore > 0.8 {
		riskColor = color.New(color.FgRed, color.Bold)
	} else if analysis.RiskScore > 0.6 {
		riskColor = color.New(color.FgRed)
	} else if analysis.RiskScore > 0.4 {
		riskColor = color.New(color.FgYellow)
	}
	
	fmt.Printf("🎯 Risk Score: %s (Confidence: %.1f%%)\n", 
		riskColor.Sprintf("%.2f/1.00", analysis.RiskScore),
		analysis.Confidence*100)
	
	// Statistical Scores
	fmt.Println("\n📈 Statistical Analysis:")
	fmt.Printf("   • Benford's Law Score: %.2f\n", analysis.StatisticalScores.BenfordScore)
	fmt.Printf("   • Velocity Score: %.2f\n", analysis.StatisticalScores.VelocityScore)
	fmt.Printf("   • Entropy Score: %.2f\n", analysis.StatisticalScores.EntropyScore)
	fmt.Printf("   • Clustering Score: %.2f\n", analysis.StatisticalScores.ClusteringScore)
	fmt.Printf("   • Temporal Anomaly: %.2f\n", analysis.StatisticalScores.TemporalAnomaly)
	
	// Behavioral Flags
	if len(analysis.BehavioralFlags) > 0 {
		fmt.Println("\n🚩 Behavioral Patterns Detected:")
		// Sort by severity
		sort.Slice(analysis.BehavioralFlags, func(i, j int) bool {
			return analysis.BehavioralFlags[i].Severity > analysis.BehavioralFlags[j].Severity
		})
		
		for _, flag := range analysis.BehavioralFlags {
			severityColor := color.New(color.FgYellow)
			if flag.Severity > 0.8 {
				severityColor = color.New(color.FgRed, color.Bold)
			} else if flag.Severity > 0.6 {
				severityColor = color.New(color.FgRed)
			}
			
			fmt.Printf("   • %s [Severity: %s]\n", 
				flag.Description,
				severityColor.Sprintf("%.2f", flag.Severity))
		}
	}
	
	// Real-time Flags
	if len(analysis.RealTimeFlags) > 0 {
		fmt.Println("\n⚡ Real-time Indicators:")
		for _, flag := range analysis.RealTimeFlags {
			fmt.Printf("   • %s\n", flag)
		}
	}
	
	// Transaction Details
	fmt.Println("\n📋 Transaction Summary:")
	fmt.Printf("   • Total Transactions: %d\n", analysis.DetailedAnalysis["transaction_count"])
	fmt.Printf("   • First Transaction: %s\n", analysis.DetailedAnalysis["first_tx_time"])
	fmt.Printf("   • Last Transaction: %s\n", analysis.DetailedAnalysis["last_tx_time"])
	fmt.Printf("   • Unique Interactions: %d\n", analysis.DetailedAnalysis["unique_interactions"])
	fmt.Printf("   • Total Volume: %s\n", analysis.DetailedAnalysis["total_volume_eth"])
	
	// Recommendations
	if len(analysis.Recommendations) > 0 {
		fmt.Println("\n💡 Recommendations:")
		for _, rec := range analysis.Recommendations {
			fmt.Printf("   %s\n", rec)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 80))
}

// Structs from the original code that were referenced
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
