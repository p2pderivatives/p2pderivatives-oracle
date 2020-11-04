package api

import (
	"fmt"
	"math"
	"net/http"
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/datafeed"
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/oracle"
	"time"

	"github.com/cryptogarageinc/server-common-go/pkg/database/orm"
	"github.com/cryptogarageinc/server-common-go/pkg/utils/iso8601"

	"github.com/sirupsen/logrus"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// URLParamTagTime Tag to use as date parameter in route
	URLParamTagTime = "time"
	// RouteGETAssetConfig relative GET route to retrieve asset configuration
	RouteGETAssetConfig = "/config"
	// RouteGETAssetRvalue relative GET route to retrieve asset rvalue
	RouteGETAssetRvalue = "/rvalue/:" + URLParamTagTime
	// RouteGETAssetSignature relative GET route to retrieve asset signature
	RouteGETAssetSignature = "/signature/:" + URLParamTagTime
)

// AssetController represents the asset api Controller
type AssetController struct {
	assetID string
	config  AssetConfig
}

// NewAssetController creates a new Controller structure with the given parameters.
func NewAssetController(assetID string, config AssetConfig) Controller {
	return &AssetController{
		assetID: assetID,
		config:  config,
	}
}

// Routes list and binds all routes to the router group provided
func (ct *AssetController) Routes(route *gin.RouterGroup) {
	route.GET(RouteGETAssetRvalue, ct.GetAssetRvalue)
	route.GET(RouteGETAssetSignature, ct.GetAssetSignature)
	route.GET(RouteGETAssetConfig, ct.GetConfiguration)
}

// GetConfiguration handler returns the asset configuration
func (ct *AssetController) GetConfiguration(c *gin.Context) {
	ginlogrus.SetCtxLoggerHeader(c, "request-header", "Get Asset Configuration")
	c.JSON(http.StatusOK, &AssetConfigResponse{
		StartDate: ct.config.StartDate,
		Frequency: iso8601.EncodeDuration(ct.config.Frequency),
		RangeD:    iso8601.EncodeDuration(ct.config.RangeD),
	})
}

// GetAssetRvalue handler returns the stored Rvalue related to the asset and time
// if not present and future time, it will generates a new one using the config start date as reference
func (ct *AssetController) GetAssetRvalue(c *gin.Context) {
	ginlogrus.SetCtxLoggerHeader(c, "request-header", "Get Asset Rvalue")
	logger := ginlogrus.GetCtxLogger(c)
	_, requestedDate, err := validateAssetAndTime(c, ct.assetID)
	if err != nil {
		c.Error(err)
		return
	}
	publishDate, err := calculatePublishDate(*requestedDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}

	oracleInstance := c.MustGet(ContextIDOracle).(*oracle.Oracle)
	db := c.MustGet(ContextIDOrm).(*orm.ORM).GetDB()
	crypto := c.MustGet(ContextIDCryptoService).(dlccrypto.CryptoService)
	dlcData, err := findOrCreateDLCData(logger, db, crypto, ct.assetID, *publishDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, NewDLCDataResponse(oracleInstance.PublicKey, dlcData))
}

// GetAssetSignature handler returns the stored signature and asset value related to the asset and time
// or if not present, it will generate a new one using the config start date as reference
func (ct *AssetController) GetAssetSignature(c *gin.Context) {
	ginlogrus.SetCtxLoggerHeader(c, "request-header", "Get Asset Signature")
	logger := ginlogrus.GetCtxLogger(c)
	_, requestedDate, err := validateAssetAndTime(c, ct.assetID)
	if err != nil {
		c.Error(err)
		return
	}
	publishDate, err := calculatePublishDate(*requestedDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}

	// check the signature has been published
	if publishDate.After(time.Now().UTC()) {
		cause := errors.Errorf("Oracle cannot sign a value not yet known, retry after %s", publishDate.String())
		c.Error(NewBadRequestError(InvalidTimeTooEarlyBadRequestErrorCode, cause, requestedDate.String()))
		return
	}

	oracleInstance := c.MustGet(ContextIDOracle).(*oracle.Oracle)
	db := c.MustGet(ContextIDOrm).(*orm.ORM).GetDB()
	crypto := c.MustGet(ContextIDCryptoService).(dlccrypto.CryptoService)
	dlcData, err := findOrCreateDLCData(logger, db, crypto, ct.assetID, *publishDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}
	if !dlcData.IsSigned() {
		logger.Debug("Computing Signature")
		asset, currency := ParseAssetID(ct.assetID)
		feed := c.MustGet(ContextIDDataFeed).(datafeed.DataFeed)
		value, err := feed.FindPastAssetPrice(asset, currency, dlcData.PublishedDate)
		if err != nil {
			c.Error(NewUnknownDataFeedError(err))
			return
		}

		// round datafeed price to neareast integer
		valueMessage := fmt.Sprintf("%d", int(math.Round(*value)))
		oracleInstance := c.MustGet(ContextIDOracle).(*oracle.Oracle)
		kvalue, err := dlccrypto.NewPrivateKey(dlcData.Kvalue)
		if err != nil {
			c.Error(NewUnknownCryptoServiceError(err))
			return
		}
		sig, err := crypto.ComputeSchnorrSignature(oracleInstance.PrivateKey, kvalue, valueMessage)
		if err != nil {
			c.Error(NewUnknownCryptoServiceError(err))
			return
		}

		dlcData, err = entity.UpdateDLCDataSignatureAndValue(
			db,
			dlcData.AssetID,
			dlcData.PublishedDate,
			sig.EncodeToString(),
			valueMessage)

		if err != nil {
			c.Error(NewUnknownDBError(err))
			return
		}
	}

	c.JSON(http.StatusOK, NewDLCDataResponse(oracleInstance.PublicKey, dlcData))
}

func findOrCreateDLCData(logger *logrus.Entry, db *gorm.DB, oracle dlccrypto.CryptoService, assetID string, publishDate time.Time, config AssetConfig) (*entity.DLCData, error) {
	dlcData, err := entity.FindDLCDataPublishedAt(db, assetID, publishDate)
	if err == nil {
		logger.Debug("Found a matching DLC Data in db")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, NewUnknownDBError(err)
	}

	// if record is not found, need to create the record in db
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Debug("Generating new DLC data Rvalue")
		signingK, rvalue, err := oracle.GenerateSchnorrKeyPair()
		if err != nil {
			return nil, NewUnknownCryptoServiceError(err)
		}
		dlcData, err = entity.CreateDLCData(
			db,
			assetID,
			publishDate,
			signingK.EncodeToString(),
			rvalue.EncodeToString())
		if err != nil {
			// need to retry to be sure a concurrent didn't try to create same DLCData
			inDb, errFind := entity.FindDLCDataPublishedAt(db, assetID, publishDate)
			if errFind != nil {
				return nil, NewUnknownDBError(err)
			}
			dlcData = inDb
		}
	}

	return dlcData, nil
}

func validateAssetAndTime(c *gin.Context, assetID string) (*entity.Asset, *time.Time, error) {
	timestampStr := c.Param(URLParamTagTime)
	db := c.MustGet(ContextIDOrm).(*orm.ORM).GetDB()
	asset, err := entity.FindAsset(db, assetID)
	if err != nil {
		return nil, nil, NewRecordNotFoundDBError(err, assetID)
	}
	requestedPublishDate, err := ParseTime(timestampStr)
	if err != nil {
		return asset, requestedPublishDate, NewBadRequestError(InvalidTimeFormatBadRequestErrorCode, err, timestampStr)
	}
	return asset, requestedPublishDate, err
}

func calculatePublishDate(requestDate time.Time, config AssetConfig) (*time.Time, error) {
	// date to use as publish date reference
	from := config.StartDate

	// calculate the difference between the requested date and the reference
	// round up to the frequency
	durationDiff := requestDate.Sub(from)
	frequencyMultiple := durationDiff.Round(config.Frequency)
	publishDate := from.Add(frequencyMultiple)
	// if round below (floor) then add one frequency duration
	if publishDate.Before(requestDate) {
		publishDate = publishDate.Add(config.Frequency)
	}

	// check publish date in range
	upTo := time.Now().UTC().Add(config.RangeD)
	if publishDate.After(upTo) {
		cause := errors.Errorf(
			"Requested Date not in oracle range, you cannot request a DLC Data that will be published after %s",
			upTo.String())
		return nil, NewBadRequestError(InvalidTimeTooLateBadRequestErrorCode, cause, publishDate.String())
	}
	return &publishDate, nil
}

// ParseTime will try to parse a string using ISO8691 format and convert it to a time.Time
func ParseTime(timeParam string) (*time.Time, error) {
	timestamp, err := time.Parse(TimeFormatISO8601, timeParam)
	if err != nil {
		err = errors.WithMessagef(err, "Invalid time format ! You should use ISO8601 ex: %s", TimeFormatISO8601)
		return nil, err
	}
	utc := timestamp.UTC()
	return &utc, nil
}

// ParseAssetID will return the asset and currency related to the asset id
func ParseAssetID(assetID string) (asset string, currency string) {
	return assetID[0:3], assetID[3:6]
}
