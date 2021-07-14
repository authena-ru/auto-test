package server

import (
	"context"

	"github.com/authena-ru/auto-test/internal/service"
)

type Server struct {
	UnimplementedAutoTestServiceServer
	checkAttempt CheckAttemptFunction
}

type CheckAttemptFunction func(attempt service.Attempt) service.AttemptCheckingResult

func New(checkAttempt CheckAttemptFunction) *Server {
	return &Server{checkAttempt: checkAttempt}
}

func (s *Server) CheckAttemptToPassTestingTask(_ context.Context, req *CheckAttemptRequest) (*CheckAttemptResponse, error) {
	attempt := unmarshalAttempt(req)
	result := s.checkAttempt(attempt)

	return marshalCheckAttemptResponse(result), nil
}

func unmarshalAttempt(req *CheckAttemptRequest) service.Attempt {
	return service.Attempt{
		TestPoints: marshalTestPoints(req.TestPoints),
		GradeScale: marshalGradeScale(req.GradeScale),
	}
}

func marshalTestPoints(tps []*TestPoint) []service.TestPoint {
	points := make([]service.TestPoint, len(tps))

	for _, tp := range tps {
		if tp == nil {
			continue
		}

		points = append(points, service.NewTestPoint(
			tp.CorrectVariantNumbers,
			tp.ChosenVariantNumbers,
		))
	}

	return points
}

func marshalGradeScale(scale *GradeScale) service.GradeScale {
	if scale == nil {
		return service.GradeScale{}
	}

	return service.NewGradeScale(
		scale.GetExcellentLowerBound(),
		scale.GetExcellentLowerBound(),
		scale.GetSatisfactoryLowerBound(),
	)
}

func marshalCheckAttemptResponse(result service.AttemptCheckingResult) *CheckAttemptResponse {
	return &CheckAttemptResponse{
		Grade:   marshalGrade(result.Grade),
		Percent: int32(result.Percent),
	}
}

func marshalGrade(grade service.Grade) CheckAttemptResponse_Grade {
	switch grade {
	case service.Excellent:
		return CheckAttemptResponse_EXCELLENT
	case service.Good:
		return CheckAttemptResponse_GOOD
	case service.Satisfactory:
		return CheckAttemptResponse_SATISFACTORY
	case service.Unsatisfactory:
		return CheckAttemptResponse_UNSATISFACTORY
	case service.NoGrade:
		return CheckAttemptResponse_NO_GRADE
	}

	return CheckAttemptResponse_NO_GRADE
}
