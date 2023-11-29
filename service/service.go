package service

import (
	"context"
	"fmt"
	imgKeeperv1 "github.com/1azar/imgKeeper-api-contracts/gen/go/imgKeeper"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

const (
	UploadFile   = "upload"
	DownloadFile = "download"
	ListFiles    = "list"
)

type ClientService struct {
	addr      string
	filePath  string
	batchSize int
	client    imgKeeperv1.ImgKeeperClient
	method    string
}

func New(addr string, filePath string, batchSize int, method string) *ClientService {
	//fmt.Println(method)
	switch method {
	case UploadFile, DownloadFile, ListFiles:
		return &ClientService{
			addr:      addr,
			filePath:  filePath,
			batchSize: batchSize,
			method:    method,
		}
	default:
		fmt.Println(method)
		fmt.Println(1111111111111)
		panic(fmt.Sprintf("method flag should be one of: %s, %s, %s", UploadFile, DownloadFile, ListFiles))
	}
}

//func (s *ClientService) SendFile() error {
//	log.Println(s.addr, s.filePath)
//	conn, err := grpc.Dial(s.addr, grpc.WithInsecure())
//	if err != nil {
//		return err
//	}
//	defer conn.Close()
//	s.client = imgKeeperv1.NewImgKeeperClient(conn)
//	interrupt := make(chan os.Signal, 1)
//	shutdownSignals := []os.Signal{
//		os.Interrupt,
//		syscall.SIGTERM,
//		syscall.SIGINT,
//		syscall.SIGQUIT,
//	}
//	signal.Notify(interrupt, shutdownSignals...)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	go func(s *ClientService) {
//		if err = s.upload(ctx, cancel); err != nil {
//			log.Fatal(err)
//			cancel()
//		}
//	}(s)
//
//	select {
//	case killSignal := <-interrupt:
//		log.Println("Got ", killSignal)
//		cancel()
//	case <-ctx.Done():
//	}
//	return nil
//}

func (s *ClientService) TransferFile() error {
	log.Println(s.addr, s.filePath)
	conn, err := grpc.Dial(s.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	s.client = imgKeeperv1.NewImgKeeperClient(conn)
	interrupt := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
	signal.Notify(interrupt, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(s *ClientService) {
		switch s.method {
		case UploadFile:
			if err = s.upload(ctx, cancel); err != nil {
				log.Fatal(err)
				cancel()
			}
		case DownloadFile:
			if err = s.download(ctx, cancel, s.filePath); err != nil {
				log.Fatal(err)
				cancel()
			}
		case ListFiles:

		}
	}(s)

	select {
	case killSignal := <-interrupt:
		log.Println("Got ", killSignal)
		cancel()
	case <-ctx.Done():
	}
	return nil
}

func (s *ClientService) upload(ctx context.Context, cancel context.CancelFunc) error {
	stream, err := s.client.UploadImg(ctx)
	if err != nil {
		return err
	}
	file, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	buf := make([]byte, s.batchSize)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := buf[:num]

		//if err := stream.Send(&uploadpb.FileUploadRequest{FileName: s.filePath, Chunk: chunk}); err != nil {
		if err := stream.Send(&imgKeeperv1.ImgUploadReq{FileName: s.filePath, Chunk: chunk}); err != nil {
			return err
		}
		log.Printf("Sent - batch #%v - size - %v\n", batchNumber, len(chunk))
		batchNumber += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("Sent - %v bytes - %s\n", res.GetSize(), res.GetFileName())
	cancel()
	return nil
}

func (s *ClientService) download(ctx context.Context, cancel context.CancelFunc, fileName string) error {
	stream, err := s.client.DownloadImg(ctx, &imgKeeperv1.ImgDownloadReq{FileName: fileName})
	if err != nil {
		return err
	}

	outpuFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outpuFile.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = outpuFile.Write(chunk.Chunk)
		if err != nil {
			return err
		}
	}

	fmt.Printf("file %s has been downloaded.\n", fileName)

	return nil
}
