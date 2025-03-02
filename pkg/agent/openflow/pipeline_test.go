// Copyright 2022 Antrea Authors
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

package openflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"antrea.io/antrea/pkg/agent/config"
	oftest "antrea.io/antrea/pkg/agent/openflow/testing"
)

func pipelineDefaultFlows(externalNodeEnabled, isEncap, isIPv4 bool) []string {
	if externalNodeEnabled {
		return []string{
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=200,ip actions=goto_table:ConntrackZone",
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=190 actions=goto_table:NonIP",
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=0 actions=drop",
			"cookie=0x1000000000000, table=ConntrackZone, priority=0 actions=goto_table:ConntrackState",
			"cookie=0x1000000000000, table=ConntrackState, priority=0 actions=goto_table:EgressSecurityClassifier",
			"cookie=0x1000000000000, table=EgressSecurityClassifier, priority=0 actions=goto_table:AntreaPolicyEgressRule",
			"cookie=0x1000000000000, table=AntreaPolicyEgressRule, priority=0 actions=goto_table:EgressRule",
			"cookie=0x1000000000000, table=EgressRule, priority=0 actions=goto_table:EgressDefaultRule",
			"cookie=0x1000000000000, table=EgressDefaultRule, priority=0 actions=goto_table:EgressMetric",
			"cookie=0x1000000000000, table=EgressMetric, priority=0 actions=goto_table:L3Forwarding",
			"cookie=0x1000000000000, table=L3Forwarding, priority=0 actions=goto_table:EgressMark",
			"cookie=0x1000000000000, table=EgressMark, priority=0 actions=goto_table:L2ForwardingCalc",
			"cookie=0x1000000000000, table=L2ForwardingCalc, priority=0 actions=goto_table:IngressSecurityClassifier",
			"cookie=0x1000000000000, table=IngressSecurityClassifier, priority=0 actions=goto_table:AntreaPolicyIngressRule",
			"cookie=0x1000000000000, table=AntreaPolicyIngressRule, priority=0 actions=goto_table:IngressRule",
			"cookie=0x1000000000000, table=IngressRule, priority=0 actions=goto_table:IngressDefaultRule",
			"cookie=0x1000000000000, table=IngressDefaultRule, priority=0 actions=goto_table:IngressMetric",
			"cookie=0x1000000000000, table=IngressMetric, priority=0 actions=goto_table:ConntrackCommit",
			"cookie=0x1000000000000, table=ConntrackCommit, priority=0 actions=goto_table:Output",
			"cookie=0x1000000000000, table=Output, priority=0 actions=drop",
			"cookie=0x1000000000000, table=NonIP, priority=0 actions=drop",
		}
	}

	var flows []string
	if isEncap {
		flows = []string{
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=0 actions=drop",
			"cookie=0x1000000000000, table=Classifier, priority=0 actions=drop",
			"cookie=0x1000000000000, table=SpoofGuard, priority=0 actions=drop",
			"cookie=0x1000000000000, table=UnSNAT, priority=0 actions=goto_table:ConntrackZone",
			"cookie=0x1000000000000, table=ConntrackZone, priority=0 actions=goto_table:ConntrackState",
			"cookie=0x1000000000000, table=ConntrackState, priority=0 actions=goto_table:PreRoutingClassifier",
			"cookie=0x1000000000000, table=PreRoutingClassifier, priority=0 actions=goto_table:SessionAffinity",
			"cookie=0x1000000000000, table=SessionAffinity, priority=0 actions=goto_table:ServiceLB",
			"cookie=0x1000000000000, table=ServiceLB, priority=0 actions=goto_table:EndpointDNAT",
			"cookie=0x1000000000000, table=EndpointDNAT, priority=0 actions=goto_table:AntreaPolicyEgressRule",
			"cookie=0x1000000000000, table=AntreaPolicyEgressRule, priority=0 actions=goto_table:EgressRule",
			"cookie=0x1000000000000, table=EgressRule, priority=0 actions=goto_table:EgressDefaultRule",
			"cookie=0x1000000000000, table=EgressDefaultRule, priority=0 actions=goto_table:EgressMetric",
			"cookie=0x1000000000000, table=EgressMetric, priority=0 actions=goto_table:L3Forwarding",
			"cookie=0x1000000000000, table=L3Forwarding, priority=0 actions=goto_table:EgressMark",
			"cookie=0x1000000000000, table=EgressMark, priority=0 actions=goto_table:L3DecTTL",
			"cookie=0x1000000000000, table=L3DecTTL, priority=0 actions=goto_table:SNATMark",
			"cookie=0x1000000000000, table=SNATMark, priority=0 actions=goto_table:SNAT",
			"cookie=0x1000000000000, table=SNAT, priority=0 actions=goto_table:L2ForwardingCalc",
			"cookie=0x1000000000000, table=L2ForwardingCalc, priority=0 actions=goto_table:TrafficControl",
			"cookie=0x1000000000000, table=TrafficControl, priority=0 actions=goto_table:IngressSecurityClassifier",
			"cookie=0x1000000000000, table=IngressSecurityClassifier, priority=0 actions=goto_table:AntreaPolicyIngressRule",
			"cookie=0x1000000000000, table=AntreaPolicyIngressRule, priority=0 actions=goto_table:IngressRule",
			"cookie=0x1000000000000, table=IngressRule, priority=0 actions=goto_table:IngressDefaultRule",
			"cookie=0x1000000000000, table=IngressDefaultRule, priority=0 actions=goto_table:IngressMetric",
			"cookie=0x1000000000000, table=IngressMetric, priority=0 actions=goto_table:ConntrackCommit",
			"cookie=0x1000000000000, table=ConntrackCommit, priority=0 actions=goto_table:Output",
			"cookie=0x1000000000000, table=Output, priority=0 actions=drop",
		}
		if isIPv4 {
			flows = append(flows,
				"cookie=0x1000000000000, table=PipelineRootClassifier, priority=200,arp actions=goto_table:ARPSpoofGuard",
				"cookie=0x1000000000000, table=PipelineRootClassifier, priority=200,ip actions=goto_table:Classifier",
				"cookie=0x1000000000000, table=ARPSpoofGuard, priority=0 actions=drop",
				"cookie=0x1000000000000, table=ARPResponder, priority=0 actions=drop",
				"cookie=0x1000000000000, table=PipelineIPClassifier, priority=210,ip,nw_dst=224.0.0.0/4 actions=resubmit:MulticastEgressRule",
				"cookie=0x1000000000000, table=PipelineIPClassifier, priority=0 actions=goto_table:UnSNAT",
				"cookie=0x1000000000000, table=MulticastEgressRule, priority=0 actions=goto_table:MulticastEgressMetric",
				"cookie=0x1000000000000, table=MulticastEgressMetric, priority=0 actions=goto_table:MulticastEgressPodMetric",
				"cookie=0x1000000000000, table=MulticastEgressPodMetric, priority=0 actions=goto_table:MulticastRouting",
				"cookie=0x1000000000000, table=MulticastRouting, priority=0 actions=goto_table:MulticastIngressRule",
				"cookie=0x1000000000000, table=MulticastIngressRule, priority=0 actions=goto_table:MulticastIngressMetric",
				"cookie=0x1000000000000, table=MulticastIngressMetric, priority=0 actions=goto_table:MulticastIngressPodMetric",
				"cookie=0x1000000000000, table=MulticastIngressPodMetric, priority=0 actions=goto_table:MulticastOutput",
				"cookie=0x1000000000000, table=MulticastOutput, priority=0 actions=drop",
			)
		} else {
			flows = append(flows,
				"cookie=0x1000000000000, table=PipelineRootClassifier, priority=200,ipv6 actions=goto_table:Classifier",
				"cookie=0x1000000000000, table=IPv6, priority=0 actions=goto_table:UnSNAT",
			)
		}
	} else {
		flows = []string{
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=200,arp actions=goto_table:ARPSpoofGuard",
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=0 actions=drop",
			"cookie=0x1000000000000, table=PipelineRootClassifier, priority=200,ip actions=goto_table:Classifier",
			"cookie=0x1000000000000, table=ARPSpoofGuard, priority=0 actions=drop",
			"cookie=0x1000000000000, table=ARPResponder, priority=0 actions=drop",
			"cookie=0x1000000000000, table=Classifier, priority=0 actions=drop",
			"cookie=0x1000000000000, table=SpoofGuard, priority=0 actions=drop",
			"cookie=0x1000000000000, table=PipelineIPClassifier, priority=210,ip,nw_dst=224.0.0.0/4 actions=resubmit:MulticastEgressRule",
			"cookie=0x1000000000000, table=PipelineIPClassifier, priority=0 actions=goto_table:UnSNAT",
			"cookie=0x1000000000000, table=UnSNAT, priority=0 actions=goto_table:ConntrackZone",
			"cookie=0x1000000000000, table=ConntrackZone, priority=0 actions=goto_table:ConntrackState",
			"cookie=0x1000000000000, table=ConntrackState, priority=0 actions=goto_table:PreRoutingClassifier",
			"cookie=0x1000000000000, table=PreRoutingClassifier, priority=0 actions=goto_table:SessionAffinity",
			"cookie=0x1000000000000, table=SessionAffinity, priority=0 actions=goto_table:ServiceLB",
			"cookie=0x1000000000000, table=ServiceLB, priority=0 actions=goto_table:EndpointDNAT",
			"cookie=0x1000000000000, table=EndpointDNAT, priority=0 actions=goto_table:AntreaPolicyEgressRule",
			"cookie=0x1000000000000, table=AntreaPolicyEgressRule, priority=0 actions=goto_table:EgressRule",
			"cookie=0x1000000000000, table=EgressRule, priority=0 actions=goto_table:EgressDefaultRule",
			"cookie=0x1000000000000, table=EgressDefaultRule, priority=0 actions=goto_table:EgressMetric",
			"cookie=0x1000000000000, table=EgressMetric, priority=0 actions=goto_table:L3Forwarding",
			"cookie=0x1000000000000, table=L3Forwarding, priority=0 actions=goto_table:EgressMark",
			"cookie=0x1000000000000, table=EgressMark, priority=0 actions=goto_table:L3DecTTL",
			"cookie=0x1000000000000, table=L3DecTTL, priority=0 actions=goto_table:SNATMark",
			"cookie=0x1000000000000, table=SNATMark, priority=0 actions=goto_table:SNAT",
			"cookie=0x1000000000000, table=SNAT, priority=0 actions=goto_table:L2ForwardingCalc",
			"cookie=0x1000000000000, table=L2ForwardingCalc, priority=0 actions=goto_table:TrafficControl",
			"cookie=0x1000000000000, table=TrafficControl, priority=0 actions=goto_table:IngressSecurityClassifier",
			"cookie=0x1000000000000, table=IngressSecurityClassifier, priority=0 actions=goto_table:AntreaPolicyIngressRule",
			"cookie=0x1000000000000, table=AntreaPolicyIngressRule, priority=0 actions=goto_table:IngressRule",
			"cookie=0x1000000000000, table=IngressRule, priority=0 actions=goto_table:IngressDefaultRule",
			"cookie=0x1000000000000, table=IngressDefaultRule, priority=0 actions=goto_table:IngressMetric",
			"cookie=0x1000000000000, table=IngressMetric, priority=0 actions=goto_table:ConntrackCommit",
			"cookie=0x1000000000000, table=ConntrackCommit, priority=0 actions=goto_table:VLAN",
			"cookie=0x1000000000000, table=VLAN, priority=0 actions=goto_table:Output",
			"cookie=0x1000000000000, table=Output, priority=0 actions=drop",
			"cookie=0x1000000000000, table=MulticastEgressRule, priority=0 actions=goto_table:MulticastEgressMetric",
			"cookie=0x1000000000000, table=MulticastEgressMetric, priority=0 actions=goto_table:MulticastEgressPodMetric",
			"cookie=0x1000000000000, table=MulticastEgressPodMetric, priority=0 actions=goto_table:MulticastRouting",
			"cookie=0x1000000000000, table=MulticastRouting, priority=0 actions=goto_table:MulticastIngressRule",
			"cookie=0x1000000000000, table=MulticastIngressRule, priority=0 actions=goto_table:MulticastIngressMetric",
			"cookie=0x1000000000000, table=MulticastIngressMetric, priority=0 actions=goto_table:MulticastIngressPodMetric",
			"cookie=0x1000000000000, table=MulticastIngressPodMetric, priority=0 actions=goto_table:MulticastOutput",
			"cookie=0x1000000000000, table=MulticastOutput, priority=0 actions=drop",
		}
	}
	return flows
}

func Test_client_defaultFlows(t *testing.T) {
	testCases := []struct {
		name             string
		enableIPv4       bool
		enableIPv6       bool
		nodeType         config.NodeType
		trafficEncapMode config.TrafficEncapModeType
		clientOptions    []clientOptionsFn
		expectedFlows    []string
	}{
		{
			name:             "IPv4,Encap,K8s Node",
			enableIPv4:       true,
			nodeType:         config.K8sNode,
			trafficEncapMode: config.TrafficEncapModeEncap,
			clientOptions:    []clientOptionsFn{enableTrafficControl, enableMulticast, enableMulticluster},
			expectedFlows:    pipelineDefaultFlows(false, true, true),
		},
		{
			name:             "IPv4,NoEncap,K8s Node",
			enableIPv4:       true,
			nodeType:         config.K8sNode,
			trafficEncapMode: config.TrafficEncapModeNoEncap,
			clientOptions:    []clientOptionsFn{enableTrafficControl, enableMulticast, enableConnectUplinkToBridge},
			expectedFlows:    pipelineDefaultFlows(false, false, true),
		},
		{
			name:             "IPv6,K8s Node",
			enableIPv6:       true,
			trafficEncapMode: config.TrafficEncapModeEncap,
			clientOptions:    []clientOptionsFn{enableTrafficControl},
			nodeType:         config.K8sNode,
			expectedFlows:    pipelineDefaultFlows(false, true, false),
		},
		{
			name:          "IPv4,ExternalNode Node",
			enableIPv4:    true,
			nodeType:      config.ExternalNode,
			expectedFlows: pipelineDefaultFlows(true, false, false),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := oftest.NewMockOFEntryOperations(ctrl)

			fc := newFakeClient(m, tc.enableIPv4, tc.enableIPv6, tc.nodeType, tc.trafficEncapMode, tc.clientOptions...)
			defer resetPipelines()

			assert.ElementsMatch(t, tc.expectedFlows, getFlowStrings(fc.defaultFlows()))
		})
	}
}
