package data

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	DataCollector sarama.SyncProducer
	DeadLetter    sarama.AsyncProducer
}

func (s *Server) collectSensorData(topic string, deadLetterTopic string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal
		var msg SensorData
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Do some processing
		if err = validateData(msg); err != nil {
			invalid, err := json.Marshal(msg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// We are not setting a message key, which means that all messages will
			// be distributed randomly over the different partitions.
			_, _, err = s.DataCollector.SendMessage(&sarama.ProducerMessage{
				Topic: deadLetterTopic,
				Value: sarama.StringEncoder(invalid),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		validMsg, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// We are not setting a message key, which means that all messages will
		// be distributed randomly over the different partitions.
		_, _, err = s.DataCollector.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(validMsg),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (s *Server) Handler(topic string, deadLetter string) http.Handler {
	return s.collectSensorData(topic, deadLetter)
}

func (s *Server) Run(addr string, topic string, deadLetter string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.Handler(topic, deadLetter),
	}
	log.Printf("Listening for requests on %s...\n", addr)
	return httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	if err := s.DataCollector.Close(); err != nil {
		log.Println("Failed to shut down data collector cleanly", err)
	}
	if err := s.DeadLetter.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
	return nil
}
