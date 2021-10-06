package halstead

import "math"

type HalsteadInfo struct {
	operators map[string]uint
	operands  map[string]uint
}

func (info *HalsteadInfo) getN1Distinct() uint {
	return uint(len(info.operators))
}

func (info *HalsteadInfo) getN2Distinct() uint {
	return uint(len(info.operands))
}

// Not universal key type because using interfaces is nasty
// and generics (with 1.18 version) are not released yet
func sumMap(targetMap map[string]uint) uint {
	var total uint = 0
	for _, count := range targetMap {
		total += count
	}
	return total
}

func (info *HalsteadInfo) getN1Total() uint {
	sum := sumMap(info.operands)
	return sum
}

func (info *HalsteadInfo) getN2Total() uint {
	sum := sumMap(info.operators)
	return sum
}

func (info *HalsteadInfo) Vocabuary() uint {
	n1Dist, n2Dist := info.getN1Distinct(), info.getN2Distinct()
	return n1Dist + n2Dist
}

func (info *HalsteadInfo) Length() uint {
	n1Tot, n2Tot := info.getN1Total(), info.getN2Total()
	return n1Tot + n2Tot
}

func (info *HalsteadInfo) Volume() float64 {
	nTot := info.Length()
	nDist := info.Vocabuary()
	return float64(nTot) * math.Log2(float64(nDist))
}

func (info *HalsteadInfo) Difficulty() float64 {
	n1Dist, n2Dist := info.getN1Distinct(), info.getN2Distinct()
	n2Tot := info.getN2Total()
	return (float64(n1Dist) * float64(n2Tot)) / (2 * float64(n2Dist))
}

func (info *HalsteadInfo) Effort() float64 {
	// Do not use other functions to avoid summing maps multiple times
	n1Dist, n2Dist := float64(info.getN1Distinct()), float64(info.getN2Distinct())
	n1Tot, n2Tot := float64(info.getN1Total()), float64(info.getN2Total())
	D := (n1Dist * n2Tot) / (2 * n2Dist)
	V := (n1Tot + n2Tot) * math.Log2(n1Dist + n2Dist)
	return D * V
}
