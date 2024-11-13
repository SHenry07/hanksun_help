package ssh

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

type DatabaseCreds struct {
	SSHHost     string // SSH Server Hostname/IP
	SSHPort     int    // SSH Port
	SSHUser     string // SSH Username
	SSHPassword string // SSH Password
	SSHKeyFile  string // SSH Key file location
	DBUser      string // DB username
	DBPass      string // DB Password
	DBHost      string // DB Hostname/IP
	DBName      string // Database name
}

// func main() {
// 	db, sshConn, err := ConnectToDB(DatabaseCreds{
// 		SSHHost:    "123.123.123.123",
// 		SSHPort:    22,
// 		SSHUser:    "root",
// 		SSHKeyFile: "sshkeyfile.pem",
// 		DBUser:     "root",
// 		DBPass:     "password",
// 		DBHost:     "localhost:3306",
// 		DBName:     "dname",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer sshConn.Close()
// 	defer db.Close()

// 	if rows, err := db.Query("SELECT 1=1"); err == nil {
// 		for rows.Next() {
// 			var result string
// 			rows.Scan(&result)
// 			fmt.Printf("Result: %s\n", result)
// 		}
// 		rows.Close()
// 	} else {
// 		fmt.Printf("Failure: %s", err.Error())
// 	}
// }

// ConnectToDB will accept the db and ssh credientials (DatabaseCreds) and
// form a connection with the database (handling any errors that might arise).
func ConnectToDB(dbCreds DatabaseCreds) (*sql.DB, *ssh.Client, error) {
	// Make SSH client: establish a connection to the local ssh-agent
	var agentClient agent.Agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close()
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            dbCreds.SSHUser,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if dbCreds.SSHKeyFile != "" {
		pemBytes, err := os.ReadFile(dbCreds.SSHKeyFile)
		if err != nil {
			return nil, nil, err
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, nil, err
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))

	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	if dbCreds.SSHPassword != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return dbCreds.SSHPassword, nil
		}))
	}

	// Connect to the SSH Server
	sshConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", dbCreds.SSHHost, dbCreds.SSHPort), sshConfig)
	if err != nil {
		return nil, nil, err
	}

	// Now we register the ViaSSHDialer with the ssh connection as a parameter
	mysql.RegisterDialContext("mysql+tcp", func(_ context.Context, addr string) (net.Conn, error) {
		dialer := &ViaSSHDialer{sshConn}
		return dialer.Dial(addr)
	})

	// And now we can use our new driver with the regular mysql connection string tunneled through the SSH connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", dbCreds.DBUser, dbCreds.DBPass, dbCreds.DBHost, dbCreds.DBName))
	if err != nil {
		return nil, sshConn, err
	}

	log.Println("Successfully connected to the db")

	return db, sshConn, err
}
