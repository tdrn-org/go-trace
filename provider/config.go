//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package provider

import (
	_ "embed" // to embed default config

	"github.com/tdrn-org/go-config-toml"
)

//go:embed defaults.toml
var defaultsData []byte

// DefaultConfig gets the default tracing configuration.
func DefaultConfig() (*TracingConfig, error) {
	return config.Defaults(&TracingConfig{}, defaultsData)
}

// LoadConfig loads a tracing configuration from the
// given file. In case a configuration attribute is not set,
// the default configuration is used.
func LoadConfig(path string) (*TracingConfig, error) {
	return config.Load(&TracingConfig{}, path, defaultsData, false)
}

// TracingConfig represents a tracing configuration.
type TracingConfig struct {
	// Enabled determines whether tracing is enabled or not.
	Enabled bool `toml:"enabled"`
	// ServiceName defines the service name used for tracing.
	ServiceName string `toml:"service_name"`
	// Protocol defines the tracing protocol to use.
	Protocol TracingProtocol `toml:"protocol"`
	// EndpointURL defines the URL of the tracing endpoint to use.
	EndpointURL config.URLSpec `toml:"endpoint_url"`
	// BatchTimeout defines the batch timeout value to use.
	BatchTimeout config.DurationSpec `toml:"batch_timeout"`
	// ExportTimeout defines the export timeout value to use.
	ExportTimeout config.DurationSpec `toml:"export_timeout"`
	// Headers defines additional headers to be sent to the tracing endpoint.
	Headers map[string]string `toml:"headers"`
}

// TracingProtocol type.
type TracingProtocol string

const (
	// TracingProtocolHttp sets the tracing protocol to HTTP.
	TracingProtocolHttp TracingProtocol = "http"
	// TracingProtocolGRPC sets the tracing protocol to gRPC.
	TracingProtocolGRPC TracingProtocol = "gRPC"
)

var tracingProtocolMarshalMap map[TracingProtocol]string = map[TracingProtocol]string{
	TracingProtocolHttp: string(TracingProtocolHttp),
	TracingProtocolGRPC: string(TracingProtocolGRPC),
}

var tracingProtocolUnmarshalMap map[string]TracingProtocol = map[string]TracingProtocol{
	string(TracingProtocolHttp): TracingProtocolHttp,
	string(TracingProtocolGRPC): TracingProtocolGRPC,
}

func (p TracingProtocol) String() string {
	return string(p)
}

func (p TracingProtocol) MarshalText() ([]byte, error) {
	return config.MarshalEnum(p, tracingProtocolMarshalMap)
}

func (p *TracingProtocol) UnmarshalText(text []byte) error {
	TracingProtocol, err := config.UnmarshalEnum(tracingProtocolUnmarshalMap, text)
	if err != nil {
		return err
	}
	*p = TracingProtocol
	return nil
}
