package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Ternary[T any](expression bool, trueResult T, falseResult T) T {
	if expression {
		return trueResult
	} else {
		return falseResult
	}
}
func TransposeMatrix[T any](matrix [][]T) [][]T {
	// Get the dimensions of the matrix
	numRows := len(matrix)
	numCols := len(matrix[0])
	// Create a new transposed matrix with swapped dimensions
	transposed := make([][]T, numCols)
	for i := range transposed {
		transposed[i] = make([]T, numRows)
	}
	// Fill in the transposed matrix by swapping rows and columns
	for i, row := range matrix {
		for j, val := range row {
			transposed[j][i] = val
		}
	}
	return transposed
}
func ExtractHostPort(addr string) (string, int, error) {
	arr := strings.Split(addr, ":")
	if len(arr) != 2 {
		return "", 0, errors.New("malformed addr string")
	}
	port, err := strconv.Atoi(arr[1])
	if err != nil {
		return "", 0, fmt.Errorf("cannot parse port: " + err.Error())
	}
	return arr[0], port, nil
}
