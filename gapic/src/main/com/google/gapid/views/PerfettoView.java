/*
 * Copyright (C) 2019 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.google.gapid.views;

import static com.google.gapid.widgets.Widgets.createButton;
import static com.google.gapid.widgets.Widgets.createCheckbox;
import static com.google.gapid.widgets.Widgets.createComposite;
import static com.google.gapid.widgets.Widgets.createLabel;
import static com.google.gapid.widgets.Widgets.withLayoutData;

import com.google.gapid.models.Capture;
import com.google.gapid.models.Models;
import com.google.gapid.perfetto.TraceView;
import com.google.gapid.proto.service.path.Path;
import com.google.gapid.util.Loadable.Message;
import com.google.gapid.widgets.Widgets;

import org.eclipse.swt.SWT;
import org.eclipse.swt.layout.GridData;
import org.eclipse.swt.layout.GridLayout;
import org.eclipse.swt.widgets.Button;
import org.eclipse.swt.widgets.Composite;
import org.eclipse.swt.widgets.Control;

public class PerfettoView extends Composite implements Tab, Capture.Listener {
  private final Models models;
  private final Button ovViewport, ovFragmentShader, ovVertexCount, ovSampling, ovTextureSize;
  private final Button goButton;

  public PerfettoView(Composite parent, Models models, Widgets widgets) {
    super(parent, SWT.NONE);
    this.models = models;

    setLayout(new GridLayout(1, false));

    Composite input = createComposite(this, new GridLayout(8, false));

    withLayoutData(createLabel(input, "Overrides:"), new GridData(SWT.LEFT, SWT.CENTER, false, false));
    ovViewport = withLayoutData(createCheckbox(input, "1x1 Viewport:", false), new GridData(SWT.LEFT, SWT.CENTER, false, false));
    ovFragmentShader = withLayoutData(createCheckbox(input, "Constant Fragment Shader:", false), new GridData(SWT.LEFT, SWT.CENTER, false, false));
    ovVertexCount = withLayoutData(createCheckbox(input, "1 Vertex/Draw:", false), new GridData(SWT.LEFT, SWT.CENTER, false, false));
    ovSampling = withLayoutData(createCheckbox(input, "Nearest Sampling:", false), new GridData(SWT.LEFT, SWT.CENTER, false, false));
    ovTextureSize = withLayoutData(createCheckbox(input, "1x1 Textures:", false), new GridData(SWT.LEFT, SWT.CENTER, false, false));
    goButton = withLayoutData(createButton(input, "Go", e -> getReport()),
        new GridData(SWT.LEFT, SWT.CENTER, false, false));
    goButton.setEnabled(false);

    withLayoutData(new TraceView(this, models, widgets, false),
        new GridData(SWT.FILL, SWT.FILL, true, true));

    models.capture.addListener(this);
    addListener(SWT.Dispose, e -> {
      models.capture.removeListener(this);
    });

    if (models.capture.isGraphics()) {
      onCaptureLoaded(null);
    }
  }

  @Override
  public void onCaptureLoadingStart(boolean maintainState) {
    goButton.setEnabled(false);
  }

  @Override
  public void onCaptureLoaded(Message error) {
    if (models.capture.isGraphics()) {
      goButton.setEnabled(true);
    }
  }

  private void getReport() {
    models.perfetto.captureReplay(models.capture.getData().path,
        models.devices.getReplayDevicePath(), Path.OverrideConfig.newBuilder()
            .setViewportSize(ovViewport.getSelection())
            .setFragmentShader(ovFragmentShader.getSelection())
            .setVertexCount(ovVertexCount.getSelection())
            .setSampling(ovSampling.getSelection())
            .setTextureSize(ovTextureSize.getSelection())
            .build());
  }

  @Override
  public Control getControl() {
    return this;
  }

  @Override
  public void reinitialize() {
    onCaptureLoadingStart(false);
    if (models.capture.isLoaded()) {
      onCaptureLoaded(null);
    }
  }
}
