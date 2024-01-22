//Variables Globales

//Abonnements

//Fonctions
function demarrer(){
    document.getElementById("bouton").addEventListener("click", exploser)
}
function exploser(){
    document.getElementById("p1").innerHTML="Bouton nucléaire pressé !"
    document.getElementById("p2").innerHTML="Il ne fallait surtout pas y toucher !"
    let body1=document.getElementById("mainSection")
    body1.style.backgroundImage="url('boum.jpg')"
    new Audio("boumloin.wav").play
    document.getElementById("bouton").removeEventListener("click", exploser)
}
//Corps
//Attends le chargement de la page pour exécuter la fonction demarrer
window.addEventListener("load", demarrer)

