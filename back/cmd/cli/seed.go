package main

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"

	"gpe_project/internal/app/adapter/postgresql"
	"gpe_project/internal/app/adapter/service"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("fatal error during dir finding: %w", err))
	}

	service.LoadConfig(dir)

	// Init database connection
	db := service.GetPostgresqlDB()

	seeder := postgresql.NewSeeder(db)

	flags()

	switch viper.GetString("seed") {
	case "notification":
		fmt.Println("seeding notifications")
		seeder.PurgeNotifications()
		seeder.AddUserNotifications(1)
	case "barometer":
		fmt.Println("seeding barometers")
		seeder.PurgeBarometers()
		seeder.AddBarometers(2)
	}
	fmt.Println("done")
}

func flags() {
	flag.String("seed", "", "Launch seed action (notification, barometers)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println("cannot bind flags")
		os.Exit(1)
	}
}
