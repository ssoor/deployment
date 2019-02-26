package main

import (
	"golang.org/x/crypto/ssh"
)

func SSH(user, password, ip_port string) (*ssh.Client, error) {
	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{User: user, Auth: PassWd}

	return ssh.Dial("tcp", ip_port, &Conf)
}
