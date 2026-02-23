package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type XrayConfig struct {
	Inbounds []struct {
		Settings struct {
			Clients []struct {
				ID    string `json:"id"`
				Email string `json:"email"`
			} `json:"clients"`
		} `json:"settings"`
	} `json:"inbounds"`
}

func main() {
	var client_filename = "/client/client-template.json"
	var server_filename = "/server/config.json"

	// Prepare client template
	data, err := os.ReadFile(client_filename)
	if err != nil {
		log.Fatal(err)
	}
	client_tmp := string(data)
	if !strings.Contains(client_tmp, "${UUID}") {
		log.Fatal("${UUID} placeholder must be in %s", client_filename)
	}

	// Extract UUIDs and emails from server config
	data, err = os.ReadFile(server_filename)
	if err != nil {
		log.Fatal(err)
	}
	var server_conf XrayConfig
	if err := json.Unmarshal(data, &server_conf); err != nil {
		log.Fatal(err)
	}
	users := make(map[string]string)
	// Only checks for the first inbound
	for _, client := range server_conf.Inbounds[0].Settings.Clients {
		users[client.ID] = client.Email
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uuid := strings.TrimPrefix(r.URL.Path, "/")
		email, found := users[uuid]
		if uuid == "" || !found {
			http.NotFound(w, r)
			return
		}

		conf := strings.ReplaceAll(client_tmp, "${UUID}", uuid)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(conf))
		log.Printf("Served config for %s\n", email)
	})

	log.Println("Subscription server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
