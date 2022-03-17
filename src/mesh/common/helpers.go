package common

// Contains checks if slice contains given value
func Contains[K comparable](s []K, e K) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SubsetOf checks if slice s is a subset of slice e
func SubsetOf[K comparable](s []K, e []K) bool {
	for _, a := range s {
		if !Contains(e, a) {
			return false
		}
	}
	return true
}
