package token

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"
)

// Payload includes serialization rules for access token
type Payload struct {
	Subject string    `json:"subject"`
	Issued  time.Time `json:"issued"`
	Expires time.Time `json:"expires"`
}

// ExpirationDays is the number of days when access token is valid
const ExpirationDays = 3

// Row serializes Payload
func (p Payload) Row() (string, error) {
	var buf bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)

	err := json.NewEncoder(encoder).Encode(p)
	if err != nil {
		return "", err
	}

	err = encoder.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateToken generates token from subject, uses current time for Payload.Issued parameter
func GenerateToken(subject string) Payload {
	issued := time.Now()
	expires := issued.AddDate(0, 0, ExpirationDays)
	return Payload{Subject: subject, Issued: issued, Expires: expires}
}

// GenerateMACToken generates access token where Payload.Subject is joined slice of available MAC addresses
// uses GenerateToken for background
func GenerateMACToken() (Payload, error) {
	macAddr, err := getMacAddr()
	if err != nil {
		return Payload{}, err
	}
	return GenerateToken(strings.Join(macAddr, ",")), nil
}

// ParseToken parses access token, returns Payload structure
func ParseToken(token []byte) (Payload, error) {
	reader := strings.NewReader(string(token))
	decoder := base64.NewDecoder(base64.StdEncoding, reader)

	payload := Payload{}
	err := json.NewDecoder(decoder).Decode(&payload)
	if err != nil {
		return Payload{}, err
	}

	return payload, err
}

// ValidateToken receives Payload structure, returns error when time.Now > Payload.Expires
func ValidateToken(payload Payload) error {
	if time.Now().Second() > payload.Expires.Second() {
		return errors.New("token expired")
	}
	return nil
}

// getMacAddr returns slice of available on current machine MAC addresses
func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
