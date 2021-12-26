package wallet

import (
	"reflect"
	"testing"

	"github.com/FiruzMurodov/wallet/pkg/types"
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

func TestService_Reject_success (t *testing.T){
	s:= Service{}

	phone := types.Phone("000000001")
	account,err:= s.RegisterAccount(phone)
	if err!= nil {
		t.Errorf("Rejected(): can't register account, error =%v",err)
		return
	}

	err= s.Deposit(account.ID,1000_00)
	if err!= nil {
		t.Errorf("Rejected(): can't deposit account, error =%v",err)
		return
	}

	payment,err:= s.Pay(account.ID, 500_00,"auto")
	if err!= nil {
		t.Errorf("Rejected(): can't create payment, error =%v",err)
		return
	}

	err = s.Reject(payment.ID)
	if err!= nil {
		t.Errorf("Rejected(): error =%v",err)
		return
	}

}

func TestService_Reject_faild (t *testing.T){
	s:= Service{}

	phone := types.Phone("000000001")
	account,err:= s.RegisterAccount(phone)
	if err!= nil {
		t.Errorf("Rejected(): can't register account, error =%v",err)
		return
	}

	err= s.Deposit(account.ID,1000_00)
	if err!= nil {
		t.Errorf("Rejected(): can't deposit account, error =%v",err)
		return
	}

	_,_= s.Pay(account.ID, 500_00,"auto")

	err = s.Reject("11")
	if err== nil {
		t.Errorf("\n got - %v \n want - %v", err,ErrPaymentNotFound)
		return
	}

}

func TestService_FindPaymentByID_success (t *testing.T){
	s:= Service{}

	phone := types.Phone("000000001")
	account,err:= s.RegisterAccount(phone)
	if err!= nil {
		t.Errorf("Rejected(): can't register account, error =%v",err)
		return
	}

	err= s.Deposit(account.ID,1000_00)
	if err!= nil {
		t.Errorf("Rejected(): can't deposit account, error =%v",err)
		return
	}

	payment,err:= s.Pay(account.ID, 500_00,"auto")
	if err!= nil {
		t.Errorf("Rejected(): can't create payment, error =%v",err)
		return
	}

	got, err:= s.FindPaymentByID(payment.ID)
	if err!= nil {
		t.Errorf("FindPaymentByID(): error =%v",err)
		return
	}

	if !reflect.DeepEqual(payment,got) {
		t.Errorf("FindPaymentByID(): wrong payment return =%v",err)
		return
	}

}

func TestService_FindPaymentByID_faild (t *testing.T){
	s:= Service{}

	phone := types.Phone("000000001")
	account,err:= s.RegisterAccount(phone)
	if err!= nil {
		t.Errorf("Rejected(): can't register account, error =%v",err)
		return
	}

	err= s.Deposit(account.ID,1000_00)
	if err!= nil {
		t.Errorf("Rejected(): can't deposit account, error =%v",err)
		return
	}

	_,_= s.Pay(account.ID, 500_00,"auto")
	if err!= nil {
		t.Errorf("Rejected(): can't create payment, error =%v",err)
		return
	}

	_,err= s.FindPaymentByID("1")
	if err== nil {
		t.Errorf("\n got - %v \n want - %v", err,ErrPaymentNotFound)
		return
	}

	
}
