package wallet

import (
	"errors"
	"fmt"

	"github.com/FiruzMurodov/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughtBalance = errors.New("account not enought balance")
var ErrPaymentNotFound = errors.New("payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

type testService struct {
	*Service
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000001",
	balance: 1_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 500_00, category: "auto"},
	},
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

// Deposit method
func (s *Service) Deposit(accountID int64, amount types.Money) error {

	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, ac := range s.accounts {
		if ac.ID == accountID {
			account = ac
			break
		}
	}
	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount

	return nil
}

//Pay method
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {

	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, ac := range s.accounts {
		if ac.ID == accountID {
			account = ac
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughtBalance
	}
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) addAcount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			return acc, nil
		}

	}

	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {

	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {

	payment,err:= s.FindPaymentByID(paymentID)

	if err != nil {
		return ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil

}

func (s *Service) Repeat(paymentID string) (*types.Payment,error)  {
	
	targetPayment,err:= s.FindPaymentByID(paymentID)

	if err != nil {
		return nil,ErrPaymentNotFound
	}

	repeatPay,err := s.Pay(targetPayment.AccountID,targetPayment.Amount,targetPayment.Category)

	if err != nil {
		return nil,err
	}

	return repeatPay,nil
}