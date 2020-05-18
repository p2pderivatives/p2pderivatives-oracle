// +build integration

package assetapi_test

import (
	"p2pderivatives-oracle/internal/api"
	"strings"
	"time"
)

func GetRouteAssetRvalue(assetID string, requestedDate time.Time) string {
	route := api.AssetBaseRoute + "/" + assetID + api.RouteGETAssetRvalue
	return strings.Replace(route, ":"+api.URLParamTagTime, requestedDate.Format(api.TimeFormatISO8601), 1)
}

func GetRouteAssetSignature(assetID string, requestedDate time.Time) string {
	route := api.AssetBaseRoute + "/" + assetID + api.RouteGETAssetSignature
	return strings.Replace(route, ":"+api.URLParamTagTime, requestedDate.Format(api.TimeFormatISO8601), 1)
}
