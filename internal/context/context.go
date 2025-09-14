// Copyright 2025 Roma Hlushko
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

package context

import (
	"context"

	"github.com/roma-glushko/frens/internal/store"
)

type AppContext struct {
	JournalDir string
	Store      store.Store
}

type ctxKey struct{}

func WithCtx(ctx context.Context, j *AppContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, j)
}

func FromCtx(ctx context.Context) *AppContext {
	if val := ctx.Value(ctxKey{}); val != nil {
		if jCtx, ok := val.(*AppContext); ok && jCtx != nil {
			return jCtx
		}
	}

	return nil
}
