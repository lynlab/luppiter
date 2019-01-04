package services

import (
	"fmt"
	"testing"
)

func TestCheckPermission(t *testing.T) {
	// checkPermission should pass for valid api key
	_, err := checkPermission(testAPIKey, testService)
	if err != nil {
		fmt.Printf("checkPermission returns error: %v\n", err)
		t.Fail()
		return
	}

	// checkPermission should return ErrUnauthorized for invalid key
	_, err = checkPermission(testAPIKey+"_invalid", testService)
	if err == nil || err != ErrUnauthorized {
		fmt.Printf("checkPermission should return ErrUnauthorized: %v\n", err)
		t.Fail()
		return
	}

	// checkPermission should return ErrInsuffPermission for non existing permission
	_, err = checkPermission(testAPIKey, testService+"_invalid")
	if err == nil || err != ErrInsuffPermission {
		fmt.Printf("checkPermission should return ErrUnauthorized: %v\n", err)
		t.Fail()
		return
	}
}
