# Gas Price Anomaly Patterns in Blockchain Security

## Overview

Abnormally high gas prices are a critical indicator of various malicious activities on the blockchain. Our enhanced analyzer detects multiple gas-related attack patterns that go beyond simple threshold checks.

## Gas Patterns and Their Implications

### 1. Front-Running Attacks
**Pattern Characteristics:**
- Gas price 10-100x above network average
- Single transaction with extreme priority
- Often targets DEX trades, token launches, or NFT mints

**Detection Method:**
```
if gasPrice > p90 * 10 && (targetsDEX || hasVictimNearby) {
    flag as front_running
}
```

**Real Example:**
- Normal gas: 30 Gwei
- Front-runner gas: 500-3000 Gwei
- Target: Uniswap large buy order
- Profit: Price difference on subsequent sell

### 2. Sandwich Attacks (MEV)
**Pattern Characteristics:**
- Two high-gas transactions surrounding a normal-gas transaction
- Same attacker address for both high-gas txs
- Buy-sell sequence targeting same pool

**Detection Sequence:**
1. TX1: High gas (300 Gwei) - Attacker buys
2. TX2: Normal gas (30 Gwei) - Victim buys (price rises)
3. TX3: High gas (300 Gwei) - Attacker sells

**Detection Method:**
```
if (tx1.gas > p90*5 && tx3.gas > p90*5 && tx2.gas < p90*2) &&
   (tx1.from == tx3.from) && (same target contract) {
    flag as sandwich_attack
}
```

### 3. Gas Wars (Competition)
**Pattern Characteristics:**
- Multiple transactions in same block
- Escalating gas prices
- High failure rate
- Common during NFT drops or token launches

**Example Escalation:**
- TX1: 100 Gwei (fails)
- TX2: 200 Gwei (fails)
- TX3: 500 Gwei (fails)
- TX4: 1000 Gwei (succeeds)

### 4. Exploit Execution
**Pattern Characteristics:**
- Single transaction with extreme gas (20x+ p99)
- Large value extraction
- Complex contract interaction (large input data)
- Time-sensitive execution

**Why High Gas for Exploits:**
- Ensure execution before patch
- Prevent other hackers from front-running
- Guarantee inclusion despite network congestion

### 5. Censorship Evasion
**Pattern Characteristics:**
- Consistently high gas across all transactions
- Often combined with mixer interactions
- Rapid sequential transactions
- Used by hackers moving stolen funds

**Detection:**
```
if (highGasTransactions / totalTransactions > 0.6) {
    flag as censorship_evasion
}
```

## Integration with Behavioral Analysis

Our analyzer combines gas patterns with other behavioral indicators:

```go
// Example combined detection
if (highGas && rapidDrainage && mixerInteraction) {
    riskScore = 0.95 // Very high risk
    flagAs("hack_fund_movement")
}
```

## Risk Scoring Impact

Gas anomalies contribute to the overall risk score:
- Front-running: +0.85 severity
- Sandwich attack: +0.9 severity
- Exploit execution: +0.95 severity
- Gas wars: +0.6 severity
- Censorship evasion: +0.8 severity

## Advanced Detection Features

### 1. Contextual Analysis
The analyzer considers:
- Target contract type (DEX, NFT, Token)
- Transaction timing and sequencing
- Historical gas patterns for the address
- Network congestion at time of transaction

### 2. MEV Bot Detection
- Identifies known MEV bot addresses
- Detects bundle patterns
- Analyzes profit extraction methods

### 3. Statistical Thresholds
- P50 (median): Baseline gas price
- P90: High gas threshold
- P99: Extreme gas threshold
- Dynamic adjustment based on network conditions

## Example Output

When high gas anomalies are detected:

```
🚩 Behavioral Patterns Detected:
   • Potential front-running: 0x123...abc used 2500 Gwei gas [Severity: 0.85]
     Evidence: {
       "gas_price_gwei": 2500,
       "target_is_dex": true,
       "victim_nearby": true,
       "normal_gas": 30
     }
   
   • Sandwich attack pattern detected [Severity: 0.90]
     Evidence: {
       "front_tx": "0xabc...",
       "victim_tx": "0xdef...",
       "back_tx": "0x123...",
       "profit_estimate": "0.5 ETH"
     }
```

## Why This Matters

1. **Early Warning**: High gas often precedes major exploits
2. **Attack Attribution**: Gas patterns help identify attack types
3. **Risk Assessment**: Quantifies threat level based on gas behavior
4. **Forensics**: Helps trace fund flows post-hack

## Configuration

You can adjust gas thresholds in the analyzer:

```go
riskThresholds: RiskThresholds{
    GasAnomalyMultiplier: 3.0,  // 3x average = anomaly
    // Add custom thresholds
    FrontRunMultiplier: 10.0,    // 10x p90 = front-run
    SandwichMultiplier: 5.0,     // 5x p90 = sandwich
    ExploitMultiplier: 20.0,     // 20x p99 = exploit
}
```

## Best Practices for Detection

1. **Combine Indicators**: Don't rely on gas alone
2. **Consider Context**: Network congestion affects baselines
3. **Update Thresholds**: Adjust based on network evolution
4. **Monitor Patterns**: Track new MEV strategies

## Future Enhancements

1. **Machine Learning**: Train models on historical gas patterns
2. **Real-time Alerts**: Webhook notifications for gas spikes
3. **Cross-chain Analysis**: Detect coordinated attacks
4. **Predictive Analysis**: Forecast attacks based on gas trends
