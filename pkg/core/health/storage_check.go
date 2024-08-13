/*
Copyright The Karpor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package health

import (
	"context"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
)

type StorageCheck struct {
	storage.Storage
}

func NewStorageCheckHandler(sg storage.Storage) Check {
	return &StorageCheck{
		Storage: sg,
	}
}

func (s *StorageCheck) Pass(ctx context.Context) bool {
	err := s.Storage.CheckStorageHealth(ctx)
	return err == nil
}

func (s *StorageCheck) Name() string {
	return "StorageCheck"
}
