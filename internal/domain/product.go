package domain

type Product struct {
	Id          int     `json:"id" example:"1"`
	Name        string  `json:"name" example:"Pineapple" binding:"required"`
	Quantity    int     `json:"quantity" example:"100" binding:"required"`
	CodeValue   string  `json:"code_value" example:"COD123" binding:"required"`
	IsPublished bool    `json:"is_published" example:"true"`
	Expiration  string  `json:"expiration" example:"25/08/2030" binding:"required"`
	Price       float64 `json:"price" example:"299" binding:"required" format:"float64"`
}

type ProductRequest struct {
	Name        string  `json:"name,omitempty" example:"Pineapple"`
	Quantity    int     `json:"quantity,omitempty" example:"100"`
	CodeValue   string  `json:"code_value,omitempty" example:"COD123"`
	IsPublished bool    `json:"is_published,omitempty" example:"true"`
	Expiration  string  `json:"expiration,omitempty" example:"25/08/2030"`
	Price       float64 `json:"price,omitempty" example:"299" format:"float64"`
}
