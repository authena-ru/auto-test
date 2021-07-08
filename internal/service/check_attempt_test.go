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
				GradeScale: service.GradeScale{
					ExcellentLowerBound:    75,
					GoodLowerBound:         50,
					SatisfactoryLowerBound: 25,
				},
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
					service.NewTestPoint([]int{0, 1}, []int{0}),
					service.NewTestPoint([]int{2}, []int{2}),
					service.NewTestPoint([]int{3}, []int{3}),
					service.NewTestPoint([]int{0, 2}, []int{0, 2}),
				},
				GradeScale: service.GradeScale{
					ExcellentLowerBound:    75,
					GoodLowerBound:         60,
					SatisfactoryLowerBound: 30,
				},
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
					service.NewTestPoint([]int{1, 3}, []int{1, 3}),
					service.NewTestPoint([]int{0}, []int{1}),
					service.NewTestPoint([]int{2, 3}, []int{2, 4}),
					service.NewTestPoint([]int{1}, []int{1}),
					service.NewTestPoint([]int{3}, []int{3}),
				},
				GradeScale: service.GradeScale{
					ExcellentLowerBound:    90,
					GoodLowerBound:         60,
					SatisfactoryLowerBound: 40,
				},
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
					service.NewTestPoint([]int{0}, []int{0}),
					service.NewTestPoint([]int{1, 0}, []int{0, 1}),
					service.NewTestPoint([]int{2}, []int{3, 2}),
					service.NewTestPoint([]int{1}, []int{2}),
					service.NewTestPoint([]int{0}, []int{2}),
				},
				GradeScale: service.GradeScale{
					ExcellentLowerBound:    80,
					GoodLowerBound:         60,
					SatisfactoryLowerBound: 40,
				},
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
					service.NewTestPoint([]int{1}, []int{1}),
					service.NewTestPoint([]int{2}, []int{3}),
					service.NewTestPoint([]int{4}, []int{5}),
					service.NewTestPoint([]int{0}, []int{5}),
					service.NewTestPoint([]int{1}, []int{2}),
					service.NewTestPoint([]int{1}, []int{3, 4}),
				},
				GradeScale: service.GradeScale{
					ExcellentLowerBound:    90,
					GoodLowerBound:         50,
					SatisfactoryLowerBound: 18,
				},
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
