package tracer

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var tracer trace.Tracer

func InitTracer(enable bool) (trace.Tracer, func()) {
	if !enable {
		tracer = noop.NewTracerProvider().Tracer("") // No-op tracer
		return tracer, func() {}                     // 返回空的清理函数
	}

	// 使用 stdout 作为示例导出器
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logrus.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes("", attribute.String("service.name", "cli"))),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("cli-tracer")

	return tracer, func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}
}

// func main() {
// 	cleanup := initTracer()
// 	defer cleanup() // 程序结束时关闭 TracerProvider

// 	var mysqlAddress string
// 	var enableSSH bool
// 	var sshHosts string
// 	var username string
// 	var password string

// 	rootCmd := &cobra.Command{
// 		Use:   "cli",
// 		Short: "CLI tool for managing MySQL with optional SSH",
// 	}

// 	mysqlCmd := &cobra.Command{
// 		Use:   "mfwmysql",
// 		Short: "Connect to MySQL",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			ctx, span := tracer.Start(context.Background(), "MySQLConnection")
// 			defer span.End()

// 			// 添加一些属性到span
// 			span.SetAttributes(
// 				attribute.String("mysql.address", mysqlAddress),
// 				attribute.Bool("enable_ssh", enableSSH),
// 				attribute.String("ssh_hosts", sshHosts),
// 			)

// 			logrus.Info("Starting MySQL connection...")
// 			fmt.Println("MySQL Address:", mysqlAddress)
// 			fmt.Println("Enable SSH:", enableSSH)
// 			fmt.Println("SSH Hosts:", sshHosts)
// 			fmt.Println("Username:", username)
// 			fmt.Println("Password:", password)

// 			// 模拟连接步骤
// 			simulateMySQLConnection(ctx)
// 		},
// 	}

// 	// 绑定环境变量
// 	viper.AutomaticEnv()
// 	viper.SetEnvPrefix("cli")

// 	mysqlCmd.Flags().StringVarP(&mysqlAddress, "mysql_address", "h", viper.GetString("mysql_address"), "MySQL address")
// 	mysqlCmd.Flags().BoolVarP(&enableSSH, "enable_ssh", "e", viper.GetBool("enable_ssh"), "Enable SSH tunneling")
// 	mysqlCmd.Flags().StringVarP(&sshHosts, "ssh_hosts", "s", viper.GetString("ssh_hosts"), "SSH hosts")
// 	mysqlCmd.Flags().StringVarP(&username, "username", "u", "", "MySQL username")
// 	mysqlCmd.Flags().StringVarP(&password, "password", "p", "", "MySQL password")

// 	viper.BindPFlag("mysql_address", mysqlCmd.Flags().Lookup("mysql_address"))
// 	viper.BindPFlag("enable_ssh", mysqlCmd.Flags().Lookup("enable_ssh"))
// 	viper.BindPFlag("ssh_hosts", mysqlCmd.Flags().Lookup("ssh_hosts"))

// 	rootCmd.AddCommand(mysqlCmd)
// 	if err := rootCmd.Execute(); err != nil {
// 		logrus.WithError(err).Fatal("CLI execution failed")
// 	}
// }

// func simulateMySQLConnection(ctx context.Context) {
// 	_, span := tracer.Start(ctx, "simulateMySQLConnection")
// 	defer span.End()
// 	// 模拟数据库连接的逻辑
// 	span.AddEvent("Connecting to MySQL server")
// }
