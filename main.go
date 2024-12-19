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
	payload := [][]byte{
		[]byte("key\n"), []byte("value\n"),
		[]byte("uper key start\n"), []byte("super key value\n"),
	}

	superRider, err := CreateReq(payload)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ProcessII(superRider)
	// fmt.Println(string(superRider))
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
			temp := pointer
			if pointer < 4 {
				copy(header[temp:], buf[:4-temp])
				pointer = 4 - pointer
				temp = 0
			}

			end := client.DecodeLength(header)
			copy(slabBloc[temp:], buf[pointer:end])
			active = false
		}

		for {
			var end int
			if !active {
				if pointer+4 > n {
					copy(header, buf[pointer:])
					pointer = n - pointer
					active = true
					break
				}

				end = pointer + client.DecodeLength(buf[pointer:pointer+4]) + 4

				slabBloc = make([]byte, 64)
				if end > n {
					copy(slabBloc, buf[pointer:])
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

func Req(buf []byte) {
	// fmt.Println(string(buf))
}

func CreateReq(msg [][]byte) ([]byte, error) {
	if len(msg)%2 != 0 {
		return nil, ErrWrongInput
	}

	superRider := make([]byte, 0, 10*1024)

	for index := 0; index < len(msg); index += 2 {
		payload, err := client.Set(msg[index], msg[index+1], 0)
		fmt.Println("[payload length]", len(payload))
		if err != nil {
			return nil, err
		}

		fmt.Printf("[encoded value] ")
		fmt.Println(client.Decode(payload))

		// copy(superRider[pointer:], payload)

		// pointer += len(payload)
		superRider = append(superRider, payload...)
	}

	return superRider, nil
}
