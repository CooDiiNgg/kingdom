package comms

import (
	"testing"

	commstypes "kingdom/internal/comms/comms_types"
	"kingdom/test/unit"

	"github.com/stretchr/testify/require"
)

func TestEncodeInternal(t *testing.T) {
	for _, test := range unit.CommsTestCases_EncodeInternal {
		t.Run(test.Name, func(t *testing.T) {
			result, err := encode_internal(test.Input)
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Expected, result)
			}
		})
	}
}

func TestDecodeInternal(t *testing.T) {
	for _, test := range unit.CommsTestCases_DecodeInternal {
		t.Run(test.Name, func(t *testing.T) {
			var result any
			var err error
			switch test.Type.(type) {
			case *commstypes.Request:
				result, err = decode_internal[*commstypes.Request](test.Input)
			case *commstypes.Task:
				result, err = decode_internal[*commstypes.Task](test.Input)
			case *commstypes.TaskResult:
				result, err = decode_internal[*commstypes.TaskResult](test.Input)
			default:
				t.Fatalf("Unsupported type: %T", test.Type)
			}
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Expected, result)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	for _, test := range unit.CommsTestCases_Encode {
		t.Run(test.Name, func(t *testing.T) {
			result, key, err := Encode(test.Input, test.Key)
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Expected, result)
				require.Equal(t, test.Key, key)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	for _, test := range unit.CommsTestCases_Decode {
		t.Run(test.Name, func(t *testing.T) {
			var result any
			var err error
			switch test.Type.(type) {
			case *commstypes.Request:
				result, err = Decode[*commstypes.Request](test.Input, test.Key)
			case *commstypes.Task:
				result, err = Decode[*commstypes.Task](test.Input, test.Key)
			case *commstypes.TaskResult:
				result, err = Decode[*commstypes.TaskResult](test.Input, test.Key)
			default:
				t.Fatalf("Unsupported type: %T", test.Type)
			}
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Expected, result)
			}
		})
	}
}
