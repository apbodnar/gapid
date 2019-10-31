// Copyright (C) 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package replay

import (
	"context"

	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/os/device/bind"
	"github.com/google/gapid/gapir"
	"github.com/google/gapid/gapis/api"
	"github.com/google/gapid/gapis/api/transform"
	"github.com/google/gapid/gapis/replay/builder"
	"github.com/google/gapid/gapis/service"
)

type WaitForFence struct{
	TraceOptions *service.TraceOptions
}

func (t *WaitForFence) Transform(ctx context.Context, id api.CmdID, cmd api.Cmd, out transform.Writer) {
	if id == 0 {
		t.AddStartPerfetto(ctx, id, out)
	}
	out.MutateAndWrite(ctx, id, cmd)
}

func (t *WaitForFence) Flush(ctx context.Context, out transform.Writer) {
	t.AddStopPerfetto(ctx, out)
}

func (t *WaitForFence) PreLoop(ctx context.Context, out transform.Writer)  {}
func (t *WaitForFence) PostLoop(ctx context.Context, out transform.Writer) {}

func (t *WaitForFence) AddStartPerfetto(ctx context.Context, id api.CmdID, out transform.Writer) {
	out.MutateAndWrite(ctx, id, Custom{T: 0, F: func(ctx context.Context, s *api.GlobalState, b *builder.Builder) error {
		fenceID := uint32(id)
		b.Wait(fenceID)
		callback := func(p *gapir.FenceReadyRequest, device bind.Device) {
			//TODO(apbodnar) Start perfetto here
		}
		return b.RegisterFenceReadyRequestCallback(fenceID, callback)
	}})
}

func (t *WaitForFence) AddStopPerfetto(ctx context.Context, out transform.Writer) {
	out.MutateAndWrite(ctx, api.CmdNoID, Custom{T: 0, F: func(ctx context.Context, s *api.GlobalState, b *builder.Builder) error {
		fenceID := uint32(0xdeaaad)
		b.Wait(fenceID)
		callback := func(p *gapir.FenceReadyRequest, device bind.Device) {
			//TODO(apbodnar) Stop perfetto here
		}
		return b.RegisterFenceReadyRequestCallback(fenceID, callback)
	}})
}
