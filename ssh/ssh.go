package ssh

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	Host           string
	User           string
	Password       string
	PrivateKeyPath string
	*ssh.Client
}

func (this *Client) Connect() error {
	conf := ssh.ClientConfig{
		User:            this.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if this.Password != "" {
		conf.Auth = append(conf.Auth, ssh.Password(this.Password))
	} else {
		privateKey, err := getPrivateKey(this.PrivateKeyPath)
		if err != nil {
			return err
		}

		conf.Auth = append(conf.Auth, privateKey)
	}

	client, err := ssh.Dial("tcp", this.Host, &conf)
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}

	this.Client = client

	return nil
}

// Close the connection
func (this *Client) Close() {
	this.Client.Close()
}

// Get the private key for current user
func getPrivateKey(privateKeyPath string) (ssh.AuthMethod, error) {
	if privateKeyPath == "" {
		privateKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa")
	}

	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %v", err)
	}

	return ssh.PublicKeys(signer), nil
}

func CreateTerminalModes() *ssh.TerminalModes {
	return &ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
}
