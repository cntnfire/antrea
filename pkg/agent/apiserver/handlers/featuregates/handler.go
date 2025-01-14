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

package featuregates

import (
	"encoding/json"
	"net/http"
	"sort"

	"k8s.io/klog/v2"

	"antrea.io/antrea/pkg/features"
)

type Response struct {
	Component string `json:"component,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
	Version   string `json:"version,omitempty"`
}

// HandleFunc returns the function which can handle queries issued by 'antctl get featuregates' command.
// The handler function populates Antrea Agent feature gates information to the response.
func HandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var featureGates []Response
		for df := range features.DefaultAntreaFeatureGates {
			if features.AgentGates.Has(df) {
				featureGates = append(featureGates, Response{
					Component: "agent",
					Name:      string(df),
					Status:    getStatus(features.DefaultFeatureGate.Enabled(df)),
					Version:   string(features.DefaultAntreaFeatureGates[df].PreRelease),
				})
			}
		}
		sort.Slice(featureGates, func(i, j int) bool {
			return featureGates[i].Name < featureGates[j].Name
		})
		err := json.NewEncoder(w).Encode(featureGates)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			klog.ErrorS(err, "Error when encoding FeatureGates to json")
		}
	}
}

func getStatus(status bool) string {
	if status {
		return "Enabled"
	}
	return "Disabled"
}
