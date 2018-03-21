// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

// AuthMethod specifies a mask of authentication methods.
type AuthMethod uint

const (
	// AuthMethodAppID indicates the app-id authentication method
	AuthMethodAppID AuthMethod = 0x01
	// AuthMethodAppRole indicates the approle authentication method
	AuthMethodAppRole AuthMethod = 0x02
)

// IsEnabled returns true if the given specific authentication method is contained in the given mask.
func (mask AuthMethod) IsEnabled(specific AuthMethod) bool {
	return mask&specific == specific
}
