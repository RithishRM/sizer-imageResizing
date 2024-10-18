package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

// Save the original image in its original format.
func saveOriginalImage(src image.Image, outputPath string) {
	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", outputPath, err)
	}
	defer output.Close()

	switch img := src.(type) {
	case *image.NRGBA:
		err = png.Encode(output, img) // Save as PNG
	case *image.YCbCr:
		err = jpeg.Encode(output, img, nil) // Save as JPEG
	case *image.GIF:
		err = gif.Encode(output, img, nil) // Save as GIF
	default:
		log.Fatalf("unsupported image format: %T", img)
	}
	if err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}
}

// Resize the image to specified dimensions.
func resizeImageCustom(src image.Image, width, height int) image.Image {
	return imaging.Resize(src, width, height, imaging.Lanczos)
}

// Handle image upload and save the original image.
func GetHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err),
			http.StatusInternalServerError)
		return
	}
	defer file.Close()

	src, err := imaging.Decode(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding image: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Save the original image
	saveOriginalImage(src, "Images/Original/original.png") // You can adjust the path based on the format

	fmt.Fprintf(w, "Image uploaded successfully.")
}

// Handle image resizing based on the requested format.
func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	width := r.FormValue("width")  // Get width from request
	height := r.FormValue("height") // Get height from request
	format := r.FormValue("format") // Get output format

	srcFile, err := os.Open("Images/Original/original.png")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening original image file: %v", err),
			http.StatusInternalServerError)
		return
	}
	defer srcFile.Close()

	src, err := imaging.Decode(srcFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding original image file: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Resize the image to specified dimensions
	customImage := resizeImageCustom(src, 300, 200) // Replace with width/height parsing

	// Set content type and encode image based on format choice
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
		err = png.Encode(w, customImage)
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, customImage, nil)
	case "gif":
		w.Header().Set("Content-Type", "image/gif")
		err = gif.Encode(w, customImage, nil)
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding image: %v", err),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Image resized successfully.")
}

func main() {
	fileServer := http.FileServer(http.Dir("./api"))
	http.Handle("/", fileServer)
	http.HandleFunc("/Upload", GetHandler)
	http.HandleFunc("/Resize", ResizeHandler)

	fmt.Printf("Server started at Port number : 8000\n")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
