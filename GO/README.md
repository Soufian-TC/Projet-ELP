Ce dossier contient les codes client et serveur pour que le serveur effectue un flou gaussien sur une image envoyé par le client 

Pour tester le code :

Dans un terminal :

cd Projet-ELP/GO

go run serveur.go


Dans un autre :

cd Projet-ELP/GO/client

go run client.go

On peut ensuite renseigner le chemin de l'image à floutter, son nom de sauvegarde et l'intensité du flou a appliquer.

Par exemple :

../images/newYork600-400.jpg

test

10

L'image flouttée est sauvegardée dans le dossier client

