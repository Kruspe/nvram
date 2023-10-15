package problem_test

import (
	"github.com/kruspe/nvram/problem"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type matrixSuite struct {
	suite.Suite
}

func Test_MatrixSuite(t *testing.T) {
	suite.Run(t, &matrixSuite{})
}

func (s *matrixSuite) Test_Multiply_1() {
	a := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	b := [][]int{
		{7, 8},
		{9, 10},
		{11, 12},
	}
	result, err := problem.Multiply(a, b)

	require.NoError(s.T(), err)
	require.Equal(s.T(), [][]int{
		{58, 64},
		{139, 154},
	}, result)
}

func (s *matrixSuite) Test_Multiply_2() {
	a := [][]int{
		{2, 5},
		{1, 0},
	}
	b := [][]int{
		{3, 4},
		{1, 1},
	}
	result, err := problem.Multiply(a, b)

	require.NoError(s.T(), err)
	require.Equal(s.T(), [][]int{
		{11, 13},
		{3, 4},
	}, result)
}

func (s *matrixSuite) Test_Multiply_3() {
	a := [][]int{
		{1, -3, 7},
		{-2, 3, 1},
		{3, 5, 5},
	}
	b := [][]int{
		{2, 1, 4},
		{0, 2, 2},
		{4, 3, 2},
	}
	result, err := problem.Multiply(a, b)

	require.NoError(s.T(), err)
	require.Equal(s.T(), [][]int{
		{30, 16, 12},
		{0, 7, 0},
		{26, 28, 32},
	}, result)
}

func (s *matrixSuite) Test_Multiply_Error() {
	a := [][]int{
		{1, 2},
		{3, 4},
	}
	b := [][]int{
		{5, 6},
		{7, 8},
		{9, 10},
	}
	_, err := problem.Multiply(a, b)

	require.Error(s.T(), err)
}
