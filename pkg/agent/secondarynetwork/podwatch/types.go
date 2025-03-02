// Copyright 2021 Antrea Authors
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

package podwatch

import "antrea.io/antrea/pkg/agent/secondarynetwork/cnipodcache"

type RouteInfo struct {
	Dst string `json:"dst,omitempty"`
}

type IPAMConfig struct {
	Type       string    `json:"type,omitempty"`
	Subnet     string    `json:"subnet,omitempty"`
	RangeStart string    `json:"rangeStart,omitempty"`
	RangeEnd   string    `json:"rangeEnd,omitempty"`
	Routes     RouteInfo `json:"routes,omitempty"`
	Gateway    string    `json:"gateway,omitempty"`
}

const (
	sriovNetworkType cnipodcache.NetworkType = "sriov"
	vlanNetworkType  cnipodcache.NetworkType = "vlan"
)

type SecondaryNetworkConfig struct {
	CNIVersion string `json:"cniVersion,omitempty"`
	Name       string `json:"name,omitempty"`
	// Set type to "antrea"
	Type string `json:"type,omitempty"`
	// Set networkType to "sriov"
	NetworkType cnipodcache.NetworkType `json:"networkType,omitempty"`

	MTU  int32      `json:"mtu,omitempty"`
	VLAN int32      `json:"vlan,omitempty"`
	IPAM IPAMConfig `json:"ipam,omitempty"`
}
