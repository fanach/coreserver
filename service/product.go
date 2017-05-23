package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/fanach/coreserver/entity"
)

const (
	pathProductsJSON = "products.json"
)

var (
	defaultProducts = []entity.Product{
		entity.Product{
			Name:         "Free",
			Description:  "Free account",
			Price:        0,
			PriceUnit:    entity.PriceUnitRMB,
			DataFlow:     1024,
			DataFlowUnit: entity.DataFlowUnitMB,
			Expire:       1,
			ExpireUnit:   entity.ExpireUnitMonth,
		},
		entity.Product{
			Name:         "1元包月",
			Description:  "1元包月, 10GB",
			Price:        1,
			PriceUnit:    entity.PriceUnitRMB,
			DataFlow:     10 * 1024,
			DataFlowUnit: entity.DataFlowUnitMB,
			Expire:       1,
			ExpireUnit:   entity.ExpireUnitMonth,
		},
	}
)

// GetProducts returns products
func GetProducts() (products *[]entity.Product, err error) {
	// if products.json not found, will use default static values
	products = &defaultProducts

	if _, err = os.Stat(pathProductsJSON); err != nil {
		log.Printf("stat file %s error: %v\n", pathProductsJSON, err)
		return
	}
	content, err := ioutil.ReadFile(pathProductsJSON)
	if err != nil {
		log.Printf("read file %s error: %v\n", pathProductsJSON, err)
		return
	}

	err = json.Unmarshal(content, products)
	if err != nil {
		log.Printf("unmarshal json to struct error: %v\n", err)
		return
	}

	return
}
