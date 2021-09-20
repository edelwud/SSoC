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

type Payload struct {
	Subject string    `json:"subject"`
	Issued  time.Time `json:"issued"`
	Expires time.Time `json:"expires"`
}

const ExpirationDays = 3

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

func GenerateToken(subject string) Payload {
	issued := time.Now()
	expires := issued.AddDate(0, 0, ExpirationDays)
	return Payload{Subject: subject, Issued: issued, Expires: expires}
}

func GenerateMACToken() (Payload, error) {
	macAddr, err := getMacAddr()
	if err != nil {
		return Payload{}, err
	}
	return GenerateToken(strings.Join(macAddr, ",")), nil
}

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

func ValidateToken(payload Payload) error {
	if time.Now().Second() > payload.Expires.Second() {
		return errors.New("token expired")
	}
	return nil
}

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
