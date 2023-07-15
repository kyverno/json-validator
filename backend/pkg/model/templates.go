package model

var DefaultAdmissionReview = AdmissionReview{
	Kind:       "AdmissionReview",
	APIVersion: "admission.k8s.io/v1",
	Request: Request{
		UID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
		Kind: Kind{
			Kind:    "MyJson",
			Group:   "testing.io",
			Version: "v1",
		},
		Resource: Resource{
			Resource: "myjsons",
			Group:    "testing.io",
			Version:  "v1",
		},
		RequestKind: Kind{
			Kind:    "MyJson",
			Group:   "testing.io",
			Version: "v1",
		},
		RequestResource: Resource{
			Resource: "myjsons",
			Group:    "testing.io",
			Version:  "v1",
		},
		Name:      "testing",
		Namespace: "default",
		UserInfo:  map[string]string{},
		Object: Object{
			Kind:       "MyJson",
			APIVersion: "testing.io/v1",
			Metadata: Metadata{
				Name:      "testing",
				Namespace: "default",
			},
		},
	},
}
