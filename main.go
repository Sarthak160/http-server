package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	companyID string = "testCompany"
	startTime int64  = 1742479649
	endTime   int64  = 1745071649
)

func main() {
	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "createdBy.cid", Value: companyID},
				{Key: "createdAt", Value: bson.D{
					{Key: "$gte", Value: startTime},
					{Key: "$lte", Value: endTime},
				}},
				{Key: "type", Value: "prMetrics"},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalCovered", Value: bson.D{
					{Key: "$sum", Value: "$eventDetails.coveredByKeploy"},
				}},
			}},
		},
	}

	// Configure MongoDB client options with retryWrites and retryReads enabled
	clientOptions := options.Client().
		ApplyURI("mongodb://staging:itAaK8L8QPonPLG4dOraYGAPM8u4iwHIykuOMPPLRCMI848Xwgu3feutla3sqAdvKMRjm4dZUPP7ACDbyrw3PA%3D%3D@staging.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=true&maxIdleTimeMS=120000&appName=@staging@").
		SetRetryWrites(false).
		SetRetryReads(true)

	ctx := context.Background()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// Create a new GraphQL server
	dbName := "api-server"
	if dbName == "" {
		dbName = "keploy"
	}
	mongoDb := client.Database(dbName)
	event := mongoDb.Collection("eventsV2")
	// populate(event)
	cursor, err := event.Aggregate(ctx, pipeline)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			panic(err)
		}
	}()

	var result struct {
		TotalCovered int `bson:"totalCovered"`
	}

	if cursor.Next(ctx) {
		println("Cursor found")
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}
		print("Total Covered: ", result.TotalCovered)
	}

}

type EventType string

type CreatorInfo struct {
	Email     string `json:"email" bson:"email"`
	UserId    string `json:"userId" bson:"userId"`
	CompanyId string `json:"cid" bson:"cid"`
}

const (
	PRMetricsType       EventType = "prMetrics"
	ListDirectoriesType EventType = "listDirectories"
	// add more events to track if any
)

type Event[T any] struct {
	ID        string      `json:"id" bson:"_id"`
	CreatedAt int64       `json:"createdAt" bson:"createdAt"`
	SessionID string      `json:"sessionId" bson:"sessionId"`
	CreatedBy CreatorInfo `json:"createdBy" bson:"createdBy"`
	Module    ModuleName  `json:"module" bson:"module"`
	Action    string      `json:"action" bson:"action"`
	Data      T           `json:"eventDetails" bson:"eventDetails"`
	IP        string      `json:"ip" bson:"ip"`
}
type ModuleName string

const (
	PRAgentModule ModuleName = "prAgent"
	VSCodeModule  ModuleName = "vsCode"
)

type FileMetrics struct {
	Name              string  `json:"name" bson:"name"`
	TotalFileCoverage float64 `json:"totalFileCoverage" bson:"totalFileCoverage"`
	DiffLinesCoverage float64 `json:"diffLinesCoverage" bson:"diffLinesCoverage"`
}

type PRMetrics struct {
	Type                  EventType     `json:"type" bson:"type"`
	PRNumber              int           `json:"prNumber" bson:"prNumber"`
	RepoName              string        `json:"repoName" bson:"repoName"`
	NumberOfFiles         int           `json:"numberOfFiles" bson:"numberOfFiles"`
	ImprovedCoverageFiles int           `json:"numberOfFilesOfIncCov" bson:"numberOfFilesOfIncCov"`
	FileMetrics           []FileMetrics `json:"fileMetrics" bson:"fileMetrics"`
	PRCoverage            int           `json:"prCoverage" bson:"prCoverage"`
	InitCovPercent        float64       `json:"initCovPercent" bson:"initCovPercent"`
	CovIncPercent         float64       `json:"covIncrPercent" bson:"covIncrPercent"`
	CoveredByKeploy       int           `json:"coveredByKeploy" bson:"coveredByKeploy"`
}
