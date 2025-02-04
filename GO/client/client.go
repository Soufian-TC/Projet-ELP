package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Demander à l'utilisateur de fournir le chemin du fichier à télécharger
	fmt.Print("Entrez le chemin du fichier image à envoyer au serveur: ")
	var imagePath string
	fmt.Scanln(&imagePath)

	fmt.Print("Entrez le nom de sauvegarde:")
	var nomImage string
	fmt.Scanln(&nomImage)

	//valeur de sigma a envoyer au serveur
	fmt.Print("Entrer l'intensite du flou:")
	var Intensite string
	fmt.Scanln(&Intensite)

	// Ouvrir l'image locale
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de l'image:", err)
		return
	}
	defer file.Close()

	// Créer un buffer pour écrire les données multipart
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Ajouter le fichier image dans le formulaire multipart
	part, err := writer.CreateFormFile("image", imagePath)
	if err != nil {
		fmt.Println("Erreur lors de la création du formulaire:", err)
		return
	}

	// Copier le contenu du fichier dans le champ du formulaire
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Erreur lors de la copie du fichier:", err)
		return
	}

	// Ajouter le champ pour l'intensité
	err = writer.WriteField("intensite", Intensite)
	if err != nil {
		fmt.Println("Erreur lors de l'ajout du champ intensité:", err)
		return
	}
	// Fermer le writer pour signaler la fin de l'écriture du formulaire
	err = writer.Close()
	if err != nil {
		fmt.Println("Erreur lors de la fermeture du writer:", err)
		return
	}

	// Créer une requête POST pour envoyer l'image au serveur
	url := "http://localhost:8080/upload"
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête:", err)
		return
	}

	// Ajouter l'en-tête Content-Type pour indiquer que c'est une requête multipart
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Créer un client HTTP et envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête:", err)
		return
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("./%s.jpg", nomImage)

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier:", err)
		return
	}
	defer outFile.Close()

	// Copier le contenu de la réponse (image) dans le fichier local
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image:", err)
		return
	}

	fmt.Println("Image téléchargée et enregistrée avec succès")
}
