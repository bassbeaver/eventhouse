// Taken from https://github.com/opentracing-contrib/go-grpc

package opentracing

import (
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	GrpcMetadata = opentracing.HTTPHeaders
)

type GrpcMetadataCarrier metadata.MD

func (c GrpcMetadataCarrier) Set(key, val string) {
	// The GRPC HPACK implementation rejects any uppercase keys here.
	//
	// As such, since the HTTP_HEADERS format is case-insensitive anyway, we
	// blindly lowercase the key (which is guaranteed to work in the
	// Inject/Extract sense per the OpenTracing spec).
	key = strings.ToLower(key)
	c[key] = append(c[key], val)
}

func (c GrpcMetadataCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range c {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
