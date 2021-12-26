package wallet

import (
	"testing"
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