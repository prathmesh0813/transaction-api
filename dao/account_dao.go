package dao

import (
	"errors"
	"transaction-api/config"
	"transaction-api/models"

	"gorm.io/gorm"
)

// GetAllAccounts retrieves all accounts using GORM
func GetAllAccounts() ([]models.Account, error) {
	var accounts []models.Account
	result := config.DB.Find(&accounts)
	return accounts, result.Error
}

// GetAccountByID retrieves an account by its ID
func GetAccountByID(id int) (models.Account, error) {
	var account models.Account
	result := config.DB.First(&account, id)
	return account, result.Error
}

// DeleteAccount deletes an account by ID
func DeleteAccount(id int) error {
	result := config.DB.Delete(&models.Account{}, id)
	return result.Error
}

// CreateAccount inserts a new account
func CreateAccount(account models.Account) error {
	result := config.DB.Create(&account)
	return result.Error
}

// UpdateAccount updates account name and balance
func UpdateAccount(id int, account models.Account) error {
	result := config.DB.Model(&models.Account{}).Where("id = ?", id).
		Updates(models.Account{Name: account.Name, Balance: account.Balance})
	return result.Error
}

// TransferAmount performs a fund transfer between accounts with transaction handling
func TransferAmount(fromID, toID int, amount float64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var fromAccount models.Account
		if err := tx.First(&fromAccount, fromID).Error; err != nil {
			return err
		}

		if fromAccount.Balance < amount {
			return errors.New("insufficient funds")
		}

		var toAccount models.Account
		if err := tx.First(&toAccount, toID).Error; err != nil {
			return err
		}

		// Deduct and add balances
		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-amount).Error; err != nil {
			return err
		}
		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+amount).Error; err != nil {
			return err
		}

		// Insert transaction record
		transaction := models.Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// MiniStatement returns the last 5 transactions for a given account
func MiniStatement(accountID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := config.DB.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Order("created_at DESC").Limit(5).Find(&transactions).Error
	return transactions, err
}
