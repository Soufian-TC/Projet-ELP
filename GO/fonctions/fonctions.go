package fonctions

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sync"
)

// FlouGaussienOptimise applique un flou gaussien à une image donnée en optimisant le traitement
// selon les dimensions de l'image (largeur ou hauteur) et en utilisant des goroutines.
// Paramètres :
// - fImage : image.Image, l'image d'entrée à flouter.
// - sigma : float64, l'écart-type utilisé pour générer le noyau gaussien.
// Retourne : *image.RGBA, l'image floutée.
func FlouGaussienOptimise(fImage image.Image, sigma float64) *image.RGBA {
	bounds := fImage.Bounds()
	largeur := bounds.Dx()
	hauteur := bounds.Dy()
	kernel := Noyeau(10, sigma)
	blurred := image.NewRGBA(bounds)

	var wg sync.WaitGroup

	if largeur >= hauteur {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			wg.Add(1)
			go func(y int) {
				defer wg.Done()
				FlouGaussienUneLigne(fImage, kernel, y, blurred)
			}(y)
		}
	} else {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			wg.Add(1)
			go func(x int) {
				defer wg.Done()
				FlouGaussienUneColonne(fImage, kernel, x, blurred)
			}(x)
		}
	}

	wg.Wait()
	return blurred
}

// OuvrirImage ouvre un fichier image à partir d'un chemin donné et le décode.
// Paramètres :
// - cheminImage : string, le chemin du fichier image à ouvrir.
// Retourne :
// - image.Image : l'image décodée en cas de succès.
// - error : une erreur en cas d'échec.
func OuvrirImage(cheminImage string) (image.Image, error) {
	file, err := os.Open(cheminImage)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de l'image : %v", err)
	}

	return img, nil
}

// Noyeau génère un noyau de flou gaussien basé sur un sigma donné.
func Noyeau(size int, sigma float64) [][]float64 {
	kernel := make([][]float64, size)
	sum := 0.0
	mid := size / 2

	for i := 0; i < size; i++ {
		kernel[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			x, y := float64(i-mid), float64(j-mid)
			value := (1 / (2 * math.Pi * sigma * sigma)) * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
			kernel[i][j] = value
			sum += value
		}
	}

	// Normalisation
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			kernel[i][j] /= sum
		}
	}

	return kernel
}

// FlouGaussienUneLigne applique un flou gaussien à une seule ligne de l'image.
func FlouGaussienUneLigne(fImage image.Image, kernel [][]float64, y int, blurred *image.RGBA) {
	bounds := fImage.Bounds()
	kSize := len(kernel)
	kMid := kSize / 2

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		var rSum, gSum, bSum float64

		for ky := 0; ky < kSize; ky++ {
			for kx := 0; kx < kSize; kx++ {
				offsetX := x + kx - kMid
				offsetY := y + ky - kMid

				if offsetX < bounds.Min.X || offsetX >= bounds.Max.X || offsetY < bounds.Min.Y || offsetY >= bounds.Max.Y {
					continue
				}

				r, g, b, _ := fImage.At(offsetX, offsetY).RGBA()
				weight := kernel[ky][kx]

				rSum += weight * float64(r>>8)
				gSum += weight * float64(g>>8)
				bSum += weight * float64(b>>8)
			}
		}

		blurred.Set(x, y, color.RGBA{
			R: uint8(rSum),
			G: uint8(gSum),
			B: uint8(bSum),
			A: 255,
		})
	}
}

// FlouGaussienUneColonne applique un flou gaussien à une seule colonne de l'image.
func FlouGaussienUneColonne(fImage image.Image, kernel [][]float64, x int, blurred *image.RGBA) {
	bounds := fImage.Bounds()
	kSize := len(kernel)
	kMid := kSize / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var rSum, gSum, bSum float64

		for ky := 0; ky < kSize; ky++ {
			for kx := 0; kx < kSize; kx++ {
				offsetX := x + kx - kMid
				offsetY := y + ky - kMid

				if offsetX < bounds.Min.X || offsetX >= bounds.Max.X || offsetY < bounds.Min.Y || offsetY >= bounds.Max.Y {
					continue
				}

				r, g, b, _ := fImage.At(offsetX, offsetY).RGBA()
				weight := kernel[ky][kx]

				rSum += weight * float64(r>>8)
				gSum += weight * float64(g>>8)
				bSum += weight * float64(b>>8)
			}
		}

		blurred.Set(x, y, color.RGBA{
			R: uint8(rSum),
			G: uint8(gSum),
			B: uint8(bSum),
			A: 255,
		})
	}
}
