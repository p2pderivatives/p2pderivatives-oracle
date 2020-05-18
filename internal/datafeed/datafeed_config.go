package datafeed

// Config corresponding to the datafeed configuration, contains any abstract config used by datafeed
// ex: `cryptocompare: {...}`
type Config struct {
	DataFeedConfigs map[string]interface{} `configkey:"datafeed"`
}