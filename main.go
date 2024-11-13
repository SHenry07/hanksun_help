package main

import (
	"context"
	"fmt"

	"jsdf-helper/utils/log"
	"jsdf-helper/utils/tracer"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var enableTrace bool

func main() {
	// init logger

	var mysqlAddress string
	var enableSSH bool
	var sshHosts string
	var username string
	var password string
	var debug bool
	var traceFlag bool

	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI tool for managing MySQL with optional SSH",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 初始化日志
			log.InitLogger(debug)

			// 初始化 Tracer
			tracer, cleanup := tracer.InitTracer(traceFlag)
			defer cleanup()
		},
	}

	mysqlCmd := &cobra.Command{
		Use:   "mfwmysql",
		Short: "Connect to MySQL",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, span := tracer.Start(context.Background(), "MySQLConnection")
			defer span.End()

			logrus.WithFields(logrus.Fields{
				"mysql_address": mysqlAddress,
				"enable_ssh":    enableSSH,
				"ssh_hosts":     sshHosts,
			}).Debug("Debug information: starting MySQL connection")

			fmt.Println("MySQL Address:", mysqlAddress)
			fmt.Println("Enable SSH:", enableSSH)
			fmt.Println("SSH Hosts:", sshHosts)
			fmt.Println("Username:", username)
			fmt.Println("Password:", password)

			simulateMySQLConnection(ctx)
		},
	}

	// 设置全局标志
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&traceFlag, "trace", false, "Enable tracing")

	// 设置命令标志
	mysqlCmd.Flags().StringVarP(&mysqlAddress, "mysql_address", "h", viper.GetString("mysql_address"), "MySQL address")
	mysqlCmd.Flags().BoolVarP(&enableSSH, "enable_ssh", "e", viper.GetBool("enable_ssh"), "Enable SSH tunneling")
	mysqlCmd.Flags().StringVarP(&sshHosts, "ssh_hosts", "s", viper.GetString("ssh_hosts"), "SSH hosts")
	mysqlCmd.Flags().StringVarP(&username, "username", "u", "", "MySQL username")
	mysqlCmd.Flags().StringVarP(&password, "password", "p", "", "MySQL password")

	viper.BindPFlag("mysql_address", mysqlCmd.Flags().Lookup("mysql_address"))
	viper.BindPFlag("enable_ssh", mysqlCmd.Flags().Lookup("enable_ssh"))
	viper.BindPFlag("ssh_hosts", mysqlCmd.Flags().Lookup("ssh_hosts"))

	rootCmd.AddCommand(mysqlCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("CLI execution failed")
	}
}

func simulateMySQLConnection(ctx context.Context) {
	_, span := tracer.Start(ctx, "simulateMySQLConnection")
	defer span.End()
	span.AddEvent("Connecting to MySQL server")
	logrus.Debug("Simulated MySQL connection established")
}
