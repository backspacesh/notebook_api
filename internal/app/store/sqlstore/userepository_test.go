package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"rest_api/internal/app/model"
	"rest_api/internal/app/store/sqlstore"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseString)
	defer teardown()

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseString)
	defer teardown()

	s := sqlstore.New(db)
	user1 := model.TestUser(t)
	s.User().Create(user1)
	user2, err := s.User().FindByEmail(user1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user2)
}