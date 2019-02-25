package account

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rhymond/go-money"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/event"
)

func TestCreateEmpty_MethodCalled_ShouldReturnNewInstance(t *testing.T) {
	t.Parallel()

	// act
	acc := Create(uuid.New())

	// assert
	assert.IsType(t, &Account{}, acc)
}

func TestAccountOpenAccount_Method_Called_ValidEventReturned(t *testing.T) {
	t.Parallel()

	// arrange
	ID := uuid.New()
	acc := Create(ID)

	expected := event.NewAccountOpened(ID, "123", *money.New(10000, "USD"))
	// act
	obtained, err := acc.OpenAccount(command.OpenAccount{
		AggID:   ID,
		Number:  "123",
		Balance: *money.New(10000, "USD"),
	})

	require.NoError(t, err)
	accountOpened := obtained[0].(event.AccountOpened)

	// assert
	assert.Equal(t, expected.AggregateID(), accountOpened.AggregateID())
	assert.Equal(t, expected.Number, accountOpened.Number)
	assert.Equal(t, expected.Balance, accountOpened.Balance)
}

func TestAccountDeposit_ValidAmountGiven_ValidEventReturned(t *testing.T) {
	t.Parallel()

	// arrange
	ID := uuid.New()
	acc := &Account{
		AggID:   ID,
		Number:  Number("123"),
		Balance: *money.New(0, "USD"),
	}

	expected := event.NewMoneyDeposited(
		ID,
		*money.New(10000, "USD"),
		*money.New(10000, "USD"),
	)

	// act
	obtained, err := acc.DepositMoney(command.DepositMoney{
		AggID:  ID,
		Amount: *money.New(10000, "USD"),
	})

	require.NoError(t, err)
	moneyDeposited := obtained[0].(event.MoneyDeposited)

	// assert
	assert.Equal(t, expected.AggregateID(), moneyDeposited.AggregateID())
	assert.Equal(t, expected.Amount, moneyDeposited.Amount)
	assert.Equal(t, expected.Balance, moneyDeposited.Balance)
}

func TestAccountDeposit_AmountInADifferentCurrencyGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(0, "USD"),
	}

	// act
	_, err := acc.DepositMoney(command.DepositMoney{Amount: *money.New(10000, "RUB")})

	// assert
	assert.EqualError(t, err, "cannot deposit account 123: Currencies don't match")
}

func TestAccountTransferMoney_ValidParamsGiven_ValidEventReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	expected := event.MoneyTransferred{
		From:    "123",
		To:      "777",
		Amount:  *money.New(1000, "USD"),
		Balance: *money.New(9000, "USD"),
	}

	// act
	obtained, err := acc.TransferMoney(command.TransferMoney{To: "777", Amount: *money.New(1000, "USD")})
	require.NoError(t, err)

	moneyTransferred := obtained[0].(event.MoneyTransferred)

	// assert
	assert.Equal(t, expected.AggregateID(), moneyTransferred.AggregateID())
	assert.Equal(t, expected.From, moneyTransferred.From)
	assert.Equal(t, expected.To, moneyTransferred.To)
	assert.Equal(t, expected.Amount, moneyTransferred.Amount)
	assert.Equal(t, expected.Balance, moneyTransferred.Balance)
}

func TestAccountTransferMoney_SendingTransferToTheSameAccount_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.TransferMoney(command.TransferMoney{To: "123", Amount: *money.New(1000, "USD")})

	// assert
	assert.EqualError(t, err, "cannot transfer money to the same account: 123")
}

//
func TestAccountTransferMoney_BalanceIsNotHighEnough_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(0, "USD"),
	}

	// act
	_, err := acc.TransferMoney(
		command.TransferMoney{From: "123", To: "777", Amount: *money.New(1000, "USD")})

	// assert
	assert.EqualError(t, err, "cannot send transfer from 123 to 777: balance -$10.00 is not high enough")
}

func TestAccountTransferMoney_AmountInADifferentCurrentGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.TransferMoney(command.TransferMoney{From: "123", To: "777", Amount: *money.New(1000, "RUB")})

	// assert
	assert.EqualError(t, err, "cannot send transfer from 123 to 777: Currencies don't match")
}

func TestAccountReceiveMoney_SendingTransferToTheSameAccount_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.ReceiveMoney(command.ReceiveMoney{From: "123", Amount: *money.New(1000, "USD")})

	// assert
	assert.EqualError(t, err, "cannot transfer money to the same account: 123")
}

func TestAccountReceiveMoney_AmountInADifferentCurrencyGiven_ErrorReturned(t *testing.T) {
	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.ReceiveMoney(command.ReceiveMoney{From: "777", Amount: *money.New(1000, "RUB")})

	// assert
	assert.EqualError(t, err, "cannot receive money from 777 to 123: Currencies don't match")
}

func TestAccountReceiveMoney_ValidParamsGiven_ValidEventReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	expected := event.MoneyReceived{
		From:    "777",
		To:      "123",
		Amount:  *money.New(1000, "USD"),
		Balance: *money.New(11000, "USD"),
	}

	// act
	obtained, err := acc.ReceiveMoney(command.ReceiveMoney{From: "777", To: "123", Amount: *money.New(1000, "USD")})
	require.NoError(t, err)

	moneyReceived := obtained[0].(event.MoneyReceived)

	// assert
	assert.Equal(t, expected.AggregateID(), moneyReceived.AggregateID())
	assert.Equal(t, expected.From, moneyReceived.From)
	assert.Equal(t, expected.To, moneyReceived.To)
	assert.Equal(t, expected.Amount, moneyReceived.Amount)
	assert.Equal(t, expected.Balance, moneyReceived.Balance)
}

func TestAccountToString_MethodCalled_ValidStringReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	expected := "#123: $100.00"

	// act
	obtained := acc.ToString()

	// assert
	assert.Equal(t, expected, obtained)
}
