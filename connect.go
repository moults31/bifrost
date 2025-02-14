package bifrost

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/pkg/term"
)

type Connect struct {
	portPath   string
	baudRate   int
	port       *term.Term
	portReader *bufio.Reader
	stateChan  chan error
}

// NewConnection returns a pointer to a Connect instance
func NewConnection(portPath string, baudRate int) (*Connect, error) {
	t := term.Speed(baudRate)
	port, err := term.Open(portPath, t)
	if err != nil {
		return nil, err
	}
	port.SetRaw()
	portReader := bufio.NewReader(port)
	stateChan := make(chan error)
	return &Connect{portPath: portPath, baudRate: baudRate, port: port,
		portReader: portReader,
		stateChan:  stateChan}, nil
}

// Start initializes a read loop that attempts to reconnect
// when the connection is broken
func (c *Connect) Start() {
	go c.read()
	for {
		select {
		case err := <-c.stateChan:
			if err != nil {
				fmt.Printf("Error connecting to %s", c.portPath)
				go c.initialize()
			} else {
				fmt.Printf(" | Connection to %s reestablished!", c.portPath)
				go c.read()
			}
		}
	}
}

func (c *Connect) initialize() {
	c.port.Close()
	for {
		time.Sleep(time.Second)
		port, err := term.Open(c.portPath)
		if err != nil {
			continue
		}
		c.port = port
		c.portReader = bufio.NewReader(port)
		c.stateChan <- nil
		return
	}
}

func (c *Connect) read() {
	buf := make([]byte, 256)

	for {
		n, err := c.portReader.Read(buf)
		// report the error
		if err != nil && err != io.EOF {
			c.stateChan <- err
			return
		}
		if n > 0 {
			fmt.Print(string(buf[:n]))
		}
	}
}

func (c *Connect) Write(message []byte) {
	_, err := c.port.Write(message)
	if err != nil {
		fmt.Printf("Error writing to serial port: %v ", err)
	}
}
