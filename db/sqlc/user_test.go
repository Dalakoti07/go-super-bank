package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/utils"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: "secret",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	account2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, user1.Username, account2.Username)
	require.Equal(t, user1.FullName, account2.FullName)
	require.Equal(t, user1.HashedPassword, account2.HashedPassword)
	require.Equal(t, user1.Email, account2.Email)
	require.WithinDuration(t, user1.CreatedAt, account2.CreatedAt, time.Second)
}
