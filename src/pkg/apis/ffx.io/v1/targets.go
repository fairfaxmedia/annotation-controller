package v1

type Target struct {
	Kind string            `json:",kind"`
	Data map[string]string `json:",data"`
}
