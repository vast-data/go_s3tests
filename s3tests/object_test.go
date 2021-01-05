package s3test

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/spf13/viper"

	"fmt"
	"strings"
	"time"

	. "../Utilities"
)

func (suite *S3Suite) TestObjectWriteToNonExistantBucket() {

	/*
		Reource object : method: get
		Operation : read object
		Assertion : read contents that were never written
	*/

	assert := suite
	non_exixtant_bucket := "bucketz"

	err := PutObjectToBucket(svc, non_exixtant_bucket, "key", "content")
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchBucket", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}

	}

}

func (suite *S3Suite) TestMultiObjectDelete() {

	/*
		Resource : object, method: put
		Scenario : delete multiple objects
		Assertion: deletes multiple objects with a single call.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "echo", "bar": "lima", "baz": "golf"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	DeleteObjects(svc, bucket)

	resp, err := GetObjects(svc, bucket)
	assert.Nil(err)
	assert.Equal(0, len(resp.Contents))
}

func (suite *S3Suite) TestObjectReadNotExist() {

	/*
		Reource object : method: get
		Operation : read object
		Assertion : read contents that were never written
	*/

	assert := suite
	bucket := GetBucketName()

	err := CreateBucket(svc, bucket)
	assert.Nil(err)

	_, err = GetObject(svc, bucket, "key6")
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchKey", awsErr.Code())
			assert.Equal("", awsErr.Message())

		}
	}

}

func (suite *S3Suite) TestObjectReadFromNonExistantBucket() {

	/*
		Reource object : method: get
		Operation : read object
		Assertion : read contents that were never written
	*/
	assert := suite
	non_exixtant_bucket := "bucketz"

	_, err := GetObject(svc, non_exixtant_bucket, "key6")
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchBucket", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}

	}

}

func (suite *S3Suite) TestObjectWriteReadUpdateReadDelete() {

	assert := suite
	bucket := GetBucketName()
	key := "key1"

	err := CreateBucket(svc, bucket)
	assert.Nil(err)

	// Write object
	PutObjectToBucket(svc, bucket, key, "hello")
	assert.Nil(err)

	// Read object
	result, _ := GetObject(svc, bucket, key)
	assert.Equal("hello", result)

	//Update object
	PutObjectToBucket(svc, bucket, key, "Come on !!")
	assert.Nil(err)

	// Read object again
	result, _ = GetObject(svc, bucket, key)
	assert.Equal("Come on !!", result)

	err = DeleteObjects(svc, bucket)
	assert.Nil(err)

	// If object was well deleted, there shouldn't be an error at this point
	err = DeleteBucket(svc, bucket)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectDeleteAll() {

	// Reading content that was never written should fail
	assert := suite
	bucket := GetBucketName()
	var empty_list []*s3.Object
	key := "key5"
	key1 := "key6"

	err := CreateBucket(svc, bucket)
	assert.Nil(err)

	PutObjectToBucket(svc, bucket, key, "hello")
	PutObjectToBucket(svc, bucket, key1, "foo")
	assert.Nil(err)
	objects, err := ListObjects(svc, bucket)
	assert.Nil(err)
	assert.Equal(2, len(objects))

	err = DeleteObjects(svc, bucket)
	assert.Nil(err)

	objects, err = ListObjects(svc, bucket)
	assert.Nil(err)
	assert.Equal(empty_list, objects)

}

func (suite *S3Suite) TestObjectCopyBucketNotFound() {

	// copy from non-existent bucket

	assert := suite
	bucket := GetBucketName()
	item := "key1"
	other := GetBucketName()

	source := bucket + "/" + item

	err := CreateBucket(svc, bucket)
	assert.Nil(err)

	// Write object
	PutObjectToBucket(svc, bucket, item, "hello")
	assert.Nil(err)

	err = CopyObject(svc, other, source, item)
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal(awsErr.Code(), "NoSuchBucket")
			assert.Equal(awsErr.Message(), "")
		}

	}

}

func (suite *S3Suite) TestObjectCopyKeyNotFound() {

	assert := suite
	bucket := GetBucketName()
	item := "key1"
	other := GetBucketName()

	source := bucket + "/" + item

	err := CreateBucket(svc, bucket)
	err = CreateBucket(svc, other)
	assert.Nil(err)

	err = CopyObject(svc, other, source, item)
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal(awsErr.Code(), "NoSuchKey")
			assert.Equal(awsErr.Message(), "")
		}

	}

}

//.....................................Test Getting Ranged Objects....................................................................................................................

func (suite *S3Suite) TestRangedRequest() {

	//getting objects in a range should return correct data

	assert := suite
	bucket := GetBucketName()
	key := "key"
	content := "testcontent"

	var data string
	var resp *s3.GetObjectOutput

	err := CreateBucket(svc, bucket)
	PutObjectToBucket(svc, bucket, key, content)

	resp, data, err = GetObjectWithRange(svc, bucket, key, "bytes=4-7")
	assert.Nil(err)
	assert.Equal(content[4:8], data)
	assert.Equal("bytes", *resp.AcceptRanges)
}

func (suite *S3Suite) TestRangedRequestSkipLeadingBytes() {

	//getting objects in a range should return correct data

	assert := suite
	bucket := GetBucketName()
	key := "key"
	content := "testcontent"

	var data string
	var resp *s3.GetObjectOutput

	err := CreateBucket(svc, bucket)
	PutObjectToBucket(svc, bucket, key, content)

	resp, data, err = GetObjectWithRange(svc, bucket, key, "bytes=4-")
	assert.Nil(err)
	assert.Equal(content[4:], data)
	assert.Equal("bytes", *resp.AcceptRanges)

}

func (suite *S3Suite) TestRangedRequestReturnTrailingBytes() {

	//getting objects in a range should return correct data

	assert := suite
	bucket := GetBucketName()
	key := "key"
	content := "testcontent"

	var data string
	var resp *s3.GetObjectOutput

	err := CreateBucket(svc, bucket)
	PutObjectToBucket(svc, bucket, key, content)

	resp, data, err = GetObjectWithRange(svc, bucket, key, "bytes=-8")
	assert.Nil(err)
	assert.Equal(content[3:11], data)
	assert.Equal("bytes", *resp.AcceptRanges)
}

func (suite *S3Suite) TestRangedRequestInvalidRange() {

	//getting objects in unaccepted range returns invalid range

	assert := suite
	bucket := GetBucketName()
	key := "key"
	content := "testcontent"

	err := CreateBucket(svc, bucket)
	PutObjectToBucket(svc, bucket, key, content)

	_, _, err = GetObjectWithRange(svc, bucket, key, "bytes=40-50")
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("InvalidRange", awsErr.Code())
			assert.Equal("", awsErr.Message())

		}
	}
}

func (suite *S3Suite) TestRangedRequestEmptyObject() {

	//getting a range of an empty object returns invalid range

	assert := suite
	bucket := GetBucketName()
	key := "key"
	content := ""

	err := CreateBucket(svc, bucket)
	PutObjectToBucket(svc, bucket, key, content)

	_, _, err = GetObjectWithRange(svc, bucket, key, "bytes=40-50")
	assert.NotNil(err)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("InvalidRange", awsErr.Code())
			assert.Equal("", awsErr.Message())

		}
	}
}

func (suite *S3Suite) TestObjectSetGetMetadataNoneToGood() {

	assert := suite
	metadata := map[string]*string{"mymeta": nil}
	got := GetSetMetadata(metadata)
	assert.Equal(metadata, got)
}

func (suite *S3Suite) TestObjectSetGetMetadataNoneToEmpty() {

	assert := suite
	metadata := map[string]*string{"": nil}
	got := GetSetMetadata(metadata)
	assert.Equal(metadata, got)
}

func (suite *S3Suite) TestObjectSetGetMetadataOverwriteToGood() {

	assert := suite

	oldmetadata := map[string]*string{"meta1": nil}
	got := GetSetMetadata(oldmetadata)
	assert.Equal(oldmetadata, got)

	newmetadata := map[string]*string{"meta2": nil}
	got = GetSetMetadata(newmetadata)
	assert.Equal(newmetadata, got)
}

func (suite *S3Suite) TestObjectSetGetMetadataOverwriteToEmpty() {

	assert := suite

	oldmetadata := map[string]*string{"meta1": nil}
	got := GetSetMetadata(oldmetadata)
	assert.Equal(oldmetadata, got)

	newmetadata := map[string]*string{"": nil}
	got = GetSetMetadata(newmetadata)
	assert.Equal(newmetadata, got)
}

//..............................................SSE-C encrypted transfer....................................................

func (suite *S3Suite) TestEncryptedTransfer1B() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 1byte
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := EncryptionSSECustomerWrite(svc, 1)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestEncryptedTransfer1KB() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 1KB
		Assertion: success.
	*/
	assert := suite

	rdata, data, err := EncryptionSSECustomerWrite(svc, 1024)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestEncryptedTransfer1MB() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 1MB
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := EncryptionSSECustomerWrite(svc, 1024*1024)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestEncryptedTransfer13B() {

	// Resource : object, method: put
	// Scenario : Test SSE-C encrypted transfer 13 bytes
	// Assertion: success.

	assert := suite

	rdata, data, err := EncryptionSSECustomerWrite(svc, 13)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestEncryptionSSECPresent() {

	/*
		Resource : object, method: put
		Scenario : write encrypted with SSE-C and read without SSE-C
		Assertion: fails.
	*/
	assert := suite

	data := strings.Repeat("A", 10)
	key := "testobj"
	bucket := GetBucketName()
	sse := []string{"AES256", "pO3upElrwuEXSoFwCfnZPdSsmt/xWeFa0N9KgDijwVs=", "DWygnHRtgiJ77HCm+1rvHw=="}

	err := CreateBucket(svc, bucket)

	err = WriteSSECEcrypted(svc, bucket, key, data, sse)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		_, err = GetObjects(svc, bucket)
		assert.NotNil(err)

	}
}

func (suite *S3Suite) TestEncryptionSSECOtherKey() {

	/*
		Resource : object, method: put/get
		Scenario : write encrypted with SSE-C but read with other key
		Assertion: fails.
	*/

	assert := suite

	data := strings.Repeat("A", 10)
	key := "testobj"
	bucket := GetBucketName()
	sse0 := []string{"AES256", "pO3upElrwuEXSoFwCfnZPdSsmt/xWeFa0N9KgDijwVs=", "DWygnHRtgiJ77HCm+1rvHw=="}
	sse1 := []string{"AES256", "6b+WOZ1T3cqZMxgThRcXAQBrS5mXKdDUphvpxptl9/4=", "arxBvwY2V4SiOne6yppVPQ=="}

	_ = CreateBucket(svc, bucket)

	err := WriteSSECEcrypted(svc, bucket, key, data, sse0)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		_, err = ReadSSECEcrypted(svc, bucket, key, sse1)
		assert.NotNil(err)

	}
}

func (suite *S3Suite) TestEncryptionSSECInvalidMd5() {

	/*
		Resource : object, method: put
		Scenario : write encrypted with SSE-C, but md5 is bad
		Assertion: fails.
	*/

	assert := suite

	data := strings.Repeat("A", 10)
	key := "testobj"
	bucket := GetBucketName()
	sse := []string{"AES256", "pO3upElrwuEXSoFwCfnZPdSsmt/xWeFa0N9KgDijwVs=", "AAAAAAAAAAAAAAAAAAAAAA=="}

	err := CreateBucket(svc, bucket)

	err = WriteSSECEcrypted(svc, bucket, key, data, sse)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		_, err = GetObjects(svc, bucket)
		assert.NotNil(err)

	}
}

func (suite *S3Suite) TestEncryptionSSECNoMd5() {

	/*
		Resource : object, method: put
		Scenario : write encrypted with SSE-C, but dont provide MD5'
		Assertion: fails.
	*/

	assert := suite

	data := strings.Repeat("A", 10)
	key := "testobj"
	bucket := GetBucketName()
	sse := []string{"AES256", "pO3upElrwuEXSoFwCfnZPdSsmt/xWeFa0N9KgDijwVs=", " "}

	err := CreateBucket(svc, bucket)

	err = WriteSSECEcrypted(svc, bucket, key, data, sse)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		_, err = GetObjects(svc, bucket)
		assert.NotNil(err)

	}
}

func (suite *S3Suite) TestEncryptionSSECNoKey() {

	/*
		Resource : object, method: put
		Scenario : declare SSE-C but do not provide key'
		Assertion: fails.
	*/

	assert := suite

	data := strings.Repeat("A", 10)
	key := "testobj"
	bucket := GetBucketName()
	sse := []string{"AES256", " ", " "}

	err := CreateBucket(svc, bucket)

	err = WriteSSECEcrypted(svc, bucket, key, data, sse)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		_, err = GetObjects(svc, bucket)
		assert.NotNil(err)

	}
}

func (suite *S3Suite) TestEncryptionKeyNoSSEC() {

	/*
		Resource : object, method: put
		Scenario : 'Do not declare SSE-C but provide key and MD5
		Assertion: operation successfull, no encryption.
	*/

	assert := suite

	data := strings.Repeat("A", 10)
	key := "testobj"
	bucket := GetBucketName()
	sse := []string{" ", "pO3upElrwuEXSoFwCfnZPdSsmt/xWeFa0N9KgDijwVs=", "DWygnHRtgiJ77HCm+1rvHw=="}

	err := CreateBucket(svc, bucket)

	err = WriteSSECEcrypted(svc, bucket, key, data, sse)
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		_, err = GetObjects(svc, bucket)
		assert.Nil(err)

	}

}

//.................................SSE and KMS......................................................................

func (suite *S3Suite) TestSSEKMSbarbTransfer13B() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 13 bytes
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSkeyIdCustomerWrite(svc, 13)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}

}

func (suite *S3Suite) TestSSEKMSbarbTransfer1MB() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 13 bytes
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSkeyIdCustomerWrite(svc, 1024*1024)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}

}

func (suite *S3Suite) TestSSEKMSbarbTransfer1KB() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 13 bytes
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSkeyIdCustomerWrite(svc, 1024)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestSSEKMSbarbTransfer1B() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-C encrypted transfer 13 bytes
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSkeyIdCustomerWrite(svc, 1)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}

}

func (suite *S3Suite) TestSSEKMSTransfer13B() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-KMS encrypted transfer 13 bytes
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSCustomerWrite(svc, 13)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestSSEKMSTransfer1MB() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-KMS encrypted transfer 1 mega byte
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSCustomerWrite(svc, 1024*1024)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}
}

func (suite *S3Suite) TestSSEKMSTransfer1KB() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-KMS encrypted transfer 1 kilobyte
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSCustomerWrite(svc, 1024)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {
		assert.Nil(err)
		assert.Equal(rdata, data)

	}

}

func (suite *S3Suite) TestSSEKMSTransfer1B() {

	/*
		Resource : object, method: put
		Scenario : Test SSE-KMS encrypted transfer 1 byte
		Assertion: success.
	*/

	assert := suite

	rdata, data, err := SSEKMSCustomerWrite(svc, 1)
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		assert.Equal(rdata, data)

	}

}

func (suite *S3Suite) TestSSEKMSPresent() {

	/*
		Resource : object, method: put
		Scenario : write encrypted with SSE-KMS and read without SSE-KMS
		Assertion: success.
	*/

	assert := suite

	bucket := GetBucketName()

	err := CreateBucket(svc, bucket)

	err = WriteSSEKMSkeyId(svc, bucket, "kay1", "test", viper.GetString("s3main.SSE"), viper.GetString("s3main.kmskeyid"))
	
	if awsErr, ok := err.(awserr.Error); ok {

		assert.NotNil(awsErr)

	} else {

		assert.Nil(err)
		data, _ := GetObject(svc, bucket, "kay1")

		assert.Equal("test", data)
	}

}

func (suite *S3Suite) TestSSEKMSNoKey() {

	/*
		Resource : object, method: put
		Scenario : declare SSE-KMS but do not provide key_id'
		Assertion: fails.
	*/

	assert := suite

	bucket := GetBucketName()

	err := CreateBucket(svc, bucket)

	err = WriteSSEKMSkeyId(svc, bucket, "kay1", "test", viper.GetString("s3main.SSE"), "")
	
	if awsErr, ok := err.(awserr.Error); ok {
		assert.NotNil(awsErr)
		assert.Equal("InvalidAccessKeyId", awsErr.Code())
	} else {

		assert.NotNil(err)
	}

}

func (suite *S3Suite) TestSSEKMSNotDeclared() {

	/*
		Resource : object, method: put
		Scenario : dDo not declare SSE-KMS but provide key_id
		Assertion: fails if either the aws:kms or key_id is not declared
	*/

	assert := suite

	bucket := GetBucketName()

	err := CreateBucket(svc, bucket)

	err = WriteSSEKMSkeyId(svc, bucket, "kay1", "test", "", viper.GetString("s3main.kmskeyid"))
	if awsErr, ok := err.(awserr.Error); ok {
		
		assert.NotNil(awsErr)
	}
	err = WriteSSEKMSkeyId(svc, bucket, "kay1", "test", viper.GetString("s3main.SSE"), "")
	if awsErr, ok := err.(awserr.Error); ok {
		
		assert.NotNil(awsErr)
	} 
}

//...................................... get object with conditions....................

func (suite *S3Suite) TestGetObjectIfmatchGood() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-Match: the latest ETag
		Assertion: suceeds.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	object, err := GetObj(svc, bucket, "foo")

	got, err := GetObjectWithIfMatch(svc, bucket, "foo", *object.ETag)
	assert.Nil(err)
	assert.Equal("bar", got)

}

func (suite *S3Suite) TestGetObjectIfmatchFailed() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-Match: bogus ETag
		Assertion: fails.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	_, err = GetObjectWithIfMatch(svc, bucket, "foo", "ABCORZ")
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			assert.Equal("PreconditionFailed", awsErr.Code())
		}
	}

}

func (suite *S3Suite) TestGetObjectIfNoneMatchGood() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-None-Match: the latest ETag
		Assertion: fails.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	object, err := GetObj(svc, bucket, "foo")

	_, err = GetObjectWithIfNoneMatch(svc, bucket, "foo", *object.ETag)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NotModified", awsErr.Code())
			assert.Equal("Not Modified", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestGetObjectIfNoneMatchFailed() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-None-Match: bogus ETag
		Assertion: suceeds.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	got, err := GetObjectWithIfNoneMatch(svc, bucket, "foo", "ABCORZ")
	assert.Nil(err)
	assert.Equal("bar", got)
}

func (suite *S3Suite) TestGetObjectIfModifiedSinceGood() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-Modified-Since: before
		Assertion: suceeds.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}
	now := time.Now()
	time.Sleep(1 * time.Second)

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	_, err = GetObj(svc, bucket, "foo")

	got, err := GetObjectWithIfModifiedSince(svc, bucket, "foo", now)
	assert.Nil(err)
	assert.Equal("bar", got)
}

func (suite *S3Suite) TestGetObjectIfUnModifiedSinceGood() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-Unmodified-Since: before
		Assertion: fails.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}
	now := time.Now()
	time.Sleep(1 * time.Second)

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	_, err = GetObjectWithIfUnModifiedSince(svc, bucket, "foo", now)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("PreconditionFailed", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestGetObjectIfUnModifiedSinceFailed() {

	/*
		Resource : object, method: get
		Scenario : get w/ If-Unmodified-Since: after
		Assertion: suceeds.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}
	now := time.Now()
	future := now.Add(time.Hour * 24 * 3)

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	got, err := GetObjectWithIfUnModifiedSince(svc, bucket, "foo", future)
	assert.Nil(err)
	assert.Equal("bar", got)
}

//................put object with condition..............................................

func (suite *S3Suite) TestPutObjectIfMatchGood() {

	/*
		Resource : object, method: get
		Scenario : data re-write w/ If-Match: the latest ETag
		Assertion: replaces previous data.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	gotData, err := GetObject(svc, bucket, "foo")
	assert.Equal("bar", gotData)

	object, err := GetObj(svc, bucket, "foo")
	err = PutObjectWithIfMatch(svc, bucket, "foo", "zar", *object.ETag)
	assert.Nil(err)

	new_data, _ := GetObject(svc, bucket, "foo")
	assert.Nil(err)
	assert.Equal("zar", new_data)
}

func (suite *S3Suite) TestPutObjectIfMatchFailed() {

	/*
		Resource : object, method: get
		Scenario : data re-write w/ If-Match: outdated ETag
		Assertion: replaces previous data.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"key1": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	gotData, err := GetObject(svc, bucket, "key1")
	assert.Equal("bar", gotData)

	err = PutObjectWithIfMatch(svc, bucket, "key1", "zar", "ABCORZmmmm")

	oldData, err := GetObject(svc, bucket, "key1")
	assert.Nil(err)
	assert.Equal("zar", oldData)
}

func (suite *S3Suite) TestPutObjectIfmatchNonexistedFailed() {

	/*
		Resource : object, method: put
		Scenario : overwrite non-existing object w/ If-Match: *
		Assertion: fails
	*/

	assert := suite
	bucket := GetBucketName()

	err := CreateBucket(svc, bucket)
	assert.Nil(err)

	err = PutObjectWithIfMatch(svc, bucket, "foo", "zar", "*")
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchKey", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestPutObjectIfNonMatchGood() {

	/*
		Resource : object, method: get
		Scenario : overwrite existing object w/ If-None-Match: outdated ETag'
		Assertion: replaces previous data and metadata.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"foo": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	gotData, err := GetObject(svc, bucket, "foo")
	assert.Equal("bar", gotData)

	err = PutObjectWithIfNoneMatch(svc, bucket, "foo", "zar", "ABCORZ")
	assert.Nil(err)

	new_data, _ := GetObject(svc, bucket, "foo")
	assert.Nil(err)
	assert.Equal("zar", new_data)
}

func (suite *S3Suite) TestPutObjectIfNonMatchNonexistedGood() {

	/*
		Resource : object, method: get
		Scenario : overwrite non-existing object w/ If-None-Match: *
		Assertion: succeeds.
	*/

	assert := suite
	bucket := GetBucketName()

	err := CreateBucket(svc, bucket)

	err = PutObjectWithIfNoneMatch(svc, bucket, "key1", "bar", "*")
	assert.Nil(err)

	data, err := GetObject(svc, bucket, "key1")
	assert.Equal("bar", data)
}

func (suite *S3Suite) TestPutObjectIfNonMatchOverwriteExistedFailed() {

	/*
		Resource : object, method: get
		Scenario : overwrite existing object w/ If-None-Match: *
		Assertion: fails.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"key1": "bar"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)

	gotData, err := GetObject(svc, bucket, "key1")
	assert.Equal("bar", gotData)

	err = PutObjectWithIfNoneMatch(svc, bucket, "key1", "zar", "*")
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("PreconditionFailed", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}

	oldData, err := GetObject(svc, bucket, "key1")
	assert.Equal("bar", oldData)
}

//......................................Multipart Upload...................................................................

func (suite *S3Suite) TestAbortMultipartUploadInvalid() {

	/*
		Resource : object, method: get
		Scenario : Abort given invalid arguments.
		Assertion: fails.
	*/

	assert := suite
	bucket := GetBucketName()
	key := "mymultipart"

	err := CreateBucket(svc, bucket)

	_, err = AbortMultiPartUploadInvalid(svc, bucket, key, key)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("InvalidParameter", awsErr.Code())
			assert.Equal("1 validation error(s) found.", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestAbortMultipartUploadNotfound() {

	/*
		Resource : object, method: get
		Scenario : Abort non existant multipart upload
		Assertion: fails.
	*/

	assert := suite
	bucket := GetBucketName()
	key := "mymultipart"

	err := CreateBucket(svc, bucket)

	_, err = AbortMultiPartUpload(svc, bucket, key, key)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchUpload", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestAbortMultipartUpload() {

	/*
		Resource : object, method: get
		Scenario : Abort multipart upload
		Assertion: successful.
	*/

	assert := suite
	bucket := GetBucketName()
	bucket2 := GetBucketName()
	key := "key"
	fmtstring := fmt.Sprintf("%s/%s", bucket2, key)
	objects := map[string]string{key: "golf"}

	err := CreateBucket(svc, bucket2)
	err = CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket2, objects)

	result, err := InitiateMultipartUpload(svc, bucket, "key")
	_, err = UploadCopyPart(svc, bucket, key, fmtstring, *result.UploadId, int64(1))

	_, err = AbortMultiPartUpload(svc, bucket, key, *result.UploadId)
	assert.Nil(err)

	resp, err := Listparts(svc, bucket, key, *result.UploadId)
	assert.Equal(0, len(resp.Parts))
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchKey", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestMultipartUploadOverwriteExistingObject() {

	/*
		Resource : object, method: get
		Scenario : multi-part upload overwrites existing key
		Assertion: successful.
	*/

	assert := suite
	bucket := GetBucketName()
	num_parts := 2

	payload := strings.Repeat("12345", 1024*1024)
	key_name := "mymultipart"

	newObject := map[string]string{key_name: "payload"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, newObject)

	result, err := InitiateMultipartUpload(svc, bucket, key_name)

	resp, err := Uploadpart(svc, bucket, key_name, *result.UploadId, payload, int64(num_parts))
	assert.Nil(err)

	_, err = CompleteMultiUpload(svc, bucket, key_name, int64(num_parts), *result.UploadId, *resp.ETag)
	assert.Nil(err)

	gotData, err := GetObject(svc, bucket, key_name)
	assert.Nil(err)
	assert.Equal(payload, gotData)
}

func (suite *S3Suite) TestMultipartUploadContents() {

	/*
		Resource : object, method: get
		Scenario : check contents of multi-part upload
		Assertion: successful.
	*/
	assert := suite
	bucket := GetBucketName()
	num_parts := 2

	payload := strings.Repeat("12345", 1024*1024)
	key_name := "mymultipart"

	err := CreateBucket(svc, bucket)

	result, err := InitiateMultipartUpload(svc, bucket, key_name)

	resp, err := Uploadpart(svc, bucket, key_name, *result.UploadId, payload, int64(num_parts))
	assert.Nil(err)

	_, err = CompleteMultiUpload(svc, bucket, key_name, int64(num_parts), *result.UploadId, *resp.ETag)
	assert.Nil(err)

	gotData, err := GetObject(svc, bucket, key_name)
	assert.Nil(err)
	assert.Equal(payload, gotData)
}

func (suite *S3Suite) TestMultipartUploadInvalidPart() {

	/*
		Resource : object, method: get
		Scenario : check failure on multiple multi-part upload with invalid etag
		Assertion: fails.
	*/
	assert := suite
	bucket := GetBucketName()
	num_parts := 2

	payload := strings.Repeat("12345", 1024*1024)
	key_name := "mymultipart"

	err := CreateBucket(svc, bucket)

	result, err := InitiateMultipartUpload(svc, bucket, key_name)

	_, err = Uploadpart(svc, bucket, key_name, *result.UploadId, payload, int64(num_parts))
	assert.Nil(err)

	_, err = CompleteMultiUpload(svc, bucket, key_name, int64(num_parts), *result.UploadId, "")
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("InvalidPart", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

// func (suite *S3Suite) TestMultipartUploadNoSuchUpload() {
// 	/*
// 		Resource : object, method: get
// 		Scenario : check failure on multiple multi-part upload with invalid upload id
// 		Assertion: fails.
// 	*/
// 	assert := suite
// 	bucket := GetBucketName()
// 	num_parts := 2

// 	payload := strings.Repeat("12345", 1024*1024)
// 	key_name := "mymultipart"

// 	err := CreateBucket(svc, bucket)

// 	result, err := InitiateMultipartUpload(svc, bucket, key_name)
// 	fmt.Println("Result: ", result)

// 	resp, err := Uploadpart(svc, bucket, key_name, *result.UploadId, payload, int64(num_parts))

// 	assert.Nil(err)
// 	fmt.Println("Resp: ", resp)

// 	_, err = CompleteMultiUpload(svc, bucket, key_name, int64(num_parts), "*result.UploadId", *resp.ETag)
//
// 	assert.NotNil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal("NoSuchKey", awsErr.Code())
// 			assert.Equal("", awsErr.Message())
// 		}
// 	}
// }

func (suite *S3Suite) TestUploadPartNoSuchUpload() {

	/*
		Resource : object, method: get
		Scenario : check failure on multiple multi-part upload with invalid upload id
		Assertion: fails.
	*/
	assert := suite
	bucket := GetBucketName()
	num_parts := 2

	payload := strings.Repeat("12345", 1024*1024)
	key_name := "mymultipart"

	err := CreateBucket(svc, bucket)

	_, err = InitiateMultipartUpload(svc, bucket, key_name)

	_, err = Uploadpart(svc, bucket, key_name, "*result.UploadId", payload, int64(num_parts))
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchUpload", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

//.....................................MD5 headers..............................................................................

func (suite *S3Suite) TestObjectCreateBadMd5InvalidShort() {

	/*
		Resource : object, method: put
		Scenario : create w/invalid MD5.
		Assertion: fails.
	*/

	assert := suite
	headers := map[string]string{"Content-MD5": "YWJyYWNhZGFicmE="}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("InvalidDigest", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestObjectCreateBadMd5Bad() {

	/*
		Resource : object, method: put
		Scenario : create w/mismatched MD5.
		Assertion: fails.
	*/

	assert := suite
	headers := map[string]string{"Content-MD5": "rL0Y20zC+Fzt72VPzMSk2A=="}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("BadDigest", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestObjectCreateBadMd5Empty() {

	/*
		Resource : object, method: put
		Scenario : create w/empty MD5.
		Assertion: fails.
	*/

	assert := suite
	headers := map[string]string{"Content-MD5": " "}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("InvalidDigest", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestObjectCreateBadMd5Unreadable() {

	/*
		Resource : object, method: put
		Scenario : create w/non-graphics in MD5.
		Assertion: fails with invalid header field value
	*/

	assert := suite
	headers := map[string]string{"Content-MD5": "\x07"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("RequestError", awsErr.Code())
			assert.Equal("send request failed", awsErr.Message())
		}
	}

}

func (suite *S3Suite) TestObjectCreateBadMd5None() {

	/*
		Resource : object, method: put
		Scenario : create w/no MD5 header.
		Assertion: suceeds.
	*/

	assert := suite
	headers := map[string]string{}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)

}

//.........................................Expect Headers............................................................

func (suite *S3Suite) TestObjectCreateBadExpectMismatch() {

	/*
		Resource : object, method: put
		Scenario : create w/Expect 200.
		Assertion: garbage, but S3 succeeds!.
	*/

	assert := suite
	headers := map[string]string{"Expect": "200"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadExpectEmpty() {

	/*
		Resource : object, method: put
		Scenario : create w/empty expect.
		Assertion: succeeds...shouldnt IMHO!.
	*/

	assert := suite
	headers := map[string]string{"Expect": ""}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadExpectNone() {

	/*
		Resource : object, method: put
		Scenario : create w/no expect.
		Assertion: succeeds!.
	*/

	assert := suite
	headers := map[string]string{"Expect": ""}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadExpectUnreadable() {

	/*
		Resource : object, method: put
		Scenario : create w/non-graphic expect.
		Assertion: fails with invalid header field value
	*/

	assert := suite
	headers := map[string]string{"Expect": "\x07"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
}

//..........................................Content Length header............................................

func (suite *S3Suite) testObjectCreateBadContentlengthEmpty() {

	/*
		Resource : object, method: put
		Scenario : create w/empty content length.
		Assertion: fails!
	*/

	assert := suite
	headers := map[string]string{"Content-Length": " "}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("None", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestObjectCreateBadContentlengthNegative() {

	/*
		Resource : object, method: put
		Scenario : create w/negative content length.
		Assertion: fails.. but error message returned from SDK...I dont agree!!!
	*/

	assert := suite
	headers := map[string]string{"Content-Length": "-1"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("MissingContentLength", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestObjectCreateBadContentlengthNone() {

	/*
		Resource : object, method: put
		Scenario : create w/no content length.
		Assertion: suceeds
	*/

	assert := suite
	headers := map[string]string{"Content-Length": ""}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadContentlengthUnreadable() {

	/*
		Resource : object, method: put
		Scenario : create w/non-graphic content length.
		Assertion: fails
	*/

	assert := suite
	headers := map[string]string{"Content-Length": "\x07"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("MissingContentLength", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestObjectCreateBadContentlengthMismatchAbove() {

	/*
		Resource : object, method: put
		Scenario : create w/content length too long.
		Assertion: fails
	*/

	assert := suite
	content := "bar"
	length := fmt.Sprint(len(content) + 1)
	headers := map[string]string{"Content-Length": length}

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("MissingContentLength", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

//..................................Content-type header.........................................................

func (suite *S3Suite) TestObjectCreateBadContenttypevalid() {

	/*
		Resource : object, method: put
		Scenario : create w/content type text/plain.
		Assertion: suceeds!
	*/

	assert := suite
	headers := map[string]string{"Content-Type": "text/plain"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadContenttypeEmpty() {

	/*
		Resource : object, method: put
		Scenario : create w/empty content type.
		Assertion: suceeds!
	*/

	assert := suite
	headers := map[string]string{"Content-Type": " "}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadContenttypeNone() {

	/*
		Resource : object, method: put
		Scenario : create w/no content type.
		Assertion: suceeds!
	*/

	assert := suite
	headers := map[string]string{"Content-Type": ""}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
}

func (suite *S3Suite) TestObjectCreateBadContenttypeUnreadable() {

	/*
		Resource : object, method: put
		Scenario : create w/non-graphic content type.
		Assertion: fails with invalid header field value
	*/

	assert := suite
	headers := map[string]string{"Content-Type": "\x08"}
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.NotNil(err)
}

//..................................Authorization header.........................................................

func (suite *S3Suite) TestObjectCreateBadAuthorizationUnreadable() {

	/*
		Resource : object, method: put
		Scenario : create w/non-graphic authorization.
		Assertion: suceeds.... but should fail
			"Authorization" is in the ingnored header list, so its value does not matter
	*/

	assert := suite
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	headers := map[string]string{"Authorization": "\x01"}

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("AccessDenied", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestObjectCreateBadAuthorizationEmpty() {

	/*
		Resource : object, method: put
		Scenario :create w/empty authorization.
		Assertion: fails
	*/

	assert := suite
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	headers := map[string]string{"Authorization": " "}

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("AccessDenied", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *S3Suite) TestObjectCreateBadAuthorizationNone() {

	/*
		Resource : object, method: put
		Scenario :create w/no authorization.
		Assertion: fails
	*/

	assert := suite
	content := "bar"

	bucket := GetBucketName()
	key := "key1"
	err := CreateBucket(svc, bucket)

	headers := map[string]string{"Authorization": ""}

	err = SetupObjectWithHeader(svc, bucket, key, content, headers)
	assert.Nil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("AccessDenied", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
}

func (suite *HeadSuite) TestObjectListPrefixDelimiterPrefixDelimiterNotExist() {

	/*
		Resource : Object, method: ListObjects
		Scenario : list under prefix w/delimiter.
		Assertion: finds nothing w/unmatched prefix and delimiter.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "y"
	delimeter := "z"
	var empty_list []*s3.Object
	objects := map[string]string{"b/a/c": "echo", "b/a/g": "lima", "b/a/r": "golf", "g": "golf"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimeterAndPrefix(svc, bucket, prefix, delimeter)
	assert.Nil(errr)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal("NoSuchBucket", awsErr.Code())
			assert.Equal("", awsErr.Message())
		}
	}
	assert.Equal([]string{}, keys)
	assert.Equal([]string{}, prefixes)
	assert.Equal(empty_list, list.Contents)
}

func (suite *HeadSuite) TestObjectListPrefixDelimiterDelimiterNotExist() {

	/*
		Resource : object, method: list
		Scenario : list under prefix w/delimiter.
		Assertion: over-ridden slash ceases to be a delimiter.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "b"
	delimeter := "z"
	objects := map[string]string{"b/a/c": "echo", "b/a/g": "lima", "b/a/r": "golf", "golffie": "golfyy"}
	expectedkeys := []string{"b/a/c", "b/a/g", "b/a/r"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimeterAndPrefix(svc, bucket, prefix, delimeter)
	assert.Nil(errr)
	assert.Equal(3, len(list.Contents))
	assert.Equal(expectedkeys, keys)
	assert.Equal([]string{}, prefixes)
}

func (suite *HeadSuite) TestObjectListPrefixDelimiterPrefixNotExist() {

	/*
		Resource : object, method: list
		Scenario : list under prefix w/delimiter.
		Assertion: finds nothing w/unmatched prefix and delimiter.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "d"
	delimeter := "/"
	var empty_list []*s3.Object
	objects := map[string]string{"b/a/r": "echo", "b/a/c": "lima", "b/a/g": "golf", "g": "g"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimeterAndPrefix(svc, bucket, prefix, delimeter)
	assert.Nil(errr)
	assert.Equal([]string{}, keys)
	assert.Equal([]string{}, prefixes)
	assert.Equal(empty_list, list.Contents)
}

func (suite *HeadSuite) TestObjectListPrefixDelimiterAlt() {

	/*
		Resource : object, method: list
		Scenario : list under prefix w/delimiter.
		Assertion: non-slash delimiters.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "ba"
	delimeter := "a"
	objects := map[string]string{"bar": "echo", "bazar": "lima", "cab": "golf", "foo": "g"}
	expected_keys := []string{"bar"}
	expected_prefixes := []string{"baza"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimeterAndPrefix(svc, bucket, prefix, delimeter)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)
	assert.Equal(delimeter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)
}

func (suite *HeadSuite) TestObjectListPrefixDelimiterBasic() {

	/*
		Resource : object, method: list
		Scenario : list under prefix w/delimiter.
		Assertion: returns only objects directly under prefix.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "foo/"
	delimeter := "/"
	objects := map[string]string{"foo/bar": "echo", "foo/baz/xyzzy": "lima", "quux/thud": "golf"}
	expected_keys := []string{"foo/bar"}
	expected_prefixes := []string{"foo/baz/"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimeterAndPrefix(svc, bucket, prefix, delimeter)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(delimeter, *list.Delimiter)
	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)
}

func (suite *HeadSuite) TestObjectListPrefixUnreadable() {

	/*
		Resource : object, method: list
		Scenario : list under prefix.
		Assertion: non-printable prefix can be specified.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "\x0a"
	objects := map[string]string{"foo/bar": "echo", "foo/baz/xyzzy": "lima", "quux/thud": "golf"}
	expected_keys := []string{}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithPrefix(svc, bucket, prefix)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(expected_prefixes, prefixes)
	assert.Equal(expected_keys, keys)

}

func (suite *HeadSuite) TestObjectListPrefixNotExist() {

	/*
		Resource : object, method: List
		Scenario : list under prefix.
		Assertion: nonexistent prefix returns nothing.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "d"
	objects := map[string]string{"foo/bar": "echo", "foo/baz": "lima", "quux": "golf"}
	expected_keys := []string{}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithPrefix(svc, bucket, prefix)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListPrefixNone() {

	/*
		Resource : object, method: list
		Scenario : list under prefix.
		Assertion: unspecified prefix returns everything.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := ""
	objects := map[string]string{"foo/bar": "echo", "foo/baz": "lima", "quux": "golf"}
	expected_keys := []string{"foo/bar", "foo/baz", "quux"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithPrefix(svc, bucket, prefix)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)
}

func (suite *HeadSuite) TestObjectListPrefixEmpty() {

	/*
		Resource : object, method: list
		Scenario : list under prefix.
		Assertion: empty prefix returns everything.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := ""
	objects := map[string]string{"foo/bar": "echo", "foo/baz": "lima", "quux": "golf"}
	expected_keys := []string{"foo/bar", "foo/baz", "quux"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithPrefix(svc, bucket, prefix)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListPrefixAlt() {

	/*
		Resource : object, method: list
		Scenario : list under prefix.
		Assertion: prefixes w/o delimiters.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "ba"
	objects := map[string]string{"bar": "echo", "baz": "lima", "foo": "golf"}
	expected_keys := []string{"bar", "baz"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithPrefix(svc, bucket, prefix)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListPrefixBasic() {

	/*
		Resource : bucket, method: get
		Scenario : list under prefix.
		Assertion: returns only objects under prefix.
	*/

	assert := suite
	bucket := GetBucketName()
	prefix := "foo/"
	objects := map[string]string{"foo/bar": "echo", "foo/baz": "lima", "quux": "golf"}
	expected_keys := []string{"foo/bar", "foo/baz"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithPrefix(svc, bucket, prefix)
	assert.Nil(errr)
	assert.Equal(prefix, *list.Prefix)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterNotExist() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: unused delimiter is not found.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := "/"
	objects := map[string]string{"bar": "echo", "baz": "lima", "cab": "golf", "foo": "golf"}
	expected_keys := []string{"bar", "baz", "cab", "foo"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterNone() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: unspecified delimiter defaults to none.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := " "
	objects := map[string]string{"bar": "echo", "baz": "lima", "cab": "golf", "foo": "golf"}
	expected_keys := []string{"bar", "baz", "cab", "foo"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterEmpty() {

	// Resource : object, method: list
	// Scenario : list .
	// Assertion: empty delimiter can be specified.

	assert := suite
	bucket := GetBucketName()
	delimiter := " "
	objects := map[string]string{"bar": "echo", "baz": "lima", "cab": "golf", "foo": "golf"}
	expected_keys := []string{"bar", "baz", "cab", "foo"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterUnreadable() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: non-printable delimiter can be specified.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := "\x0a"
	objects := map[string]string{"bar": "echo", "baz": "lima", "cab": "golf", "foo": "golf"}
	expected_keys := []string{"bar", "baz", "cab", "foo"}
	expected_prefixes := []string{}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterDot() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: dot delimiter characters.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := "."
	objects := map[string]string{"b.ar": "echo", "b.az": "lima", "c.ab": "golf", "foo": "golf"}
	expected_keys := []string{"foo"}
	expected_prefixes := []string{"b.", "c."}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(2, len(prefixes))
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterPercentage() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: percentage delimiter characters.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := "%"
	objects := map[string]string{"b%ar": "echo", "b%az": "lima", "c%ab": "golf", "foo": "golf"}
	expected_keys := []string{"foo"}
	expected_prefixes := []string{"b%", "c%"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(2, len(prefixes))
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterWhiteSpace() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: whitespace delimiter characters.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := " "
	objects := map[string]string{"b ar": "echo", "b az": "lima", "c ab": "golf", "foo": "golf"}
	expected_keys := []string{"foo"}
	expected_prefixes := []string{"b ", "c "}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(2, len(prefixes))
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterAlt() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: non-slash delimiter characters.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := "a"
	objects := map[string]string{"bar": "echo", "baz": "lima", "cab": "golf", "foo": "golf"}
	expected_keys := []string{"foo"}
	expected_prefixes := []string{"ba", "ca"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(delimiter, *list.Delimiter)

	assert.Equal(expected_keys, keys)
	assert.Equal(len(prefixes), 2)
	assert.Equal(prefixes, expected_prefixes)

}

func (suite *HeadSuite) TestObjectListDelimiterBasic() {

	/*
		Resource : object, method: list
		Scenario : list .
		Assertion: prefixes in multi-component object names.
	*/

	assert := suite
	bucket := GetBucketName()
	delimiter := "/"
	objects := map[string]string{"foo/bar": "echo", "foo/baz/xyzzy": "lima", "quux/thud": "golf", "asdf": "golf"}
	expected_keys := []string{"asdf"}
	expected_prefixes := []string{"foo/", "quux/"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	list, keys, prefixes, errr := ListObjectsWithDelimiter(svc, bucket, delimiter)
	assert.Nil(errr)
	assert.Equal(*list.Delimiter, delimiter)

	assert.Equal(keys, expected_keys)
	assert.Equal(2, len(prefixes))
	assert.Equal(expected_prefixes, prefixes)

}

func (suite *HeadSuite) TestObjectListMaxkeysNone() {

	/*
		Resource : Object, Method: list
		Operation : List all keys
		Assertion : pagination w/o max_keys.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"key1": "echo", "key2": "lima", "key3": "golf"}
	ExpectedKeys := []string{"key1", "key2", "key3"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, err := GetObjects(svc, bucket)
	assert.Nil(err)

	keys := []string{}
	for _, key := range resp.Contents {
		keys = append(keys, *key.Key)
	}
	assert.Equal(ExpectedKeys, keys)
	assert.Equal(int64(1000), *resp.MaxKeys)
	assert.Equal(false, *resp.IsTruncated)
}

func (suite *HeadSuite) TestObjectListMaxkeysZero() {

	/*
		Resource : object, method: get
		Operation : List all keys .
		Assertion: pagination w/max_keys=0.
	*/

	assert := suite
	bucket := GetBucketName()
	maxkeys := int64(0)
	objects := map[string]string{"key1": "echo", "key2": "lima", "key3": "golf"}
	ExpectedKeys := []string(nil)

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMaxKeys(svc, bucket, maxkeys)
	assert.Nil(errr)
	assert.Equal(ExpectedKeys, keys)
	assert.Equal(false, *resp.IsTruncated)
}

func (suite *HeadSuite) TestObjectListMaxkeysOne() {

	/*
		Resource : bucket, method: get
		Operation : List keys all keys.
		Assertion: pagination w/max_keys=1, marker.
	*/

	assert := suite
	bucket := GetBucketName()
	maxkeys := int64(1)
	objects := map[string]string{"key1": "echo", "key2": "lima", "key3": "golf"}
	EKeysMaxkey := []string{"key1"}
	EKeysMarker := []string{"key2", "key3"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMaxKeys(svc, bucket, maxkeys)
	assert.Nil(errr)
	assert.Equal(EKeysMaxkey, keys)
	assert.Equal(true, *resp.IsTruncated)

	resp, keys, errs := GetKeysWithMarker(svc, bucket, EKeysMaxkey[0])
	assert.Nil(errs)
	assert.Equal(false, *resp.IsTruncated)
	assert.Equal(EKeysMarker, keys)

}

//............................................Test Get object with marker...................................

func (suite *HeadSuite) TestObjectListMarkerBeforeList() {

	/*
		Resource : object, method: get
		Scenario : list all objects.
		Assertion: marker before list.
	*/

	assert := suite
	bucket := GetBucketName()
	marker := "aaa"
	objects := map[string]string{"bar": "echo", "baz": "lima", "quux": "golf"}
	expected_keys := []string{"bar", "baz", "quux"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMarker(svc, bucket, marker)
	assert.Nil(errr)
	assert.Equal(marker, *resp.Marker)
	assert.Equal(expected_keys, keys)
	assert.Equal(false, *resp.IsTruncated)

	err = DeleteObjects(svc, bucket)
	err = DeleteBucket(svc, bucket)
	assert.Nil(err)

}

func (suite *HeadSuite) TestObjectListMarkerAfterList() {

	/*
		Resource : object, method: get
		Scenario : list all objects.
		Assertion: marker after list.
	*/

	assert := suite
	bucket := GetBucketName()
	marker := "zzz"
	objects := map[string]string{"bar": "echo", "baz": "lima", "quux": "golf"}
	expected_keys := []string(nil)

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMarker(svc, bucket, marker)
	assert.Nil(errr)
	assert.Equal(marker, *resp.Marker)
	assert.Equal(false, *resp.IsTruncated)
	assert.Equal(expected_keys, keys)

}

func (suite *HeadSuite) TestObjectListMarkerNotInList() {

	/*
		Resource : object, method: get
		Scenario : list all objects.
		Assertion: marker not in list.
	*/

	assert := suite
	bucket := GetBucketName()
	marker := "blah"
	objects := map[string]string{"bar": "echo", "baz": "lima", "quux": "golf"}
	expected_keys := []string{"quux"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMarker(svc, bucket, marker)
	assert.Nil(errr)
	assert.Equal(marker, *resp.Marker)
	assert.Equal(expected_keys, keys)
}

func (suite *HeadSuite) TestObjectListMarkerUnreadable() {

	/*
		Resource : object, method: get
		Scenario : list all objects.
		Assertion: non-printing marker.
	*/

	assert := suite
	bucket := GetBucketName()
	marker := "\x0a"
	objects := map[string]string{"bar": "echo", "baz": "lima", "quux": "golf"}
	expected_keys := []string{"bar", "baz", "quux"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMarker(svc, bucket, marker)
	assert.Nil(errr)
	assert.Equal(marker, *resp.Marker)
	assert.Equal(false, *resp.IsTruncated)
	assert.Equal(expected_keys, keys)

}

func (suite *HeadSuite) TestObjectListMarkerEmpty() {

	/*
		Resource : object, method: get
		Scenario : list all objects.
		Assertion: no pagination, empty marker.
	*/

	assert := suite
	bucket := GetBucketName()
	marker := ""
	objects := map[string]string{"bar": "echo", "baz": "lima", "quux": "golf"}
	expected_keys := []string{"bar", "baz", "quux"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMarker(svc, bucket, marker)
	assert.Nil(errr)
	assert.Equal(marker, *resp.Marker)
	assert.Equal(false, *resp.IsTruncated)
	assert.Equal(expected_keys, keys)

}

func (suite *HeadSuite) TestObjectListMarkerNone() {

	/*
		Resource : object, method: get
		Scenario : list all objects.
		Assertion: no pagination, no marker.
	*/

	assert := suite
	bucket := GetBucketName()
	marker := ""
	objects := map[string]string{"bar": "echo", "baz": "lima", "quux": "golf"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, _, errr := GetKeysWithMarker(svc, bucket, marker)
	assert.Nil(errr)
	assert.Equal(marker, *resp.Marker)

}

func (suite *HeadSuite) TestObjectListMany() {

	/*
		Resource : object, method: list
		Scenario : list all keys
		Assertion: pagination w/max_keys=2, no marker.
	*/

	assert := suite
	bucket := GetBucketName()
	maxkeys := int64(2)
	keys := []string{}
	objects := map[string]string{"foo": "echo", "bar": "lima", "baz": "golf"}
	expected_keys := []string{"bar", "baz"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, keys, errr := GetKeysWithMaxKeys(svc, bucket, maxkeys)
	assert.Nil(errr)
	assert.Equal(2, len(resp.Contents))
	assert.Equal(true, *resp.IsTruncated)
	assert.Equal(expected_keys, keys)

	resp, keys, errs := GetKeysWithMarker(svc, bucket, expected_keys[1])
	assert.Nil(errs)
	assert.Equal(1, len(resp.Contents))
	assert.Equal(false, *resp.IsTruncated)
	expected_keys = []string{"foo"}

}

func (suite *HeadSuite) TestObjectHeadZeroBytes() {

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{"bar": ""}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)

	resp, err := GetObject(svc, bucket, "bar")
	assert.Nil(err)
	assert.Equal(0, len(resp))
}

func (suite *HeadSuite) TestObjectCreateUnreadable() {

	/*
		Resource : object, method: put
		Scenario : write to non-printing key
		Assertion: passes.
	*/

	assert := suite
	bucket := GetBucketName()
	objects := map[string]string{string('\x0a'): "echo"}

	err := CreateBucket(svc, bucket)
	err = CreateObjects(svc, bucket, objects)
	assert.Nil(err)
}
