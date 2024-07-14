package postgres

import (
	"database/sql"
	"fmt"
	pb "product_service/genproto/proto"
	"time"
)

type ProdutRepo struct {
	Db *sql.DB
}

func NewProdutRepo(db *sql.DB) *ProdutRepo {
	return &ProdutRepo{Db: db}
}

func (p *ProdutRepo) CreateProduct(req *pb.CreateProductRequest) (*pb.ProductInfo, error) {
	resp := pb.ProductInfo{}
	query := `insert into products(
						name, 
						description,  
						price, 
						category_id, 
						quantity, 
						artisan_id,
						created_at)
				values($1, $2, $3, $4, $5, $6, $7)  
				returing id, name, description, price, quantity, artisan_id, created_at `
	err := p.Db.QueryRow(query, req.Name, req.Description,
		req.Price, req.CategoryId, req.Quantity, req.ArtisanId, time.Now()).Scan(
		&resp.Id, &req.Name, &resp.Description, &resp.Price, &resp.CategoryId, &resp.Quantity, &resp.ArtisanId, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (p *ProdutRepo) UpdateProduct(req *pb.UpdateProductRequest) (*pb.ProductInfo, error) {
	query := `
		update 
			products 
		set 
			name = $1, price = $2, product_id = $3, updated_at = $5
		where 
			id = $4 and deleted_at is null `
	_, err := p.Db.Exec(query, req.Name, req.Price, req.ProductId, req.Id, time.Now())
	if err != nil {
		return nil, err
	}

	resp, err := p.GetByIdProduct(&pb.Id{ProductId: req.Id})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *ProdutRepo) GetByIdProduct(req *pb.Id) (*pb.ProductInfo, error) {
	resp := pb.ProductInfo{}
	query := `
		select 
			  id, name, description, price, quantity, artisan_id, created_at, updated_at 
		from 
			  products
		where 
			  id = $1 and deleted_at is null `

	err := p.Db.QueryRow(query, req.ProductId).Scan(&resp.Id, &resp.Name,
		&resp.Description, &resp.Price, &resp.CategoryId, &resp.Quantity,
		&resp.ArtisanId, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (p *ProdutRepo) DeleteProduct(req *pb.Id) (*pb.DeleteProductResponse, error) {
	q := `
		update
			products
		set
			deleted_at = $1
		where 
			id = $2`
	_, err := p.Db.Exec(q, time.Now(), req.ProductId)
	if err != nil {
		return &pb.DeleteProductResponse{Message: "product id was not deleted"}, nil
	}

	return &pb.DeleteProductResponse{Message: "product id deleted"}, nil
}




func (p *ProdutRepo) ListProducts(req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	resp := pb.ListProductsResponse{}

	m, err := p.GetAllProduct()
	if err != nil {
		return nil, err
	}
	total := len(m)
	totalPage := total/int(req.Limit)
	if totalPage < int(req.Page) {
		return nil, fmt.Errorf("totalPage dan page katta bo`lib qoldi!")
	}
	start_row := 1
	end_row := 0
	for i := 1; i < int(req.Page); i++ {
		if int(req.Page) > 1 {
			start_row += int(req.Limit)
			end_row += start_row + int(req.Limit)
		} else {
			start_row = 1
			end_row = int(req.Limit)
		}
	}

	b := []pb.ProductInfo{}

	q := `WITH NumberedRows AS (
    	SELECT 
			*, ROW_NUMBER() OVER (ORDER BY id) AS row_num
   	 	FROM 
			users
		)
		SELECT 
			* 
		FROM 
			NumberedRows
		WHERE 
			row_num BETWEEN $1 AND $2 `

	rows, err := p.Db.Query(q, start_row, end_row).Scan()
	
}

func (p *ProdutRepo) GetAllProduct() ([]pb.ProductInfo, error) {
	resps := []pb.ProductInfo{}

	q := `select 
				id, name, description, price, quantity, updated_at, created_at, category_id, artisan_id
		  from 
		  		products
		  where 
		  		deleted_at is null`

	rows, err := p.Db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		resp := pb.ProductInfo{}
		err = rows.Scan(&resp.Id, &resp.Name, &resp.Description, &resp.Price,
			&resp.Quantity, &resp.UpdatedAt, &resp.CreatedAt, &resp.CategoryId, &resp.ArtisanId)
		if err != nil {
			return nil, err
		}

		resps = append(resps, resp)
	}

	return resps, nil
}

