
 ## S3 compatibility tests

This is a set of integration tests for the S3 (AWS) interface of [RGW](http://docs.ceph.com/docs/mimic/radosgw/). 

It might also be useful for people implementing software that exposes an S3-like API.

The test suite only covers the REST interface and uses [GO amazon SDK](https://aws.amazon.com/sdk-for-go/) version 1.11.364 and [Golang Environment setup](https://golang.org/doc/install).

### Get the source code

Clone the repository

	git clone https://github.com/adamyanova/go_s3tests

### Edit Configuration

	cd go_s3tests
	cp config.yaml.sample config.yaml

The config file should look like this:

	DEFAULT :

		host : localhost 
		port : 8000
		is_secure : false

	fixtures :

		bucket_prefix : test-

	s3main :

		access_key : 0555b35654ad1656d804
		access_secret : h7GhxuBLTrlhVUyxSPUKUV8r/2EI4ngqJxD7iBdBYLhwluN30JaT3Q==
		bucket : bucket1
		region : us-east-1
		endpoint : localhost:8000
		host : localhost
		port : 8000
		display_name :
		email : someone@gmail.com
		is_secure : false
		SSE : aws:kms 
		kmskeyid : testkey-1 


	s3alt :

		access_key : NOPQRSTUVWXYZABCDEFG
		access_secret : nopqrstuvwxyzabcdefghijklmnabcdefghijklm
		bucket : bucket1
		region : us-east-1
		endpoint : localhost:8000
		display_name :
		email : someone@gmail.com
		SSE : your SSE
		kmskeyid : testkey-1
		is_secure : false

The credentials match the default S3 test users created by RGW.

#### RGW

The tests connect to the Ceph RGW, therefore one shoud start RGW beforehand and use the provided credentials. Details on building Ceph and starting RGW can be found in the [ceph repository](https://github.com/ceph/ceph).

The **s3tests.teuth.config.yaml** files is required for the Ceph test framework [Teuthology](http://docs.ceph.com/teuthology/docs/README.html). 
It is irrelevant for standalone testing.

### Install prerequisits
#### Golang
The **boostrap.sh** script will install **golang**.

The GOPATH variable should beset before running. Details on setting up Go environments can be found [here](https://golang.org/doc/install)
	
	export GOPATH=$HOME/go

#### Test dependencies
	cd 
	go get -v -d ./...
	go get -v github.com/stretchr/testify

### Run the Tests

To run all tests:

	cd s3tests
	go test -v  

To run a specific test e.g. TestSignWithBodyReplaceRequestBody():
	
	cd s3tests
	go test -v -run TestSuite/TestSignWithBodyReplaceRequestBody

To run all tests with "TestSSEKMS" in their name:

	cd s3tests
	go test -v -run TestSuite/TestSSEKMS

**Using SSL**
The server certificate must be present in the cetrificate pool of the system on which the tests are executed.
