package service

import (
	"context"
	"database/sql"
	"product_service/config"
	user "product_service/genproto/auth"
	pbu "product_service/genproto/products"
	"product_service/pkg"
	"product_service/storage/postgres"
)

type ProductService struct {
	pbu.UnimplementedProductServiceServer
	Auth    user.AuthClient
	Product *postgres.ProdutRepo
}

func NewProductService(db *sql.DB, cfg *config.Config) *ProductService {
	return &ProductService{
		Product: postgres.NewProdutRepo(db),
		Auth: pkg.CreateAuthClient(cfg),
	}
}

func (h *ProductService) CreateProduct(ctx context.Context, req *pbu.CreateProductRequest) (*pbu.ProductInfo, error) {
	exist, err := h.Auth.ValidateUserId(ctx, &user.Id{Id: req.ArtisanId})
	if !exist.Exist || err != nil {
		return nil, err
	}
	
	resp, err := h.Product.CreateProduct(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *ProductService) UpdateProduct(ctx context.Context, req *pbu.UpdateProductRequest) (*pbu.ProductInfo, error) {
	resp, err := h.Product.UpdateProduct(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *ProductService) DeleteProduct(ctx context.Context, req *pbu.Id) (*pbu.DeleteProductResponse, error) {
	resp, err := h.Product.DeleteProduct(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *ProductService) ListProducts(ctx context.Context, req *pbu.ListProductsRequest) (*pbu.ListProductsResponse, error) {
	resp, err := h.Product.ListProducts(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *ProductService) GetProduct(ctx context.Context, req *pbu.Id) (*pbu.ProductInfo, error) {
	resp, err := h.Product.GetByIdProduct(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *ProductService) SearchProducts(ctx context.Context, req *pbu.SearchProductsRequest) (*pbu.SearchProductsResponse, error) {
	resp, err := h.Product.SearchProduct(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *ProductService) AddProductRating(ctx context.Context, req *pbu.AddProductRatingRequest) (*pbu.RatingInfo, error) {
	exist, err := h.Auth.ValidateUserId(ctx, &user.Id{Id: req.UserId})
	if !exist.Exist || err != nil {
		return nil, err
	}

	resp, err := h.Product.AddProductRating(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (h *ProductService) GetProductRatings(ctx context.Context, req *pbu.Id) (*pbu.GetProductRatingsResponse, error) {
	resp, err := h.Product.GetProductRatings(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}














