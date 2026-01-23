package database

import "model-manager/backend/models"

// ResetAllPendingClientFiles deletes all ClientFile records with status 'pending'.
// This is used at startup to clear any stuck states from previous runs or crashes.
func ResetAllPendingClientFiles() error {
	return DB.Unscoped().Where("status = ?", "pending").Delete(&models.ClientFile{}).Error
}

// ResetPendingClientFilesForClient deletes ClientFile records with status 'pending'
// for a specific client. This is used when a client disconnects to prevent stuck models.
func ResetPendingClientFilesForClient(clientID string) error {
	return DB.Unscoped().Where("client_id = ? AND status = ?", clientID, "pending").Delete(&models.ClientFile{}).Error
}
