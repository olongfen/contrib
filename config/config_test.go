package config

import "testing"

type c struct {
	Config `yaml:"-"`
	Name  string
	Age  int
}

var d = &c{
	Name:   "Test",
	Age:    20,
}

// TestLoadConfiguration
func TestLoadConfiguration(t *testing.T) {
	var(
		err error
	)
	if err = LoadConfiguration("test.yaml",d,d);err!=nil{
		t.Fatal(err)
	}
}

func TestConfig_Save(t *testing.T) {
	var(
		err error
	)
	if err = d.Save(nil);err!=nil{
		t.Fatal(err)
	}
}