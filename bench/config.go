package bench

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	StageInit    = "init"
	StageWrite   = "write"
	StageRead    = "read"
	StageClean   = "clean"
	StageDispose = "dispose"
)

type WorkConfig struct {
	Stage        string `json:"stage"`
	Concurrent   int    `json:"concurrent"`
	File         string `json:"file"`
	Filesize     int    `json:"filesize"`
	BucketPrefix string `json:"bucketPrefix"`
	BucketStart  int    `json:"bucketStart"` // bucket name start number
	BucketEnd    int    `json:"bucketEnd"`   // bucket name end number
	ObjectPrefix string `json:"ObjectPrefix"`
	ObjectStart  int    `json:"ObjectStart"` // object name start number
	ObjectEnd    int    `json:"ObjectEnd"`   // object name end number
	HashCheck    bool   `json:"hashCheck"`   // only for write|read work stage
	Enabled      bool   `json:"enabled"`     // does this work enabled?
}

func (wc *WorkConfig) BucketNumber() int {
	if wc.BucketEnd < wc.BucketStart {
		return 0
	}

	return wc.BucketEnd - wc.BucketStart + 1
}

func (wc *WorkConfig) ObjectNumber() int {
	if wc.ObjectEnd < wc.ObjectStart {
		return 0
	}

	return wc.ObjectEnd - wc.ObjectStart + 1
}

func (wc *WorkConfig) ObjectTotals() int {
	return wc.BucketNumber() * wc.ObjectNumber()
}

type Config struct {
	Host            string `json:"host"`
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`

	Workflow []*WorkConfig `json:"workflow"`
}

func NewConfig(fpath string) (cfg *Config, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &cfg)
	return
}
