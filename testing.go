package main

func TestNewConfig(url string) *Config {
	return &Config{
		APIKey:       "testapikey",
		APISecret:    "testapisecret",
		BaseEndpoint: url,
		Debug:        true,
	}
}
