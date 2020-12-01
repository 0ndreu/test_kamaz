package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	URL    = "http://127.0.0.1:8000"
	HEALTH = "/health"
	COUNT  = 200
	ADD    = "/api/v1/goods/add"
)

type good struct {
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
	Amount   int64 `json:"amount"`
	Object   int64 `json:"object"`
	Method   int64 `json:"method"`
}

func main() {
	client := &http.Client{}
	err := ping(client)
	if err != nil {
		log.Fatalf("Failed ping: %s", err)
	}

	goods := genGoods(COUNT)

	status, err := sendGoods(client, goods)
	if err != nil {
		log.Fatalf("Failed post: %s\nStatus code: %d", err, &status)
	}
	log.Println(*status)
}

func genGoods(maxCount int) []good {
	rand.Seed(time.Now().UTC().UnixNano())
	var itemList []good

	intCh := make(chan good)
	countRand := rand.Intn(maxCount)
	for i := 0; i < countRand; i++ {
		go func() {
			item := good{
				Price:    int64(rand.Intn(maxCount)),
				Quantity: int64(rand.Intn(maxCount)),
				Amount:   int64(rand.Intn(maxCount)),
				Object:   int64(rand.Intn(maxCount)),
				Method:   int64(rand.Intn(maxCount)),
			}
			intCh <- item
		}()
		item := <-intCh
		itemList = append(itemList, item)
	}
	return itemList
}

func ping(client *http.Client) error {
	resp, err := client.Get(fmt.Sprintf("%s%s", URL, HEALTH))
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	resp.Body.Close()
	return nil
}

func sendGoods(client *http.Client, goods []good) (*int, error) {
	goodsJSON, err := json.Marshal(goods)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(fmt.Sprintf("%s%s", URL, ADD), "application/json", bytes.NewBuffer(goodsJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &resp.StatusCode, err
}
