package main

import (
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var pottid *string

// Crypto algs
var kex []string
var enc []string
var mac []string

// Handler for password auth, return true to be able to send data to client
func passauthHandler(ctx ssh.Context, password string) bool {
	constr := strings.Split(ctx.RemoteAddr().String(), ":")
	if len(constr) > 1 {
		fmt.Printf("{\"timestamp\": %q, \"id\": %q, \"user\": %q, \"clientip\": %q, \"srcport\": %q, \"password\": %q, \"clientversion\": %q }\n", time.Now().Format(time.RFC3339), *pottid, ctx.User(), constr[0], constr[1], password, ctx.ClientVersion())
	}
	return true
}

// Handler for publickey auth, always return false to proceed to passwordauth
func pubauthHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	authorizedKey := gossh.MarshalAuthorizedKey(key)
	authorizedKey = authorizedKey[:len(authorizedKey)-1]
	constr := strings.Split(ctx.RemoteAddr().String(), ":")
	if len(constr) > 1 {
		fmt.Printf("{\"timestamp\": %q, \"id\": %q, \"user\": %q, \"clientip\": %q, \"srcport\": %q, \"publickey\": %q, \"clientversion\": %q }\n", time.Now().Format(time.RFC3339), *pottid, ctx.User(), constr[0], constr[1], authorizedKey, ctx.ClientVersion())
	}
	return false
}

// Callback function to set custom server settings
func serverconfigHandler(ctx ssh.Context) *gossh.ServerConfig {
	var sconf *gossh.ServerConfig
	sconf = new(gossh.ServerConfig)
	sconf.Config.KeyExchanges = kex
	sconf.Config.Ciphers = enc
	sconf.Config.MACs = mac
	return sconf
}


func main() {

	// Custom usage func
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "potter ssh honeypot.\n\nUsage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n\nExample:\npotter -p 9000 -s \"OpenSSH_8.9\" -k ./potter.key -m \"Go away\" -i mysensor\n\n")
	}

	// Get commandline options
	ssh_port := flag.String("p", "2222", "ssh port")
	ssh_string := flag.String("s", "OpenSSH_4.5", "ssh versionstring")
	ssh_ip := flag.String("l", "0.0.0.0", "listen address")
	ssh_hostkey := flag.String("k", "potter.key", "hostkey filename")
	denymsg := flag.String("m", "Access denied.", "Custom deny message")
	pottid = flag.String("i", "ssh-pott", "custom tag in JSON")
	// Options for custom alg.

	// See https://cs.opensource.google/go/x/crypto/+/master:ssh/kex.go for valid values
	kex_flag := flag.String("kex","curve25519-sha256@libssh.org","Custom kex algs.")

	// https://cs.opensource.google/go/x/crypto/+/master:ssh/keys.go
	enc_flag := flag.String("enc","aes256-ctr,aes128-ctr","Custom enc. algs.")
	
	// https://cs.opensource.google/go/x/crypto/+/master:ssh/mac.go
	mac_flag := flag.String("mac","hmac-sha2-256-etm@openssh.com,hmac-sha2-512-etm@openssh.com","Custom mac. algs.")
	flag.Parse()

	// Set custom algs
	kex = strings.Split(*kex_flag,",")
	enc = strings.Split(*enc_flag,",")
	mac = strings.Split(*mac_flag,",")

	// Open hostkey file
	keyb, err := ioutil.ReadFile(*ssh_hostkey)
	if err != nil {
		log.Fatalf("Failed to open hostkey: %s.", *ssh_hostkey)
	}
	hkey, err := gossh.ParsePrivateKey(keyb)
	if err != nil {
		log.Fatal("Failed to parse hostkey.")
	}

	// SSH server config
	s := &ssh.Server{
		Addr:                 *ssh_ip + ":" + *ssh_port,
		Version:              *ssh_string,
		IdleTimeout:          time.Duration(10 * time.Second),
		MaxTimeout:           time.Duration(15 * time.Second),
		PasswordHandler:      passauthHandler,
		PublicKeyHandler:     pubauthHandler,
		ServerConfigCallback: serverconfigHandler,
	}

	// Add the hostkey to config
	s.AddHostKey(hkey)

	// Handler to write output to accepted SSH sessions
	s.Handle(func(s ssh.Session) {
		io.WriteString(s, *denymsg+"\n")
		return
	})

	// Start the SSH server
	log.Printf("Potter SSH honeypot starting. port: %q, ip: %q, sshstring: %q\n", *ssh_port, *ssh_ip, *ssh_string)
	log.Fatal(s.ListenAndServe())
}
