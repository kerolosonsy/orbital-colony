package main

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var ErrSectorMissing = errors.New("sector missing")

type PressureLookup func(context.Context, string) (int64, error)
type inspectionService struct {
	UnimplementedBerthControlServer
	lookup PressureLookup
}

func sectorName(sector string) bool {
	if len(sector) < 1 || len(sector) > 16 || sector[0] < 'a' || sector[0] > 'z' {
		return false
	}
	for _, b := range []byte(sector) {
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-') {
			return false
		}
	}
	return true
}
func (s inspectionService) Inspect(ctx context.Context, q *BerthQuery) (*BerthReading, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	if !sectorName(q.GetSector()) {
		return nil, status.Error(codes.InvalidArgument, "invalid sector")
	}
	pressure, err := s.lookup(ctx, q.GetSector())
	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}
	if errors.Is(err, ErrSectorMissing) {
		return nil, status.Error(codes.NotFound, "sector missing")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "inspection unavailable")
	}
	return &BerthReading{Sector: q.GetSector(), Pressure: pressure}, nil
}
func inspectionGuard(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (reply any, err error) {
		defer func() {
			if recover() != nil {
				reply = nil
				err = status.Error(codes.Internal, "inspection unavailable")
			}
		}()
		if ctx.Err() != nil {
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		md, _ := metadata.FromIncomingContext(ctx)
		keys := md.Get("x-colony-key")
		if len(keys) != 1 || keys[0] != key {
			return nil, status.Error(codes.Unauthenticated, "access denied")
		}
		return next(ctx, req)
	}
}
func NewInspectionServer(key string, lookup PressureLookup) *grpc.Server {
	srv := grpc.NewServer(grpc.UnaryInterceptor(inspectionGuard(key)))
	RegisterBerthControlServer(srv, inspectionService{lookup: lookup})
	return srv
}
func main() {}
