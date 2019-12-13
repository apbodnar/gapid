// Copyright (C) 2017 Google Inc.
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

package vulkan

import (
  "context"

  "github.com/google/gapid/core/log"
  "github.com/google/gapid/gapis/api"
  "github.com/google/gapid/gapis/api/transform"
)

var renderStagesLayer = "VkRenderStagesProducer"

type ProfilingLayers struct {
}

func (t *ProfilingLayers) Transform(ctx context.Context, id api.CmdID, cmd api.Cmd, out transform.Writer) {
  ctx = log.Enter(ctx, "ProfilingLayers")

  s := out.State()
  cb := CommandBuilder{Thread: cmd.Thread(), Arena: out.State().Arena}
  l := s.MemoryLayout
  allocated := []api.AllocResult{}
  defer func() {
    for _, d := range allocated {
      d.Free()
    }
  }()
  mustAlloc := func(ctx context.Context, v ...interface{}) api.AllocResult {
    res := s.AllocDataOrPanic(ctx, v...)
    allocated = append(allocated, res)
    return res
  }

  switch cmd := cmd.(type) {
  case *VkCreateInstance:
    cmd.Extras().Observations().ApplyReads(s.Memory.ApplicationPool())
    info := cmd.PCreateInfo().MustRead(ctx, cmd, s, nil)
    layerCount := info.EnabledLayerCount()
    layers := []Charᶜᵖ{}
    hasProfilingLayer := false
    for _, l := range info.PpEnabledLayerNames().Slice(0, uint64(layerCount), l).MustRead(ctx, cmd, s, nil) {
      if renderStagesLayer == l.String() {
        hasProfilingLayer = true
      }
    }

    if !hasProfilingLayer {
      d := mustAlloc(ctx, renderStagesLayer)
      layers = append(layers, NewCharᶜᵖ(d.Ptr()))
    }

    layersData := mustAlloc(ctx, layers)
    extCount := info.EnabledExtensionCount()
    exts := info.PpEnabledExtensionNames().Slice(0, uint64(extCount), l).MustRead(ctx, cmd, s, nil)
    extsData := mustAlloc(ctx, exts)

    info.SetEnabledLayerCount(uint32(len(layers)))
    info.SetPpEnabledLayerNames(NewCharᶜᵖᶜᵖ(layersData.Ptr()))
    info.SetEnabledExtensionCount(uint32(len(exts)))
    info.SetPpEnabledExtensionNames(NewCharᶜᵖᶜᵖ(extsData.Ptr()))
    infoData := mustAlloc(ctx, info)

    newCmd := cb.VkCreateInstance(infoData.Ptr(), cmd.PAllocator(), cmd.PInstance(), cmd.Result())
    newCmd.AddRead(
      infoData.Data(),
    ).AddRead(
      layersData.Data(),
    )
    // Also add back all the other read/write observations of the original vkCreateInstance
    for _, r := range cmd.Extras().Observations().Reads {
      newCmd.AddRead(r.Range, r.ID)
    }
    for _, w := range cmd.Extras().Observations().Writes {
      newCmd.AddWrite(w.Range, w.ID)
    }
    out.MutateAndWrite(ctx, id, newCmd)

  default:
    out.MutateAndWrite(ctx, id, cmd)

  }
}

func (t *ProfilingLayers) PreLoop(ctx context.Context, out transform.Writer)  {}
func (t *ProfilingLayers) PostLoop(ctx context.Context, out transform.Writer) {}
func (t *ProfilingLayers) Flush(ctx context.Context, out transform.Writer) {}
