package wallet

import (
	"reflect"
	"testing"

	"github.com/FiruzMurodov/wallet/pkg/types"
	"github.com/google/uuid"
)

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