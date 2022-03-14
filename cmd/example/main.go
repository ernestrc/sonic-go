package main

import (
	"context"
	"flag"
	"os"

	sonic "github.com/ernestrc/sonic-go"
	log "github.com/sirupsen/logrus"
)

var (
	sonicAddr = flag.String("addr", "127.0.0.1:10001", "Sonic endpoint address")
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

func main() {
	flag.Parse()

	c := sonic.NewClient(*sonicAddr)
	defer c.Close()

	query := sonic.Query{Query: "10000", Config: map[string]interface{}{
		"class":          "SyntheticSource",
		"seed":           1000,
		"progress-delay": 10,
	}}
	rx, err := c.Stream(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	var schema sonic.StreamSchema

	for msg := range rx {
		switch res := msg.(type) {
		case sonic.StreamStarted:
			log.Println("stream started ... ")
		case sonic.StreamSchema:
			schema = res
			log.Println(res)
		case sonic.StreamProgress:
			log.Println(res)
		case sonic.StreamOutput:
			log.Println(res.UnmarshalRaw(schema))
		case sonic.StreamCompleted:
			log.Printf("stream completed: error: %s\n", msg.(sonic.StreamCompleted).Err)
		}
	}
}
