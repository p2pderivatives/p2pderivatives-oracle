package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"p2pderivatives-oracle/src/api"
	"p2pderivatives-oracle/test/integration/helper"
)

func TestMain(m *testing.M) {
	helper.InitHelper()
	os.Exit(m.Run())
}

func Test_GetHello_Returns_HelloWorld(t *testing.T) {
	client := helper.CreateDefaultClient()
	resp, err := client.R().Get(api.HelloBaseRoute)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	var worldArr []string
	for i := 0; i < helper.APIConfig.HelloCount; i++ {
		worldArr = append(worldArr, fmt.Sprintf("world %d", (i+1)))
	}
	expectedBody, err := json.Marshal(gin.H{"hello": worldArr})
	if err != nil {
		t.Error(err)
	}
	actualBody := resp.String()
	assert.JSONEq(t, string(expectedBody), actualBody)
}
