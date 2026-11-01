//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package provider_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-trace/provider"
)

func TestLoadConfig(t *testing.T) {
	_, err := provider.LoadConfig("testdata/tracing.toml")
	require.NoError(t, err)
}
