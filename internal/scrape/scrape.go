package scrape

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leviharrison/fios-exporter/internal/metrics"
)

var hash string

// Init sets up the scraper
func Init(host string) {
	retrieved, err := getHash(host)
	if err != nil {
		log.Fatalf("Could not connect to router: %v", err)
	}

	hash = retrieved
}

// Scrape starts a scraping loop
func Scrape(host, password string) {
	retried := false
	cookies, err := login(host, password)
	if err != nil {
		log.Fatal(err)
	}

	for {
		recieved, status, err := getData(host, cookies)
		if err != nil {
			log.Printf("Error requesting resource: %v", err)
			if retried {
				log.Fatalf("Back to back requests with errors, exiting: %v", err)
			}
			retried = true
			time.Sleep(15 * time.Second)
			continue
		}

		if status != 200 {
			if status == 401 {
				log.Printf("Relogging in")
				cookies, err = login(host, password)
				if err != nil {
					if retried {
						log.Fatalf("Back to back requests with errors, exiting: %v", err)
					}
					retried = true
					log.Printf("Failed relogging in")
				}
				time.Sleep(15 * time.Second)
				continue
			}

			if retried {
				log.Fatalf("Back to back requests with status code %d, exiting", status)
			}
			retried = true
			time.Sleep(15 * time.Second)
			continue
		}

		metrics.TXMinute1.Set(float64(recieved.Bandwidth.MinutesTX[0] * 8))
		metrics.RXMinute1.Set(float64(recieved.Bandwidth.MinutesRX[0] * 8))

		retried = false
		time.Sleep(15 * time.Second)
		continue
	}
}

type data struct {
	Bandwidth bandwidth
}

type bandwidth struct {
	HoursRX   []int
	HoursTX   []int
	MinutesRX []int
	MinutesTX []int
}

func getData(host string, cookies []string) (data, int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", host+"/api/network/1", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Cookie", cookies[0]+cookies[1])
	req.Header.Add("X-XSRF-TOKEN", cookies[1][11:len(cookies[1])-1])

	res, err := client.Do(req)
	if err != nil {
		return data{}, res.StatusCode, err
	}

	if res.StatusCode != 200 {
		return data{}, res.StatusCode, nil
	}

	defer res.Body.Close()

	result := data{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result, res.StatusCode, nil
}

type loginPayload struct {
	Password string `json:"password"`
}

func login(host, password string) ([]string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	hasher := sha512.New()
	passwordBytes := append([]byte(password), []byte(hash)...)
	hasher.Write(passwordBytes)
	result := hex.EncodeToString(hasher.Sum(nil))

	payload, err := json.Marshal(loginPayload{result})
	if err != nil {
		log.Fatal(err)
	}

	b := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", host+"/api/login", b)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		return []string{res.Header["Set-Cookie"][0][:len(res.Header["Set-Cookie"][0])-7], res.Header["Set-Cookie"][1][:len(res.Header["Set-Cookie"][1])-7]}, nil
	}

	return nil, fmt.Errorf("Invalid password or error")
}

type invalidLoginResponse struct {
	PasswordSalt string
}

func getHash(host string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", host+"/api", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 401 {
		defer res.Body.Close()

		result := invalidLoginResponse{}
		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		return result.PasswordSalt, nil
	}

	return "", fmt.Errorf("Could not get password hash")
}
