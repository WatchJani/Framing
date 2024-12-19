package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"root/client"
)

var (
	ErrWrongInput = errors.New("wrong input")
)

func main() {
	data, err := os.ReadFile("./key_test.bin")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	payload := make([][]byte, 0)

	var start int
	for index := range data {
		if data[index] == '\n' {
			payload = append(payload, data[start:index])
			start = index + 1
		}
	}

	superRider, err := CreateReq(payload)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ProcessII(superRider)
}

// stage I
func ProcessI(buf []byte) {
	var (
		pointer int
		n       = len(buf)
	)

	for {
		if pointer+4 > n {
			break
		}

		end := pointer + client.DecodeLength(buf[pointer:pointer+4]) + 4
		Req(buf[pointer+4 : end])
		pointer = end
	}
}

func ProcessII(all []byte) {
	var (
		pointer  int
		active   bool
		header   = make([]byte, 4)
		slabBloc []byte
	)

	for _, buf := range Read(all) {
		n := len(buf)

		if active {
			temp := pointer //koliko je zapisano podataka u prenosu
			pointer = 0

			if temp < 4 {
				pointer = 4 - temp
				temp += copy(header[temp:], buf[:pointer])
				temp = 0
			}

			end := client.DecodeLength(header)
			copy(slabBloc[temp:], buf[pointer:end-temp+4])
			pointer += end - temp
		}

		for {
			if !active {
				if pointer+4 > n {
					copy(header, buf[pointer:])
					pointer = n - pointer
					active = true
					break
				} else {
					copy(header, buf[pointer:pointer+4])
				}

				end := pointer + client.DecodeLength(buf[pointer:pointer+4]) + 4

				slabBloc = make([]byte, 4096)
				pointer += 4
				if end > n {
					copy(slabBloc, buf[pointer:])
					pointer = n - pointer
					active = true
					break
				}
				copy(slabBloc, buf[pointer:end])
				pointer = end
			}

			Req(slabBloc)
			active = false
		}
	}

	data, err := os.ReadFile("./key_test.bin")
	if err != nil {
		log.Println(err)
	}

	payload := make([][]byte, 0, 2)

	var start int
	for index := range data {
		if data[index] == '\n' {
			payload = append(payload, data[start:index])
			start = index + 1
		}
	}

	testRealData := make([]string, 0, len(payload)/2)
	for index := 0; index < len(payload); index += 2 {
		testRealData = append(testRealData, string(payload[index])+string(payload[index+1]))
	}

	for index, value := range testRealData {
		if value != test[index] {
			fmt.Println(index, "|", value, "|", test[index])
		}
	}
}

func Read(all []byte) [][]byte {
	const chunkSize = 4096
	f := make([][]byte, 0, (len(all)+chunkSize-1)/chunkSize)

	for index := 0; index < len(all); index += chunkSize {
		end := index + chunkSize
		if end > len(all) {
			end = len(all)
		}

		f = append(f, all[index:end])
	}

	return f
}

// var counter int

var test []string

func Req(buf []byte) {
	_, key, _, body := client.Decode(buf)
	headerSize := 10

	test = append(test, string(buf[headerSize:headerSize+int(key)+int(body)]))
}

func CreateReq(msg [][]byte) ([]byte, error) {
	if len(msg)%2 != 0 {
		return nil, ErrWrongInput
	}

	superRider := make([]byte, 0, 10*1024)

	for index := 0; index < len(msg); index += 2 {
		payload, err := client.Set(msg[index], msg[index+1], 0)
		if err != nil {
			return nil, err
		}

		superRider = append(superRider, payload...)
	}

	return superRider, nil
}
