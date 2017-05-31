package api

// TokenInfo contains information about an access token.
type TokenInfo struct {
	Type         string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccountID    string `json:"account_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Devices contains a list of pills and senses connected to the account.
type Devices struct {
	Pills  []Pill  `json:"pills"`
	Senses []Sense `json:"senses"`
}

// DeviceInfo contains common information about a device.
type DeviceInfo struct {
	ID              string `json:"id"`
	State           string `json:"state"`
	Color           string `json:"color"`
	FirmwareVersion string `json:"firmware_version"`
	LastUpdated     int64  `json:"last_updated"`
}

// Pill contains information about a sleeping pill.
type Pill struct {
	DeviceInfo
	BatteryLevel int    `json:"battery_level"`
	BatteryType  string `json:"battery_type"`
}

// Sense contains information about a sense device.
type Sense struct {
	DeviceInfo
	HardwareVersion string   `json:"hw_version"`
	WifiInfo        WifiInfo `json:"wifi_info"`
}

// WifiInfo contains information about a Wifi a Sense is connected to.
type WifiInfo struct {
	SSID        string `json:"ssid"`
	RSSI        int    `json:"rssi"`
	Condition   string `json:"condition"`
	LastUpdated int64  `json:"last_updated"`
}

// RoomInfo contains information about the sensors in the sense device.
type RoomInfo struct {
	Humidity     SensorData `json:"humidity"`
	Light        SensorData `json:"light"`
	Particulates SensorData `json:"particulates"`
	Sound        SensorData `json:"sound"`
	Temperature  SensorData `json:"temperature"`
}

// SensorData contains information about a single sense sensor.
type SensorData struct {
	Unit            string  `json:"unit"`
	Value           float64 `json:"value"`
	LastUpdated     int64   `json:"last_updated_utc"`
	Condition       string  `json:"condition"`
	IdealConditions string  `json:"ideal_conditions"`
	Message         string  `json:"message"`
}
