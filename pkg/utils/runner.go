package utils

import "github.com/uptrace/bun"

type QueryRunner interface {
	NewSelect() *bun.SelectQuery
	NewInsert() *bun.InsertQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
}

func GetQueryRunner(tx *bun.Tx, db *bun.DB) QueryRunner {
	if tx != nil {
		return tx
	}
	return db
}
