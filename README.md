
 ## S3 compatibility tests

This is a set of completely unofficial Amazon AWS S3 compatibility
tests, that will hopefully be useful to people implementing software
that exposes an S3-like API.

The tests only cover the REST interface.


### Download and setup the tests

Clone the repository

	git clone https://github.com/nanjekyejoannah/go_s3tests
	cd go_s3tests
	./bootstrap.sh

### Configuration

Copy the sample configuration.

	cp config.yaml.sample config.yaml

Edit the config.yaml.sample file to your needs.  

The config file should look  like this:

	
	DEFAULT :

    		host : s3.amazonaws.com
    		port : 8080
    		is_secure : yes

	fixtures :
    		bucket_prefix : yournamehere-

	s3main :
    		access_key : 0555b35654ad1656d804
    		access_secret : h7GhxuBLTrlhVUyxSPUKUV8r/2EI4ngqJxD7iBdBYLhwluN30JaT3Q==
    		bucket : bucket1
    		region : us-east-1
    		endpoint : localhost:8000
    		host : localhost
    		port : 8000
    		display_name :
    		email : tester@gmail.com
    		is_secure : false
    		SSE : AES256
    		kmskeyid : barbican_key_id

	s3alt :
    		access_key : NOPQRSTUVWXYZABCDEFG
    		access_secret : nopqrstuvwxyzabcdefghijklmnabcdefghijklm
    		bucket : bucket1
    		region : us-east-1
    		endpoint : localhost:8000
    		display_name :
    		email : johndoe@gmail.com
    		SSE : your SSE
    		kmskeyid : barbican_key_id
    		is_secure : false


### RGW

The tests connect to the Ceph RGW ,therefore you shoud have started your RGW and use the credentials you get. Details on building Ceph and starting RGW can be found in the [ceph repository](https://github.com/ceph/ceph).

### Gopath and Dependencies

You need to set your GoPath . Details on setting up Go environments can be found [here](https://golang.org/doc/install)
	
	export GOPATH=$HOME/go

#### Installing dependencies

You should be in the project root folder to run this.

	 go get -d ./...

#### To run the tests
	
	cd s3tests
	go test -v  

#### To Do

+ Host Style 
+ Versioning 			 	
