package client

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/todaatsushi/basic_tcp/internal/encoding"
)

func Send(msg string, port int, translator encoding.Translator) {
	if len(msg) == 0 {
		log.Fatal("Can't send empty message.")
	}

	log.Printf("Encoding '%s' and sending to server %d times.(length: %d).", msg, 10, len(msg))

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Artificial stream creation - it should send as one message
	encoded := translator.Encode(msg)

	for range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		_, err := conn.Write(encoded)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Stdin(port int, translator encoding.Translator) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for stdinScanner := bufio.NewScanner(os.Stdin); stdinScanner.Scan(); {
		log.Printf("sent: %s\n", stdinScanner.Text())

		msg := stdinScanner.Text()

		encoded := translator.Encode(msg)
		if _, err := conn.Write(encoded); err != nil {
			log.Fatalf("error writing to %s: %v", conn.RemoteAddr(), err)
		}

		if stdinScanner.Err() != nil {
			log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
		}
	}
}
