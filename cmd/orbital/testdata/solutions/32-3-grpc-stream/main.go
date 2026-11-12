package main

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DriftAdjust func(context.Context, int64, int64) (int64, error)
type driftService struct {
	UnimplementedDriftControlServer
	adjust DriftAdjust
}

func (s driftService) step(ctx context.Context, total int64, q *DriftStep) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, status.FromContextError(err).Err()
	}
	next, err := s.adjust(ctx, total, q.GetDelta())
	if ctx.Err() != nil {
		return 0, status.FromContextError(ctx.Err()).Err()
	}
	if err != nil {
		return 0, status.Error(codes.Internal, "calibration unavailable")
	}
	return next, nil
}
func (s driftService) Calibrate(stream grpc.BidiStreamingServer[DriftStep, DriftAck]) error {
	var total int64
	for sequence := uint32(1); ; sequence++ {
		q, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if sequence > 64 {
			return status.Error(codes.ResourceExhausted, "step limit")
		}
		if q.GetSequence() != sequence || q.GetDelta() < -1000 || q.GetDelta() > 1000 {
			return status.Error(codes.InvalidArgument, "invalid step")
		}
		total, err = s.step(stream.Context(), total, q)
		if err != nil {
			return err
		}
		if err := stream.Send(&DriftAck{Sequence: sequence, Total: total}); err != nil {
			return err
		}
	}
}
func NewDriftServer(adjust DriftAdjust) *grpc.Server {
	srv := grpc.NewServer()
	RegisterDriftControlServer(srv, driftService{adjust: adjust})
	return srv
}
func main() {}
