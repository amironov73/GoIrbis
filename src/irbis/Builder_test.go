package irbis

import "testing"

func TestNeedWrap_1(t *testing.T) {
	if !NeedWrap("") {
		t.Fail()
	}

	if NeedWrap("(Hello)") {
		t.Fail()
	}
}

func TestKeyword_1(t *testing.T) {
	text := Keyword("1").String()
	if text != "K=1" {
		t.Fail()
	}
	text = Keyword("1", "2").String()
	if text != "(K=1 + K=2)" {
		t.Fail()
	}
	text = Keyword("1", "2", "3").String()
	if text != "(K=1 + K=2 + K=3)" {
		t.Fail()
	}
}

func TestKeyword_2(t *testing.T) {
	text := Keyword("1 1").String()
	if text != "\"K=1 1\"" {
		t.Fail()
	}
	text = Keyword("1 1", "2 2").String()
	if text != "(\"K=1 1\" + \"K=2 2\")" {
		t.Fail()
	}
	text = Keyword("1 1", "2 2", "3 3").String()
	if text != "(\"K=1 1\" + \"K=2 2\" + \"K=3 3\")" {
		t.Fail()
	}
}

func TestKeyword_3(t *testing.T) {
	text := Keyword("1").And(Title("2")).Or(Author("3")).String()
	if text != "((K=1 * T=2) + A=3)" {
		t.Fail()
	}

	text = Keyword("1").Not(Title("2")).String()
	if text != "(K=1 ^ T=2)" {
		t.Fail()
	}
}

func TestBuilder_1(t *testing.T) {
	text := All().String()
	if text != "I=$" {
		t.Fail()
	}

	text = Number("1", "2").String()
	if text != "(IN=1 + IN=2)" {
		t.Fail()
	}

	text = Publisher("1").String()
	if text != "O=1" {
		t.Fail()
	}

	text = Place("1").String()
	if text != "MI=1" {
		t.Fail()
	}

	text = Subject("1").String()
	if text != "S=1" {
		t.Fail()
	}

	text = Language("1").String()
	if text != "J=1" {
		t.Fail()
	}

	text = Year("1").String()
	if text != "G=1" {
		t.Fail()
	}

	text = Magazine("1").String()
	if text != "TJ=1" {
		t.Fail()
	}

	text = DocumentKind("1").String()
	if text != "V=1" {
		t.Fail()
	}

	text = Udc("1").String()
	if text != "U=1" {
		t.Fail()
	}

	text = Bbk("1").String()
	if text != "BBK=1" {
		t.Fail()
	}

	text = Rzn("1").String()
	if text != "RZN=1" {
		t.Fail()
	}

	text = Mhr("1").String()
	if text != "MHR=1" {
		t.Fail()
	}
}

func TestSearch_SameField_1(t *testing.T) {
	text := Keyword("1").SameField(Keyword("2")).String()
	if text != "(K=1 (G) K=2)" {
		t.Fail()
	}
}

func TestSearch_SameRepeat_1(t *testing.T) {
	text := Keyword("1").SameRepeat(Keyword("2")).String()
	if text != "(K=1 (F) K=2)" {
		t.Fail()
	}
}
