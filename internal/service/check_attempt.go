package service

import "math"

type Attempt struct {
	TestPoints []TestPoint
	GradeScale GradeScale
}

type GradeScale struct {
	excellentLowerBound    int
	goodLowerBound         int
	satisfactoryLowerBound int
}

func NewGradeScale(excellentLowerBound, goodLowerBound, satisfactoryLowerBound int32) GradeScale {
	return GradeScale{
		excellentLowerBound:    int(excellentLowerBound),
		goodLowerBound:         int(goodLowerBound),
		satisfactoryLowerBound: int(satisfactoryLowerBound),
	}
}

type TestPoint struct {
	correctVariantNumbers map[int]bool
	chosenVariantNumbers  map[int]bool
}

func NewTestPoint(correctVariantNumbers, chosenVariantNumbers []int32) TestPoint {
	correct := make(map[int]bool, len(correctVariantNumbers))
	for _, n := range correctVariantNumbers {
		correct[int(n)] = true
	}

	chosen := make(map[int]bool, len(chosenVariantNumbers))
	for _, n := range chosenVariantNumbers {
		chosen[int(n)] = true
	}

	return TestPoint{
		correctVariantNumbers: correct,
		chosenVariantNumbers:  chosen,
	}
}

type AttemptCheckingResult struct {
	Grade   Grade
	Percent int
}

type Grade uint

const (
	NoGrade Grade = iota
	Excellent
	Good
	Satisfactory
	Unsatisfactory
)

func CheckAttemptToPassTestingTask(attempt Attempt) AttemptCheckingResult {
	percent := calculatePercent(attempt.TestPoints)
	grade := evaluate(percent, attempt.GradeScale)

	return AttemptCheckingResult{
		Grade:   grade,
		Percent: percent,
	}
}

func evaluate(percent int, gradeScale GradeScale) Grade {
	if percent < 0 {
		return NoGrade
	}

	if percent >= gradeScale.excellentLowerBound {
		return Excellent
	}

	if percent >= gradeScale.goodLowerBound {
		return Good
	}

	if percent >= gradeScale.satisfactoryLowerBound {
		return Satisfactory
	}

	return Unsatisfactory
}

func calculatePercent(points []TestPoint) int {
	passed := computePassedTestPoints(points)

	pointsNumber := len(points)
	if pointsNumber == 0 {
		return -1
	}

	nonRounded := float64(passed) / float64(pointsNumber) * 100

	return int(math.Round(nonRounded))
}

func computePassedTestPoints(points []TestPoint) int {
	passedTestPoints := 0

	for _, tp := range points {
		if testPointCorrect(tp) {
			passedTestPoints++
		}
	}

	return passedTestPoints
}

func testPointCorrect(point TestPoint) bool {
	if len(point.correctVariantNumbers) != len(point.chosenVariantNumbers) {
		return false
	}

	for chosen := range point.chosenVariantNumbers {
		if !point.correctVariantNumbers[chosen] {
			return false
		}
	}

	return true
}
