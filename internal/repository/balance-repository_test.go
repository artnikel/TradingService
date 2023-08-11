package repository

import (
	"context"
	"testing"

	bproto "github.com/artnikel/BalanceService/proto"
	"github.com/artnikel/BalanceService/proto/mocks"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(100.9),
	}
)

func TestBalanceOperation(t *testing.T) {
	client := new(mocks.BalanceServiceClient)
	strOperation := testBalance.Operation.String()
	client.On("BalanceOperation", mock.Anything, mock.Anything).
		Return(&bproto.BalanceOperationResponse{Operation: strOperation}, nil)
	rep := NewBalanceRepository(client)
	_, err := rep.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	client := new(mocks.BalanceServiceClient)
	client.On("GetBalance", mock.Anything, mock.Anything).
		Return(&bproto.GetBalanceResponse{Money: testBalance.Operation.InexactFloat64()}, nil)
	rep := NewBalanceRepository(client)
	_, err := rep.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	client.AssertExpectations(t)
}
