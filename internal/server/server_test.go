package server_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/authena-ru/auto-test/internal/server"
	"github.com/authena-ru/auto-test/internal/service"
)

func TestServer_CheckAttemptToPassTestingTask(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name         string
		CheckAttempt func(attempt service.Attempt) service.AttemptCheckingResult
		Request      *server.CheckAttemptRequest
		Response     *server.CheckAttemptResponse
	}{
		{
			Name: "no_grade_when_empty_test_points_in_request",
			CheckAttempt: func(_ service.Attempt) service.AttemptCheckingResult {
				return service.AttemptCheckingResult{
					Grade:   service.NoGrade,
					Percent: -1,
				}
			},
			Request: &server.CheckAttemptRequest{
				GradeScale: &server.GradeScale{
					ExcellentLowerBound:    70,
					GoodLowerBound:         50,
					SatisfactoryLowerBound: 30,
				},
			},
			Response: &server.CheckAttemptResponse{
				Grade:   server.CheckAttemptResponse_NO_GRADE,
				Percent: -1,
			},
		},
		{
			Name: "excellent_when_test_points_presents_in_request",
			CheckAttempt: func(_ service.Attempt) service.AttemptCheckingResult {
				return service.AttemptCheckingResult{
					Grade:   service.Excellent,
					Percent: 100,
				}
			},
			Request: &server.CheckAttemptRequest{
				GradeScale: &server.GradeScale{
					ExcellentLowerBound:    90,
					GoodLowerBound:         70,
					SatisfactoryLowerBound: 50,
				},
				TestPoints: []*server.TestPoint{
					{CorrectVariantNumbers: []int32{0, 1}, ChosenVariantNumbers: []int32{1, 0}},
					{CorrectVariantNumbers: []int32{2}, ChosenVariantNumbers: []int32{2}},
					{CorrectVariantNumbers: []int32{1}, ChosenVariantNumbers: []int32{1}},
				},
			},
			Response: &server.CheckAttemptResponse{
				Grade:   server.CheckAttemptResponse_EXCELLENT,
				Percent: 100,
			},
		},
		{
			Name: "good_when_test_points_presents_in_request",
			CheckAttempt: func(_ service.Attempt) service.AttemptCheckingResult {
				return service.AttemptCheckingResult{
					Grade:   service.Good,
					Percent: 80,
				}
			},
			Request: &server.CheckAttemptRequest{
				GradeScale: &server.GradeScale{
					ExcellentLowerBound:    90,
					GoodLowerBound:         80,
					SatisfactoryLowerBound: 60,
				},
			},
			Response: &server.CheckAttemptResponse{
				Grade:   server.CheckAttemptResponse_GOOD,
				Percent: 80,
			},
		},
		{
			Name: "satisfactory_when_test_points_in_request",
			CheckAttempt: func(_ service.Attempt) service.AttemptCheckingResult {
				return service.AttemptCheckingResult{
					Grade:   service.Satisfactory,
					Percent: 60,
				}
			},
			Request: &server.CheckAttemptRequest{
				GradeScale: &server.GradeScale{
					ExcellentLowerBound:    90,
					GoodLowerBound:         70,
					SatisfactoryLowerBound: 60,
				},
			},
			Response: &server.CheckAttemptResponse{
				Grade:   server.CheckAttemptResponse_SATISFACTORY,
				Percent: 60,
			},
		},
		{
			Name: "unsatisfactory_when_test_points_in_request",
			CheckAttempt: func(_ service.Attempt) service.AttemptCheckingResult {
				return service.AttemptCheckingResult{
					Grade:   service.Unsatisfactory,
					Percent: 30,
				}
			},
			Request: &server.CheckAttemptRequest{
				GradeScale: &server.GradeScale{
					ExcellentLowerBound:    80,
					GoodLowerBound:         60,
					SatisfactoryLowerBound: 40,
				},
			},
			Response: &server.CheckAttemptResponse{
				Grade:   server.CheckAttemptResponse_UNSATISFACTORY,
				Percent: 30,
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			conn := newServerConnection(t, c.CheckAttempt)
			defer conn.Close()

			client := server.NewAutoTestServiceClient(conn)

			resp, err := client.CheckAttemptToPassTestingTask(context.Background(), c.Request)
			require.NoError(t, err)

			require.Equal(t, c.Response.Grade, resp.Grade)
			require.Equal(t, c.Response.Percent, resp.Percent)
		})
	}
}

func newServerConnection(
	t *testing.T,
	checkAttempt func(attempt service.Attempt) service.AttemptCheckingResult,
) *grpc.ClientConn {
	t.Helper()

	const bufSize = 1024 * 1024

	lis := bufconn.Listen(bufSize)
	servReg := grpc.NewServer()
	serv := server.New(checkAttempt)

	server.RegisterAutoTestServiceServer(servReg, serv)

	go func() {
		err := servReg.Serve(lis)
		require.NoError(t, err)
	}()

	ctx := context.Background()
	withDialer := grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
		return lis.Dial()
	})

	conn, err := grpc.DialContext(ctx, "bufnet", withDialer, grpc.WithInsecure())
	require.NoError(t, err)

	return conn
}
