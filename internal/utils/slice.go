// Copyright 2026 Roma Hlushko
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

package utils

func Unique[T comparable](s []T) []T {
	m := make(map[T]struct{}, len(s))

	for _, v := range s {
		m[v] = struct{}{}
	}

	unique := make([]T, 0, len(m))

	for k := range m {
		unique = append(unique, k)
	}

	return unique
}
