// Copyright The Karbour Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resourcegroup

import (
	"context"
)

type ResourceGroupManager struct{}

func NewResourceGroupManager() *ResourceGroupManager {
	return &ResourceGroupManager{}
}

func (c *ResourceGroupManager) GetResourceGroup(ctx context.Context) error {
	panic("need to implement")
}

func (c *ResourceGroupManager) ListResourceGroup(ctx context.Context) error {
	panic("need to implement")
}

func (c *ResourceGroupManager) CreateResourceGroup(ctx context.Context) error {
	panic("need to implement")
}

func (c *ResourceGroupManager) UpdateResourceGroup(ctx context.Context) error {
	panic("need to implement")
}

func (c *ResourceGroupManager) DeleteResourceGroup(ctx context.Context) error {
	panic("need to implement")
}
