package adaptor

import (
	"io"
	"log"
	proto "mc_reverse_proxy/src/control/controlProto"
	service "mc_reverse_proxy/src/control/service"
	"net"

	"google.golang.org/grpc"
)

type GRPCControlCenter struct {
	EventService *service.EventService

	proto.UnimplementedMetricServiceServer
	proto.UnimplementedCommandServiceServer

	address string
}

func (s *GRPCControlCenter) Serve() error {
	server := grpc.NewServer()

	proto.RegisterMetricServiceServer(server, s)
	proto.RegisterCommandServiceServer(server, s)

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		// log.Fatalf("Failed to listen on %s: %v", s.address, err)
		return err

	}

	log.Printf("[GRPC control] Start server on %s", s.address)
	if err := server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *GRPCControlCenter) Metric(stream proto.MetricService_MetricServer) error {
	for {
		req, err := stream.Recv()
		log.Printf(req.ServerID)
		if err == io.EOF {
			// End of stream, send response
			return stream.SendAndClose(&proto.Placeholder{})
		}
		if err != nil {
			log.Panicf(err.Error())
			return err
		}

		s.EventService.Publish("metric", service.EventData{MetricData: req})
	}
}

func (s *GRPCControlCenter) Command(req *proto.Placeholder, stream proto.CommandService_CommandServer) error {
	channel := s.EventService.Subscribe("command")
	for {
		select {
		case data := <-channel:
			log.Printf("%v", data)
			// switch data.CommandData.Command {
			// case proto.CommandEnum_TIMESET:
			// 	stream.Send(&proto.CommandData{Command: data.CommandData.Command, TimesetData: data.CommandData.TimesetData})
			// }
			if data.CommandData.TimesetData != nil {
				stream.Send(&proto.CommandData{Command: data.CommandData.Command, TimesetData: data.CommandData.TimesetData})
			}
		}
	}
}

func NewGRPCControlCenter(address string, eventService *service.EventService) *GRPCControlCenter {
	return &GRPCControlCenter{
		address:      address,
		EventService: eventService,
	}
}
