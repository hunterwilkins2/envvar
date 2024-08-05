package envvar_test

import (
	"os"
	"testing"

	"github.com/hunterwilkins2/envvar"
)

func TestGetString(t *testing.T) {
	os.Setenv("STRING_VAR", "abc")
	result := envvar.Get[string]("STRING_VAR")
	if result != "abc" {
		t.Errorf("Expected abc; got %s", result)
	}
}

func TestGetBool(t *testing.T) {
	os.Setenv("BOOL_VAR", "true")
	result := envvar.Get[bool]("BOOL_VAR")
	if result != true {
		t.Errorf("Expected true; got false")
	}
}

func TestGetInt(t *testing.T) {
	os.Setenv("INT_VAR", "42")
	result := envvar.Get[int]("INT_VAR")
	if result != 42 {
		t.Errorf("Expected 42; got %d", result)
	}
}

func TestGetUint(t *testing.T) {
	os.Setenv("UINT_VAR", "1")
	result := envvar.Get[uint]("UINT_VAR")
	if result != 1 {
		t.Errorf("Expected 1; got %d", result)
	}
}

func TestGetFloat(t *testing.T) {
	os.Setenv("FLOAT_VAR", "3.14")
	result := envvar.Get[float64]("FLOAT_VAR")
	if result != 3.14 {
		t.Errorf("Expected 3.14; got %f", result)
	}
}

func TestGetWithDefaultNotSet(t *testing.T) {
	result := envvar.GetWithDefault("NOT_SET", "abc")
	if result != "abc" {
		t.Errorf("Expected abc; got %s", result)
	}
}

func TestGetWithDefaultCannotParse(t *testing.T) {
	os.Setenv("INT_VAR", "abc")
	result := envvar.GetWithDefault("INT_VAR", 10)
	if result != 10 {
		t.Errorf("Expected 10; got %d", result)
	}
}

func TestGetWithDefault(t *testing.T) {
	os.Setenv("INT_VAR", "3")
	result := envvar.GetWithDefault("INT_VAR", 20)
	if result != 3 {
		t.Errorf("Expected 3; got %d", result)
	}
}

func TestGetAndValidateNotSet(t *testing.T) {
	result, err := envvar.GetAndValidate("NOT_SET", envvar.GreaterThan(0))
	expectedErr := "could not parse \"NOT_SET\": \"NOT_SET\" is not set"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0; gpt %d", result)
	}
}

func TestGetAndValidateCannotParse(t *testing.T) {
	os.Setenv("INT_VAR", "abc")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.GreaterThan(0))
	expectedErr := "could not parse \"INT_VAR\": strconv.Atoi: parsing \"abc\": invalid syntax"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0; gpt %d", result)
	}
}

func TestGetAndValidateGreaterThan(t *testing.T) {
	os.Setenv("INT_VAR", "3")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.GreaterThan(2))
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
	if result != 3 {
		t.Errorf("Expected 3; got %d", result)
	}
}

func TestGetAndValidateGreaterThanFailed(t *testing.T) {
	os.Setenv("INT_VAR", "5")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.GreaterThan(10))
	expectedErr := "\"INT_VAR\" failed validation: must be greater than 10"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0; gpt %d", result)
	}
}

func TestGetAndValidateLessThan(t *testing.T) {
	os.Setenv("INT_VAR", "3")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.LessThan(4))
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
	if result != 3 {
		t.Errorf("Expected 3; got %d", result)
	}
}

func TestGetAndValidateLessThanFailed(t *testing.T) {
	os.Setenv("INT_VAR", "5")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.LessThan(3))
	expectedErr := "\"INT_VAR\" failed validation: must be less than 3"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0; gpt %d", result)
	}
}

func TestGetAndValidateBetween(t *testing.T) {
	os.Setenv("INT_VAR", "3")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.Between(0, 10))
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
	if result != 3 {
		t.Errorf("Expected 3; got %d", result)
	}
}

func TestGetAndValidateBetweenMinFailed(t *testing.T) {
	os.Setenv("INT_VAR", "-1")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.Between(0, 10))
	expectedErr := "\"INT_VAR\" failed validation: must be between 0 and 10"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0; gpt %d", result)
	}
}

func TestGetAndValidateBetweenMaxFailed(t *testing.T) {
	os.Setenv("INT_VAR", "11")
	result, err := envvar.GetAndValidate("INT_VAR", envvar.Between(0, 10))
	expectedErr := "\"INT_VAR\" failed validation: must be between 0 and 10"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0; gpt %d", result)
	}
}

func TestGetAndValidateWithin(t *testing.T) {
	os.Setenv("STRING_VAR", "development")
	result, err := envvar.GetAndValidate("STRING_VAR", envvar.Within([]string{"development", "testing", "staging", "production"}))
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
	if result != "development" {
		t.Errorf("Expected development; got %s", result)
	}
}

func TestGetAndValidateWithinFailed(t *testing.T) {
	os.Setenv("STRING_VAR", "develop")
	result, err := envvar.GetAndValidate("STRING_VAR", envvar.Within([]string{"development", "testing", "staging", "production"}))
	expectedErr := "\"STRING_VAR\" failed validation: develop not within [development testing staging production]"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != "" {
		t.Errorf("Expected 0; gpt %s", result)
	}
}

func TestGetAndValidateEmail(t *testing.T) {
	os.Setenv("STRING_VAR", "test@abc.com")
	result, err := envvar.GetAndValidate("STRING_VAR", envvar.ValidEmail())
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
	if result != "test@abc.com" {
		t.Errorf("Expected test@abc.com; got %s", result)
	}
}

func TestGetAndValidateEmailFailed(t *testing.T) {
	os.Setenv("STRING_VAR", "not-valid")
	result, err := envvar.GetAndValidate("STRING_VAR", envvar.ValidEmail())
	expectedErr := "\"STRING_VAR\" failed validation: not-valid is not a valid email"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %s", expectedErr, err.Error())
	}
	if result != "" {
		t.Errorf("Expected 0; gpt %s", result)
	}
}

func TestGetAndValidateUrl(t *testing.T) {
	os.Setenv("STRING_VAR", "http://localhost:8080/abc/efg")
	result, err := envvar.GetAndValidate("STRING_VAR", envvar.ValidUrl())
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
	if result != "http://localhost:8080/abc/efg" {
		t.Errorf("Expected http://localhost:8080/abc/efg; got %s", result)
	}
}

func TestGetAndValidateUrlFailed(t *testing.T) {
	os.Setenv("STRING_VAR", "http://192.168.0.%31/")
	result, err := envvar.GetAndValidate("STRING_VAR", envvar.ValidUrl())
	expectedErr := "\"STRING_VAR\" failed validation: http://192.168.0.%31/ is not a valid url"
	if err.Error() != expectedErr {
		t.Errorf("Expected %s; got %v", expectedErr, err)
	}
	if result != "" {
		t.Errorf("Expected \"\"; got %s", result)
	}
}
