// Copyright © 2017 Dennis Walters
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"os"
)

const (
	envVar = "MOCKABLE"
)

// Mocked is true if the MOCKABLE environment variable is set, but is false
// otherwise.
func Mocked() bool {
	_, present := os.LookupEnv(envVar)

	return present
}

// Enable sets the MOCKABLE environment variable to a non-null value. This is
// really only handy for use within test suites.
func Enable() {
	if !Mocked() {
		os.Setenv(envVar, "1")
	}
}

// Disable deletes the MOCKABLE environment variable from the environment. This
// is really only handy for use within test suites.
func Disable() {
	if Mocked() {
		os.Unsetenv(envVar)
	}
}
