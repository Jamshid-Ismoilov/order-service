CREATE TABLE IF NOT EXISTS orders(
    id varchar(256),
    book_id varchar(64),
    description varchar(64),
    created_at timestamp default current_timestamp,
    updated_at timestamp default null,
    deleted_at timestamp default null
);