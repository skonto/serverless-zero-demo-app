package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/skonto/serverless-zero-demo-app/pkg/data"
)

var (
	addr       = flag.String("addr", ":8080", "The address to bind to")
	topic      = flag.String("topic", os.Getenv("KAFKA_TOPIC"), "The topic name to use for valid data")
	deadLetter = flag.String("dead", os.Getenv("DEAD_TOPIC"), "The topic name to use for invalid data")
	brokers    = flag.String("brokers", os.Getenv("KAFKA_BROKERS"), "The Kafka brokers to connect to, as a comma separated list")
	verbose    = flag.Bool("verbose", false, "Turn on Sarama logging")
	certFile   = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile    = flag.String("key", "", "The optional key file for client authentication")
	caFile     = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl  = flag.Bool("verify", false, "Optional verify ssl certificates chain")
)

func main() {
	flag.Parse()

	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	server := &data.Server{
		DataCollector: data.NewDataCollector(brokerList, certFile, keyFile, caFile, verifySsl),
		DeadLetter:    data.NewDeadLetterProducer(brokerList, certFile, keyFile, caFile, verifySsl),
	}
	defer func() {
		if err := server.Close(); err != nil {
			log.Println("Failed to close server", err)
		}
	}()

	log.Fatal(server.Run(*addr, *topic, *deadLetter))
}
