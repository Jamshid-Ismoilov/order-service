package postgres

import (
	"database/sql"
	"time"
	"fmt"
	"log"
	// // "context"

	// // "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	"github.com/jmoiron/sqlx"

	pb "github.com/Jamshid-Ismoilov/order-service/genproto/order_service"
	// pb2 "github.com/Jamshid-Ismoilov/order-service/genproto/catalog_service"
)

type orderRepo struct {
	db *sqlx.DB
}

// NewOrderRepo ...
func NewOrderRepo(db *sqlx.DB) *orderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order pb.Order) (pb.Book, error) {
	var id string
	err := r.db.QueryRow(`
        INSERT INTO orders(id, book_id, description)
        VALUES ($1, $2, $3) returning id`, order.Id, order.BookId, order.Description).Scan(&id)
	if err != nil {
		return pb.Book{}, err
	}

	// conn, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("Did not connect %v", err)
	// }
	// client := pb2.NewCatalogServiceClient(conn)
	
	// input := pb2.ByIdReq{
	// 	Id: order.BookId,
	// }

	// res, err := client.BookGet(context.Background(), &input)
	// if err != nil {
	// 	log.Printf("failed to get user", err)
	// }

	// if err != nil {
	// 	return pb.Book{}, err
	// }
	// var result pb.Book
	// book := *res
	// result.Id, result.Name, result.CreatedAt, result.UpdateAt = book.Id, book.Name, book.CreatedAt, book.UpdateAt
	return pb.Book{}, nil
}

func (r *orderRepo) Get(id string) (pb.Order, error) {
	var order pb.Order
	var updated sql.NullString

	err := r.db.QueryRow(`
        SELECT id, book_id, description, created_at, updated_at FROM orders
        WHERE id=$1 and deleted_at is null`, id).Scan(
			&order.Id, 
			&order.BookId, 
			&order.Description, 
			&order.CreatedAt, 
			&updated,
		)
	order.UpdatedAt = updated.String

	if err != nil {
		return pb.Order{}, err
	}
	fmt.Println("GET function")

	return order, nil
}

func (r *orderRepo) List(page, limit int64) ([]*pb.Order, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT id, book_id, description, created_at, updated_at FROM orders WHERE deleted_at is null LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		orders []*pb.Order
		order  pb.Order
		count int64
	)
	for rows.Next() {
		var updated sql.NullString
		err = rows.Scan(&order.Id, &order.BookId, &order.Description, &order.CreatedAt, &updated)
		order.UpdatedAt = updated.String
	if err != nil {
			return nil, 0, err
		}
		orders = append(orders, &order)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM orders`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(orders)
	return orders, count, nil
}

func (r *orderRepo) Update(order pb.Order) (pb.Order, error) {
	log.Println(order)
	result, err := r.db.Exec(`UPDATE orders SET book_id=$1, description=$2, updated_at=current_timestamp WHERE id=$3`,
		order.BookId, order.Description, order.Id)
	if err != nil {
		return pb.Order{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Order{}, sql.ErrNoRows
	}

	order, err = r.Get(order.Id)
	if err != nil {
		return pb.Order{}, err
	}

	return order, nil
}



func (r *orderRepo) Delete(id string) error {
	result, err := r.db.Exec(`UPDATE orders SET deleted_at = $2 WHERE id=$1`, id, time.Now())
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

