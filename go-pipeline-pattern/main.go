package main

import (
	"fmt"
	imageprocessing "go-pipeline-pattern/image_processing"
	"image"
	"strings"
)

type Job struct {
	InputPath  string
	Image      image.Image
	OutputPath string
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)

	go func() {
		for _, p := range paths {
			job := Job{
				InputPath:  p,
				OutputPath: strings.Replace(p, "images/", "images/output/", 1),
			}

			job.Image = imageprocessing.ReadImage(p)

			out <- job
		}

		close(out)
	}()

	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)

	go func() {
		for job := range input {
			job.Image = imageprocessing.Resize(job.Image)

			out <- job
		}

		close(out)
	}()

	return out
}

func convertToGrayScale(input <-chan Job) <-chan Job {
	out := make(chan Job)

	go func() {
		for job := range input {
			job.Image = imageprocessing.GrayScale(job.Image)

			out <- job
		}

		close(out)
	}()

	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)

	go func() {
		for job := range input {
			imageprocessing.WriteImage(job.OutputPath, job.Image)
			out <- true
		}

		close(out)
	}()

	return out
}

func main() {
	imagePaths := []string{
		"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayScale(channel2)
	writeResults := saveImage(channel3)

	for writeResult := range writeResults {
		if writeResult == true {
			fmt.Println("Success")
		} else {
			fmt.Println("failed")
		}
	}

}
