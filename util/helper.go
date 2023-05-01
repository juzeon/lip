package util

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
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
func WriteTable(matrix [][]string, writer io.Writer, reverse bool) {
	firstRowAsHeader := true
	table := tablewriter.NewWriter(writer)
	if reverse {
		matrix = TransposeMatrix(matrix)
		firstRowAsHeader = false
	}
	if firstRowAsHeader {
		table.SetHeader(matrix[0])
		table.AppendBulk(matrix[1:])
	} else {
		table.AppendBulk(matrix)
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
}
func Map[T any, E any](arr []T, function func(value T) E) []E {
	var result []E
	for _, item := range arr {
		result = append(result, function(item))
	}
	return result
}
func FilterKeep[T any](arr []T, function func(value T) bool) []T {
	var result []T
	for _, item := range arr {
		if function(item) {
			result = append(result, item)
		}
	}
	return result
}

// GetProxyDialer returns a proxy.Dialer with timeout.
// Note that we should convert it into proxy.ContextDialer to ensure a timeout.
func GetProxyDialer(proxyStringOrEmpty string, uncertainTimeout time.Duration) (proxy.Dialer, error) {
	var dialer proxy.Dialer
	if proxyStringOrEmpty != "" {
		proxyURL, err := url.Parse(proxyStringOrEmpty)
		if err != nil {
			return &net.Dialer{}, err
		}
		dialer, err = proxy.FromURL(proxyURL, &net.Dialer{Timeout: uncertainTimeout})
		if err != nil {
			return &net.Dialer{}, err
		}
	} else {
		dialer = proxy.FromEnvironment()
	}
	return dialer, nil
}
