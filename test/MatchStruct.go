package test

type Match struct {
	Input string
	InputMatch string
	Response string
	Score float64
}

type ByScore []Match

func (s ByScore) Len() int {
    return len(s)
}
func (s ByScore) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s ByScore) Less(i, j int) bool {
    return s[i].Score > s[j].Score
}
