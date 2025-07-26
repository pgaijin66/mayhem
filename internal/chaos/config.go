package chaos

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

// ChaosConfig represents the configuration for chaos injection
type ChaosConfig struct {
	DelayEnabled       bool     `json:"delay_enabled"`
	DelayMin           Duration `json:"delay_min"`
	DelayMax           Duration `json:"delay_max"`
	DelayProbability   float64  `json:"delay_probability"`
	ErrorEnabled       bool     `json:"error_enabled"`
	ErrorCodes         []int    `json:"error_codes"`
	ErrorProbability   float64  `json:"error_probability"`
	ErrorMessage       string   `json:"error_message"`
	TimeoutEnabled     bool     `json:"timeout_enabled"`
	TimeoutDuration    Duration `json:"timeout_duration"`
	TimeoutProbability float64  `json:"timeout_probability"`
}

// NewConfigFromFlags creates a new configuration from command line flags
func NewConfigFromFlags(delayMin, delayMax time.Duration, delayProb, errorProb float64,
	errorCodes, errorMsg string, timeoutDur time.Duration, timeoutProb float64) (*ChaosConfig, error) {

	var codes []int
	for _, code := range strings.Split(errorCodes, ",") {
		if c, err := strconv.Atoi(strings.TrimSpace(code)); err == nil {
			codes = append(codes, c)
		}
	}

	return &ChaosConfig{
		DelayEnabled:       true,
		DelayMin:           Duration{delayMin},
		DelayMax:           Duration{delayMax},
		DelayProbability:   delayProb,
		ErrorEnabled:       true,
		ErrorCodes:         codes,
		ErrorProbability:   errorProb,
		ErrorMessage:       errorMsg,
		TimeoutEnabled:     true,
		TimeoutDuration:    Duration{timeoutDur},
		TimeoutProbability: timeoutProb,
	}, nil
}

// LoadFromFile loads configuration from a JSON file
func (c *ChaosConfig) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}
