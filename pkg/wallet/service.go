package wallet

import (
	"errors"

	"github.com/FiruzMurodov/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughtBalance = errors.New("account not enought balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites 	  []*types.Favorite
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

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite,error)  {
	
	payment,err:= s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,ErrPaymentNotFound
	}

	favoriteID:=uuid.New().String()
	favorite:= &types.Favorite{
		ID: favoriteID,
		AccountID: payment.AccountID,
		Name: name,
		Amount: payment.Amount,
		Category: payment.Category,
	}

	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite,error){
 
	for _,i:=range s.favorites{
	  if i.ID==favoriteID{
		return i,nil
	  }
	}
	 return nil,ErrFavoriteNotFound
  }

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment,error)  {
	
	favorite,err:=s.FindFavoriteByID(favoriteID)

	if err != nil {
		return nil, err
	}

	return s.Pay(favorite.AccountID,favorite.Amount,favorite.Category)
		
}