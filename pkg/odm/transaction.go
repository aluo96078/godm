package odm

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// WithTransaction 為需要原子性操作的業務提供事務支持。
// WithTransaction provides transaction support for operations that require atomicity.
func (o *GODM) WithTransaction(callback func(sessCtx mongo.SessionContext) error) error {
	session, err := MongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("start session error: %w", err)
	}
	defer session.EndSession(o.getContext())

	err = mongo.WithSession(o.getContext(), session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("start transaction error: %w", err)
		}
		if err := callback(sessCtx); err != nil {
			if abortErr := session.AbortTransaction(sessCtx); abortErr != nil {
				return fmt.Errorf("abort transaction error: %w", abortErr)
			}
			return fmt.Errorf("transaction error: %w", err)
		}
		if err := session.CommitTransaction(sessCtx); err != nil {
			return fmt.Errorf("commit transaction error: %w", err)
		}
		return nil
	})
	return err
}
