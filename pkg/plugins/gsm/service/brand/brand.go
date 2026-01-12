package brand

type Brand struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Slug            string `json:"slug,omitempty"`
	NumberOfDevices int    `json:"number_of_devices,omitempty"`
}
