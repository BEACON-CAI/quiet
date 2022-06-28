package plugins

import (
	"fmt"
	"net"
	"time"

	"quiet/crackpwd/models"
	"quiet/vars"

	"golang.org/x/crypto/ssh"
)

// SSH password cracking
func CrackSsh(s models.Service) (r models.CrackResult, err error) {
	r.Service = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout: time.Duration(vars.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// SSH dial
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.IP, s.Port), config)
	if err != nil {
		return r, err
	}

	// func (c *Client) NewSession() (*Session, error)
	// NewSession opens a new Session for this client.
	// (A session is a remote execution of a program.)
	session, err := client.NewSession()
	if err != nil {
		return r, err

	}
	// func (s *Session) Run(cmd string) error
	// Run runs cmd on the remote host.
	// Typically, the remote server passes cmd to the shell for interpretation.
	// A Session only accepts one call to Run, Start, Shell, Output, or CombinedOutput.
	err = session.Run("echo quiet")
	if err != nil {
		return r, err
	}

	r.Result = true

	defer func() {
		if client != nil {
			_ = client.Close()
		}
		if session != nil {
			_ = session.Close()
		}
	}()

	return r, err
}
