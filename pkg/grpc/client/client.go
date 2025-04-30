package client

import (
	"fmt"
	"grpc-gateway/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func New(clientGRPC *config.ClientGRPC) (*grpc.ClientConn, error) {
	const op = "grpc_client.New"

	if clientGRPC.MaxMsgSize < 4 {
		clientGRPC.MaxMsgSize = 4
	}
	if clientGRPC.MaxMsgSize > 32 {
		clientGRPC.MaxMsgSize = 32
	}

	switch clientGRPC.NegotiationType {
	case "plaintext":
		var dialOptions []grpc.DialOption
		dialOptions = append(dialOptions,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		dialOptions = append(dialOptions,
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(clientGRPC.MaxMsgSize*1024*1024),
				grpc.MaxCallSendMsgSize(clientGRPC.MaxMsgSize*1024*1024),
			),
		)

		conn, err := grpc.NewClient(clientGRPC.Connect, dialOptions...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return conn, nil

	case "tls":
		credential, err := credentials.NewClientTLSFromFile(clientGRPC.Cert, "")
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var dialOptions []grpc.DialOption
		dialOptions = append(dialOptions,
			grpc.WithTransportCredentials(credential),
		)
		dialOptions = append(dialOptions,
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(clientGRPC.MaxMsgSize*1024*1024),
				grpc.MaxCallSendMsgSize(clientGRPC.MaxMsgSize*1024*1024),
			),
		)

		conn, err := grpc.NewClient(clientGRPC.Connect, dialOptions...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return conn, nil

	default:
		return nil, fmt.Errorf("%s: not supported negotiation-type: %s", op, clientGRPC.NegotiationType)
	}
}
