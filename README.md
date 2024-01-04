# file-server-go
An application that exposes an API accepting and serving files.

## Getting Started

### Prerequisites

- Go installed on your machine
- Docker and Docker Compose for running Minio locally

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/huangrpablo/file-server-go.git
   cd file-server-go
   
2. Start Minio using Docker Compose:

    ```bash
   docker-compose up -d

Minio will be available at http://127.0.0.1:9000 with default credentials `minioadmin:minioadmin`.

3. Make a bucket in Minio for file storage. The default bucket is `taurus-file-store`, which can be configured in `config/config.yaml`.

4. Build and run the File Server App:

    ```bash
   go run main.go

The File Server App should now be running at http://127.0.0.1:8080.

## API definition

### Upload a file

* Endpoint: `POST /v1/files`
* Request:
    ```bash
    curl -X POST -F "file=@/path/to/your/file.txt" http://127.0.0.1:8080/v1/files

* Response:
    ```json
    {"message":"File uploaded successfully"}
  
### Download a file

* Endpoint: `GET /v1/file/{fileName}`
* Request:
    ```bash
    curl http://127.0.0.1:8080/v1/file/{fileName} -o myfile
* Response:
  File contents streamed to the `myfile`

## Configuration
The configuration can be made in `config/config.yaml`.

* `app.chunkSize`: the chunkSize for multipart uploading to minio (default: 5242880)
* `app.bucketName`: the bucket where files are stored in minio (default: taurus-file-store)

## Design choices


## TODO
- [ ] Build the file server app into a docker image
- [ ] Find out why minio actions are so slow: for a 500MB file, it takes 4 ~ 9 seconds to upload (it could be the limited resources allocated)
- [ ] Write more unit tests, some with dependency injection
- [ ] More things can be configured, such as the cryptography algorithm

