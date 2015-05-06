package util

type Comparable interface {
	AreEqual(int, int) bool
	Len() int
}

func GetUnique(l Comparable) []int {
	uniques := []int{}
	for i := 0; i < l.Len(); i++ {
		found := false
		for j := i + 1; j < l.Len(); j++ {
			if l.AreEqual(i, j) {
				found = true
				break
			}
		}
		if !found {
			uniques = append(uniques, i)
		}
	}

	return uniques
}
