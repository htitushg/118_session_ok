//Fonctions
function demarrer(){
    document.getElementById("bouton").addEventListener("click", affichemdp)
}
function affichemdp(){
    if (document.getElementById("pwd1").type=="text")
    {
        document.getElementById("pwd1").type="password"
        document.getElementById("pwd2").type="password"
    }
    else
    {
        document.getElementById("pwd1").type="text"
        document.getElementById("pwd2").type="text"
    }
    //document.getElementById("bouton").removeEventListener("click", exploser)
}
//Corps
//Attends le chargement de la page pour exécuter la fonction demarrer
window.addEventListener("load", demarrer)
//____________________________________________________________



//____________________________________________________________
/*
Saisissez le mot de passe [7 à 15 caractères contenant uniquement des caractères, des chiffres numériques, un trait de soulignement et le premier caractère doit être une lettre]
*/
function CheckPassword1(inputtxt) 
{ 
var passw=  /^[A-Za-z]\w{7,15}$/;
if(inputtxt.value.match(passw)) 
{ 
alert('Correct, try another...')
return true;
}
else
{ 
alert('Wrong...!')
return false;
}
}

/*
Saisissez votre mot de passe [6 à 20 caractères contenant au moins un chiffre numérique, une lettre majuscule et une lettre minuscule]
*/
function CheckPassword2(inputtxt, inputtxt2) 
{ 
var passw = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,20}$/;
if(inputtxt.value.match(passw)) 
{ 
//alert('Correct, try another...')
if(inputtxt.value == inputtxt2.value)
    alert('Correct, les 2 entrées sont bien identiques')
else{
    alert('Incorrect, les 2 entrées doivent être identiques') 
}
return true;
}
else
{ 
alert('Wrong...Vous devez entrer au moins une majuscule, une minuscule et un chiffre!')
return false;
}
}
/*
Saisissez votre mot de passe [7 à 15 caractères contenant au moins un chiffre numérique, une lettre et un caractère spécial]
*/
function CheckPassword3(inputtxt) 
{ 
var paswd=  /^(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{7,15}$/;
if(inputtxt.value.match(paswd)) 
{ 
alert('Correct, try another...')
return true;
}
else
{ 
alert('Wrong...!')
return false;
}
} 
/*
Saisissez le mot de passe [8 à 20 caractères contenant au moins une lettre minuscule, une lettre majuscule, un chiffre numérique et un caractère spécial]
*/
function CheckPassword4(inputtxt, inputtxt2) 
{ 
var decimal=  /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z0-9])(?!.*\s).{8,20}$/;
if(inputtxt.value.match(decimal)) 
{ 
//alert('Correct, try another...')
if(inputtxt.value == inputtxt2.value)
    alert('Correct, les 2 entrées sont bien conformes et identiques')
else{
    alert('Incorrect, les 2 entrées doivent être conformes et identiques') 
}
return true;
}
else
{ 
alert('8 à 20 caractères, au moins 1 chiffre, 1 majuscule, 1 minuscule, 1 caractère spécial !')
return false;
}
} 
//______________________________________________________________________________________________

function ValidateEmail(uemail)
{
var mailformat = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
if(uemail.value.match(mailformat))
{
return true;
}
else
{
alert("l'adresse courriel est invalide !");
uemail.focus();
return false;
}
}
//______________________________________________________________________________________________
