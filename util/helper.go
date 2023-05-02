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
	h, p := RemoveProtocol(addr)
	arr := strings.Split(h, ":")
	if len(arr) != 2 && p == -1 ||
		len(arr) < 1 || len(arr) > 2 {
		return "", 0, errors.New("malformed addr string")
	}
	if len(arr) == 2 {
		port, err := strconv.Atoi(arr[1])
		if err != nil {
			return "", 0, fmt.Errorf("cannot parse port: " + err.Error())
		}
		p = port
	}
	return arr[0], p, nil
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
func RemoveProtocol(addr string) (host string, port int) {
	if strings.HasPrefix(addr, "http://") {
		return addr[7:], 80
	} else if strings.HasPrefix(addr, "https://") {
		return addr[8:], 443
	} else if strings.HasPrefix(addr, "ftp://") {
		return addr[6:], 21
	} else if strings.HasPrefix(addr, "socks5://") {
		return addr[9:], 1080
	} else if strings.HasPrefix(addr, "socks5h://") {
		return addr[10:], 1080
	} else {
		return addr, -1
	}
}
