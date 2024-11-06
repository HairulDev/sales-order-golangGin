package models

type ItemOrder struct {
	Id_Item   string  `json:"id_item"`
	Id_Order  string  `json:"id_order"`
	Item_Name string  `json:"item_name"`
	Qty      int     `json:"qty"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total"`
}
