package model

type AdmissionReview struct {
	Kind       string  `json:"kind"`
	APIVersion string  `json:"apiVersion"`
	Request    Request `json:"request"`
	OldObject  any     `json:"oldObject"`
	DryRun     bool    `json:"dryRun"`
	Options    any     `json:"options"`
}

type Kind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

type Resource struct {
	Group    string `json:"group"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Object struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       any      `json:"spec"`
}

type Request struct {
	UID             string   `json:"uid"`
	Kind            Kind     `json:"kind"`
	Resource        Resource `json:"resource"`
	RequestKind     Kind     `json:"requestKind"`
	RequestResource Resource `json:"requestResource"`
	Name            string   `json:"name"`
	Namespace       string   `json:"namespace"`
	Operation       string   `json:"operation"`
	UserInfo        any      `json:"userInfo"`
	Roles           any      `json:"roles"`
	ClusterRoles    any      `json:"clusterRoles"`
	Object          Object   `json:"object"`
	OldObject       any      `json:"oldObject"`
	DryRun          bool     `json:"dryRun"`
	Options         any      `json:"options"`
}

type Response struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Request    struct {
		UID  string `json:"uid"`
		Kind struct {
			Group   string `json:"group"`
			Version string `json:"version"`
			Kind    string `json:"kind"`
		} `json:"kind"`
		Resource struct {
			Group    string `json:"group"`
			Version  string `json:"version"`
			Resource string `json:"resource"`
		} `json:"resource"`
		RequestKind struct {
			Group   string `json:"group"`
			Version string `json:"version"`
			Kind    string `json:"kind"`
		} `json:"requestKind"`
		RequestResource struct {
			Group    string `json:"group"`
			Version  string `json:"version"`
			Resource string `json:"resource"`
		} `json:"requestResource"`
		Name      string   `json:"name"`
		Namespace string   `json:"namespace"`
		Operation string   `json:"operation"`
		UserInfo  struct{} `json:"userInfo"`
		Object    struct {
			APIVersion string `json:"apiVersion"`
			Kind       string `json:"kind"`
			Metadata   struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec any
		} `json:"object"`
		OldObject any  `json:"oldObject"`
		DryRun    bool `json:"dryRun"`
		Options   any  `json:"options"`
	} `json:"request"`
	Response struct {
		UID     string `json:"uid"`
		Allowed bool   `json:"allowed"`
		Status  struct {
			Metadata struct{} `json:"metadata"`
			Status   string   `json:"status"`
			Message  string   `json:"message"`
		} `json:"status"`
	} `json:"response"`
}
