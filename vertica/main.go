package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/vertica/vertica-sql-go"
	"github.com/vertica/vertica-sql-go/logger"
)

func main() {
	// Have our logger output INFO and above.
	logger.SetLogLevel(logger.INFO)

	var testLogger = logger.New("samplecode")

	ctx := context.Background()

	// Create a connection to our database. Connection is lazy and won't
	// happen until it's used.
	connDB, err := sql.Open("vertica", "vertica://dbadmin:@localhost:5433/dbadmin")

	if err != nil {
		testLogger.Fatal(err.Error())
		os.Exit(1)
	}

	defer connDB.Close()

	// Ping the database connnection to force it to attempt to connect.
	if err = connDB.PingContext(ctx); err != nil {
		testLogger.Fatal(err.Error())
		os.Exit(1)
	}

	// Query a standard metric table in Vertica.
	rows, err := connDB.QueryContext(ctx, "SELECT * FROM v_monitor.cpu_usage LIMIT 5")

	if err != nil {
		testLogger.Fatal(err.Error())
		os.Exit(1)
	}

	defer rows.Close()

	// Iterate over the results and print them out.
	for rows.Next() {
		var nodeName string
		var startTime string
		var endTime string
		var avgCPU float64

		if err = rows.Scan(&nodeName, &startTime, &endTime, &avgCPU); err != nil {
			testLogger.Fatal(err.Error())
			os.Exit(1)
		}

		testLogger.Info("%s\t%s\t%s\t%f", nodeName, startTime, endTime, avgCPU)
	}

	testLogger.Info("Test complete")

	os.Exit(0)
}
