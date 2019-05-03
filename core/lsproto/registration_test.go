package lsproto

import "testing"

func TestParseRegistration1(t *testing.T) {
	s := "go,.go,tcp,goexec"
	reg, err := ParseRegistration(s)
	if err != nil {
		t.Fatal(err)
	}
	s2 := RegistrationString(reg)
	if s2 != "go,.go,tcp,goexec" {
		t.Fatal(s2)
	}
}

func TestParseRegistration2(t *testing.T) {
	s := "c/c++,.c .h .hpp,tcp,\"cexec opt1\""
	reg, err := ParseRegistration(s)
	if err != nil {
		t.Fatal(err)
	}
	s2 := RegistrationString(reg)
	if s2 != "c/c++,\".c .h .hpp\",tcp,\"cexec opt1\"" {
		t.Fatal(s2)
	}
}