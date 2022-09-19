package newrelic

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Sender struct {
	Region     string
	LicenseKey string
	SendTimeout int
}

func CreateSenderFromConfig(config string, sendTimeout int) (Sender) {

	s := strings.Split(config, ":")

	var sender Sender
	sender.Region = s[0]
	sender.LicenseKey = s[1]
	sender.SendTimeout = sendTimeout
	return sender
}



func getLogType(logtype int) string{
	switch logtype {
	case 0:
		return "EMERGENCY"
	case 1:
		return "ALERT"
	case 2:
		return "CRITICAL"
	case 3:
		return "ERROR"
	case 4:
		return "WARNING"
	case 5:
		return "NOTICE"
	case 6:
		return "INFO"
	case 7:
		return "DEBUG"
	default:
		return "INFO"
}
}

func endpointUrl(region string) string {
	switch region {
	case "eu":
		return "https://log-api.eu.newrelic.com/log/v1"
	default:
		return "https://log-api.newrelic.com/log/v1"
}
}

func (s *Sender) Send(data []byte) (err error) {
	var result map[string]interface{}
	if err = json.Unmarshal(data, &result); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	result["logtype"] = getLogType(result["level"].(int))
	result["message"] = result["short_message"].(string) + result["full_message"].(string)
	result["timestamp"] = int(result["timestamp"].(float64))
	delete (result, "short_message")
	delete (result, "full_message")
	delete (result, "level")

	output, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	req, err := http.NewRequest("POST", endpointUrl(s.Region), bytes.NewReader(output))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", s.LicenseKey)

	client := &http.Client{
		Timeout: time.Duration(s.SendTimeout) * time.Millisecond,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		err = errors.New(resp.Status)
	}
	return
}
