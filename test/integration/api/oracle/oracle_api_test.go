// +build integration

package oracleapi_test

import (
	"encoding/json"
	"net/http"
	"os"
	"p2pderivatives-oracle/internal/api"
	"p2pderivatives-oracle/internal/dlccrypto"
	helper "p2pderivatives-oracle/test/integration"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helper.InitHelper()
	os.Exit(m.Run())
}

func TestGetOraclePublicKey_Returns_CorrectValue(t *testing.T) {
	// arrange
	client := helper.CreateDefaultClient()

	// act
	resp, err := client.R().Get(api.OracleBaseRoute + api.RouteGETOraclePublicKey)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	actual := &api.OraclePublicKeyResponse{}
	err = json.Unmarshal(resp.Body(), actual)
	assert.NoError(t, err)
	actualPubkey, err := dlccrypto.NewPublicKey(actual.PublicKey)
	assert.NoError(t, err)
	assert.Equal(t, helper.ExpectedOracle.PublicKey, actualPubkey)
}
