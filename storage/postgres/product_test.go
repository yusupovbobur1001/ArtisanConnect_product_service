package postgres

import (
	"fmt"
	pb "product_service/genproto/products"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)

	req := pb.CreateProductRequest{
		Name:        "olma",
		Description: "alfaj",
		Price:       28.52,
		CategoryId:  "c6ff9cf9-63fb-42e4-b3a8-227b4aa51e9c",
		Quantity:    20,
		ArtisanId:   "880e8400-e29b-41d4-a716-446655440001",
	}
	fmt.Println(1111)
	resp, err := product.CreateProduct(&req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", resp)
}

func TestUpdateProduct(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)

	p := pb.UpdateProductRequest{
		Name:  "banana",
		Price: 55.2,
		Id:    "6bf19dd3-08b1-4b2b-bdb2-7680474b7f67",
	}

	res, err := product.UpdateProduct(&p)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", res)
}

func TestDeleteProduct(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)

	req := pb.Id{
		ProductId: "6bf19dd3-08b1-4b2b-bdb2-7680474b7f67",
	}

	resp, err := product.DeleteProduct(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestListProducts(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)

	req := pb.ListProductsRequest{
		Page:  2,
		Limit: 2,
	}

	resp, err := product.ListProducts(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestGetProduct(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)
	req := pb.Id{
		ProductId: "a69f363b-d9da-4e6e-a19d-1fa70073f0e7",
	}
	reap, err := product.GetByIdProduct(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(reap)
}

func TestSearchProducts(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)

	req := pb.SearchProductsRequest{
		Query:    "olma",
		Category: "da4b6097-8460-4996-ba06-99217c331cb4",
		MinPrice: 1,
		MaxPrice: 700,
		Page:     1,
		Limit:    5,
	}

	resp, err := product.SearchProduct(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestAddProductRating(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	product := NewProdutRepo(db)

	req := pb.AddProductRatingRequest{
		ProductId: "a82de4d2-2278-46fd-bb2f-bf65486db4e4",
		UserId:    "e72e99af-4f0f-4eca-b3d8-a5891afd1e98",
		Rating:    5,
		Comment:   "nima gap nima nima gap",
	}
	resp, err := product.AddProductRating(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func TestGetProductRatings(t *testing.T) {
	db, err := ConnDB()
	if err != nil {
		panic(err)
	}

	product := NewProdutRepo(db)

	req := pb.Id{
		ProductId: "f7999ed7-58e5-4697-b91c-2cddc6e1f1e1",
	}
	resp, err := product.GetProductRatings(&req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
