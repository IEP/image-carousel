package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"sync"
)

type StaticImageServer struct {
	Metadata   map[string][]Image
	BucketSize map[string]int
	BucketList []string
	BasePath   string

	mu sync.Mutex
}

// NewStaticImageServer based on metadata.json
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
		size[k] = len(v)
		list = append(list, k)
	}

	return &StaticImageServer{
		Metadata:   data,
		BucketSize: size,
		BucketList: list,
		BasePath:   basePath,
	}, nil
}

var _ ImageServer = &StaticImageServer{}

func (s *StaticImageServer) GetBucketsName() []string {
	return s.BucketList
}

func (s *StaticImageServer) GetRandomImage(bucketName string) Image {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := rand.Intn(s.BucketSize[bucketName])
	img := s.Metadata[bucketName][idx]
	img.PhotoPath = path.Join(s.BasePath, bucketName, img.PhotoPath)

	return img
}