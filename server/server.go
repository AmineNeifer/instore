package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fake.com/instore/storage"
	"fake.com/instore/storepb"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// initializeing filePath to for later use
var _, b, _, _ = runtime.Caller(0)

// basepath = "/path/to/server" the directory in which we find server.go
var basepath = filepath.Dir(b)

// filePath = /path/to/server/csvFiles/ the directory in which we find
// csv file to be used by server. ie where the data is stored
var filePath = basepath + "/csvFiles/" + filename

const (
	bSvr     string = "[svr]  "
	filename string = "data.csv"
)

var collection *mongo.Collection

type ItemBD struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type server struct{}

func main() {

	// helps with detecting line in which we got an error <3 in case...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection = client.Database("mydb").Collection("data")

	// server authentication SSL/TLS
	// "attempt lol"
	//
	// certFile := "ssl/server.crt"
	// keyFile := "ssl/server.pem"
	// creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	// if sslErr != nil {
	// 	fmt.Printf(bSvr+"%v", sslErr)
	// 	os.Exit(1)
	// }
	// s := grpc.NewServer(grpc.Creds(creds))

	s := grpc.NewServer()
	// listen on tcp
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Printf(bSvr+"%v\n", err)
		os.Exit(1)
	}
	// Registers a services and its implementations to the gRPC server.
	storepb.RegisterStoreServiceServer(s, &server{})
	storepb.RegisterStoreDbServiceServer(s, &server{})

	go func() {
		fmt.Println(bSvr + "Oh you chose to run me :D!")
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(bSvr + "Good choice ;)")
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(bSvr + "waiting for the client...")
		err := s.Serve(lis)
		if err != nil {
			fmt.Printf(bSvr+"%v\n", err)
			os.Exit(1)
		}
	}()

	// wait for ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is recieved
	<-ch

	s.Stop()
	lis.Close()
	fmt.Println(" ------Interrupted with Ctrl-C------")
	fmt.Println(bSvr + "Thank you for your stay! :D")
	fmt.Println(bSvr + "Have a good day!")
}

// UseCsv is a funcion that changes storing mode to CSV
func (*server) UseCsv(ctx context.Context, req *storepb.UseCsvRequest) (*storepb.UseCsvResponse, error) {
	fmt.Println(bSvr + "UseCsv function was invoked")
	storage.StoreType = ""
	result := bSvr + "Using CSV mode now..."
	res := &storepb.UseCsvResponse{
		Result: result,
	}
	return res, nil
}

// AddCsv is a funcion that adds key-value pair to the in-memory store
func (*server) AddCsv(ctx context.Context, req *storepb.AddCsvRequest) (*storepb.AddCsvResponse, error) {
	fmt.Printf(bSvr + "AddCsv function was invoked with %v\n", req)
	m := storage.LoadCsv(filePath)
	key := req.GetKey()
	value := req.GetValue()

	m.Put(key, value)

	result := bSvr + "Added: {'" + key + "': '" + value + "'}"
	res := &storepb.AddCsvResponse{
		Result: result,
	}
	storage.SaveCsv(m, filePath)
	return res, nil
}

// GetvCsv is a funcion that gets values of the key in the in-memory store
func (*server) GetvCsv(ctx context.Context, req *storepb.GetvCsvRequest) (*storepb.GetvCsvResponse, error) {
	fmt.Printf(bSvr+"GetvCsv function was invoked with %v\n", req)
	m := storage.LoadCsv(filePath)

	key := req.GetKey()

	values, found := m.Get(key)
	if !found {
		result := bSvr + "Key '" + key + "' entered isn't found in database"
		res := &storepb.GetvCsvResponse{
			Result: result,
		}
		return res, nil
	}
	var list []string
	for _, v := range values {
		list = append(list, v.(string))
	}
	result := bSvr + "list of the values: [" + strings.Join(list, " ") + "]"
	res := &storepb.GetvCsvResponse{
		Result: result,
	}
	storage.SaveCsv(m, filePath)
	return res, nil
}

// GetkCsv is a funcion that gets keys from value in the in-memory store
func (*server) GetkCsv(ctx context.Context, req *storepb.GetkCsvRequest) (*storepb.GetkCsvResponse, error) {
	fmt.Printf(bSvr+"GetkCsv function was invoked with %v\n", req)
	m := storage.LoadCsv(filePath)

	value := req.GetValue()

	if !m.ContainsValue(value) {
		result := bSvr + "Value '" + value + "' entered isn't found in database"
		res := &storepb.GetkCsvResponse{
			Result: result,
		}
		return res, nil
	}
	var list []string
	for _, k := range m.KeySet() {
		values, _ := m.Get(k)
		for _, v := range values {
			if v == value {
				list = append(list, k.(string))
			}
		}
	}

	result := bSvr + "list of the keys: [" + strings.Join(list, " ") + "]"
	res := &storepb.GetkCsvResponse{
		Result: result,
	}
	storage.SaveCsv(m, filePath)
	return res, nil
}

// GetAllCsv is a function that get all key-value pairs in the in-memory store
func (*server) GetAllCsv(req *storepb.GetAllCsvRequest, stream storepb.StoreService_GetAllCsvServer) error {
	fmt.Println(bSvr + "GetAllCsv function was invoked")
	m := storage.LoadCsv(filePath)

	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			res := &storepb.GetAllCsvResponse{
				Key:   k.(string),
				Value: v.(string),
			}
			stream.Send(res)
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

// AddCsvFromFile is a funcion that adds key-value pairs imported from a csv file to the in-memory store
func (*server) AddCsvFromFile(stream storepb.StoreService_AddCsvFromFileServer) error {
	fmt.Printf(bSvr + "AddCsvFromFile function was invoked with a streaming request\n")
	m := storage.LoadCsv(filePath)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			storage.SaveCsv(m, filePath)
			return stream.SendAndClose(&storepb.AddCsvFromFileResponse{
				Result: bSvr + "Csv pairs added with success",
			})
		}
		if err != nil {
			fmt.Printf(bSvr+"%v", err)
			os.Exit(1)
		}
		key := req.GetKey()
		value := req.GetValue()
		m.Put(key, value)
		fmt.Println(m)
	}
}

// RemoveCsv is a funcion that removes key-value pair from the in-memory store
func (*server) RemoveCsv(ctx context.Context, req *storepb.RemoveCsvRequest) (*storepb.RemoveCsvResponse, error) {
	fmt.Printf(bSvr+"RemoveCsv function was invoked with %v\n", req)
	m := storage.LoadCsv(filePath)
	key := req.GetKey()
	value := req.GetValue()

	m.Remove(key, value)

	result := bSvr + "Removed: {'" + key + "': '" + value + "'}"
	res := &storepb.RemoveCsvResponse{
		Result: result,
	}
	storage.SaveCsv(m, filePath)
	return res, nil
}

// RemoveCsvFromFile is a funcion that removes key-value pairs imported from a csv file to the in-memory store
func (*server) RemoveCsvFromFile(stream storepb.StoreService_RemoveCsvFromFileServer) error {
	fmt.Printf(bSvr + "RemoveCsvFromFile function was invoked with a streaming request\n")
	m := storage.LoadCsv(filePath)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			storage.SaveCsv(m, filePath)
			return stream.SendAndClose(&storepb.RemoveCsvFromFileResponse{
				Result: bSvr + "Csv pairs removed with success",
			})
		}
		if err != nil {
			fmt.Printf(bSvr+"%v", err)
			os.Exit(1)
		}
		key := req.GetKey()
		value := req.GetValue()
		m.Remove(key, value)
	}
}

// RemoveAllCsv is a funcion that removes all key-value pairs from the in-memory store
func (*server) RemoveAllCsv(ctx context.Context, req *storepb.RemoveAllCsvRequest) (*storepb.RemoveAllCsvResponse, error) {
	fmt.Println(bSvr + "RemoveAllCsv function was invoked")
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		res := &storepb.RemoveAllCsvResponse{
			Result: bSvr + "Already empty",
		}
		return res, nil
	}
	os.Remove(filePath)
	res := &storepb.RemoveAllCsvResponse{
		Result: bSvr + "All removed",
	}
	return res, nil
}

// *******************************************************************
// 			 MONGODB MODE FUNCTIONS (StoreDBService)
// *******************************************************************

// UseDb is a funcion that adds key-value pair to the in-memory store
func (*server) UseDb(ctx context.Context, req *storepb.UseDbRequest) (*storepb.UseDbResponse, error) {
	fmt.Println(bSvr + "UseDb function was invoked")
	storage.StoreType = "db"
	result := bSvr + "Using MongoDB now..."
	res := &storepb.UseDbResponse{
		Result: result,
	}
	return res, nil
}

// AddDb is a funcion that adds key-value pair to the in-memory store
func (*server) AddDb(ctx context.Context, req *storepb.AddDbRequest) (*storepb.AddDbResponse, error) {
	fmt.Printf(bSvr+"AddDb function was invoked with %v\n", req)

	pair := req.GetData()
	data := ItemBD{
		Key:   pair.GetKey(),
		Value: pair.GetValue(),
	}
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v", err),
		)
	}
	result := bSvr + "Added: {'" + data.Key + "': '" + data.Value + "'}"

	return &storepb.AddDbResponse{
		Result: result,
	}, nil
}

// GetvDb is a funcion that gets values of the key in the in-memory store
func (*server) GetvDb(ctx context.Context, req *storepb.GetvDbRequest) (*storepb.GetvDbResponse, error) {
	fmt.Printf(bSvr+"GetvDb function was invoked with %v\n", req)
	key := req.GetKey()
	// var pair ItemBD
	var pairs []ItemBD
	cursor, err := collection.Find(context.Background(), bson.M{"key": key})
	if err != nil {
		fmt.Println("Error finding all documents: ", err)
		defer cursor.Close(ctx)
	} else {
		for cursor.Next(ctx) {
			var r ItemBD
			err := cursor.Decode(&r)
			if err != nil {
				fmt.Printf(bSvr+"%v\n", err)
				os.Exit(1)
			} else {
				pairs = append(pairs, r)
			}
		}
	}
	var values []string
	for _, pair := range pairs {
		values = append(values, pair.Value)
	}

	var result string
	if len(values) == 0 {
		result = bSvr + "Key '" + key + "' entered isn't found in database"
	} else {
		result = bSvr + "list of the values: [" + strings.Join(values, " ") + "]"
	}

	res := &storepb.GetvDbResponse{
		Result: result,
	}
	return res, nil
}

// GetkDb is a funcion that gets keys of the value in the in-memory store
func (*server) GetkDb(ctx context.Context, req *storepb.GetkDbRequest) (*storepb.GetkDbResponse, error) {
	fmt.Printf(bSvr+"GetkDb function was invoked with %v\n", req)
	value := req.GetValue()
	// var pair ItemBD
	var pairs []ItemBD
	cursor, err := collection.Find(context.Background(), bson.M{"value": value})
	if err != nil {
		fmt.Printf(bSvr+"%v\n", err)
		defer cursor.Close(ctx)
	} else {
		for cursor.Next(ctx) {
			var r ItemBD
			err := cursor.Decode(&r)
			if err != nil {
				fmt.Printf(bSvr+"%v\n", err)
				os.Exit(1)
			} else {
				pairs = append(pairs, r)
			}
		}
	}
	var keys []string
	for _, pair := range pairs {
		keys = append(keys, pair.Key)
	}

	var result string
	if len(keys) == 0 {
		result = bSvr + "Value '" + value + "' entered isn't found in database"
	} else {
		result = bSvr + "list of the keys: [" + strings.Join(keys, " ") + "]"
	}
	res := &storepb.GetkDbResponse{
		Result: result,
	}
	return res, nil
}

// RemoveDb is a funcion that adds key-value pair to the in-memory store
func (*server) RemoveDb(ctx context.Context, req *storepb.RemoveDbRequest) (*storepb.RemoveDbResponse, error) {
	fmt.Printf(bSvr+"RemoveDb function was invoked with %v\n", req)
	data := ItemBD{
		Key:   req.GetKey(),
		Value: req.GetValue(),
	}
	_, err := collection.DeleteOne(context.Background(), bson.M{"key": data.Key})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v", err),
		)
	}
	result := bSvr + "Removed: {'" + req.GetKey() + "': '" + req.GetValue() + "'}"
	return &storepb.RemoveDbResponse{
		Result: result,
	}, nil
}
