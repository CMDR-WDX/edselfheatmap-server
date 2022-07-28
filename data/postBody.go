package data

type RequestBody struct {
	SystemName *string  `json:"systemName"`
	X          *float32 `json:"x"`
	Y          *float32 `json:"y"`
	Z          *float32 `json:"z"`
}

type PixelEntry struct {
	SystemName string
	X          int
	Y          int
}

func MakePixelEntry(r RequestBody) PixelEntry {
	return PixelEntry{
		SystemName: *r.SystemName,
		X:          int(*r.X) / 10,
		Y:          int(*r.Z) / 10,
	}
}
