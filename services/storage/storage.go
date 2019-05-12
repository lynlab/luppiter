package storage

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/h2non/filetype"
)

const awsBucketName = "luppiter.lynlab.co.kr"

var sess *session.Session
var cachePath string

var (
	StorageErr          = errors.New("StorageErr")
	StorageErrNotExists = errors.New("StorageErrNotExists")
)

// ReadFile 함수는 주어진 namespace와 key에 해당하는 파일 및 컨텐츠 타입을 반환합니다.
func ReadFile(namespace, key string) (io.Reader, string, error) {
	// 우선 로컬 캐시 디렉토리에서 가져온다.
	file, err := ioutil.ReadFile(cachePath + "/" + namespace + "/" + key)
	if err == nil {
		kind, _ := filetype.Match(file)
		return bytes.NewReader(file), kind.MIME.Value, nil
	}

	// 없으면 S3에서 가져온다.
	output, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(namespace + "/" + key),
	})
	if err != nil {
		switch err.(awserr.Error).Code() {
		case s3.ErrCodeNoSuchKey:
			return nil, "", StorageErrNotExists
		default:
			return nil, "", StorageErr
		}
	}
	var buffer bytes.Buffer
	io.Copy(&buffer, output.Body)
	body := buffer.Bytes()

	// 가져온 파일을 로컬 캐시 디렉토리에 저장한다.
	os.MkdirAll(cachePath+"/"+namespace, os.ModePerm)
	newFile, err := os.Create(cachePath + "/" + namespace + "/" + key)
	if err == nil {
		io.Copy(newFile, bytes.NewReader(body))
	}

	return bytes.NewReader(body), *output.ContentType, nil
}

// WriteFile 함수는 주어진 namespace와 key에 해당하는 파일을 저장합니다.
func WriteFile(namespace, key string, reader io.Reader) error {
	var buffer bytes.Buffer
	io.Copy(&buffer, reader)
	data := buffer.Bytes()

	kind, _ := filetype.Match(data)
	_, err := s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(awsBucketName),
		Key:         aws.String(namespace + "/" + key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(kind.MIME.Value),
	})
	return err
}

func init() {
	newSess, err := session.NewSession(&aws.Config{Region: aws.String("ap-northeast-2")})
	if err != nil {
		panic(err)
	}

	cachePath = os.Getenv("LUPPITER_STORAGE_CACHE_PATH")
	if cachePath == "" {
		cachePath = "/tmp"
	}

	sess = newSess
}
