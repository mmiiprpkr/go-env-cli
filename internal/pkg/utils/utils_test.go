package utils

import (
	"testing"
)

func TestParseKeyValuePair(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedKey    string
		expectedValue  string
		expectedErrMsg string
	}{
		{
			name:          "Basic key value pair",
			input:         "KEY=value",
			expectedKey:   "KEY",
			expectedValue: "value",
		},
		{
			name:          "Key value pair with spaces",
			input:         "  KEY  =  value  ",
			expectedKey:   "KEY",
			expectedValue: "value",
		},
		{
			name:          "Key value with empty value",
			input:         "KEY=",
			expectedKey:   "KEY",
			expectedValue: "",
		},
		{
			name:           "Invalid input without equals sign",
			input:          "KEY",
			expectedErrMsg: "invalid format: KEY (expected key=value)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, value, err := ParseKeyValuePair(tc.input)

			if tc.expectedErrMsg != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tc.expectedErrMsg)
					return
				}
				if err.Error() != tc.expectedErrMsg {
					t.Errorf("expected error %q, got %q", tc.expectedErrMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if key != tc.expectedKey {
				t.Errorf("expected key %q, got %q", tc.expectedKey, key)
			}

			if value != tc.expectedValue {
				t.Errorf("expected value %q, got %q", tc.expectedValue, value)
			}
		})
	}
}
