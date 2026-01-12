package specification

type SpecificationOverviewPopularity struct {
	Percentage string `json:"percentage,omitempty"`
	Hits       string `json:"hits,omitempty"`
}

type SpecificationOverviewGeneralInfo struct {
	Launched string `json:"launched,omitempty"`
	Body     string `json:"body,omitempty"`
	OS       string `json:"os,omitempty"`
	Storage  string `json:"storage,omitempty"`
}

type SpecificationOverviewDisplay struct {
	Touchscreen string `json:"touchscreen,omitempty"`
	Size        string `json:"size,omitempty"`
	Resolution  string `json:"resolution,omitempty"`
}

type SpecificationOverviewCamera struct {
	Photo string `json:"photo,omitempty"`
	Video string `json:"video,omitempty"`
}

type SpecificationOverviewExpansion struct {
	RAM     string `json:"ram,omitempty"`
	Chipset string `json:"chipset,omitempty"`
}

type SpecificationOverviewBattery struct {
	Capacity   string `json:"capacity,omitempty"`
	Technology string `json:"technology,omitempty"`
}

type SpecificationOverview struct {
	GeneralInfo SpecificationOverviewGeneralInfo `json:"general_info,omitempty"`
	Display     SpecificationOverviewDisplay     `json:"display,omitempty"`
	Camera      SpecificationOverviewCamera      `json:"camera,omitempty"`
	Expansion   SpecificationOverviewExpansion   `json:"expansion,omitempty"`
	Battery     SpecificationOverviewBattery     `json:"battery,omitempty"`
}
