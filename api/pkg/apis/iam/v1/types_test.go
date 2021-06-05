package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestShowcase(t *testing.T) {

	// organizationScope represents the scope for a given organization
	organizationScope := OrganizationScope{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nrc",
		},
	}

	// ugandaScope represents the NRC Uganda operation context
	ugandaScope := LocalScope{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "uganda",
			Namespace: "nrc",
		},
	}

	// kenyaScope represents the NRC Kenya operation context
	kenyaScope := LocalScope{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kenya",
			Namespace: "nrc",
		},
	}

	// firstNameAttribute represents a user attribute that is
	// understood throughout NRC
	firstNameAttribute := OrganizationUserAttribute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "nrc",
			Name:      "firstName",
		},
		Spec: OrganizationUserAttributeSpec{
			Type: "String",
		},
	}

	// lastNameAttribute represents a user attribute that is
	// understood throughout NRC
	lastNameAttribute := OrganizationUserAttribute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "nrc",
			Name:      "lastName",
		},
		Spec: OrganizationUserAttributeSpec{
			Type: "String",
		},
	}

	// ugandaGenderBasedViolence represents a user attribute
	// that is understood in Kenya
	ugandaGenderBasedViolence := Claim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "genderBasedViolence",
			Namespace: "nrc",
		},
		Spec: UserAttributeSpec{
			Type: "String",
			ScopeRef: ScopeRef{
				Kind:     "LocalScope",
				APIGroup: "iam.nrc.no",
				Name:     ugandaScope.Name,
			},
		},
	}

}
