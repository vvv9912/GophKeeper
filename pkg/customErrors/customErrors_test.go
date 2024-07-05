package customErrors

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCustomError_Error(t *testing.T) {

	err := NewCustomError(nil, 11, "test")
	require.Equal(t, "[Custom Error Code: 11], [custom msg: test],[Error: ]", err.Error())
}
func TestError_HandlesEmptyMessage(t *testing.T) {
	err := errors.New("sample error")

	customErr := NewCustomError(err, 400, "")

	require.Equal(t, "[Custom Error Code: 400], [custom msg: ],[Error: sample error]", customErr.Error())

}
