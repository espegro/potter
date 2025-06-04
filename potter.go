package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var pottid *string

// Crypto algs
var kex []string
var enc []string
var mac []string

// Function to calculate HASSH https://github.com/salesforce/hassh
func calculateHASSH(algorithms []string) string {
	concatenated := strings.Join(algorithms, ";")
	hash := md5.Sum([]byte(concatenated)) // Use MD5
	return hex.EncodeToString(hash[:])
}

// Handler for password auth, return true to be able to send data to client
func passauthHandler(ctx ssh.Context, password string) bool {
	host, port, err := net.SplitHostPort(ctx.RemoteAddr().String())
	if err == nil {
		hassh := calculateHASSH([]string{kex[0], enc[0], mac[0]}) // Simplified, consider all algorithms
		fmt.Printf("{\"timestamp\": %q, \"id\": %q, \"user\": %q, \"clientip\": %q, \"srcport\": %q, \"password\": %q, \"clientversion\": %q, \"hassh\": %q }\n", time.Now().Format(time.RFC3339), *pottid, ctx.User(), host, port, password, ctx.ClientVersion(), hassh)
	}
	return true
}

// Handler for publickey auth, always return false to proceed to passwordauth
func pubauthHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	authorizedKey := gossh.MarshalAuthorizedKey(key)
	authorizedKey = authorizedKey[:len(authorizedKey)-1]
	host, port, err := net.SplitHostPort(ctx.RemoteAddr().String())
	if err == nil {
		hassh := calculateHASSH([]string{kex[0], enc[0], mac[0]}) // Simplified, consider all algorithms
		fmt.Printf("{\"timestamp\": %q, \"id\": %q, \"user\": %q, \"clientip\": %q, \"srcport\": %q, \"publickey\": %q, \"clientversion\": %q, \"hassh\": %q }\n", time.Now().Format(time.RFC3339), *pottid, ctx.User(), host, port, authorizedKey, ctx.ClientVersion(), hassh)
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
	kex_flag := flag.String("kex", "curve25519-sha256@libssh.org", "Custom kex algs.")
	enc_flag := flag.String("enc", "aes256-ctr,aes128-ctr", "Custom enc. algs.")
	mac_flag := flag.String("mac", "hmac-sha2-256-etm@openssh.com,hmac-sha2-512-etm@openssh.com", "Custom mac. algs.")
	flag.Parse()

	// Set custom algs
	kex = strings.Split(*kex_flag, ",")
	enc = strings.Split(*enc_flag, ",")
	mac = strings.Split(*mac_flag, ",")

	// Open hostkey file
	keyb, err := os.ReadFile(*ssh_hostkey)
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
