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

// Add is a funcion that adds key-value pair to the in-memory store
func Add(key string, value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.StoreRequest{
		Key:   key,
		Value: value,
	}
	res, err := c.Store(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetV is a funcion that gets values of the key in the in-memory store
func GetV(key string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.GetVRequest{
		Key: key,
	}
	res, err := c.GetV(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v\n", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetK is a funcion that gets keys from value in the in-memory store
func GetK(value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.GetKRequest{
		Value: value,
	}
	res, err := c.GetK(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// GetAll is a function that get all key-value pairs in the in-memory store
func GetAll() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)

	req := &storepb.GetAllRequest{
		Msg: "",
	}

	resStream, err := c.GetAll(context.Background(), req)
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
		fmt.Println("{'" + msg.GetKey()+ "': '" + msg.GetValue() + "'}")
	}
}

// AddCsv is a funcion that adds key-value pairs imported from a csv file to the in-memory store
func AddCsv(filename string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	m := storage.LoadCsv(filename)
	var requests []*storepb.StoreCsvRequest
	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			request := &storepb.StoreCsvRequest{
				Key:   k.(string),
				Value: v.(string),
			}
			requests = append(requests, request)
		}
	}
	stream, err := c.StoreCsv(context.Background())
	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	fmt.Println(requests)
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// Remove is a funcion that removes key-value pair from the in-memory store
func Remove(key string, value string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.RemoveRequest{
		Key:   key,
		Value: value,
	}
	res, err := c.Remove(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"Error while removing pair: %v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// RemoveCsv is a funcion that removes key-value pairs imported from a csv file to the in-memory store
func RemoveCsv(filename string) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	m := storage.LoadCsv(filename)
	var requests []*storepb.RemoveCsvRequest
	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			request := &storepb.RemoveCsvRequest{
				Key:   k.(string),
				Value: v.(string),
			}
			requests = append(requests, request)
		}
	}
	stream, err := c.RemoveCsv(context.Background())
	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}

// RemoveAll is a funcion that removes all key-value pairs from the in-memory store
func RemoveAll() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		fmt.Printf(bClt+"%v", err)
		os.Exit(1)
	}
	defer cc.Close()

	c := storepb.NewStoreServiceClient(cc)
	req := &storepb.RemoveAllRequest{
		Msg: "",
	}
	res, err := c.RemoveAll(context.Background(), req)

	if err != nil {
		fmt.Printf(bClt+"Error while removing all: %v", err)
		os.Exit(1)
	}
	fmt.Println(res.Result)
}
