package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var tcpingCmd = &cobra.Command{
	Use:   "tcping <1.1.1.1:443 or 1.1.1.1 443>",
	Short: "TCPing an address with a port",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		var host string
		var port int
		if len(args) == 1 {
			h, p, err := extractHostPort(args[0])
			if err != nil {
				log.Fatalln("malformed argument: " + err.Error())
			}
			host = h
			port = p
		} else {
			host = args[0]
			p, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalln("malformed arguments: " + err.Error())
			}
			port = p
		}
		address := host + ":" + strconv.Itoa(port)
		if tcpingFlags.Continuously {
			for {
				res := doTcpingOnce(address)
				res.WriteOut(os.Stdout)
				if tcpingFlags.ExitOnce {
					break
				}
			}
		} else if tcpingFlags.Second == 0 {
			for i := 1; i <= tcpingFlags.Count; i++ {
				res := doTcpingOnce(address)
				res.WriteOut(os.Stdout)
				if tcpingFlags.ExitOnce {
					break
				}
			}
		} else {
			t := time.Now()
			for {
				res := doTcpingOnce(address)
				res.WriteOut(os.Stdout)
				if time.Since(t) > time.Duration(tcpingFlags.Second)*time.Second || tcpingFlags.ExitOnce {
					break
				}
			}
		}
	},
}

type tcpingFlagStruct struct {
	Continuously bool
	Count        int
	Second       int
	Timeout      int
	ExitOnce     bool
	Wait         float64
}

var tcpingFlags = tcpingFlagStruct{}

func init() {
	rootCmd.AddCommand(tcpingCmd)

	tcpingCmd.Flags().BoolVarP(&tcpingFlags.Continuously, "continuously", "t", false,
		"ping continuously until stopped via control-c")
	tcpingCmd.Flags().IntVarP(&tcpingFlags.Count, "count", "c", 4,
		"send pings for n times")
	tcpingCmd.Flags().IntVarP(&tcpingFlags.Second, "second", "s", 0,
		"send pings for n seconds. Overwrite `count` if exists")
	tcpingCmd.Flags().IntVarP(&tcpingFlags.Timeout, "timeout", "u", 5,
		"connection timeout for every ping")
	tcpingCmd.Flags().BoolVarP(&tcpingFlags.ExitOnce, "exit", "e", false,
		"automatically exit on a successful ping")
	tcpingCmd.Flags().Float64VarP(&tcpingFlags.Wait, "wait", "w", 0.5,
		"wait n(float64) seconds between pings")
}

func extractHostPort(addr string) (string, int, error) {
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

const (
	tcpingOpen = iota
	tcpingNoResponse
	tcpingRST
	tcpingUnknown
)

type tcpingResult struct {
	Addr string
	Num  int
	Dur  time.Duration
}

func (t tcpingResult) GetStatus() string {
	var tcpingStatus = map[int]string{
		tcpingOpen:       "Port is open (open)",
		tcpingNoResponse: "No response (filtered)",
		tcpingRST:        "Respond with RST (closed)",
		tcpingUnknown:    "Unknown error",
	}
	return tcpingStatus[t.Num]
}
func (t tcpingResult) WriteOut(writer io.Writer) {
	_, err := fmt.Fprintln(writer, "Probing "+t.Addr+"/tcp - "+t.GetStatus()+" - time="+t.Dur.String())
	if err != nil {
		log.Println("cannot write tcping result: " + err.Error())
	}
}

func doTcpingOnce(address string) tcpingResult {
	t := time.Now()
	conn, err := net.DialTimeout("tcp", address, time.Duration(tcpingFlags.Timeout)*time.Second)
	if err != nil {
		if strings.Contains(err.Error(), "actively refused") ||
			strings.Contains(err.Error(), "connection refused") {
			return tcpingResult{Num: tcpingRST, Dur: time.Since(t), Addr: address}
		} else if strings.Contains(err.Error(), "timeout") {
			return tcpingResult{Num: tcpingNoResponse, Dur: time.Since(t), Addr: address}
		} else {
			return tcpingResult{Num: tcpingUnknown, Dur: time.Since(t), Addr: address}
		}
	}
	defer time.Sleep(time.Duration(tcpingFlags.Wait * float64(time.Second)))
	defer conn.Close()
	return tcpingResult{Num: tcpingOpen, Dur: time.Since(t), Addr: address}
}
