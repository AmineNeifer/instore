package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"fake.com/instore/storage"
	"fake.com/instore/storepb"
	"google.golang.org/grpc"
)

const bClt string = "[clt] "

// UseCsv is a funcion that adds key-value pair to the in-memory store
func UseCsv() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.UseCsvRequest{
		Msg: "",
	}
	res, err := c.UseCsv(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// AddCsv is a funcion that adds key-value pair to the in-memory store
func AddCsv(key string, value string) {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.AddCsvRequest{
		Key:   key,
		Value: value,
	}
	res, err := c.AddCsv(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetvCsv is a funcion that gets values of the key in the in-memory store
func GetvCsv(key string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.GetvCsvRequest{
		Key: key,
	}
	res, err := c.GetvCsv(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetkCsv is a funcion that gets keys from value in the in-memory store
func GetkCsv(value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.GetkCsvRequest{
		Value: value,
	}
	res, err := c.GetkCsv(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetAllCsv is a function that get all key-value pairs in the in-memory store
func GetAllCsv() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)

	req := &storepb.GetAllCsvRequest{
		Msg: "",
	}

	resStream, err := c.GetAllCsv(context.Background(), req)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(bClt + "Getting all key-value pairs...")
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println("{'" + msg.GetKey() + "': '" + msg.GetValue() + "'}")
	}
}

// AddCsvFromFile is a funcion that adds key-value pairs imported from a csv file to the in-memory store
func AddCsvFromFile(filename string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	m := storage.LoadCsv(filename)
	var requests []*storepb.AddCsvFromFileRequest
	// request{key, value} and request is a slice of request
	// fill requests with key-value pairs from setmultimap `m`
	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			request := &storepb.AddCsvFromFileRequest{
				Key:   k.(string),
				Value: v.(string),
			}
			requests = append(requests, request)
		}
	}
	stream, err := c.AddCsvFromFile(context.Background())
	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// RemoveCsv is a funcion that removes key-value pair from the in-memory store
func RemoveCsv(key string, value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.RemoveCsvRequest{
		Key:   key,
		Value: value,
	}
	res, err := c.RemoveCsv(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"Error while removing pair: %v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// RemoveCsvFromFile is a funcion that removes key-value pairs imported from a csv file to the in-memory store
func RemoveCsvFromFile(filename string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	m := storage.LoadCsv(filename)
	var requests []*storepb.RemoveCsvFromFileRequest
	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			request := &storepb.RemoveCsvFromFileRequest{
				Key:   k.(string),
				Value: v.(string),
			}
			requests = append(requests, request)
		}
	}
	stream, err := c.RemoveCsvFromFile(context.Background())
	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// RemoveAllCsv is a funcion that removes all key-value pairs from the in-memory store
func RemoveAllCsv() {
	// Attenpt to secure connection here
	//
	// certFile := "ssl/ca.crt"
	// creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	// if sslErr != nil {
	// 	fmt.Printf(bClt+"%v", sslErr)
	// 		os.Exit(1)
	// }
	// cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.RemoveAllCsvRequest{
		Msg: "",
	}
	res, err := c.RemoveAllCsv(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"Error while removing all: %v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}
// *******************************************************************
// 			 MONGODB MODE FUNCTIONS (StoreDBService)
// *******************************************************************
// AddDb is a funcion that adds key-value pair to the in-memory store
func AddDb(key string, value string) {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreDbServiceClient(cc)

	pair := &storepb.Data{
		Key:   key,
		Value: value,
	}

	req := &storepb.AddDbRequest{
		Data: pair,
	}
	res, err := c.AddDb(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// UseDb is a funcion that adds key-value pair to the in-memory store
func UseDb() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreDbServiceClient(cc)
	req := &storepb.UseDbRequest{
		Msg: "",
	}
	res, err := c.UseDb(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetvDb is a funcion that gets values of the key in the in-memory store
func GetvDb(key string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreDbServiceClient(cc)
	req := &storepb.GetvDbRequest{
		Key: key,
	}
	res, err := c.GetvDb(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetkDb is a funcion that gets keys of the value in the in-memory store
func GetkDb(value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreDbServiceClient(cc)
	req := &storepb.GetkDbRequest{
		Value: value,
	}
	res, err := c.GetkDb(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}


// RemoveDb is a funcion that adds key-value pair to the in-memory store
func RemoveDb(key string, value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreDbServiceClient(cc)
	req := &storepb.RemoveDbRequest{
		Key:   key,
		Value: value,
	}
	res, err := c.RemoveDb(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}
