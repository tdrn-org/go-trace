//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package provider

import "context"

type noopProvider struct{}

func (p *noopProvider) Shutdown(_ context.Context) error {
	return nil
}
