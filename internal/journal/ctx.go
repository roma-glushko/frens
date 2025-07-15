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

package journal

import "context"

type journalCtxKey struct{}

func WithCtx(ctx context.Context, j *Journal) context.Context {
	return context.WithValue(ctx, journalCtxKey{}, j)
}

func FromCtx(ctx context.Context) *Journal {
	if val := ctx.Value(journalCtxKey{}); val != nil {
		if j, ok := val.(*Journal); ok && j != nil {
			return j
		}
	}

	return nil
}
