package postgres

import "database/sql"

type ProdutRepo struct {
	Db *sql.DB
}

func NewProdutRepo(db *sql.DB) *ProdutRepo {
	return &ProdutRepo{Db: db}
} 

func (p *ProdutRepo) CreateProduct() () {
	
}