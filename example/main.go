package main

import (
	"bufio"
	"log"
	"os"

	"github.com/pkg/errors"
	sdb "github.com/streamsdb/driver/go/sdb"
)

func main() {
	// create streamsdb connection
	conn := sdb.MustOpen(`sdb://USERHERE:PASSWORDHERE@sdb03.streamsdb.io:443?tls=1`)
	db := conn.DB("logs")
	defer conn.Close()

	// create a channel to get notified from any errors
	errs := make(chan error)

	streamName := "test2"
	// read user input from stdin and append it to the stream
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := db.Append(streamName, sdb.MessageInput{Value: scanner.Bytes()})
			if err != nil {
				errs <- errors.Wrap(err, "append error")
				return
			}
		}
	}()

	// streamName = "lastnoted_dev_events"

	// watch the inputs streams for messages and print them
	go func() {
		watch := db.Watch(streamName, 0, 10)
		for slice := range watch.Slices {
			for _, msg := range slice.Messages {
				println("received: ", string(msg.Value))
			}
		}

		errs <- errors.Wrap(watch.Err(), "watch error")
	}()

	log.Fatalf((<-errs).Error())
}
