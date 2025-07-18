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

package version

import (
	"fmt"
	"runtime"
)

const AppName = "frens"

var Version = "dev"

// GitCommit is the commit hash of the current version.
var GitCommit = "unknown"

// BuildDate is the date when the binary was built.
var BuildDate = "unknown"

var FullVersion string

func init() {
	FullVersion = fmt.Sprintf(
		"%s (commit: %s, built at: %s, runtime: %s)",
		Version,
		GitCommit,
		BuildDate,
		runtime.Version(),
	)
}
