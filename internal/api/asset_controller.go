package api

import (
	"net/http"
	"p2pderivatives-oracle/internal/database/entity"
	"p2pderivatives-oracle/internal/database/orm"
	"p2pderivatives-oracle/internal/datafeed"
	"p2pderivatives-oracle/internal/dlccrypto"
	"p2pderivatives-oracle/internal/oracle"
	"strconv"
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	URLParamTagTime        = "time"
	RouteGETAssetConfig    = "/config"
	RouteGETAssetRvalue    = "/rvalue/:" + URLParamTagTime
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
	logger := ginlogrus.GetCtxLogger(c)
	logger.Info("Accessing Asset Configuration")
	c.JSON(http.StatusOK, NewResponse(ct.config))
}

// GetAssetRvalue handler returns the stored Rvalue related to the asset and time
// if not present and future time, it will generates a new one and store it in db
func (ct *AssetController) GetAssetRvalue(c *gin.Context) {
	ginlogrus.SetCtxLoggerHeader(c, "request-header", "Get Asset Rvalue")
	logger := ginlogrus.GetCtxLogger(c)
	logger.Info("Validating assetID and time parameters")
	_, requestedPublishDate, err := validateAssetAndTime(c, ct.assetID, ct.config)
	if err != nil {
		c.Error(err)
		return
	}
	db := c.MustGet("orm").(*orm.ORM).GetDB()
	crypto := c.MustGet(ContextIDCryptoService).(dlccrypto.CryptoService)

	logger.Info("Recovering DLC Data")
	dlcData, err := findOrCreateDLCData(db, crypto, ct.assetID, requestedPublishDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, NewResponse(dlcData))
}

// GetAssetSignature handler returns the stored signature and asset value related to the asset and time
// or if not present, it will generate a new one using previous Rvalue and store it in db
func (ct *AssetController) GetAssetSignature(c *gin.Context) {
	ginlogrus.SetCtxLoggerHeader(c, "request-header", "Get Asset Signature")
	logger := ginlogrus.GetCtxLogger(c)
	logger.Info("Validating Asset and Time")
	_, requestedNearDate, err := validateAssetAndTime(c, ct.assetID, ct.config)
	if err != nil {
		c.Error(err)
		return
	}

	db := c.MustGet(ContextIDOrm).(*orm.ORM).GetDB()
	logger.Info("Calculating Publish date")
	publishDate, err := calculatePublishDate(db, ct.assetID, requestedNearDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}
	if publishDate.After(time.Now().UTC()) {
		cause := errors.Errorf("Oracle cannot sign a value not yet known, retry after %s", publishDate.String())
		c.Error(NewBadRequestError(InvalidTimeTooEarlyBadRequestErrorCode, cause, requestedNearDate.String()))
		return
	}
	crypto := c.MustGet(ContextIDCryptoService).(dlccrypto.CryptoService)
	logger.Info("Recovering DLC data")
	dlcData, err := findOrCreateDLCData(db, crypto, ct.assetID, requestedNearDate, ct.config)
	if err != nil {
		c.Error(err)
		return
	}
	if !dlcData.IsSigned() {
		logger.Info("Computing Signature")
		asset, currency := ParseAssetID(ct.assetID)
		feed := c.MustGet(ContextIDDataFeed).(datafeed.DataFeed)
		value, err := feed.FindPastAssetPrice(asset, currency, dlcData.PublishedDate)
		if err != nil {
			c.Error(NewUnknownDataFeedError(err))
			return
		}

		dlcData.Value = strconv.FormatUint(value, 10)
		oracleInstance := c.MustGet(ContextIDOracle).(*oracle.Oracle)
		kvalue, err := dlccrypto.NewPrivateKey(dlcData.Kvalue)
		if err != nil {
			c.Error(NewUnknownCryptoServiceError(err))
			return
		}
		sig, err := crypto.ComputeSchnorrSignature(oracleInstance.PrivateKey, kvalue, dlcData.Value)
		if err != nil {
			c.Error(NewUnknownCryptoServiceError(err))
			return
		}

		dlcData.Signature = sig.EncodeToString()
		err = entity.UpdateDLCDataSignatureAndValue(db, dlcData.AssetID, dlcData.PublishedDate, dlcData.Signature, dlcData.Value)
		if err != nil {
			c.Error(NewUnknownDBError(err))
			return
		}
	}

	c.JSON(http.StatusOK, NewResponse(dlcData))
}

func findOrCreateDLCData(db *gorm.DB, oracle dlccrypto.CryptoService, assetID string, nearDate time.Time, config AssetConfig) (*entity.DLCData, error) {
	dlcData, err := entity.FindDLCDataPublishedNear(db, assetID, nearDate, config.Frequency)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, NewUnknownDBError(err)
	}
	// if record is not found, need to create the record in db
	if err != nil && gorm.IsRecordNotFoundError(err) {
		publishDate, err := calculatePublishDate(db, assetID, nearDate, config)
		if err != nil {
			return nil, err
		}
		signingK, err := oracle.GenerateKvalue()
		if err != nil {
			return nil, NewUnknownCryptoServiceError(err)
		}
		rvalue, err := oracle.ComputeRvalue(signingK)
		if err != nil {
			return nil, NewUnknownCryptoServiceError(err)
		}
		dlcData, err = entity.CreateDLCData(db, assetID, publishDate, signingK.EncodeToString(), rvalue.EncodeToString(), true)
		if err != nil {
			return nil, NewUnknownDBError(err)
		}
	}

	return dlcData, nil
}

func validateAssetAndTime(c *gin.Context, assetID string, config AssetConfig) (*entity.Asset, time.Time, error) {
	timestampStr := c.Param(URLParamTagTime)
	db := c.MustGet("orm").(*orm.ORM).GetDB()
	asset, err := entity.FindAsset(db, assetID)
	if err != nil {
		return nil, time.Time{}, NewRecordNotFoundDBError(err, assetID)
	}
	requestedPublishDate, err := ParseTime(timestampStr)
	if err != nil {
		return asset, requestedPublishDate, NewBadRequestError(InvalidTimeFormatBadRequestErrorCode, err, timestampStr)
	}
	upTo := time.Now().Add(config.RangeD)
	if requestedPublishDate.After(upTo) {
		cause := errors.Errorf("Invalid time, you cannot request a value that will be sign after %s", upTo.String())
		err = NewBadRequestError(InvalidTimeTooLateBadRequestErrorCode, cause, requestedPublishDate.String())
	}
	return asset, requestedPublishDate, err
}

func calculatePublishDate(db *gorm.DB, assetID string, requestDate time.Time, config AssetConfig) (time.Time, error) {
	var from time.Time
	mostRecent, err := entity.FindDLCDataPublishedBefore(db, assetID, requestDate)
	// no data found
	if err != nil && gorm.IsRecordNotFoundError(err) {
		year, month, day := time.Now().Date()
		from = time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(-1 * config.RangeD)
	}
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return time.Time{}, NewUnknownDBError(err)
	}
	if err == nil {
		from = mostRecent.PublishedDate
	}

	bufD := requestDate.Sub(from)
	div := (bufD / config.Frequency) * config.Frequency
	q := div.Round(config.Frequency)
	publishDate := from.Add(q)
	if publishDate.Before(requestDate) {
		publishDate = publishDate.Add(config.Frequency)
	}
	return publishDate, nil
}

// ParseTime will try to parse a string using RFC3339 format and convert it to a time.Time
func ParseTime(timeParam string) (time.Time, error) {
	// RFC3339 is close to ISO8601
	timestamp, err := time.Parse(TimeFormatISO8601, timeParam)
	if err != nil {
		err = errors.WithMessagef(err, "invalid time format ! You should use ISO8601 ex: %s", TimeFormatISO8601)
	}
	return timestamp.UTC(), err
}

// ParseAssetID will return the asset and currency related to the asset id
func ParseAssetID(assetID string) (asset string, currency string) {
	return assetID[0:3], assetID[3:6]
}
