package component

import (
	"testing"
)

type SKU struct {
	SkuCode      string `json:"sku_code" table:"label:sku名称;hover:true;click:true;"`
	ProviderName string `json:"provider_name" table:"label:供应商名称"`
	OmitCol      string
}

func TestNewTable(t *testing.T) {
}
