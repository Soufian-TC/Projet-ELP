package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"time"

	"github.com/Soufian-TC/Projet-ELP/GO/fonctions"
)

// Fonction pour traiter l'upload de l'image, appliquer le flou et renvoyer l'image floutée
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée, utilisez POST", http.StatusMethodNotAllowed)
		return
	}

	// Log de la requête reçue
	fmt.Printf("Requête reçue de : %s\n", r.RemoteAddr)

	// Parse la requête multipart pour récupérer les fichiers
	err := r.ParseMultipartForm(10 << 20) // Limite à 10 Mo
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données", http.StatusBadRequest)
		return
	}

	// Récupère le fichier envoyé avec le champ "image"
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ouvrir l'image
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Erreur lors de l'ouverture de l'image", http.StatusInternalServerError)
		return
	}

	// Récupérer la valeur de "intensite"
	intensite := r.FormValue("intensite")
	fmt.Println("Intensité reçue:", intensite)

	intensitec, err := strconv.ParseFloat(intensite, 64)
	if err != nil {
		http.Error(w, "Erreur lors de la conversion de l'intensité en float64", http.StatusBadRequest)
		return
	}

	// Appliquer un flou gaussien sur l'image (choisir la taille du noyau et sigma)
	start := time.Now()
	blurredImg := fonctions.FlouGaussienOptimise(img, intensitec)
	end := time.Now()
	duration := end.Sub(start)

	// Renvoyer l'image floutée directement dans la réponse HTTP
	w.Header().Set("Content-Disposition", "attachment; filename=blurred_image.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, blurredImg, nil) // Utiliser la réponse HTTP comme le fichier de sortie
	fmt.Printf("Durée de l'opération : %.2f secondes\n", duration.Seconds())
	if err != nil {
		http.Error(w, "Erreur lors de l'envoi de l'image floutée", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Configure le handler pour l'endpoint /upload
	http.HandleFunc("/upload", handleUpload)

	// Démarre le serveur sur le port 8080 pour toutes les interfaces réseau
	port := "8080"
	fmt.Printf("Serveur démarré sur http://0.0.0.0:%s\n", port)
	err := http.ListenAndServe("0.0.0.0:"+port, nil) // Écoute sur toutes les interfaces
	if err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur : %v\n", err)
	}
}
