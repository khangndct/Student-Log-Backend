package services

import (
	"backend/internal/models"
	"errors"
	"gorm.io/gorm"
)

func DeleteAccountWithCascade(db *gorm.DB, accountID int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Convert accountID (int64) to uint for comparison with OwnerID and WriterID
		accountIDUint := uint(accountID)

		// Remove the account ID from WriterIDList arrays in all LogHeads
		// This handles the case where the account is a writer but not the owner
		// PostgreSQL's array_remove function removes all occurrences of the value
		if err := tx.Exec(
			"UPDATE log_heads SET writer_id_list = array_remove(writer_id_list, ?) WHERE ? = ANY(writer_id_list)",
			accountID, accountID,
		).Error; err != nil {
			return err
		}

		// Delete all LogHeads where OwnerID matches the account ID
		// This will cascade delete their LogContents via the existing constraint
		if err := tx.Where("owner_id = ?", accountIDUint).Delete(&models.LogHead{}).Error; err != nil {
			return err
		}

		// Delete any remaining LogContents where WriterID matches the account ID
		// These are log contents in log heads owned by other users
		if err := tx.Where("writer_id = ?", accountIDUint).Delete(&models.LogContent{}).Error; err != nil {
			return err
		}

		// Finally, delete the Account itself
		result := tx.Delete(&models.Account{}, accountID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("account not found")
		}

		return nil
	})
}

