package utilities

type Counter struct {
	Count uint    `json:"count"`
	Avg   float32 `json:"avg"`
	ID    uint    `json:"id"`
}

func RemoveDuplicates(intSlice []uint) []uint {
	keys := make(map[uint]bool)
	list := make([]uint, 0)
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
