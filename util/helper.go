package util

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
