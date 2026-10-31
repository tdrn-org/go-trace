//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package config

import (
	_ "embed"

	"github.com/tdrn-org/go-config-toml"
)

//go:embed defaults.toml
var defaultsData []byte

func Default() (*TracingConfig, error) {
	return config.Defaults(&TracingConfig{}, defaultsData)
}

func Load(path string) (*TracingConfig, error) {
	return config.Load(&TracingConfig{}, path, defaultsData, false)
}

type TracingConfig struct {
	Enabled       bool                `toml:"enabled"`
	ServiceName   string              `toml:"service_name"`
	Protocol      TracingProtocol     `toml:"protocol"`
	EndpointURL   config.URLSpec      `toml:"endpoint_url"`
	BatchTimeout  config.DurationSpec `toml:"batch_timeout"`
	ExportTimeout config.DurationSpec `toml:"export_timeout"`
	Headers       map[string]string   `toml:"headers"`
}

type TracingProtocol string

const (
	TracingProtocolHttp TracingProtocol = "http"
	TracingProtocolGRPC TracingProtocol = "gRPC"
)

var TracingProtocolMarshalMap map[TracingProtocol]string = map[TracingProtocol]string{
	TracingProtocolHttp: string(TracingProtocolHttp),
	TracingProtocolGRPC: string(TracingProtocolGRPC),
}

var TracingProtocolUnmarshalMap map[string]TracingProtocol = map[string]TracingProtocol{
	string(TracingProtocolHttp): TracingProtocolHttp,
	string(TracingProtocolGRPC): TracingProtocolGRPC,
}

func (l TracingProtocol) String() string {
	return TracingProtocolMarshalMap[l]
}

func (l TracingProtocol) MarshalText() ([]byte, error) {
	return config.MarshalEnum(l, TracingProtocolMarshalMap)
}

func (l *TracingProtocol) UnmarshalText(text []byte) error {
	TracingProtocol, err := config.UnmarshalEnum(TracingProtocolUnmarshalMap, text)
	if err != nil {
		return err
	}
	*l = TracingProtocol
	return nil
}
