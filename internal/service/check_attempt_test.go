package service_test

import (
	"testing"

	"github.com/authena-ru/auto-test/internal/service"
	"github.com/stretchr/testify/require"
)

func TestCheckAttemptToPassTestingTask(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                  string
		Attempt               service.Attempt
		AttemptCheckingResult service.AttemptCheckingResult
	}{
		{
			Name: "no_grade_when_empty_test_points",
			Attempt: service.Attempt{
				GradeScale: service.NewGradeScale(75, 50, 25),
			},
			AttemptCheckingResult: service.AttemptCheckingResult{
				Grade:   service.NoGrade,
				Percent: -1,
			},
		},
		{
			Name: "excellent_grade_when_better_than_lower_bound",
			Attempt: service.Attempt{
				TestPoints: []service.TestPoint{
					service.NewTestPoint([]int32{0, 1}, []int32{0}),
					service.NewTestPoint([]int32{2}, []int32{2}),
					service.NewTestPoint([]int32{3}, []int32{3}),
					service.NewTestPoint([]int32{0, 2}, []int32{0, 2}),
				},
				GradeScale: service.NewGradeScale(75, 60, 30),
			},
			AttemptCheckingResult: service.AttemptCheckingResult{
				Grade:   service.Excellent,
				Percent: 75,
			},
		},
		{
			Name: "good_grade_when_better_than_lower_bound",
			Attempt: service.Attempt{
				TestPoints: []service.TestPoint{
					service.NewTestPoint([]int32{1, 3}, []int32{1, 3}),
					service.NewTestPoint([]int32{0}, []int32{1}),
					service.NewTestPoint([]int32{2, 3}, []int32{2, 4}),
					service.NewTestPoint([]int32{1}, []int32{1}),
					service.NewTestPoint([]int32{3}, []int32{3}),
				},
				GradeScale: service.NewGradeScale(90, 60, 40),
			},
			AttemptCheckingResult: service.AttemptCheckingResult{
				Grade:   service.Good,
				Percent: 60,
			},
		},
		{
			Name: "satisfactory_grade_when_better_than_lower_bound",
			Attempt: service.Attempt{
				TestPoints: []service.TestPoint{
					service.NewTestPoint([]int32{0}, []int32{0}),
					service.NewTestPoint([]int32{1, 0}, []int32{0, 1}),
					service.NewTestPoint([]int32{2}, []int32{3, 2}),
					service.NewTestPoint([]int32{1}, []int32{2}),
					service.NewTestPoint([]int32{0}, []int32{2}),
				},
				GradeScale: service.NewGradeScale(80, 60, 40),
			},
			AttemptCheckingResult: service.AttemptCheckingResult{
				Grade:   service.Satisfactory,
				Percent: 40,
			},
		},
		{
			Name: "unsatisfactory_when_worse_than_satisfactory_lower_bound",
			Attempt: service.Attempt{
				TestPoints: []service.TestPoint{
					service.NewTestPoint([]int32{1}, []int32{1}),
					service.NewTestPoint([]int32{2}, []int32{3}),
					service.NewTestPoint([]int32{4}, []int32{5}),
					service.NewTestPoint([]int32{0}, []int32{5}),
					service.NewTestPoint([]int32{1}, []int32{2}),
					service.NewTestPoint([]int32{1}, []int32{3, 4}),
				},
				GradeScale: service.NewGradeScale(90, 50, 18),
			},
			AttemptCheckingResult: service.AttemptCheckingResult{
				Grade:   service.Unsatisfactory,
				Percent: 17,
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			givenResult := service.CheckAttemptToPassTestingTask(c.Attempt)

			require.Equal(t, c.AttemptCheckingResult, givenResult)
		})
	}
}
