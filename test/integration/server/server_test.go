package server_test

import (
	"net/http"
	"os"
	"p2pderivatives-oracle/test/integration/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helper.InitHelper()
	os.Exit(m.Run())
}

func Test_Server_Is_Reachable(t *testing.T) {
	resp, err := helper.CreateDefaultClient().R().Get("")

	assert.Nil(t, err)
	// Reachable as base url will return not found status code
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
}
