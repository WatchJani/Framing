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
	data, err := os.ReadFile("./test_data.bin")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// msg := string([]byte{66, 99, 77, 87, 55, 80, 79, 108, 69, 116, 97, 111, 101, 77, 73, 118, 72, 83, 66, 52})
	// fmt.Println(msg)

	payload := make([][]byte, 0, 2)

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
		// fmt.Println(end)
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
			temp := 4096 - pointer
			if pointer < 4 {
				copy(header[temp:], buf[:4-temp])
				pointer = 4 - pointer
				temp = 0
			} else {
				pointer = 0
			}

			end := client.DecodeLength(header)
			// fmt.Println(buf[pointer : end-temp+4])
			copy(slabBloc[:temp], buf[pointer:end-temp+4])
			pointer += end - temp + 4
			active = false
		}

		//30 0 0 0 83 8 0 0 0 0 12 0 0 0 66 99 77 87 55 80 79 108 69 116 97 111 101 77 73 118

		for {
			var end int
			if !active {
				if pointer+4 > n {
					copy(header, buf[pointer:])
					pointer = n - pointer
					active = true
					break
				} else {
					copy(header, buf[pointer:pointer+4])
				}

				end = pointer + client.DecodeLength(buf[pointer:pointer+4]) + 4

				slabBloc = make([]byte, 64)
				if end > n {
					copy(slabBloc, buf[pointer:])
					active = true
					break
				}

				copy(slabBloc, buf[pointer+4:end])
			}

			Req(slabBloc)
			pointer = end
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

func Req(buf []byte) {
	_, key, _, body := client.Decode(buf)
	headerSize := 10
	fmt.Println(string(buf[headerSize : headerSize+int(key)+int(body)]))
}

func CreateReq(msg [][]byte) ([]byte, error) {
	if len(msg)%2 != 0 {
		return nil, ErrWrongInput
	}

	superRider := make([]byte, 0, 10*1024)

	for index := 0; index < len(msg); index += 2 {
		payload, err := client.Set(msg[index], msg[index+1], 0)
		// fmt.Println("[payload length]", len(payload))
		if err != nil {
			return nil, err
		}

		// fmt.Printf("[encoded value] ")
		// fmt.Println(client.Decode(payload))

		// copy(superRider[pointer:], payload)

		// pointer += len(payload)
		superRider = append(superRider, payload...)
	}

	return superRider, nil
}
