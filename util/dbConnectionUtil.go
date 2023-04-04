package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// CreateConnection creates database connection using pgxpool
func CreateConnection() *pgxpool.Pool {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://"+dbUsername+":"+dbPassword+"@"+dbHost+":"+strconv.Itoa(dbPort)+"/"+dbName+"?sslmode=disable")
	// dbpool, err := pgxpool.Connect(context.Background(), "postgres://ykakirfxfsrqyr:5225c8b97bf4ce844ae9fd6e49775d84724504b84da0fc6c898e5cef5d8b1aac@ec2-35-172-73-125.compute-1.amazonaws.com:5432/d3fkc9srt7jkbk")
	// dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:msaha@localhost:5432/satcom?sslmode=disable")
	// dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:omur@localhost:5432/satcom?sslmode=disable")
	// dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/satcom?sslmode=disable")
	// dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:satcom@localhost:5432/satcom?sslmode=disable")
	if err != nil {
		fmt.Println("Unable to connect to database: ")
		fmt.Println(err)
		// os.Exit(1)
	}
	return dbpool
}

// CreatePayrollConnectionUsingGorm creates database connection using gorm
func CreatePayrollConnectionUsingGorm() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
		},
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ViperReturnStringConfigVariableFromLocalConfigJSON("payroll_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// panic("failed to connect database")
		panic(err)
	} else {
		return db
	}
}

// CreateConnectionUsingGorm2 creates database connection using gorm
func CreateConnectionUsingGorm2() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ViperReturnStringConfigVariableFromLocalConfigJSON("advertisement_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// panic("failed to connect database")
		panic(err)
	} else {
		return db
	}
}

// CreateConnectionUsingGorm creates database connection using gorm
func CreateConnectionUsingGorm222222() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "test" + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// panic("failed to connect database")
		panic(err)
	} else {
		return db
	}
}

// ViperReturnStringConfigVariableFromLocalConfigJSON returns values of string variable from local-config.json
func ViperReturnStringConfigVariableFromLocalConfigJSON(key string) string {
	// viper.SetConfigFile("local-config.json")
	viper.SetConfigName("local-config") // name of config file (without extension)
	viper.SetConfigType("json")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config")     // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	// viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
		return ""
	}
	return value
}

// ViperReturnIntegerConfigVariableFromLocalConfigJSON returns values of int variable from local-config.json
func ViperReturnIntegerConfigVariableFromLocalConfigJSON(key string) int {
	// viper.SetConfigFile("local-config.json")
	viper.SetConfigName("local-config") // name of config file (without extension)
	viper.SetConfigType("json")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config")     // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value := viper.GetInt(key)
	return value
}

// CreateConnectionUsingGorm2 creates database connection using gorm
func CreateConnectionToCirculationSchemaUsingGorm() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
		},
	)
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ViperReturnStringConfigVariableFromLocalConfigJSON("circulation_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// panic("failed to connect database")
		panic(err)
	} else {
		return db
	}
}

// Create Connection To Accounts Schema Using Gorm
func CreateConnectionToAccountsSchemaUsingGorm() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
		},
	)
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ViperReturnStringConfigVariableFromLocalConfigJSON("accounts_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// panic("failed to connect database")
		panic(err)
	} else {
		return db
	}
}

// Create Connection To Accounts Schema Using Gorm

// CreatePayrollConnectionUsingGorm creates database connection using gorm
func CreateConnectionToPFAccountsSchemaUsingGorm() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
		},
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ViperReturnStringConfigVariableFromLocalConfigJSON("pfaccounts_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// panic("failed to connect database")
		panic(err)
	} else {
		return db
	}
}
