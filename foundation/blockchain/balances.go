package blockchain

import "fmt"

// BalanceSheet represents the data representation to maintain address balances.
type BalanceSheet map[string]uint

// newBalanceSheet constructs a new balance sheet for use.
func newBalanceSheet() BalanceSheet {
	return make(BalanceSheet)
}

// replace updates the balance sheet for a given address.
func (bs BalanceSheet) replace(address string, value uint) {
	bs[address] = value
}

// =============================================================================

// copyBalanceSheet makes a copy of the specified balance sheet.
func copyBalanceSheet(org BalanceSheet) BalanceSheet {
	balanceSheet := newBalanceSheet()
	for acct, value := range org {
		balanceSheet.replace(acct, value)
	}
	return balanceSheet
}

// applyMiningRewardToBalance gives the beneficiary address a reward for mining a block.
func applyMiningRewardToBalance(balanceSheet BalanceSheet, beneficiary string, reward uint) {
	balanceSheet[beneficiary] += reward
}

// applyMiningFeeToBalance gives the beneficiary address a fee for mining the block.
func applyMiningFeeToBalance(balanceSheet BalanceSheet, beneficiary string, tx Tx) {
	fee := tx.Gas + tx.Tip

	balanceSheet[beneficiary] += fee
	balanceSheet[tx.From] -= fee
}

// applyTransactionToBalance performs the business logic for applying a
// transaction to the balance sheet.
func applyTransactionToBalance(balanceSheet BalanceSheet, tx Tx) error {
	if tx.Status == TxStatusError {
		return nil
	}

	if tx.Data == TxDataReward {
		balanceSheet[tx.To] += tx.Value
		return nil
	}

	if tx.From == tx.To {
		return fmt.Errorf("invalid transaction, do you mean to give a reward, from %s, to %s", tx.From, tx.To)
	}

	if tx.Value > balanceSheet[tx.From] {
		return fmt.Errorf("%s has an insufficient balance", tx.From)
	}

	balanceSheet[tx.From] -= tx.Value
	balanceSheet[tx.To] += tx.Value

	if tx.Tip > 0 {
		balanceSheet[tx.From] -= tx.Tip
	}

	return nil
}
