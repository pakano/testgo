package rs256

import (
	"fmt"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	user := User{UserID: 1, Username: "aa", GrantScope: "x"}

	token, err := GenerateTokenUsingRS256(user)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token = ", token)

	time.Sleep(time.Second * 2)

	my_claim, err := ParseTokenRs256(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("my claim = ", my_claim)
}
