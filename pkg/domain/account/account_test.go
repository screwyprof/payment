package account

import (
	"testing"

	"github.com/rhymond/go-money"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/pkg/event"
)

func TestCreateEmpty_MethodCalled_ShouldReturnNewInstance(t *testing.T) {
	t.Parallel()

	// act
	acc := CreateEmpty()

	// assert
	assert.IsType(t, &Account{}, acc)
}

func TestAccountOpenAccount_Method_Called_ValidEventReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := CreateEmpty()

	numberGenerator = func() Number {
		return Number("123")
	}

	expected := event.AccountOpened{
		Number:  "123",
		Balance: *money.New(10000, "USD"),
	}

	// act
	obtained := acc.OpenAccount(*money.New(10000, "USD"))

	// assert
	assert.Equal(t, expected, obtained)
}

func TestAccountDeposit_ValidAmountGiven_ValidEventReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(0, "USD"),
	}

	expected := event.MoneyDeposited{
		Number:  "123",
		Amount:  *money.New(10000, "USD"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	obtained, err := acc.Deposit(*money.New(10000, "USD"))
	require.NoError(t, err)

	// assert
	assert.Equal(t, expected, obtained)
}

func TestAccountDeposit_AmountInADifferentCurrentGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(0, "USD"),
	}

	// act
	_, err := acc.Deposit(*money.New(10000, "RUB"))

	// assert
	assert.EqualError(t, err, "cannot deposit account 123: Currencies don't match")
}

func TestAccountSendTransferTo_ValidParamsGiven_ValidEventReturned(t *testing.T) {
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
	obtained, err := acc.SendTransferTo(Number("777"), *money.New(1000, "USD"))
	require.NoError(t, err)

	// assert
	assert.Equal(t, expected, obtained)
}

func TestAccountSendTransferTo_SendingTransferToTheSameAccount_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.SendTransferTo(Number("123"), *money.New(1000, "USD"))

	// assert
	assert.EqualError(t, err, "cannot transfer money to the same account: 123")
}

func TestAccountSendTransferTo_BalanceIsNotHighEnough_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(0, "USD"),
	}

	// act
	_, err := acc.SendTransferTo(Number("777"), *money.New(1000, "USD"))

	// assert
	assert.EqualError(t, err, "cannot send transfer from 123 to 777: balance -$10.00 is not high enough")
}

func TestAccountSendTransferTo_AmountInADifferentCurrentGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.SendTransferTo(Number("777"), *money.New(1000, "RUB"))

	// assert
	assert.EqualError(t, err, "cannot send transfer from 123 to 777: Currencies don't match")
}

func TestAccountReceiveMoneyFrom_SendingTransferToTheSameAccount_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.ReceiveMoneyFrom(Number("123"), *money.New(1000, "USD"))

	// assert
	assert.EqualError(t, err, "cannot transfer money to the same account: 123")
}

func TestAccountReceiveMoneyFrom_AmountInADifferentCurrentGiven_ErrorReturned(t *testing.T) {
	// arrange
	acc := &Account{
		Number:  Number("123"),
		Balance: *money.New(10000, "USD"),
	}

	// act
	_, err := acc.ReceiveMoneyFrom(Number("777"), *money.New(1000, "RUB"))

	// assert
	assert.EqualError(t, err, "cannot receive money from 777 to 123: Currencies don't match")
}

func TestAccountReceiveMoneyFrom_ValidParamsGiven_ValidEventReturned(t *testing.T) {
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
	obtained, err := acc.ReceiveMoneyFrom(Number("777"), *money.New(1000, "USD"))
	require.NoError(t, err)

	// assert
	assert.Equal(t, expected, obtained)
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
