package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CheckHTTPErr(resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		respData, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code - %d; body - %s", resp.StatusCode, respData)
	}

	return nil
}

func ParseHTTPResponse[T any](resp *http.Response) (*T, error) {
	var result T

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parseHTTPResponse.Decode: %w", err)
	}

	return &result, nil
}
