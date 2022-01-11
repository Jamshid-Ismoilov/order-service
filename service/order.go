package service

import (
	"context"
	"time"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Jamshid-Ismoilov/order-service/genproto/order_service"
	l "github.com/Jamshid-Ismoilov/order-service/pkg/logger"
	"github.com/Jamshid-Ismoilov/order-service/storage"
	"github.com/gofrs/uuid"
)

// OrderService is an object that implements user interface.
type OrderService struct {
	storage storage.IStorage
	logger  l.Logger
}

// NewOrderService ...
func NewOrderService(storage storage.IStorage, log l.Logger) *OrderService {
	return &OrderService{
		storage: storage,
		logger:  log,
	}
}

func (s *OrderService) Create(ctx context.Context, req *pb.Order) (*pb.Book, error) {
	
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	if req.GetCreatedAt() == "" {
		req.CreatedAt = time.Now().String()
	}
	book, err := s.storage.Order().Create(*req)
	if err != nil {
		s.logger.Error("failed to create order", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	return &book, nil
}

func (s *OrderService) Get(ctx context.Context, req *pb.ByIdReq) (*pb.Order, error) {
	order, err := s.storage.Order().Get(req.GetId())
	if err != nil {
		s.logger.Error("failed to get order", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get order")
	}

	return &order, nil
}

func (s *OrderService) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	orders, count, err := s.storage.Order().List(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list tasks", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list tasks")
	}

	return &pb.ListResp{
		Orders: orders,
		Count: count,
	}, nil
}

func (s *OrderService) Update(ctx context.Context, req *pb.Order) (*pb.Order, error) {
	order, err := s.storage.Order().Update(*req)
	if err != nil {
		s.logger.Error("failed to update order", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update order")
	}

	return &order, nil
}

func (s *OrderService) Delete(ctx context.Context, req *pb.ByIdReq) (*pb.EmptyRes, error) {
	err := s.storage.Order().Delete(req.Id)
	if err != nil {
		s.logger.Error("failed to delete order", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete order")
	}

	return &pb.EmptyRes{}, nil
}
