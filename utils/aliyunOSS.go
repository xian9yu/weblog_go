package utils

//
//import (
//	"fmt"
//	"io"
//
//	"github.com/aliyun/aliyun-oss-go-sdk/oss"
//)
//
//var (
//	v               = GetConfig().Oss // 读取内存中的配置文件
//	endpoint        = v.EndPoint
//	accessKeyID     = v.AccessKeyID
//	accessKeySecret = v.AccessKeySecret
//	bucketName      = v.BucketName
//)
//
//func OSSClient() OSS {
//	return OSS{}
//}
//
//type OSS struct {
//}
//
//func (OSS) aliOSS() *oss.Client {
//	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
//	if err != nil {
//		panic("OSS 连接失败: " + err.Error())
//	}
//	return client
//}
//
//// 上传文件
//func (o OSS) Put(objectKey string, reader io.Reader) error {
//	bucket, err := o.aliOSS().Bucket(bucketName)
//	if err != nil {
//		return err
//	}
//
//	err = bucket.PutObject(objectKey, reader)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// 下载文件
//func (o OSS) Get(objectKey, filePath string) error {
//	bucket, err := o.aliOSS().Bucket(bucketName)
//	if err != nil {
//		return err
//	}
//
//	err = bucket.GetObjectToFile(objectKey, filePath)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// 文件列表
//func (o OSS) List() error {
//	bucket, err := o.aliOSS().Bucket(bucketName)
//	if err != nil {
//		return err
//	}
//
//	lsRes, err := bucket.ListObjects()
//	if err != nil {
//		return err
//	}
//
//	for _, object := range lsRes.Objects {
//		fmt.Println("Objects:", object.Key)
//	}
//	return nil
//}
//
//// 删除文件
//func (o OSS) Delete(objectKey string) error {
//	bucket, err := o.aliOSS().Bucket(bucketName)
//	if err != nil {
//		return err
//	}
//
//	err = bucket.DeleteObject(objectKey)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
