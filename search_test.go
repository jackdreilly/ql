package main

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	fmt.Println(GetLyrics(GoogleSearch("beat up old guitar")))
}
