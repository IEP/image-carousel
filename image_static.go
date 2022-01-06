package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"sync"
)

// StaticImageServer implements ImageServer
type StaticImageServer struct {
	Metadata   map[string][]Image
	BucketSize map[string]int
	BucketList []string
	BasePath   string

	bucketListWithoutRandom []string

	bucketOnce sync.Once
}

// NewStaticImageServer based on metadata.json with basePath expressing the location of metadata.json
func NewStaticImageServer(basePath string) (ImageServer, error) {
	metadataPath := path.Join(basePath, "metadata.json")
	content, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %w", err)
	}

	var data map[string][]Image
	if err = json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	size := make(map[string]int)
	list := make([]string, 0)

	for k, v := range data {
		if k == "random" {
			return nil, errors.New("bucket name 'random' is reserved")
		}
		size[k] = len(v)
		list = append(list, k)
	}

	list = append(list, "random") // random pseudo-bucket

	return &StaticImageServer{
		Metadata:   data,
		BucketSize: size,
		BucketList: list,
		BasePath:   basePath,
	}, nil
}

// check StaticImageServer implementation whether it satisfy ImageServer interface
var _ ImageServer = &StaticImageServer{}

// GetBucketsName list
func (s *StaticImageServer) GetBucketsName() []string {
	return s.BucketList
}

func (s *StaticImageServer) getBucketsNameWithoutRandom() []string {
	s.bucketOnce.Do(func() {
		s.bucketListWithoutRandom = make([]string, 0)
		for _, b := range s.GetBucketsName() {
			b := b
			if b == "random" {
				continue
			}
			s.bucketListWithoutRandom = append(s.bucketListWithoutRandom, b)
		}
	})

	return s.bucketListWithoutRandom
}

// GetRandomImage based on the bucket name choosen
func (s *StaticImageServer) GetRandomImage(bucketName string) Image {
	// handle random bucket
	if bucketName == "random" {
		buckets := s.getBucketsNameWithoutRandom()
		idx := rand.Intn(len(buckets))
		bucketName = buckets[idx]
	}

	// handle invalid bucket
	bucketLen := s.BucketSize[bucketName]
	if bucketLen == 0 {
		return Image{}
	}

	idx := rand.Intn(s.BucketSize[bucketName])
	img := s.Metadata[bucketName][idx]
	img.PhotoPath = path.Join(s.BasePath, bucketName, img.PhotoPath)

	return img
}
