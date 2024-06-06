package file

import (
	"bytes"
	"io"
	"net/url"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/config"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type ServiceS3 struct {
	log.Log
	ctx *config.Context
}

// NewServiceS3 NewServiceS3
func NewServiceS3(ctx *config.Context) *ServiceS3 {

	return &ServiceS3{
		Log: log.NewTLog("ServiceS3"),
		ctx: ctx,
	}
}

// UploadFile 上传文件
func (s *ServiceS3) UploadFile(filePath string, contentType string, copyFileWriter func(io.Writer) error) (map[string]interface{}, error) {
	bucketName := s.ctx.GetConfig().OSS.BucketName
	// strs := strings.Split(filePath, "/")
	// if len(strs) > 0 {
	// 	bucketName = strs[0]
	// }
	session := s.newSession()
	uploader := s3manager.NewUploader(session)
	buff := bytes.NewBuffer(make([]byte, 0))
	err := copyFileWriter(buff)
	if err != nil {
		s.Error("复制文件内容失败！", zap.Error(err))
		return nil, err
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filePath),
		Body:        buff,
		ContentType: &contentType,
	})

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *ServiceS3) DownloadURL(path string, filename string) (string, error) {
	ossCfg := s.ctx.GetConfig().OSS

	rpath, _ := url.JoinPath(ossCfg.BucketURL, path)
	return rpath, nil
}

func (s *ServiceS3) newSession() *session.Session {
	ossCfg := s.ctx.GetConfig().OSS
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(ossCfg.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(false),
		DisableSSL:       aws.Bool(false),
		Credentials: credentials.NewStaticCredentials(
			ossCfg.AccessKeyID,
			ossCfg.AccessKeySecret,
			"",
		),
	})
	return sess
}
