// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0

// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the build v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=build.dev
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "build.dev", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme is used in the generated kube code
	AddToScheme = SchemeBuilder.AddToScheme
)
