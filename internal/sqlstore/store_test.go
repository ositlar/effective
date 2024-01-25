package sqlstore_test

import (
	"database/sql"
	"testing"

	"github.com/ositlar/effective/internal/model"
	"github.com/ositlar/effective/internal/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	m := model.Test()
	db, err := sql.Open("postgres", "host=localhost dbname=effective_test user=dev password=12345 sslmode=disable")
	assert.NoError(t, err)
	store := sqlstore.NewStore(db)
	err = store.Insert(m)
	assert.NoError(t, err)
}
