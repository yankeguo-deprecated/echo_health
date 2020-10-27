package echo_health

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCheck struct {
	alive bool
	ready bool
}

func (t *testCheck) CheckLiveness() error {
	if t.alive {
		return nil
	} else {
		return errors.New("test error")
	}
}

func (t *testCheck) CheckReadiness() error {
	if t.ready {
		return nil
	} else {
		return errors.New("test error")
	}
}

func TestNew(t *testing.T) {
	ck1 := &testCheck{
		alive: true,
		ready: true,
	}
	ck2 := &testCheck{
		alive: true,
		ready: true,
	}
	e := echo.New()
	h := New(ck1, ck2)(func(c echo.Context) error {
		return c.String(http.StatusOK, "wai bi ba bu")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	err := h(e.NewContext(req, rec))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "wai bi ba bu", rec.Body.String())

	checkAlive := func() (msg string, ok bool) {
		req := httptest.NewRequest(http.MethodGet, PathLiveness, nil)
		rec := httptest.NewRecorder()
		err := h(e.NewContext(req, rec))
		require.NoError(t, err)
		msg, ok = rec.Body.String(), rec.Code == http.StatusOK
		return
	}
	checkReady := func() (msg string, ok bool) {
		req := httptest.NewRequest(http.MethodGet, PathReadiness, nil)
		rec := httptest.NewRecorder()
		err := h(e.NewContext(req, rec))
		require.NoError(t, err)
		msg, ok = rec.Body.String(), rec.Code == http.StatusOK
		return
	}

	type Case struct {
		alive1   bool
		alive2   bool
		ready1   bool
		ready2   bool
		aliveOK  bool
		aliveMsg string
		readyOK  bool
		readyMsg string
	}

	cases := []Case{
		{
			alive1:   true,
			ready1:   true,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  true,
			readyMsg: "OK",
		},
		{
			alive1:   false,
			ready1:   true,
			alive2:   false,
			ready2:   true,
			aliveOK:  false,
			aliveMsg: "test error; test error",
			readyOK:  true,
			readyMsg: "OK",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   true,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  true,
			readyMsg: "OK",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  true,
			aliveMsg: "OK",
			readyOK:  false,
			readyMsg: "test error",
		},
		{
			alive1:   true,
			ready1:   false,
			alive2:   true,
			ready2:   true,
			aliveOK:  false,
			aliveMsg: "cascading readiness failure",
			readyOK:  false,
			readyMsg: "test error",
		},
	}

	for i, c := range cases {
		ck1.alive = c.alive1
		ck2.alive = c.alive2
		ck1.ready = c.ready1
		ck2.ready = c.ready2
		aliveMsg, aliveOK := checkAlive()
		readyMsg, readyOK := checkReady()
		t.Run(fmt.Sprintf("case_%d_CK1_%t_%t__CK2_%t_%t", i, c.alive1, c.ready1, c.alive2, c.ready2), func(t *testing.T) {
			require.Equal(t, c.aliveOK, aliveOK)
			require.Equal(t, c.aliveMsg, aliveMsg)
			require.Equal(t, c.readyOK, readyOK)
			require.Equal(t, c.readyMsg, readyMsg)
		})
	}
}
