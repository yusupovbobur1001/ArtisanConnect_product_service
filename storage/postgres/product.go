package postgres

import (
	"database/sql"
	"fmt"
	pb "product_service/genproto/products"
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
				returning id, name, description, price, category_id, quantity, artisan_id, created_at `
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
			name = $1, price = $2, updated_at = $4
		where 
			id = $3 and deleted_at is null
		 `
	_, err := p.Db.Exec(query, req.Name, req.Price, req.Id, time.Now())
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

	err := p.Db.QueryRow(query, req.ProductId).Scan(
		&resp.Id,
		&resp.Name,
		&resp.Description,
		&resp.Price,
		&resp.Quantity,
		&resp.ArtisanId,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
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
	// resp := &pb.ListProductsResponse{}

	m, err := p.GetAllProduct()
	if err != nil {
		return nil, err
	}
	total := len(m)

	totalPage := (total + int(req.Limit) - 1) / int(req.Limit)
	if totalPage < int(req.Page) {
		return nil, fmt.Errorf("totalPage dan page katta bo`lib qoldi!, err: %v", err)
	}

	startRow := (int(req.Page)-1)*int(req.Limit) + 1
	endRow := startRow + int(req.Limit) - 1
	if endRow > total {
		endRow = total
	}

	b := []*pb.ProductInfo{}

	q := `
		WITH NumberedRows AS (
			SELECT 
				id, name, description, price, category_id, quantity, artisan_id, created_at, updated_at,
				ROW_NUMBER() OVER (ORDER BY id) AS row_num
			FROM 
				products
		)
		SELECT 
			id, name, description, price, category_id, quantity, artisan_id, created_at, updated_at
		FROM 
			NumberedRows
		WHERE 
			row_num BETWEEN $1 AND $2
	`

	rows, err := p.Db.Query(q, startRow, endRow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := &pb.ProductInfo{}
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CategoryId,
			&product.Quantity,
			&product.ArtisanId,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		b = append(b, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.ListProductsResponse{
		Products: b,
		Total:    int32(total),
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

func (p *ProdutRepo) GetAllProduct() ([]*pb.ProductInfo, error) {
	resps := []*pb.ProductInfo{}

	q := `SELECT 
				id, name, description, price, quantity, updated_at, created_at, category_id, artisan_id
		  FROM 
		  		products
		  WHERE 
		  		deleted_at IS NULL`

	rows, err := p.Db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resp := &pb.ProductInfo{}
		err = rows.Scan(&resp.Id, &resp.Name, &resp.Description, &resp.Price,
			&resp.Quantity, &resp.UpdatedAt, &resp.CreatedAt, &resp.CategoryId, &resp.ArtisanId)
		if err != nil {
			return nil, err
		}

		resps = append(resps, resp)
	}

	return resps, nil
}

func (p *ProdutRepo) SearchProduct(req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	b, err := p.GetAllProduct()
	if err != nil {
		fmt.Println(err, "+++++++++++++")
		return nil, err
	}
	total := len(b)

	totalPage := (total + int(req.Limit) - 1) / int(req.Limit)
	if totalPage < int(req.Page) {
		return nil, fmt.Errorf("totalPage dan page katta bo`lib qoldi!, err: %v", err)
	}

	startRow := (int(req.Page)-1)*int(req.Limit) + 1
	endRow := startRow + int(req.Limit) - 1
	if endRow > total {
		endRow = total
	}

	q := `
        WITH NumberedRows AS (
            SELECT 
                id, name, description, price, category_id, quantity, artisan_id, created_at, updated_at,
                ROW_NUMBER() OVER (ORDER BY id) AS row_num
            FROM 
                products
            WHERE 
                price > $3 AND price < $4
        )
        SELECT 
            id, name, description, price, category_id, quantity, artisan_id, created_at, updated_at
        FROM 
            NumberedRows
        WHERE 
            row_num BETWEEN $1 AND $2 
    `

	params := []interface{}{startRow, endRow, req.MinPrice, req.MaxPrice}

	if len(req.Query) > 0 {
		q += " AND name ILIKE $5"
		params = append(params, "%"+req.Query+"%")
	}

	rows, err := p.Db.Query(q, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	searchPs := []*pb.ProductInfo{}

	for rows.Next() {
		searchP := &pb.ProductInfo{}
		err = rows.Scan(
			&searchP.Id,
			&searchP.Name,
			&searchP.Description,
			&searchP.Price,
			&searchP.CategoryId,
			&searchP.Quantity,
			&searchP.ArtisanId,
			&searchP.CreatedAt,
			&searchP.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		searchPs = append(searchPs, searchP)
	}
	resp := pb.SearchProductsResponse{
		Products: searchPs,
		Total:    int32(total),
		Page:     req.Page,
		Limit:    req.Limit,
	}
	return &resp, nil
}

func (p *ProdutRepo) AddProductRating(req *pb.AddProductRatingRequest) (*pb.RatingInfo, error) {
	resp := pb.RatingInfo{}
	q := `
		insert into ratings(
								product_id,
								user_id,
								rating, 
								comment,
								created_at)
		values($1, $2, $3, $4, $5)
		returning id, product_id, user_id, rating, comment, created_at  `

	err := p.Db.QueryRow(q, req.ProductId, req.UserId, req.Rating, req.Comment, time.Now()).Scan(
							&resp.Id, 
							&resp.ProductId, 
							&resp.UserId, 
							&resp.Rating, 
							&resp.Comment, 
							&resp.CreatedAt,
					)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p *ProdutRepo) GetProductRatings(req *pb.Id) (*pb.GetProductRatingsResponse, error) {
	resp := pb.GetProductRatingsResponse{}

	q1 := `
		SELECT
    		id,
    		product_id,
    		user_id,
    		rating,
    		comment,
    		created_at
		FROM
    		ratings
		WHERE
    		product_id = $1
	`
	ratings := []*pb.RatingInfo{}
	rows, err := p.Db.Query(q1, req.ProductId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rating := &pb.RatingInfo{}
		err = rows.Scan(
			&rating.Id,
			&rating.ProductId,
			&rating.UserId,
			&rating.Rating,
			&rating.Comment,
			&rating.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		ratings = append(ratings, rating)
	}

	q2 := `
		SELECT
    		AVG(rating) as average_rating,
    		COUNT(*) as total_ratings
		FROM
    		ratings
		WHERE
    		product_id = $1

	`
	err = p.Db.QueryRow(q2, req.ProductId).Scan(
		&resp.AverageRating,
		&resp.TotalRatings,
	)

	if err != nil {
		return nil, err
	}

	return &pb.GetProductRatingsResponse{
		Ratings:       ratings,
		AverageRating: resp.AverageRating,
		TotalRatings:  resp.TotalRatings,
	}, nil
}
