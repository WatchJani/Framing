package main

import (
	"log"
	"os"
	"testing"
)

func BenchmarkParsingSingReqI(b *testing.B) {
	b.StopTimer()

	payload := [][]byte{
		[]byte("key"), []byte("value"),
		[]byte("uper key start\n\n"), []byte("super key value\n\n"),
	}

	superRider, err := CreateReq(payload)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ProcessI(superRider)
	}
}

func BenchmarkParsingSingReqII(b *testing.B) {
	b.StopTimer()

	payload := [][]byte{
		[]byte("key"), []byte("value"),
		[]byte("uper key start\n\n"), []byte("super key value\n\n"),
	}

	superRider, err := CreateReq(payload)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ProcessII(superRider)
	}
}
