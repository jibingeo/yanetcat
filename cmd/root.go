package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/jibingeo/yanetcat/pipe"
	"github.com/spf13/cobra"
)

var listen, target string
var listenAddr, targetAddr *url.URL

//Parse uri to to url.URL
func Parse(address string) (*url.URL, error) {
	tmp, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme: tmp.Scheme,
		Host:   tmp.Host + tmp.Path,
	}, nil
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&listen, "listen", "l", "", "listen address")
	RootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target address")
}

func handleRequest(lConn net.Conn) {
	log.Println("Connection opened:", lConn.RemoteAddr().String())
	defer lConn.Close()
	defer func() {
		log.Println("Connection closed:", lConn.RemoteAddr().String())
	}()

	rConn, err := net.Dial(targetAddr.Scheme, targetAddr.Host)
	if err != nil {
		log.Println("Docker API:", err)
		return
	}
	defer rConn.Close()
	pipe.Pipe(lConn, rConn)
}

//RootCmd is yanetcat main Root Command
var RootCmd = &cobra.Command{
	Use:   "yanetcat",
	Short: "Yet Another netcat",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if listen == "" {
			return errors.New("listen address missing")
		}
		if target == "" {
			return errors.New("target address missing")
		}
		_listenAddr, err := Parse(listen)
		if err != nil {
			return err
		}
		_targetAddr, err := Parse(target)
		if err != nil {
			return err
		}
		listenAddr, targetAddr = _listenAddr, _targetAddr
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		listener, err := net.Listen(listenAddr.Scheme, listenAddr.Host)
		if err != nil {
			log.Fatalln("Error listening:", err.Error())
		}
		log.Println("Started listening:", listener.Addr().String())
		defer listener.Close()

		for {
			// Listen for an incoming connection.
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				continue
			}
			go handleRequest(conn)
		}
	},
}
