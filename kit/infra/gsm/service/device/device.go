package device

type Device struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Image       string `json:"image,omitempty"`
	Description string `json:"description,omitempty"`
}

type DeviceList struct {
	Items       []Device `json:"items,omitempty"`
	CurrentPage int      `json:"current_page,omitempty"`
	TotalPage   int      `json:"total_page,omitempty"`
}

func (list *DeviceList) AddDevice(device Device) {
	items := append(list.Items, device)
	list.Items = items
	return
}
