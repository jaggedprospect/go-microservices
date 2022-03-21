package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure for an API product.
type Product struct {
	ID          int     `json:"id"`                       // 'ID' will appear as 'id' in JSON output
	Name        string  `json:"name" validate:"required"` // added a validate tag
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"` // 'CreatedOn' will be excluded
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
} // ^^ using JSON struct tags to format output of data

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) Validate() error {
	// create new validator
	validate := validator.New()

	// register custom validation function
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

// validateSKU is a custom validation function used to check the
// SKU field within a Product.
func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-abcd-abcde
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)

	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Marshal
// as it does not have to buffer the output into an in memory
// slice of bytes -- this reduces allocations and overhead of
// the service
func (p *Products) ToJSON(w io.Writer) error {
	// create encoder
	e := json.NewEncoder(w)
	// return encoded caller data
	return e.Encode(p)
}

// GetProducts returns a list of Products
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd343",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
