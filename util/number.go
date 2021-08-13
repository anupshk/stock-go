package util

func MinMax(array []float64) (min float64, max float64) {
	min = array[0]
	max = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return
}
