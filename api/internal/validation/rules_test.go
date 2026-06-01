package validation

import "testing"

func TestFieldErrors(t *testing.T) {
	t.Run("empty has no errors", func(t *testing.T) {
		errs := make(FieldErrors)
		if errs.Has() {
			t.Error("empty FieldErrors should not have errors")
		}
	})

	t.Run("add and check", func(t *testing.T) {
		errs := make(FieldErrors)
		errs.Add("email", "required")
		if !errs.Has() {
			t.Error("FieldErrors should have errors after Add")
		}
		if len(errs["email"]) != 1 || errs["email"][0] != "required" {
			t.Errorf("expected email error 'required', got %v", errs["email"])
		}
	})

	t.Run("multiple errors per field", func(t *testing.T) {
		errs := make(FieldErrors)
		errs.Add("password", "too short")
		errs.Add("password", "missing uppercase")
		if len(errs["password"]) != 2 {
			t.Errorf("expected 2 password errors, got %d", len(errs["password"]))
		}
	})

	t.Run("ToAppError", func(t *testing.T) {
		errs := make(FieldErrors)
		errs.Add("email", "required")
		appErr := errs.ToAppError()
		if appErr == nil {
			t.Error("ToAppError should return non-nil error")
		}
	})
}

func TestRequired(t *testing.T) {
	errs := make(FieldErrors)

	Required("field", "", errs)
	if !errs.Has() {
		t.Error("Required should add error for empty value")
	}

	errs = make(FieldErrors)
	Required("field", "value", errs)
	if errs.Has() {
		t.Error("Required should not add error for non-empty value")
	}
}

func TestMinLength(t *testing.T) {
	errs := make(FieldErrors)

	MinLength("password", "ab", 8, errs)
	if !errs.Has() {
		t.Error("MinLength should add error for too short value")
	}

	errs = make(FieldErrors)
	MinLength("password", "longpassword", 8, errs)
	if errs.Has() {
		t.Error("MinLength should not add error for long enough value")
	}
}

func TestMaxLength(t *testing.T) {
	errs := make(FieldErrors)

	MaxLength("name", "very long name here", 5, errs)
	if !errs.Has() {
		t.Error("MaxLength should add error for too long value")
	}

	errs = make(FieldErrors)
	MaxLength("name", "short", 5, errs)
	if errs.Has() {
		t.Error("MaxLength should not add error for short enough value")
	}
}

func TestEmail(t *testing.T) {
	errs := make(FieldErrors)

	Email("email", "invalid", errs)
	if !errs.Has() {
		t.Error("Email should add error for invalid email")
	}

	errs = make(FieldErrors)
	Email("email", "", errs) // empty should be skipped
	if errs.Has() {
		t.Error("Email should skip empty values (use Required separately)")
	}

	errs = make(FieldErrors)
	Email("email", "test@example.com", errs)
	if errs.Has() {
		t.Error("Email should not add error for valid email")
	}

	errs = make(FieldErrors)
	Email("email", "user+tag@domain.co.uk", errs)
	if errs.Has() {
		t.Error("Email should accept valid complex addresses")
	}
}

func TestOneOf(t *testing.T) {
	errs := make(FieldErrors)

	OneOf("format", "yaml", []string{"csv", "ofx"}, errs)
	if !errs.Has() {
		t.Error("OneOf should add error for value not in list")
	}

	errs = make(FieldErrors)
	OneOf("format", "csv", []string{"csv", "ofx"}, errs)
	if errs.Has() {
		t.Error("OneOf should not add error for valid value")
	}
}

func TestPasswordStrength(t *testing.T) {
	errs := make(FieldErrors)

	PasswordStrength("password", "", errs) // empty should be skipped
	if errs.Has() {
		t.Error("PasswordStrength should skip empty values")
	}

	errs = make(FieldErrors)
	PasswordStrength("password", "short", errs) // too short
	if !errs.Has() {
		t.Error("PasswordStrength should reject short passwords")
	}

	errs = make(FieldErrors)
	PasswordStrength("password", "abcdefgh", errs) // no uppercase or digit
	if !errs.Has() {
		t.Error("PasswordStrength should reject passwords without uppercase and digit")
	}

	errs = make(FieldErrors)
	PasswordStrength("password", "Abcdefg1", errs) // valid
	if errs.Has() {
		t.Error("PasswordStrength should accept valid password")
	}

	errs = make(FieldErrors)
	PasswordStrength("password", "LongValidPass123", errs) // valid long
	if errs.Has() {
		t.Error("PasswordStrength should accept long valid password")
	}
}

func TestIntRange(t *testing.T) {
	errs := make(FieldErrors)

	IntRange("page", -1, 1, 100, errs)
	if !errs.Has() {
		t.Error("IntRange should add error for below minimum")
	}

	errs = make(FieldErrors)
	IntRange("page", 101, 1, 100, errs)
	if !errs.Has() {
		t.Error("IntRange should add error for above maximum")
	}

	errs = make(FieldErrors)
	IntRange("page", 50, 1, 100, errs)
	if errs.Has() {
		t.Error("IntRange should not add error for value in range")
	}
}
