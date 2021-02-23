package utils

// Difference return uint64s from a that are not in b
func Difference(a, b []uint64) []uint64 {
	m := make(map[uint64]bool)

	for _, item := range b {
		m[item] = true
	}

	diff := []uint64{}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// Intersection return uint64s from a that are not in b
func Intersection(a, b []uint64) []uint64 {
	m := make(map[uint64]bool)

	for _, item := range b {
		m[item] = true
	}

	inner := []uint64{}
	for _, item := range a {
		if _, ok := m[item]; ok {
			inner = append(inner, item)
		}
	}
	return inner
}

// DifferenceString return strings from a that are not in b
func DifferenceString(a, b []string) []string {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	diff := []string{}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// Equals check if array a and b have same elements
func Equals(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Unique remove duplicated elements from the array
func Unique(intSlice []uint64) []uint64 {
	keys := make(map[uint64]bool)
	list := []uint64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func MergeUintArray(a []uint64, b []uint64) []uint64 {
	return append(a, Difference(b, a)...)
}
