package traceparser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TraceData struct {
	TraceId   string
	Duration  int
	Name      string
	Timestamp int
	Tags      Tag
}

type Tag struct {
	Error string
}

type Traces []TraceData

type TimeScaleData struct {
	TraceId   string
	Name      string
	Duration  float32
	Timestamp time.Time
	Error     int
	Message   string
}

func RetrieveData(spanName string, content []byte) []TimeScaleData {
	// Let's first read the `config.json` file
	var traceDataList []TimeScaleData
	var payload []Traces
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Println("Error during Unmarshal(): ", err)
		return traceDataList
	}
	for _, traces := range payload {
		for _, trace := range traces {
			if trace.Name == spanName {
				timeScaleData := TimeScaleData{
					Duration:  float32(trace.Duration) / 1000,
					Timestamp: time.Unix(0, int64(trace.Timestamp)*1000),
					Name:      spanName,
					TraceId:   trace.TraceId,
				}
				if trace.Tags.Error != "" {
					timeScaleData.Error = 1
					timeScaleData.Message = trace.Tags.Error
				}
				traceDataList = append(traceDataList, timeScaleData)
				//fmt.Printf("Trace Name: %s\n", trace.Name)
				//fmt.Printf("Trace Timestamp: %d\n\n", trace.Timestamp)
			}
		}
	}
	return traceDataList
}

func InsertData(connStr string, traceData []TimeScaleData) {
	/********************************************/
	/* Connect using Connection Pool            */
	/********************************************/
	ctx := context.Background()
	//connStr := "yourConnectionStringHere"
	dbpool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Generate data to insert
	type result struct {
		Time        time.Time
		SensorId    int
		Temperature float64
		CPU         float64
	}
	/********************************************/
	/* Batch Insert into TimescaleDB            */
	/********************************************/
	//create batch
	batch := &pgx.Batch{}
	numInserts := len(traceData)
	//load insert statements into batch queue
	queryInsertTimeseriesData := `
   INSERT INTO trace_data (timestamp, duration, name, error, trace_id, message) VALUES ($1, $2, $3, $4, $5, $6);
   `

	for i := range traceData {
		var r = traceData[i]
		batch.Queue(queryInsertTimeseriesData, r.Timestamp, r.Duration, r.Name, r.Error, r.TraceId, r.Message)
	}
	batch.Queue("select count(*) from trace_data")
	fmt.Println("Num inserts : ", numInserts)
	//send batch to connection pool
	br := dbpool.SendBatch(ctx, batch)
	//execute statements in batch queue
	_, err = br.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute statement in batch queue %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully batch inserted data")

	//Compare length of results slice to size of table
	fmt.Printf("size of results: %d\n", len(traceData))
	//check size of table for number of rows inserted
	// result of last SELECT statement

	err = br.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to closer batch %v\n", err)
		os.Exit(1)
	}

}
