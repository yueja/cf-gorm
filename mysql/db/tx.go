package db

// WithTransaction 事务支持
func WithTransaction(fn func(tx *Curd) error) error {
	var err error
	tx := DbCurd().db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
		if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = fn(&Curd{db: tx})
	if err != nil {
		return err
	}
	return nil
}

// WithTransactionSpecify 事务支持
func WithTransactionSpecify(client string, fn func(tx *Curd) error) error {
	var err error
	tx := SpecifyDbCurd(client).db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
		if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = fn(&Curd{db: tx})
	if err != nil {
		return err
	}
	return nil
}
