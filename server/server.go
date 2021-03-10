package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"fake.com/instore/storage"
	"fake.com/instore/storepb"
	"google.golang.org/grpc"
)

const (
	bSvr     string = "[srv]  "
	filename        = "data.csv"
)

func (*server) Store(ctx context.Context, req *storepb.StoreRequest) (*storepb.StoreResponse, error) {
	fmt.Printf(bSvr+"Add function was invoked with %v\n", req)
	m := storage.LoadCsv("data.csv")
	key := req.GetKey()
	value := req.GetValue()

	m.Put(key, value)

	result := bSvr + "Added: {'" + key + "': '" + value + "'}"
	res := &storepb.StoreResponse{
		Result: result,
	}
	storage.SaveCsv(m, filename)
	return res, nil
}

func (*server) StoreCsv(stream storepb.StoreService_StoreCsvServer) error {
	fmt.Printf(bSvr + "AddCsv function was invoked with a streaming request\n")
	m := storage.LoadCsv(filename)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			storage.SaveCsv(m, filename)
			return stream.SendAndClose(&storepb.StoreCsvResponse{
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
	}
}

func (*server) GetV(ctx context.Context, req *storepb.GetVRequest) (*storepb.GetVResponse, error) {
	fmt.Printf(bSvr+"GetV function was invoked with %v\n", req)
	m := storage.LoadCsv("data.csv")

	key := req.GetKey()

	values, found := m.Get(key)
	if !found {
		result := bSvr + "Key '" + key + "' entered isn't found in database"
		res := &storepb.GetVResponse{
			Result: result,
		}
		return res, nil
	}
	var list []string
	for _, v := range values {
		list = append(list, v.(string))
	}
	result := bSvr + "list of the values: [" + strings.Join(list, " ") + "]"
	res := &storepb.GetVResponse{
		Result: result,
	}
	storage.SaveCsv(m, filename)
	return res, nil
}

func (*server) GetK(ctx context.Context, req *storepb.GetKRequest) (*storepb.GetKResponse, error) {
	fmt.Printf(bSvr+"GetK function was invoked with %v\n", req)
	m := storage.LoadCsv("data.csv")

	value := req.GetValue()

	if !m.ContainsValue(value) {
		result := bSvr + "Value '" + value + "' entered isn't found in database"
		res := &storepb.GetKResponse{
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
	res := &storepb.GetKResponse{
		Result: result,
	}
	storage.SaveCsv(m, filename)
	return res, nil
}

func (*server) GetAll(req *storepb.GetAllRequest, stream storepb.StoreService_GetAllServer) error {
	fmt.Println(bSvr + "GetAll function was invoked")
	m := storage.LoadCsv(filename)

	for _, k := range m.KeySet() {
		value, _ := m.Get(k)
		for _, v := range value {
			res := &storepb.GetAllResponse{
				Key:   k.(string),
				Value: v.(string),
			}
			stream.Send(res)
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

func (*server) Remove(ctx context.Context, req *storepb.RemoveRequest) (*storepb.RemoveResponse, error) {
	fmt.Printf(bSvr+"Remove function was invoked with %v\n", req)
	m := storage.LoadCsv("data.csv")
	key := req.GetKey()
	value := req.GetValue()

	m.Remove(key, value)

	result := bSvr + "Removed: {'" + key + "': '" + value + "'}"
	res := &storepb.RemoveResponse{
		Result: result,
	}
	storage.SaveCsv(m, filename)
	return res, nil
}

func (*server) RemoveCsv(stream storepb.StoreService_RemoveCsvServer) error {
	fmt.Printf(bSvr + "RemoveCsv function was invoked with a streaming request\n")
	m := storage.LoadCsv(filename)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			storage.SaveCsv(m, filename)
			return stream.SendAndClose(&storepb.RemoveCsvResponse{
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

func (*server) RemoveAll(ctx context.Context, req *storepb.RemoveAllRequest) (*storepb.RemoveAllResponse, error) {
	fmt.Println(bSvr + "RemoveAll function was invoked")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		res := &storepb.RemoveAllResponse{
			Result: bSvr + "Already empty",
		}
		return res, nil
	}
	os.Remove(filename)
	res := &storepb.RemoveAllResponse{
		Result: bSvr + "All removed",
	}
	return res, nil
}

type server struct{}

func main() {
	fmt.Println(bSvr + "Oh you chose to run me :D!")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println(bSvr + "Good choice ;)")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println(bSvr + "waiting for the client...")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Printf(bSvr+"%v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	storepb.RegisterStoreServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		fmt.Printf(bSvr+"%v\n", err)
		os.Exit(1)
	}
}
