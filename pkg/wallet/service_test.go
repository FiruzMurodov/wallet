package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/FiruzMurodov/wallet/pkg/types"
	"github.com/google/uuid"
)

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
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

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
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

func TestService_FindAccountByID_success(t *testing.T) {
	s := Service{}
	s.RegisterAccount("000000001")

	_, err := s.FindAccountByID(1)

	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

}

func TestService_Account_NotFound(t *testing.T) {
	s := Service{}
	s.RegisterAccount("000000001")

	_, err := s.FindAccountByID(2)

	if err == nil {
		t.Errorf("\ngot > %v \nwant > %v", err, ErrAccountNotFound)
	}

}

func TestService_Reject_success(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)

	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Rejected(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Rejected(): can't find payment by ID error = %v", err)
		return
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Rejected(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Rejected(): can't find account by ID error = %v", err)
		return
	}

	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Rejected(): amount didn't returned, account = %v", savedAccount)
		return
	}

}

func TestService_Reject_faild(t *testing.T) {
	s := Service{}

	phone := types.Phone("000000001")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Rejected(): can't register account, error =%v", err)
		return
	}

	err = s.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("Rejected(): can't deposit account, error =%v", err)
		return
	}

	_, _ = s.Pay(account.ID, 500_00, "auto")

	err = s.Reject("11")
	if err == nil {
		t.Errorf("\n got - %v \n want - %v", err, ErrPaymentNotFound)
		return
	}

}

func TestService_FindPaymentByID_success(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, payment) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}

}

func TestService_FindPaymentByID_faild(t *testing.T) {

	s := newTestService()

	_, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): must returned error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must returned errPaymentNotFound, returned = %v", err)
		return
	}

}

func TestService_Repeat_success(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.Repeat(payment.ID)

	if err != nil {
		t.Errorf("Repeat(): \ngot - %v", err)
		return
	}

}

func TestService_Repeat_faild(t *testing.T) {

	s := newTestService()

	_, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.Repeat("12")

	if err == nil {
		t.Errorf("\n got - %v \n want - %v", err, ErrPaymentNotFound)
		return
	}

}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	_, err = s.FavoritePayment(payment.ID, "switch")

	if err != nil {
		t.Errorf("FavoritePayment(): \ngot - %v", err)
		return
	}

}

func TestService_FavoritePayment_faild(t *testing.T) {
	s := newTestService()

	_, _, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FavoritePayment("12", "switch")

	if err == nil {
		t.Errorf("\n got - %v \n want - %v", err, ErrPaymentNotFound)
		return
	}

}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "switch")

	if err != nil {
		t.Errorf("FavoritePayment(): \ngot - %v", favorite)
	}

	if favorite.AccountID!= payment.AccountID {
		t.Errorf("FavoritePayment(): \ngot - %v", ErrAccountNotFound)
	}

	if favorite.Amount!= payment.Amount {
		t.Errorf("FavoritePayment(): \ngot - %v", ErrPaymentNotFound)
	}

	if favorite.Category!= payment.Category {
		t.Errorf("FavoritePayment(): \ngot - %v", ErrPaymentNotFound)
	}
	
	pay_favorite, err := s.PayFromFavorite(favorite.ID)

	if err != nil {
		t.Errorf("method PayFromFavorite returned not nil error, payfromtFavorite => %v", pay_favorite)

	}

}
