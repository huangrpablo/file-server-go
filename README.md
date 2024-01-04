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
* `app.uploadInChunk`: whether a file will be uploaded in parts to minio (default: false)
* `app.bucketName`: the bucket where files are stored in minio (default: taurus-file-store)

## Design choices
1. Filename as a unique file id
   1. For simplicity, filename of the uploaded file is chosen as the object key in minio, which is decided by the user. Two files with the same name will overwrite the object in the bucket. But as the versioning can be enabled, the overwritten file can still be retrieved (not implemented via API yet). The download of the file also uses the filename.
   2. However, it is also an option to let the server generate a unique file id as the object key, and return the id in the HTTP response. The user uses the returned id to access and manipulate the resource (an API for file modification needs to be implemented).

2. Upload in chunks done by Minio API
   1. Minio supports multipart upload by specifying `PartSize` (controlled via `chunkSize` in `config/config.yaml`).
   2. However, the minimum chunk size is 5MB, while in the instruction, 1MB should be possible. 
   3. To achieve this, the chunking can be implemented by the server itself (not implemented yet).

3. Input validation done by gin binding
4. Simple AES encryption (can be improved by using more custom and complex encryption strategy).

## TODO
- [ ] Build the file server app into a docker image
- [ ] Find out why minio actions are so slow: for a 500MB file, it takes 4 ~ 9 seconds to upload (it could be the limited resources allocated)
- [ ] Write more unit tests, some with dependency injection
- [ ] More things can be configured, such as the cryptography algorithm

