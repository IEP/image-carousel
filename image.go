package main

// Image object that consists of photo path and its description
type Image struct {
	PhotoPath   string `json:"photo_path"`
	Description string `json:"description"`
}

// ImageServer interface to manage the buckets and image assets
type ImageServer interface {
	GetRandomImage(bucketName string) Image
	GetBucketsName() []string
}
