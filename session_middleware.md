# Utilisation d'un middleware pour gérer les sessions

## Question

j'utilise un middleware qui met à jour le cookie, mais quand je sort du middleware, le http.Request contient encore l'ancien cookie, je ne sais pas comment le mettre à jour avec le nouveau cookie pour pouvoir accéder aux infos situés dans la map
Afficher les suggestions

## Réponse

Le middleware ne met pas à jour le http.Request

Lorsque vous utilisez un middleware pour mettre à jour un cookie, il est important de comprendre que le middleware ne met pas à jour le http.Request lui-même. Le middleware modifie plutôt le contexte de la requête, qui est un objet distinct du http.Request.

### Conséquences

    Si vous essayez d'accéder à la valeur du cookie après le middleware, vous obtiendrez l'ancienne valeur, car le http.Request n'a pas été modifié.
    Vous ne pourrez pas accéder aux informations situées dans la map du nouveau cookie.

### Solutions

1. Accéder au contexte de la requête:

    Vous pouvez utiliser le package context pour accéder au contexte de la requête et obtenir la valeur du cookie mis à jour.
    Le package context fournit des fonctions pour obtenir et définir des valeurs dans le contexte.

### Exemple

    ```Go

package main

import (
    "context"
    "fmt"
    "net/http"
)

func main() {
    // Créer un middleware
    middleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Mettre à jour le cookie
            cookie := &http.Cookie{
                Name:  "nom_cookie",
                Value: "nouvelle_valeur",
            }
            http.SetCookie(w, cookie)

            // Accéder au contexte de la requête
            ctx := r.Context()
            updatedCookie, ok := ctx.Value("updatedCookie").(*http.Cookie)
            if !ok {
                fmt.Println("Le cookie mis à jour n'est pas disponible dans le contexte")
                return
            }

            // Accéder aux informations du cookie
            fmt.Println("Valeur du cookie mis à jour:", updatedCookie.Value)

            // Appeler le prochain handler
            next.ServeHTTP(w, r)
        })
    }

    // Créer un handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Handler")
    })

    // Enregistrer le middleware
    http.HandleFunc("/", middleware(handler))

    // Démarrer le serveur
    http.ListenAndServe(":8080", nil)
}

```

### Stocker le cookie mis à jour dans une variable

    Vous pouvez stocker le cookie mis à jour dans une variable et l'utiliser ultérieurement dans votre code.
    Cette solution est plus simple que d'accéder au contexte de la requête, mais elle est moins flexible.

### Exemple 1

    ```Go

package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Créer un middleware
    middleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Mettre à jour le cookie
            cookie := &http.Cookie{
                Name:  "nom_cookie",
                Value: "nouvelle_valeur",
            }
            http.SetCookie(w, cookie)

            // Stocker le cookie mis à jour dans une variable
            updatedCookie := cookie

            // Appeler le prochain handler
            next.ServeHTTP(w, r)
        })
    }

    // Créer un handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Accéder aux informations du cookie
        fmt.Println("Valeur du cookie mis à jour:", updatedCookie.Value)
    })

    // Enregistrer le middleware
    http.HandleFunc("/", middleware(handler))

    // Démarrer le serveur
    http.ListenAndServe(":8080", nil)
}

```

Utilisez ce code avec précaution.
Le middleware ne met pas à jour le http.Request

Lorsque vous utilisez un middleware pour mettre à jour un cookie, il est important de comprendre que le middleware ne met pas à jour le http.Request lui-même. Le middleware modifie plutôt le contexte de la requête, qui est un objet distinct du http.Request.

Conséquences:

    Si vous essayez d'accéder à la valeur du cookie après le middleware, vous obtiendrez l'ancienne valeur, car le http.Request n'a pas été modifié.
    Vous ne pourrez pas accéder aux informations situées dans la map du nouveau cookie.

Solutions:

1. Accéder au contexte de la requête:

    Vous pouvez utiliser le package context pour accéder au contexte de la requête et obtenir la valeur du cookie mis à jour.
    Le package context fournit des fonctions pour obtenir et définir des valeurs dans le contexte.

### Exemple 2

    ```Go

package main

import (
    "context"
    "fmt"
    "net/http"
)

func main() {
    // Créer un middleware
    middleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Mettre à jour le cookie
            cookie := &http.Cookie{
                Name:  "nom_cookie",
                Value: "nouvelle_valeur",
            }
            http.SetCookie(w, cookie)

            // Accéder au contexte de la requête
            ctx := r.Context()
            updatedCookie, ok := ctx.Value("updatedCookie").(*http.Cookie)
            if !ok {
                fmt.Println("Le cookie mis à jour n'est pas disponible dans le contexte")
                return
            }

            // Accéder aux informations du cookie
            fmt.Println("Valeur du cookie mis à jour:", updatedCookie.Value)

            // Appeler le prochain handler
            next.ServeHTTP(w, r)
        })
    }

    // Créer un handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Handler")
    })

    // Enregistrer le middleware
    http.HandleFunc("/", middleware(handler))

    // Démarrer le serveur
    http.ListenAndServe(":8080", nil)
}

    ```

Utilisez ce code avec précaution.

### Stocker le cookie mis à jour dans une variable 2

    Vous pouvez stocker le cookie mis à jour dans une variable et l'utiliser ultérieurement dans votre code.
    Cette solution est plus simple que d'accéder au contexte de la requête, mais elle est moins flexible.

###  Exemple 3

    ```Go
package main
import (
    "fmt"
    "net/http"
)

func main() {
    // Créer un middleware
    middleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Mettre à jour le cookie
            cookie := &http.Cookie{
                Name:  "nom_cookie",
                Value: "nouvelle_valeur",
            }
            http.SetCookie(w, cookie)

            // Stocker le cookie mis à jour dans une variable
            updatedCookie := cookie

            // Appeler le prochain handler
            next.ServeHTTP(w, r)
        })
    }

    // Créer un handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Accéder aux informations du cookie
        fmt.Println("Valeur du cookie mis à jour:", updatedCookie.Value)
    })

    // Enregistrer le middleware
    http.HandleFunc("/", middleware(handler))

    // Démarrer le serveur
    http.ListenAndServe(":8080", nil)
}

    ```

Utilisez ce code avec précaution.

### Conclusion

Lorsque vous utilisez un middleware pour mettre à jour un cookie, il est important de comprendre que le middleware ne met pas à jour le http.Request lui-même. Vous devez utiliser le package context ou stocker le cookie mis à jour dans une variable pour accéder à sa valeur ultérieurement dans votre code.
