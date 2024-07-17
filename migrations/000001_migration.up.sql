CREATE TABLE products (
    id UUID PRIMARY KEY default gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT, --+
    price DECIMAL(10, 2) NOT NULL,
    category_id UUID REFERENCES product_categories(id),
    artisan_id UUID, 
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
    
CREATE TABLE orders ( --+
    id UUID PRIMARY KEY default gen_random_uuid(),
    user_id UUID,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    shipping_address JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP        
);

CREATE TABLE order_items (--+
    id UUID PRIMARY KEY default gen_random_uuid(),
    order_id UUID REFERENCES orders(id),
    product_id UUID REFERENCES products(id),
    total_amount float NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE product_categories (--+
    id UUID PRIMARY KEY default gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ratings (  
    id UUID PRIMARY KEY default gen_random_uuid(),
    product_id UUID REFERENCES products(id),
    user_id UUID,
    rating DECIMAL(2, 1) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id UUID PRIMARY KEY default gen_random_uuid(),
    order_id UUID REFERENCES orders(id),
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    transaction_id serial,
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


