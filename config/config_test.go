//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-trace/config"
)

func TestLoadConfig(t *testing.T) {
	_, err := config.Load("testdata/tracing.toml")
	require.NoError(t, err)
}
