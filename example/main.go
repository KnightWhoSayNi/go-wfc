package main

import (
	"fmt"
	"time"

	"github.com/zfedoran/go-wfc/pkg/wfc"
)

func main() {
	input_tileset := "./example/internal/input"
	output_image := "./example/internal/output/%d.png"

	// Print the adjacency hash values for the provided tileset.
	printAdjacencyHashValues(input_tileset)

	// Generate an image from the tileset.
	collapseWave(input_tileset, output_image)
}

func collapseWave(tileset_folder, output_image string) {
	// This is just a `[]image.Image`, you can use whatever loader function you'd like
	images, err := wfc.LoadImageFolder(tileset_folder)
	if err != nil {
		panic(err)
	}

	// The random seed to use when collapsing the wave
	// (given the same seed number, the Collapse() fn would generate the same state every time)
	seed := int(time.Now().UnixNano())

	width := 8
	height := 8

	// Setup the initialized state
	wave := wfc.New(images, width, height)
	wave.Initialize(seed)

	// Force the top tiles to be sky
	sky := wfc.GetConstraintFromHex("8cbc4760")
	for i := 0; i < width; i++ {
		slot := wave.PossibilitySpace[i]
		modules := make([]*wfc.Module, 0)
		for _, m := range slot.Superposition {
			if m.Adjacencies[wfc.Up] == sky {
				modules = append(modules, m)
			}
		}
		slot.Superposition = modules
	}

	// Collapse the wave function (make up to 100 attempts)
	err = wave.Collapse(200)
	if err != nil {
		// don't panic here, we want to generate the image anyway
		fmt.Printf("unable to generate: %v", err)
	}

	// Lets generate an image
	output := wave.ExportImage()
	output_file := fmt.Sprintf(output_image, seed)

	wfc.SaveImage(output_file, output)
	fmt.Printf("Image saved to: %s\n", output_file)
}

func printAdjacencyHashValues(input_tileset string) {
	fmt.Printf("Adjacency hash values:\n\n")

	images, err := wfc.LoadImageFolder(input_tileset)
	if err != nil {
		panic(err)
	}

	// We could use pretty table to do this, but this is just a demo and I don't
	// want the extra dependency.

	fmt.Println("|-------|----------|----------|")
	fmt.Println("|Tile\t|Direction |Hash      |")
	fmt.Println("|-------|----------|----------|")
	for i, img := range images {
		for _, d := range wfc.Directions {
			fmt.Printf("|%d\t|%s\t   | %s | %dx%d\n", i, d.ToString(), wfc.GetConstraintId(img, d), img.Bounds().Max.X, img.Bounds().Max.Y)
		}
		fmt.Printf("|- - - -|- - - - - |- - - - - |\n")
	}
	fmt.Printf("|-------|----------|----------|\n\n")
}
