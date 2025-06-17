package comms

import (
	"testing"

	"github.com/stretchr/testify/require"

	"kingdom/test/unit"
)

func TestEncrypt(t *testing.T) {
	for _, test := range unit.CommsTestCases_Encrypt {
		t.Run(test.Name, func(t *testing.T) {
			result, err := encrypt(test.Input, test.Key, test.IV)
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Expected, result)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	for _, test := range unit.CommsTestCases_Decrypt {
		t.Run(test.Name, func(t *testing.T) {
			result, err := decrypt(test.Input, test.Key, test.IV)
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Expected, result)
			}
		})
	}
}
