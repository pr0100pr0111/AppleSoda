package lyrics

func ActiveIndex(lines []Line, posSec float64) int {
	if len(lines) == 0 {
		return 0
	}
	idx := 0
	for i, l := range lines {
		if l.TimeSec <= posSec {
			idx = i
		} else {
			break
		}
	}
	return idx
}
