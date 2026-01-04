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

package lang

import (
	"reflect"
	"regexp"
	"strings"
)

var propRe *regexp.Regexp

func init() {
	propRe = regexp.MustCompile(
		`\$(?P<name>[^\s:]+):(?P<value>[^\s$]+)`,
	)
}

func ExtractProps[T any](s string) (*T, error) {
	matches := propRe.FindAllStringSubmatch(s, -1)
	names := propRe.SubexpNames()

	out := new(T)
	props := map[string]string{}

	for _, match := range matches {
		var name, value string

		for i, n := range names {
			if n == "name" {
				name = match[i]
			}

			if n == "value" {
				value = strings.TrimSpace(match[i])
			}
		}

		props[name] = value
	}

	// Use reflection to set fields on the struct
	v := reflect.ValueOf(out).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		key := field.Tag.Get("frentxt")

		if val, ok := props[key]; ok {
			v.Field(i).SetString(val) // handle string only here
		}
	}

	return out, nil
}

func RenderProps[T any](props T) string {
	v := reflect.ValueOf(props)
	t := v.Type()

	var sb strings.Builder

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.IsZero() {
			continue
		}

		key := field.Tag.Get("frentxt")

		if key == "" {
			continue
		}

		sb.WriteString("$")
		sb.WriteString(key)
		sb.WriteString(":")
		sb.WriteString(value.String())
		sb.WriteString(" ")
	}

	return strings.TrimSpace(sb.String())
}

func RemoveProps(s string) string {
	return propRe.ReplaceAllString(s, "")
}
