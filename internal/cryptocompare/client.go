package cryptocompare

import (
	"fmt"
	"p2pderivatives-oracle/internal/datafeed"
	"strings"
	"time"

	"github.com/cryptogarageinc/server-common-go/pkg/log"

	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
)

const (
	priceRoute           = "/price"
	pricePastHourRoute   = "/v2/histohour"
	pricePastMinuteRoute = "/v2/histominute"
	limitPastResponse    = 1
)

// NewClient returns a new CryptoCompare Client (not initialized)
func NewClient(l *log.Log, config *Config) *Client {
	return &Client{
		config:      config,
		initialized: false,
		log:         l,
	}
}

type apiPriceResponse map[string]float64
type apiPastPriceResponse struct {
	Data struct {
		TimeFrom int64 `json:"TimeFrom"`
		TimeTo   int64 `json:"TimeTo"`
		// Data data response with time (like one element for each minute/hour)
		// the requested time should be the last element
		Data []struct {
			Time  int64   `json:"time"`
			Close float64 `json:"close"`
		} `json:"Data"`
	} `json:"Data"`
}

// Client represents a CryptoCompare REST client
type Client struct {
	datafeed.DataFeed
	config      *Config
	httpClient  *resty.Client
	initialized bool
	log         *log.Log
}

// Initialize initializes the http client
func (c *Client) Initialize() {
	c.httpClient = resty.New()
	c.httpClient.SetHostURL(c.config.APIBaseURL)
	c.httpClient.SetHeader("Accept", "application/json")
	if c.config.APIKey != "" {
		c.httpClient.SetHeader("authorization", "Apikey "+c.config.APIKey)
	}
	c.initialized = true
}

// IsInitialized returns true if the Client has been initialized
func (c *Client) IsInitialized() bool {
	return c.initialized
}

// FindCurrentAssetPrice sends a GET request to the CryptoCompare API to retrieve the current price of an asset
func (c *Client) FindCurrentAssetPrice(assetID string) (*float64, error) {
	assetConfig, ok := c.config.AssetsConfig[assetID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Could not find config for asset %v", assetID))
	}
	route := fmt.Sprintf(priceRoute+"?fsym=%s&tsyms=%s", assetConfig.Fsym, assetConfig.Tsym)
	resp, err := c.getAssetPrice(route, apiPriceResponse{})
	if err != nil {
		return nil, err
	}

	res := *(resp.Result().(*apiPriceResponse))
	val, ok := res[strings.ToUpper(assetConfig.Tsym)]

	// it should not happened if the request was well formed
	if !ok {
		return nil, errors.Errorf("error currency %s not found in cryptocompare response", assetConfig.Tsym)
	}

	return &val, nil
}

// FindPastAssetPrice sends a GET request to the CryptoCompare API to retrieve a past price of an asset
func (c *Client) FindPastAssetPrice(assetID string, date time.Time) (*float64, error) {
	now := time.Now()
	if now.Before(date) {
		return nil, errors.New("date should be before now")
	}
	var precisionRoute string
	// before seven days (minute precision are stored only seven days in cryptocompare)
	if date.Before(now.Add(-time.Hour * 168)) {
		precisionRoute = pricePastHourRoute
	} else {
		precisionRoute = pricePastMinuteRoute
	}
	var assetConfig, ok = c.config.AssetsConfig[assetID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No config found for asset %v", assetID))
	}
	route := fmt.Sprintf(
		precisionRoute+"?fsym=%s&tsym=%s&toTs=%d&limit=%d",
		assetConfig.Fsym,
		assetConfig.Tsym,
		date.Unix(),
		limitPastResponse)
	resp, err := c.getAssetPrice(route, apiPastPriceResponse{})
	if err != nil {
		return nil, err
	}

	res := resp.Result().(*apiPastPriceResponse)

	// should not happen
	if len(res.Data.Data) != limitPastResponse+1 {
		c.log.Logger.Errorln("Unexpected data for route ", route, "got response ", resp)
		return nil, errors.New("cryptocompare response did not contain the requested element")
	}

	// limitPastResponse should be the last element
	value := res.Data.Data[limitPastResponse].Close
	return &value, nil
}

func (c *Client) getAssetPrice(route string, resultType interface{}) (*resty.Response, error) {
	if !c.IsInitialized() {
		return nil, errors.New("crypto compare client is not initialized")
	}
	req := c.httpClient.R()
	req.SetResult(resultType)
	resp, err := req.Get(route)
	if err != nil {
		return nil, errors.WithMessagef(err, "error while sending a request to cryptocompare api %v", resp.String())
	}

	return resp, nil
}
