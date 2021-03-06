// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0

package validate

import (
	"context"
	"fmt"
	"sort"
	"strings"

	build "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretRef contains all required fields
// to validate a Build spec secrets definitions
type SecretRef struct {
	Build  *build.Build
	Client client.Client
}

// ValidatePath implements BuildPath interface and validates
// that all referenced secrets under spec exists
func (s SecretRef) ValidatePath(ctx context.Context) error {
	var missingSecrets []string
	secret := &corev1.Secret{}

	secretNames := s.buildSecretReferences()

	for refSecret, secretType := range secretNames {
		if err := s.Client.Get(ctx, types.NamespacedName{Name: refSecret, Namespace: s.Build.Namespace}, secret); err != nil && !apierrors.IsNotFound(err) {
			return err
		} else if apierrors.IsNotFound(err) {
			s.Build.Status.Reason = secretType
			s.Build.Status.Message = fmt.Sprintf("referenced secret %s not found", refSecret)
			missingSecrets = append(missingSecrets, refSecret)
		}
	}

	// sorts a list of secret names in increasing order
	sort.Strings(missingSecrets)

	if len(missingSecrets) > 1 {
		s.Build.Status.Reason = build.MultipleSecretRefNotFound
		s.Build.Status.Message = fmt.Sprintf("missing secrets are %s", strings.Join(missingSecrets, ","))
	}
	return nil
}

func (s SecretRef) buildSecretReferences() map[string]build.BuildReason {
	// Validate if the referenced secrets exist in the namespace
	secretRefMap := map[string]build.BuildReason{}
	if s.Build.Spec.Output.SecretRef != nil && s.Build.Spec.Output.SecretRef.Name != "" {
		secretRefMap[s.Build.Spec.Output.SecretRef.Name] = build.SpecOutputSecretRefNotFound
	}
	if s.Build.Spec.Source.SecretRef != nil && s.Build.Spec.Source.SecretRef.Name != "" {
		secretRefMap[s.Build.Spec.Source.SecretRef.Name] = build.SpecSourceSecretRefNotFound
	}
	if s.Build.Spec.BuilderImage != nil && s.Build.Spec.BuilderImage.SecretRef != nil && s.Build.Spec.BuilderImage.SecretRef.Name != "" {
		secretRefMap[s.Build.Spec.BuilderImage.SecretRef.Name] = build.SpecRuntimeSecretRefNotFound
	}
	return secretRefMap
}
