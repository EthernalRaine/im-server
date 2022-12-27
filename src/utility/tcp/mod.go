package tcp

import (
	"fmt"
	"net"
	"os"
	"phantom/utility"
	"phantom/utility/logging"
	"time"
)

type TcpConnection struct {
	server net.Listener
	client net.Conn
}

func CreateListener(port int) TcpConnection {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		logging.Fatal("TCP Listener", "Failed to start listener! (%s)", err.Error())
		os.Exit(0)
	}

	logging.Info("TCP Listener", "Listening on 0.0.0.0:%d", port)

	conn := TcpConnection{
		server: tcpListener,
	}

	return conn
}

func (tcp *TcpConnection) WriteTraffic(data string) error {
	logging.Trace("TCP/WriteTraffic", "Writing Data: %s", utility.SanitizeString(data))
	_, err := tcp.client.Write([]byte(data))
	return err
}

func (tcp *TcpConnection) ReadTraffic() (data string, err error) {
	return tcp.ExReadTraffic(0)
}

func (tcp *TcpConnection) ExReadTraffic(timeout int) (data string, err error) {
	tcp.client.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))

	buf := make([]byte, 4096)
	_, err = tcp.client.Read(buf)

	if err != nil {
		logging.Error("TCP/ReadTraffic", "Failed to read traffic! (%s)", err.Error())
		return string(buf), err
	}

	ret := utility.SanitizeString(string(buf))

	logging.Trace("TCP/ReadTraffic", "Reading Data: %s", ret)

	return ret, err
}