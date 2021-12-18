package main

import (
	"testing"
)

func TestStaticImageServer(t *testing.T) {
	s, err := NewStaticImageServer("data")
	if err != nil {
		t.Errorf("NewStaticImageServer: %s", err)
	}

	t.Run("get_buckets_name", func(t *testing.T) {
		names := s.GetBucketsName()
		t.Logf("bucket names: %+v", names)
		size := len(names)
		if size != 3 { // 3 is hardcoded
			t.Errorf("len(s.GetBucketsName) != 3: %d", size)
		}
	})

	t.Run("get_random_image", func(t *testing.T) {
		img1 := s.GetRandomImage("ugm")
		img2 := s.GetRandomImage("ugm")
		t.Logf("image 1: %+v", img1)
		t.Logf("image 2: %+v", img2)

		if img1.Description == img2.Description || img1.PhotoPath == img2.PhotoPath {
			t.Error("img1 should not equal to img2")
		}
	})
}
