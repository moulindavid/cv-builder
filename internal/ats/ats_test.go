package ats

import "testing"

func TestCheckOrder(t *testing.T) {
	if err := CheckOrder("NAME\nSUMMARY\nEXPERIENCE\nEDUCATION", "NAME", "SUMMARY", "EXPERIENCE", "EDUCATION"); err != nil {
		t.Fatal(err)
	}
	if err := CheckOrder("EDUCATION SUMMARY", "SUMMARY", "EDUCATION"); err == nil {
		t.Fatal("expected order error")
	}
}
