package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Client struct {
	Host           string
	User           string
	Password       string
	PrivateKeyPath string
	*ssh.Client
}

func (sshClient *Client) Connect() error {

	conf := ssh.ClientConfig{
		User:            sshClient.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if sshClient.Password != "" {
		conf.Auth = append(conf.Auth, ssh.Password(sshClient.Password))
	} else if sshClient.PrivateKeyPath != "" {
		privateKey, err := getPrivateKey(sshClient.PrivateKeyPath)
		if err != nil {
			return err
		}

		conf.Auth = append(conf.Auth, privateKey)
	} else {
		// if occur error "Failed to open SSH_AUTH_SOCK: dial unix: missing address",
		// execute command: eval `ssh-agent`,and enter passphrase
		conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			log.Fatalf("Failed to open SSH_AUTH_SOCK: %v", err)
		}

		agentClient := agent.NewClient(conn)
		// Use a callback rather than PublicKeys so we only consult the
		// agent once the remote server wants it.
		conf.Auth = append(conf.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}
	client, err := ssh.Dial("tcp", sshClient.Host, &conf)

	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}

	sshClient.Client = client

	return nil
}

// Close the connection
func (sshClient *Client) Close() {
	sshClient.Client.Close()
}

// Get the private key for current user
func getPrivateKey(privateKeyPath string) (ssh.AuthMethod, error) {
	if !fileExist(privateKeyPath) {
		defaultPrivateKeyPath := filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa")
		log.Printf("Warning: private key path [%s] does not exist, using default %s instead", privateKeyPath, defaultPrivateKeyPath)

		privateKeyPath = defaultPrivateKeyPath
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

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}
