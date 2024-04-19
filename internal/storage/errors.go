package storage

import "errors"

var (
	ErrStartTx    = errors.New("failed to begin transaction")
	ErrCommitTx   = errors.New("failed to commit transaction")
	ErrRollbackTx = errors.New("failed to rollback transaction")

	ErrCarExist = errors.New("car with this register number already exist")
	ErrCarNotFound = errors.New("car with this id not found")
)