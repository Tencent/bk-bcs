/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// DNSEndpoints returns a DNSEndpointInformer.
	DNSEndpoints() DNSEndpointInformer
	// Domains returns a DomainInformer.
	Domains() DomainInformer
	// IngressDNSRecords returns a IngressDNSRecordInformer.
	IngressDNSRecords() IngressDNSRecordInformer
	// ServiceDNSRecords returns a ServiceDNSRecordInformer.
	ServiceDNSRecords() ServiceDNSRecordInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// DNSEndpoints returns a DNSEndpointInformer.
func (v *version) DNSEndpoints() DNSEndpointInformer {
	return &dNSEndpointInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Domains returns a DomainInformer.
func (v *version) Domains() DomainInformer {
	return &domainInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// IngressDNSRecords returns a IngressDNSRecordInformer.
func (v *version) IngressDNSRecords() IngressDNSRecordInformer {
	return &ingressDNSRecordInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ServiceDNSRecords returns a ServiceDNSRecordInformer.
func (v *version) ServiceDNSRecords() ServiceDNSRecordInformer {
	return &serviceDNSRecordInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}