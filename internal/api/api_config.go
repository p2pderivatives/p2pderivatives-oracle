package api

// Configuration contains the API configuration
type Config struct {
	HelloCount int `configkey:"api.hello.count" default:"1"`
}
