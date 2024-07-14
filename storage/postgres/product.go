package postgres

import (
	"database/sql"
	pb "product_service/genproto/product"
)

type ProdutRepo struct {
	Db *sql.DB
}

func NewProdutRepo(db *sql.DB) *ProdutRepo {
	return &ProdutRepo{Db: db}
} 

func (p *ProdutRepo) CreateProduct(req *pb.CreateOrderRequest) () {

}