package main

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

// Gas analysis patterns and thresholds
type GasPattern struct {
	Type        string
	Description string
	Indicators  []string
	Severity    float64
}

var gasPatterns = map[string]GasPattern{
	"front_running": {
		Type:        "front_running",
		Description: "Front-running attack pattern detected",
		Indicators: []string{
			"Gas price 10x+ above average",
			"Transaction immediately before large trades",
			"Targeting DEX routers or popular contracts",
		},
		Severity: 0.85,
	},
	"sandwich_attack": {
		Type:        "sandwich_attack",
		Description: "MEV sandwich attack pattern",
		Indicators: []string{
			"Multiple high-gas transactions in sequence",
			"Buy-sell pattern around victim transactions",
			"Consistent profit extraction",
		},
		Severity: 0.9,
	},
	"gas_war": {
		Type:        "gas_war",
		Description: "Gas war competition detected",
		Indicators: []string{
			"Escalating gas prices over short period",
			"Multiple failed transactions",
			"Targeting limited resources (NFTs, tokens)",
		},
		Severity: 0.6,
	},
	"exploit_execution": {
		Type:        "exploit_execution",
		Description: "Potential exploit execution",
		Indicators: []string{
			"Single extremely high gas transaction",
			"Complex contract interaction",
			"Large value extraction",
		},
		Severity: 0.95,
	},
	"censorship_evasion": {
		Type:        "censorship_evasion",
		Description: "Censorship evasion pattern",
		Indicators: []string{
			"Consistently high gas across transactions",
			"Interactions with mixers or blacklisted addresses",
			"Rapid fund movement",
		},
		Severity: 0.8,
	},
}

// Enhanced gas analysis function
func (ba *BehavioralAnalyzer) analyzeGasAnomaliesEnhanced(txs []Transaction) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	if len(txs) < 3 {
		return flags
	}

	// Calculate gas statistics
	gasPrices := []int64{}
	gasUsage := make(map[string][]GasInfo) // timestamp -> gas info
	
	for _, tx := range txs {
		gasPrice, ok := new(big.Int).SetString(tx.GasPrice, 10)
		if !ok || gasPrice.Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		
		gasPriceInt := gasPrice.Int64()
		gasPrices = append(gasPrices, gasPriceInt)
		
		timestamp := tx.TimeStamp
		gasUsage[timestamp] = append(gasUsage[timestamp], GasInfo{
			GasPrice: gasPriceInt,
			TxHash:   tx.Hash,
			To:       tx.To,
			Value:    tx.Value,
			Input:    tx.Input,
			IsError:  tx.IsError,
		})
	}

	if len(gasPrices) == 0 {
		return flags
	}

	// Sort gas prices for percentile calculation
	sort.Slice(gasPrices, func(i, j int) bool {
		return gasPrices[i] < gasPrices[j]
	})

	// Calculate percentiles
	p50 := gasPrices[len(gasPrices)/2]
	p90 := gasPrices[int(float64(len(gasPrices))*0.9)]
	p99 := gasPrices[int(float64(len(gasPrices))*0.99)]

	// Detect patterns
	frontRunFlags := ba.detectFrontRunning(txs, p90)
	flags = append(flags, frontRunFlags...)

	sandwichFlags := ba.detectSandwichAttacks(txs, p90)
	flags = append(flags, sandwichFlags...)

	gasWarFlags := ba.detectGasWars(gasUsage, p50)
	flags = append(flags, gasWarFlags...)

	exploitFlags := ba.detectExploitGasPattern(txs, p99)
	flags = append(flags, exploitFlags...)

	censorshipFlags := ba.detectCensorshipEvasion(gasPrices, p50)
	flags = append(flags, censorshipFlags...)

	return flags
}

// Detect front-running patterns
func (ba *BehavioralAnalyzer) detectFrontRunning(txs []Transaction, p90 int64) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	// Check for transactions with gas price > 10x median targeting DEX routers
	dexRouters := map[string]bool{
		"0x7a250d5630b4cf539739df2c5dacb4c659f2488d": true, // Uniswap V2
		"0xe592427a0aece92de3edee1f18e0157c05861564": true, // Uniswap V3
		"0xd9e1ce17f2641f24ae83637ab66a2cca9c378b9f": true, // SushiSwap
	}

	for i, tx := range txs {
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		if gasPrice.Int64() > p90*10 {
			// Check if targeting DEX
			isDEX := dexRouters[strings.ToLower(tx.To)]
			
			// Check if there's a victim transaction nearby
			victimFound := false
			if i > 0 && i < len(txs)-1 {
				prevGas, _ := new(big.Int).SetString(txs[i-1].GasPrice, 10)
				nextGas, _ := new(big.Int).SetString(txs[i+1].GasPrice, 10)
				
				if prevGas.Int64() < p90 || nextGas.Int64() < p90 {
					victimFound = true
				}
			}

			if isDEX || victimFound {
				gweiPrice := float64(gasPrice.Int64()) / 1e9
				flags = append(flags, BehavioralFlag{
					Type:        "front_running_suspected",
					Severity:    0.85,
					Description: fmt.Sprintf("Potential front-running: %s used %.0f Gwei gas", tx.Hash[:10]+"...", gweiPrice),
					Evidence: map[string]interface{}{
						"gas_price_gwei": gweiPrice,
						"target_is_dex":  isDEX,
						"victim_nearby":  victimFound,
						"tx_hash":        tx.Hash,
					},
				})
			}
		}
	}

	return flags
}

// Detect sandwich attack patterns
func (ba *BehavioralAnalyzer) detectSandwichAttacks(txs []Transaction, p90 int64) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	// Look for buy-sell patterns with high gas
	for i := 0; i < len(txs)-2; i++ {
		tx1 := txs[i]
		tx2 := txs[i+1]
		tx3 := txs[i+2]
		
		gas1, _ := new(big.Int).SetString(tx1.GasPrice, 10)
		gas2, _ := new(big.Int).SetString(tx2.GasPrice, 10)
		gas3, _ := new(big.Int).SetString(tx3.GasPrice, 10)
		
		// Check if tx1 and tx3 have high gas, tx2 has normal gas
		if gas1.Int64() > p90*5 && gas3.Int64() > p90*5 && gas2.Int64() < p90*2 {
			// Check if same sender for tx1 and tx3
			if strings.EqualFold(tx1.From, tx3.From) {
				// Check if interacting with same contract
				if strings.EqualFold(tx1.To, tx2.To) && strings.EqualFold(tx2.To, tx3.To) {
					flags = append(flags, BehavioralFlag{
						Type:        "sandwich_attack",
						Severity:    0.9,
						Description: "Sandwich attack pattern detected",
						Evidence: map[string]interface{}{
							"front_tx":     tx1.Hash[:10] + "...",
							"victim_tx":    tx2.Hash[:10] + "...",
							"back_tx":      tx3.Hash[:10] + "...",
							"gas_multiple": float64(gas1.Int64()) / float64(gas2.Int64()),
						},
					})
				}
			}
		}
	}

	return flags
}

// Detect gas war patterns
func (ba *BehavioralAnalyzer) detectGasWars(gasUsage map[string][]GasInfo, median int64) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	// Look for multiple transactions in same block with escalating gas
	for timestamp, infos := range gasUsage {
		if len(infos) >= 3 {
			// Check if gas prices are escalating
			escalating := true
			highGasCount := 0
			
			for i := 1; i < len(infos); i++ {
				if infos[i].GasPrice <= infos[i-1].GasPrice {
					escalating = false
				}
				if infos[i].GasPrice > median*5 {
					highGasCount++
				}
			}
			
			if escalating && highGasCount >= 2 {
				maxGas := infos[len(infos)-1].GasPrice
				flags = append(flags, BehavioralFlag{
					Type:        "gas_war",
					Severity:    0.7,
					Description: fmt.Sprintf("Gas war detected: %d competing transactions", len(infos)),
					Evidence: map[string]interface{}{
						"block_timestamp":  timestamp,
						"competing_txs":    len(infos),
						"max_gas_gwei":     float64(maxGas) / 1e9,
						"escalation_ratio": float64(maxGas) / float64(infos[0].GasPrice),
					},
				})
			}
		}
	}

	return flags
}

// Detect potential exploit execution patterns
func (ba *BehavioralAnalyzer) detectExploitGasPattern(txs []Transaction, p99 int64) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	for _, tx := range txs {
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		value, _ := new(big.Int).SetString(tx.Value, 10)
		
		// Check for extremely high gas + large value extraction + complex input
		if gasPrice.Int64() > p99*20 && value.Cmp(big.NewInt(0)) > 0 && len(tx.Input) > 1000 {
			ethValue := new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(1e18))
			ethFloat, _ := ethValue.Float64()
			
			if ethFloat > 10.0 { // Significant value
				flags = append(flags, BehavioralFlag{
					Type:        "exploit_execution",
					Severity:    0.95,
					Description: fmt.Sprintf("Potential exploit: High gas + large extraction (%.2f ETH)", ethFloat),
					Evidence: map[string]interface{}{
						"tx_hash":          tx.Hash,
						"gas_price_gwei":   float64(gasPrice.Int64()) / 1e9,
						"value_eth":        ethFloat,
						"input_size":       len(tx.Input),
						"contract_address": tx.To,
					},
				})
			}
		}
	}

	return flags
}

// Detect censorship evasion patterns
func (ba *BehavioralAnalyzer) detectCensorshipEvasion(gasPrices []int64, median int64) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	// Count how many transactions have consistently high gas
	highGasCount := 0
	for _, gas := range gasPrices {
		if gas > median*3 {
			highGasCount++
		}
	}
	
	// If more than 60% of transactions have high gas, it's suspicious
	if float64(highGasCount)/float64(len(gasPrices)) > 0.6 {
		avgHigh := int64(0)
		for _, gas := range gasPrices {
			if gas > median*3 {
				avgHigh += gas
			}
		}
		avgHigh /= int64(highGasCount)
		
		flags = append(flags, BehavioralFlag{
			Type:        "censorship_evasion",
			Severity:    0.8,
			Description: "Consistent high gas usage suggests censorship evasion",
			Evidence: map[string]interface{}{
				"high_gas_ratio":    float64(highGasCount) / float64(len(gasPrices)),
				"avg_high_gas_gwei": float64(avgHigh) / 1e9,
				"median_gas_gwei":   float64(median) / 1e9,
			},
		})
	}

	return flags
}

// Helper struct for gas analysis
type GasInfo struct {
	GasPrice int64
	TxHash   string
	To       string
	Value    string
	Input    string
	IsError  string
}

// MEV-specific detection
func (ba *BehavioralAnalyzer) detectMEVActivity(txs []Transaction) []BehavioralFlag {
	flags := []BehavioralFlag{}
	
	// Known MEV bot addresses (in practice, this would be a larger database)
	mevBots := map[string]string{
		"0x00000000000006b2e6e3fc3e0dfd8dd7dba7b0": "MEV Bot Alpha",
		"0xa69babef1ca67a37ffaf7a485dfff3382056e78c": "Flashbots Builder",
	}
	
	// Check for MEV patterns
	for _, tx := range txs {
		// Check if from known MEV bot
		if botName, found := mevBots[strings.ToLower(tx.From)]; found {
			flags = append(flags, BehavioralFlag{
				Type:        "mev_activity",
				Severity:    0.7,
				Description: fmt.Sprintf("Transaction from known MEV bot: %s", botName),
				Evidence: map[string]interface{}{
					"bot_address": tx.From,
					"bot_name":    botName,
					"tx_hash":     tx.Hash,
				},
			})
		}
		
		// Check for bundle patterns (multiple txs in same block from same sender)
		// This would need block-level analysis in practice
	}
	
	return flags
}
