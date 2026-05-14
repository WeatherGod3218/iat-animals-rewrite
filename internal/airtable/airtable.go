package airtable

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type AirtableResponse struct {
	Records []map[string]interface{} `json:"records"`
	Offset  string                   `json:"offset"`
}

func GetAirtableURI(table string) ([]map[string]interface{}, error) {
	base := os.Getenv("AIRTABLE_BASE")
	baseURL := fmt.Sprintf(
		"https://api.airtable.com/v0/%s/%s",
		base,
		url.PathEscape(table),
	)

	var allRecords []map[string]interface{}
	var offset string

	client := &http.Client{}

	for {
		fetchURL := baseURL

		if offset != "" {
			fetchURL += "?offset=" + url.QueryEscape(offset)
		}

		req, err := http.NewRequest("GET", fetchURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+os.Getenv("AIRTABLE_API_KEY"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(
				"airtable error: %d %s",
				resp.StatusCode,
				string(body),
			)
		}

		var result AirtableResponse

		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		allRecords = append(allRecords, result.Records...)

		if result.Offset == "" {
			break
		}

		offset = result.Offset
	}

	return allRecords, nil
}
