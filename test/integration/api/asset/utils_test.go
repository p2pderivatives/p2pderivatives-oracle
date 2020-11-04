// +build integration

package assetapi_test

import (
	"p2pderivatives-oracle/internal/api"
	"strings"
	"time"
)

func GetRouteAssetAnnouncement(assetID string, requestedDate time.Time) string {
	route := api.AssetBaseRoute + "/" + assetID + api.RouteGETAssetAnnouncement
	return strings.Replace(route, ":"+api.URLParamTagTime, requestedDate.Format(api.TimeFormatISO8601), 1)
}

func GetRouteAssetAttestation(assetID string, requestedDate time.Time) string {
	route := api.AssetBaseRoute + "/" + assetID + api.RouteGETAssetAttestation
	return strings.Replace(route, ":"+api.URLParamTagTime, requestedDate.Format(api.TimeFormatISO8601), 1)
}
