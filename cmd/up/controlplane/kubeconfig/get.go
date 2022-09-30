// Copyright 2021 Upbound Inc
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

package kubeconfig

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/pterm/pterm"

	"github.com/upbound/up/internal/kube"
	"github.com/upbound/up/internal/upbound"
)

// AfterApply sets default values in command after assignment and validation.
func (c *getCmd) AfterApply(upCtx *upbound.Context) error {
	c.stdin = os.Stdin
	return nil
}

// getCmd gets kubeconfig data for an Upbound control plane.
type getCmd struct {
	stdin io.Reader

	File  string `type:"path" short:"f" help:"File to merge kubeconfig."`
	Token string `required:"" help:"API token used to authenticate."`

	ID string `arg:"" name:"control-plane-ID" required:"" help:"ID (name) of control plane."`
}

// Run executes the get command.
func (c *getCmd) Run(p pterm.TextPrinter, upCtx *upbound.Context) error {
	// TODO(hasheddan): consider implementing a custom decoder
	if c.Token == "-" {
		b, err := io.ReadAll(c.stdin)
		if err != nil {
			return err
		}
		c.Token = strings.TrimSpace(string(b))
	}

	var err error
	var kubeCurCtx string

	upCtx.ProxyEndpoint.Path = "/v1/controlPlanes"
	kubeCurCtx, err = kube.BuildControlPlaneKubeconfig(upCtx.ProxyEndpoint, path.Join(upCtx.Account, c.ID), c.Token, c.File)
	if err != nil {
		return err
	}

	p.Printfln("Current context set to %s", kubeCurCtx)
	return nil
}
