
## Les bases de donnes qui m'interessent
Anything sql related 
elastic search if I want to do something with research (that works for logs, but for that is it a good use case)

## Comment je vais implementer tout cela.
Pour le moment et histoire de bien poser les abstractions necessaires je vais travailler avec une base SQL et je vais
stocker les images sur un S3 storage. 
Pour la solution finale, je compte stocker stocker les donnees de facon locale dans l'application (avec de quoi faire 
un replica au cas ou si je perds tout cela). Je compte utiliser sqlite parce que de ce que je comprends je peux avec
embed la database dans mon application contrairement a un postgres qui tourne sur un serveur. Donc la database est 
quelque part hosted sur le cloud et je la replique localement. En ce qui concerne la BDD pour l'application mobile, 
je ne vais garder en memoire que les fichiers non closes.
Pour la base de donnee utilisateur, je sais pas si je peux stocker cela localement.
Peut etre la crypter, de telle sorte qu'une requete d'authentification, connecte au serveur pour chercher une cle de 
hachage stocke quelque part poour comparer ce qu'il y a dans le form avec ce que j'ai stocke localement.

## Je veux organiser cela pour aussi faciliter la realisation de la comptabilite


## Liste des tables que je veux faire
1. users (pour l'authentication)
2. person of interest (personne investiguee) .Pour fetch les POI d'un dossier il me faudra un join
3. case (la table associe a un dossier particulier) 
4. enquete_type (incendie, adultere etc...)
5. photo_collection (on repertories toutes nos photos)
6. client (Les clients avec qui papa bosse : avocats, assurance) -> on met le contact qu'il a avec ce client 
7. interlocuteur (la personne a qui on s'adresse chez le client)
8. witness (une personne qui aide dans l'enquete style voisin, famille etc...)
9. helper (la personne qui aide sur l'enquete)
10. helper_photos (les photos du helper dans le dossier)
11. factures
12. mandat
13. devis

<!-- NOTE: le schema de chaque table -->
1. users : (comment on gere le cote chiffrement)
    (
        id UUID PRIMARY KEY,
        username TEXT UNIQUE, 
        - [ ] password VARCHAR(16), // le mot de passe fait 16  char apres hashage (juste une idee)
        is_admin BOOLEAN (or int if they do not have that type)
    )

2. persons of interest : 
    (
        id UUID PRIMARY KEY, 
        lastname TEXT,
        firstname TEXT, 
        case_id UUID 
        city TXT
        FOREIGN KEY(case_id)
        references case(id)
    )

3. cases : 
    (
        id UUID PRIMARY KEY,
        start_date DATE
        end_date DATE
        enquete_type ENQUETE TYPE
        client_id UUID FOREIGN KEY UNIQUE
        interlocuteur_id UUID FOREIGN KEY UNIQUE
        facture_id TEXT
        mandat_id TEXT
        devis_id TEXT
        pompier_id TEXT
        huissier_id TEXT
        is_close BOOLEAN
    )

cases_interlocuteurs

4. photos:
    (
        id UUID PRIMARY KEY UNIQUE
        case_id UUID 
        todo: add metadata for the photos of locations and stuff like that 
        FOREIGN KEY (case_id) REFERENCES case(id)
    )

<!-- NOTE: Je peux mettre les deux suivants dans une seule table et faire des select avec le distinct value pour chercher  -->
<!-- le client par exemple. Je ne sais pas ce qui est plus efficace -->
5. clients : 
    (
        id INTEGER PRIMARY KEY UNIQUE,
        name TEXT (ASSU 2000 etc)
        type TEXT
    )

client_type :
    (
        id INTEGER PRIMARY KEY UNIQUE,
        type TEXT (Assurance, Avocat)
    )

6. interlocuteurs : 
    (
        id INTEGER PRIMARY KEY UNIQUE,
        lastname TEXT,
        firstname TEXT,
        phone TEXT,
        mail TEXT,
        client_id INTEGER 
        FOREIGN KEY (client_id) REFERENCES client(id)
    )

<!-- NOTE: A modelling of the many to many relationship of witness and cases -->
100. case_witnesses :
    (
        case_id INTEGER 
        witness_id INTEGER 
    )


7. witnesses : 
    (
        id INTEGER PRIMARY KEY,
        location TEXT
        occupation TEXT
        phone INT
        mail TEXT
        case_id UUID FOREIGN KEY UNIQUE
        help_description TEXT (un petit texte pour se rappeler en quoi la personne a aide)
    )

8. helpers
    (
        id INTERGER PRIMARY KEY     
        firstname TEXT
        lastname TEXT
        company_name TEXT
        workhours
        mission_description TEXT
    )

100. cases_helpers :
    (
        case_id UUID
        helper_id INTEGER
        UNIQUE(case_id, helper_id)
    )

9. helper_photos :
    (
        id UUID PRIMARY KEY // ou photo name
        case_id String 
        helper_id String 
        FOREIGN KEY (interlocuteur_id) REFERENCES interlocuteur(id)
        FOREIGN KEY (case_id) REFERENCES case(id)
    )

11. factures:
    (
        id UUID PRIMARY KEY
        tva FLOAT
        workhours INT
        workhours_paid INT
        solved BOOLEAN
    )
12. mandat:
13. devis:



Relationship (ce qui permet de se faire une idee des keys)
1 helper (intervenant_id) peut intervenir sur plusieurs enquetes (case_id)
1 person_of_interest (person_of_interest_id) peut etre le sujet de plusieurs enquetes (potentiellement)
1 witness (witness_id) can appear in different cases (case_id)
1 case (case_id) peut avoir plusieurs clients (client_id)
1 client (client_id) a plusieurs case (case_id)
1 client can have multiple interlocuteur (secretaire(s), avocat direct etc..)
1 photo belong to one case (1 to 1 relationship)
