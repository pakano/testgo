package hs256

import (
	"fmt"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	user := User{UserID: 1, Username: "aa", GrantScope: "x"}
	token, err := GenerateTokenUsingHs256(user)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token = ", token)

	time.Sleep(time.Second * 1)

	my_claim, err := ParseTokenHs256(token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("my claim = %+v", my_claim)
}
