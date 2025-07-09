package specification

const (
	Tablet = "Tablet"
	Phone  = "Phone"
)

type Specification struct {
	Brand      string                `json:"brand,omitempty"`
	DeviceName string                `json:"device_name,omitempty"`
	DeviceType string                `json:"device_type,omitempty"`
	ImageURL   string                `json:"image_url,omitempty"`
	Overview   SpecificationOverview `json:"overview,omitempty"`
	Detail     SpecificationDetail   `json:"detail,omitempty"`
}
