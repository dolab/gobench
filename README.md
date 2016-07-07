# S3API Benchmark

### Usage

```bash
// for sample
$bench -s3.config=./config.json

// enable debug
$bench -s3.config=.=/config.json -s3.debug=true
```

### Config

```json
{
    "host": "https://s3.amazonaws.com",
    "region": "us-east-1",
    "accessKeyId": "",
    "accessKeySecret": "",
    "workflow": [
        {
            "stage": "init",
            "concurrent": 1,
            "bucketPrefix": "s3bucket",
            "bucketStart": 1,
            "bucketEnd": 9,
            "enabled": false
        },
        {
            "stage": "write",
            "concurrent": 10,
            "file": "./README.md",
            "filesize": 65536,
            "bucketPrefix": "s3bucket",
            "bucketStart": 1,
            "bucketEnd": 9,
            "objectPrefix": "myobjects",
            "objectStart": 1,
            "objectEnd": 10,
            "hashCheck": true,
            "enabled": false
        },
        {
            "stage": "read",
            "concurrent": 10,
            "bucketPrefix": "s3bucket",
            "bucketStart": 1,
            "bucketEnd": 1,
            "objectPrefix": "myobjects",
            "objectStart": 1,
            "objectEnd": 10,
            "hashCheck": true,
            "enabled": true
        },
        {
            "stage": "clean",
            "concurrent": 10,
            "bucketPrefix": "s3bucket",
            "bucketStart": 1,
            "bucketEnd": 9,
            "objectPrefix": "myobjects",
            "objectStart": 1,
            "objectEnd": 10,
            "enabled": false
        },
        {
            "stage": "dispose",
            "concurrent": 1,
            "bucketPrefix": "s3bucket",
            "bucketStart": 1,
            "bucketEnd": 9,
            "enabled": false
        }
    ]
}

```
### Development

```bash
$git clone git@github.com:dolab/gobench.git
$cd gobench
$source env.sh
$make
```
