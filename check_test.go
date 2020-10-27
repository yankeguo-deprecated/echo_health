package echo_health

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheck(t *testing.T) {
	fn := func() error {
		return errors.New("bop")
	}
	c1 := CheckFunc(fn)
	require.Error(t, c1.CheckLiveness())
	require.Error(t, c1.CheckReadiness())
	c2 := ReadinessCheckFunc(fn)
	require.NoError(t, c2.CheckLiveness())
	require.Error(t, c2.CheckReadiness())
	c3 := LivenessCheckFunc(fn)
	require.Error(t, c3.CheckLiveness())
	require.NoError(t, c3.CheckReadiness())
}
