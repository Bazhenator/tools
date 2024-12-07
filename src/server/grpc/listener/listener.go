package listener

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	EnvGrpcHost = "GRPC_HOST"
	EnvGrpcPort = "GRPC_PORT"

	EnvGrpcTimeout = "CONTEXT_TIMEOUT"
	DefGrpcTimeout = 30000

	EnvMaxGrpcReceive = "MAX_GRPC_RECEIVE"
	DefMaxGrpcReceive = 1024 * 1024 * 1024 * 1

	EnvMaxGrpcSend = "MAX_GRPC_SEND"
	DefMaxGrpcSend = 1024 * 1024 * 1024 * 1
)

// GrpcConfig is a grpc listener configuration
type GrpcConfig struct {
	Host           string
	Port           string
	Timeout        int
	MaxReceiveSize int
	MaxSendSize    int
}

// NewStandardGrpcConfig creates a new instance of GrpcConfig struct
func NewStandardGrpcConfig() (*GrpcConfig, error) {
	var multiErr error

	grpcHost, ok := os.LookupEnv(EnvGrpcHost)
	if !ok {
		return nil, fmt.Errorf("field ('%s') is not defined", EnvGrpcHost)
	}
	  

	grpcPort, ok := os.LookupEnv(EnvGrpcPort)
	if !ok {
		multierr.AppendInto(&multiErr, fmt.Errorf("field ('%s') is not defined", EnvGrpcPort))
	} else if _, err := strconv.Atoi(grpcPort); err != nil {
		multierr.AppendInto(&multiErr, err)
	}

	var grpcTimeoutInt int
	_, ok = os.LookupEnv(EnvGrpcTimeout)
	if !ok {
		grpcTimeoutInt = DefGrpcTimeout
	}

	var maxGrpcReceiveInt int
	_, ok = os.LookupEnv(EnvMaxGrpcReceive)
	if !ok {
		maxGrpcReceiveInt = DefMaxGrpcReceive
	}

	var maxGrpcSendInt int
	_, ok = os.LookupEnv(EnvMaxGrpcSend)
	if !ok {
		maxGrpcSendInt = DefMaxGrpcReceive
	}

	if multiErr != nil {
		return nil, multiErr
	}

	c := &GrpcConfig{
		Host:           grpcHost,
		Port:           grpcPort,
		Timeout:        grpcTimeoutInt,
		MaxReceiveSize: maxGrpcReceiveInt,
		MaxSendSize:    maxGrpcSendInt,
	}

	return c, nil
}

// NewGrpcListener initializes new grpc listener
func NewGrpcListener(c *GrpcConfig) (net.Listener, func(net.Listener), error) {
	sl := zap.S().With("entity", "grpc")

	if err := c.checkEmptyFields(); err != nil {
		sl.Error("error while checking config", zap.Error(err))
		return nil, nil, err
	}

	sl.Infof("server has started at %s:%s", c.Host, c.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to strat listener: %w", err)
	}

	deferFunc := func(net.Listener) {
		if err = lis.Close(); err != nil {
			sl.Error("error while closing listner", zap.Error(err))
		}
	}

	return lis, deferFunc, nil
}

// checkEmptyFields checks whether necessary fields is empty and return error though
func (c GrpcConfig) checkEmptyFields() error {
	var err error

	if c.Port == "" {
		multierr.AppendInto(&err,  fmt.Errorf("field ('%s') is not defined", "port"))
	}
	if c.Timeout == 0 {
		multierr.AppendInto(&err,  fmt.Errorf("field ('%s') is not defined", "timeout"))
	}
	if c.MaxReceiveSize == 0 {
		multierr.AppendInto(&err,  fmt.Errorf("field ('%s') is not defined", "maxReceiveSize"))
	}
	if c.MaxSendSize == 0 {
		multierr.AppendInto(&err,  fmt.Errorf("field ('%s') is not defined", "maxSendSize"))
	}

	return err
}