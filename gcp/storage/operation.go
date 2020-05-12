package storage

import (
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/taisukeyamashita/test/lib/config"
)

type StorageOpeator interface {
	StorageClient() *storage.Client
	GetFromCSVBucket() (*os.File, error)
}

const (
	publicURLtmp = "https://storage.cloud.google.com/%s/%s/%s"
	csvFileName  = "c01.csv"
	bucketName   = "free_test_bucket"
	filepath     = "csv/"
)

type Config struct {
	PublicURLtmp string
	BucketName   string
	CsvFilePath  string
	CsvFileName  string
}

type Operator struct {
	Client *storage.Client
	Config Config
}

var _ StorageOpeator = &Operator{}

func ProvideStorageOpeator(client *storage.Client) StorageOpeator {
	return &Operator{
		Client: client,
		Config: Config{
			PublicURLtmp: publicURLtmp,
			CsvFileName:  csvFileName,
			BucketName:   config.MyFreeBucketName,
			CsvFilePath:  config.CsvBucketFilePath,
		},
	}
}

func (op *Operator) GetFromCSVBucket() (*os.File, error) {
	csvURL := fmt.Sprintf(publicURLtmp, bucketName, filepath, csvFileName)
	log.Print("storageURL :" + csvURL)
	file, err := os.Open(csvURL)
	if err != nil {
		return file, err
	}
	defer file.Close()
	// reader := csv.NewReader(file)
	return file, err
	// defer f.Close()
}

func (op *Operator) StorageClient() *storage.Client {
	return op.Client
}

// func GetFromCSVBucket() (*os.File, error) {
// 	csvURL := fmt.Sprintf(publicURLtmp, bucketName, filepath, csvFileName)
// 	log.Print("storageURL :" + csvURL)
// 	file, err := os.Open(csvURL)
// 	if err != nil {
// 		return file, err
// 	}
// 	defer file.Close()
// 	reader := csv.NewReader(file)
// 	return file, err
// 	// defer f.Close()
// }
