package logtimer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormatDuration(t *testing.T) {
	mustParse := func(s string) time.Duration {
		a, err := time.ParseDuration(s)
		require.NoError(t, err)
		return a
	}
	require.Equal(t, "03:58:44", FormatDuration(14324*time.Second, "%X"))
	require.Equal(t, "03:58:44.123456", FormatDuration(mustParse("3h58m44s123456789ns"), "%Xf"))
	require.Equal(t, "03:58:44.123456789", FormatDuration(mustParse("3h58m44s123456789ns"), "%Xn"))
	require.Equal(t, "03:58:44.000000789", FormatDuration(mustParse("3h58m44s789ns"), "%Xn"))
	require.Equal(t, "03:58:44.999999999", FormatDuration(mustParse("3h58m44s999999999ns"), "%Xn"))
	require.Equal(t, "03:58:45.000000001", FormatDuration(mustParse("3h58m44s1000000001ns"), "%Xn"))
}
