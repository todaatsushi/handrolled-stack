package producer

import (
	"io"
	"log"
	"net"

	"github.com/todaatsushi/queue/internal/messages"
)

func QueueTask(w io.Writer, msg string) error {
	parsed := messages.NewMessage(messages.Enqueue, msg)

	data, err := parsed.MarshalBinary()
	if err != nil {
		return err
	}

	n, err := w.Write(data)
	if err != nil {
		return err
	}

	if n != 4+len(msg)+1 {
		return err
	}
	return nil
}

func QueueTasks(port int, msgs ...string) {
	for _, msg := range msgs {
		errs := []error{}
		var err error

		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: port})
		if err != nil {
			log.Fatal(err)
		}

		err = QueueTask(conn, msg)
		if err != nil {
			errs = append(errs, err)
		}

		for _, e := range errs {
			log.SetPrefix("ERRS:" + "\t")
			log.Println(e.Error())
		}
	}
}
