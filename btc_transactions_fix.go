// Updated function to better handle self-transfers and fix confirmation count
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

	// Get current BTC price and block height
	btcPrice := getBTCPrice()
	currentBlockHeight := getCurrentBlockHeight()

	transactions := make([]Transaction, 0)
	
	for _, tx := range btcResp.Txs {
		var amount int64 = 0
		var from, to string
		var fee int64 = tx.Fee
		var txType string = "UNKNOWN"
		
		// Check if this is a self-transfer
		isSelfTransfer := false
		
		// Check all inputs and outputs
		walletInInputs := false
		walletInOutputs := false
		
		for _, input := range tx.Inputs {
			if input.PrevOut.Addr == wallet {
				walletInInputs = true
				from = wallet
				break
			}
		}
		
		for _, out := range tx.Out {
			if out.Addr == wallet {
				walletInOutputs = true
				to = wallet
				amount += out.Value
			}
		}
		
		// Determine transaction type
		if walletInInputs && walletInOutputs {
			// Self-transfer
			txType = "SELF"
			isSelfTransfer = true
			from = wallet
			to = wallet
			// For self-transfers, calculate the net change (could be negative due to fees)
			var totalIn int64 = 0
			var totalOut int64 = 0
			
			for _, input := range tx.Inputs {
				if input.PrevOut.Addr == wallet {
					totalIn += input.PrevOut.Value
				}
			}
			
			for _, out := range tx.Out {
				if out.Addr == wallet {
					totalOut += out.Value
				}
			}
			
			amount = totalOut // Show the amount received back
		} else if walletInInputs && !walletInOutputs {
			// Outgoing transaction
			txType = "OUT"
			from = wallet
			// Find the main recipient
			var maxOut int64 = 0
			for _, out := range tx.Out {
				if out.Value > maxOut {
					to = out.Addr
					amount = out.Value
					maxOut = out.Value
				}
			}
		} else if !walletInInputs && walletInOutputs {
			// Incoming transaction
			txType = "IN"
			to = wallet
			// Get the sender
			if len(tx.Inputs) > 0 && tx.Inputs[0].PrevOut.Addr != "" {
				from = tx.Inputs[0].PrevOut.Addr
			}
		}
		
		// Skip if we couldn't determine the transaction details
		if amount == 0 && !isSelfTransfer {
			continue
		}
		
		// Convert satoshis to BTC
		btcAmount := float64(amount) / 100000000
		feeAmount := float64(fee) / 100000000
		
		// Calculate proper confirmations
		confirmations := int64(0)
		if tx.BlockHeight > 0 && currentBlockHeight > 0 {
			confirmations = currentBlockHeight - tx.BlockHeight + 1
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
			Confirmations: confirmations,
			USDValue:      btcAmount * btcPrice,
			Type:          txType,
		})
	}
	
	return transactions, nil
}

// Get current block height
func getCurrentBlockHeight() int64 {
	url := "https://blockchain.info/q/getblockcount"
	
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	
	var height int64
	fmt.Sscanf(string(body), "%d", &height)
	return height
}

// Update the display function to handle self-transfers
func displayTransactionTable(txs []Transaction, myWallet string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	
	t.AppendHeader(table.Row{"#", "Type", "Time", "From → To", "Amount", "USD Value", "Confirms", "Fee"})
	
	for i, tx := range txs {
		var txType string
		var typeColor text.Color
		
		switch tx.Type {
		case "SELF":
			txType = "SELF"
			typeColor = text.FgYellow
		case "IN":
			txType = "IN"
			typeColor = text.FgGreen
		case "OUT":
			txType = "OUT"
			typeColor = text.FgRed
		default:
			txType = "???"
			typeColor = text.FgWhite
		}
		
		// Format time
		timeStr := tx.Time.Format("01/02 15:04")
		if time.Since(tx.Time) < 24*time.Hour {
			timeStr = tx.Time.Format("15:04")
		}
		
		// Format addresses
		fromAddr := truncate(tx.From)
		toAddr := truncate(tx.To)
		
		// For self-transfers, make it clear
		if tx.Type == "SELF" {
			fromAddr = "[SELF]"
			toAddr = "[SELF]"
		} else {
			if tx.From == myWallet {
				fromAddr = "[TRACKED]"
			}
			if tx.To == myWallet {
				toAddr = "[TRACKED]"
			}
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
		
		if i < len(txs)-1 {
			t.AppendSeparator()
		}
	}
	
	fmt.Println(t.Render())
}
