package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/juzeon/lip/util"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var packetCmd = &cobra.Command{
	Use:   "packet <host:port or host port> [content]",
	Short: "Send TCP/UDP packets. Default: TCP",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(3)),
	Run: func(cmd *cobra.Command, args []string) {
		var host string
		var port int
		var content string
		if len(args) == 1 {
			if !packetFlags.Interactive {
				log.Fatalln("non-interactive mode requires a content argument")
			}
			h, p, err := util.ExtractHostPort(args[0])
			if err != nil {
				log.Fatalln("cannot extract host and port: " + err.Error())
			}
			host = h
			port = p
		} else if len(args) == 2 {
			if packetFlags.Interactive {
				host, _ = util.RemoveProtocol(args[0])
				p, err := strconv.Atoi(args[1])
				if err != nil {
					log.Fatalln("cannot parse port: " + err.Error())
				}
				port = p
			} else {
				h, p, err := util.ExtractHostPort(args[0])
				if err != nil {
					log.Fatalln("cannot extract host and port: " + err.Error())
				}
				host = h
				port = p
				content = args[1]
			}
		} else {
			if packetFlags.Interactive {
				log.Fatalln("interactive mode does not accept a content argument")
			}
			host, _ = util.RemoveProtocol(args[0])
			p, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalln("cannot parse port: " + err.Error())
			}
			port = p
			content = args[2]
		}
		addr := host + ":" + strconv.Itoa(port)
		dialer, err := util.GetProxyDialer(persistentFlags.Proxy, time.Duration(packetFlags.Timeout)*time.Second)
		if err != nil {
			log.Fatalln("cannot get dialer: " + err.Error())
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(packetFlags.Timeout)*time.Second)
		defer cancel()
		conn, err := dialer.(proxy.ContextDialer).
			DialContext(ctx, util.Ternary(packetFlags.UDP, "udp", "tcp"), addr)
		if err != nil {
			log.Fatalln("network dial error: " + err.Error())
		}
		defer conn.Close()
		go func() {
			for {
				b := make([]byte, 1024*100)
				n, err := conn.Read(b)
				if err != nil && !errors.Is(err, io.EOF) {
					log.Println("read conn error: " + err.Error())
					os.Exit(1)
				}
				if err != nil {
					os.Exit(0)
				}
				_, err = os.Stdout.Write(b[:n])
				if err != nil {
					log.Println("write to stdout error: " + err.Error())
				}
			}
		}()
		if packetFlags.Interactive {
			fmt.Println("Interactive mode. Use enter to send content followed by a predefined linebreak.")
			for {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				packetWrite(strings.TrimSpace(text), conn)
			}
		} else {
			packetWrite(content, conn)
			time.Sleep(time.Duration(packetFlags.Timeout) * time.Second)
		}
	},
}

func packetWrite(v string, writer io.Writer) {
	v = strings.ReplaceAll(v, "\\r", "\r")
	v = strings.ReplaceAll(v, "\\n", "\n")
	_, err := writer.Write([]byte(v))
	if err != nil {
		log.Println("" + err.Error())
	}
	linebreak := packetFlags.GetRealLinebreak()
	if linebreak != "" {
		_, err = io.WriteString(writer, linebreak)
		if err != nil {
			log.Println("" + err.Error())
		}
	}
}

type packetFlagStruct struct {
	UDP         bool
	Interactive bool
	Timeout     int
	Linebreak   string
}

func (p packetFlagStruct) GetRealLinebreak() string {
	content := strings.ReplaceAll(p.Linebreak, "\\r", "\r")
	return strings.ReplaceAll(content, "\\n", "\n")
}

var packetFlags = packetFlagStruct{}

func init() {
	rootCmd.AddCommand(packetCmd)

	packetCmd.Flags().BoolVarP(&packetFlags.UDP, "udp", "u", false,
		"use UDP instead of TCP")
	packetCmd.Flags().BoolVarP(&packetFlags.Interactive, "interactive", "i", false,
		"enter interactive mode. Note that this mode does not accept a content argument")
	packetCmd.Flags().IntVarP(&packetFlags.Timeout, "timeout", "t", 5,
		"n seconds timeout for network dial and data read")
	packetCmd.Flags().StringVarP(&packetFlags.Linebreak, "linebreak", "b", "\r\n",
		"used for linebreak automatically added to the end of content. "+
			"Pass an empty string to disable auto linebreak")
}
