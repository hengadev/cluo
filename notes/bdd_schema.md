## Liste des tables que je veux faire
- users (for authentication)
- person_of_interests (personne investiguee)
- departments (guadeloupe, martinique etc...)
- cases (la table associe a un dossier particulier) 
- case_types (incendie, adultere etc...)
- photos (on repertorie toutes nos photos)
- clients (GROUPAMA Guadeloupe, COPPET)
- client_types (avocats, assurance)
- interlocuteurs (la personne a qui on s'adresse chez le client)
- witnesses (une personne qui aide dans l'enquete style voisin, famille etc...)
- helper (un autre detective sur l'enquete)
- factures
- mandat
- devis
- department
- city (since wee work in Guadeloupe and Martinique, add table with list of these things)

- users : (comment on gere le cote chiffrement)
    (
        id UUID PRIMARY KEY,
        username TEXT UNIQUE, 
        password TEXT UNIQUE
        is_admin BOOLEAN (or int if they do not have that type)
    )

- persons of interest : 
    (
        id UUID PRIMARY KEY, 
        lastname TEXT,
        firstname TEXT, 
        dpt TXT
        city TXT
        case_id UUID 
        FOREIGN KEY(case_id)
        references case(id)
    )

- cases : 
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

- cases_interlocuteurs :
    (
        case_id UUID PRIMARY KEY,
        interlocuteur_id UUID PRIMARY KEY,
    )

- photos:
    (
        id UUID PRIMARY KEY UNIQUE
        case_id UUID 
        user_id TEXT (who took the photo)
        todo: add metadata for the photos of locations and stuff like that 
        todo: add fields to help with filigrane if we want to add it.
        FOREIGN KEY (case_id) REFERENCES case(id)
    )

- clients : 
    (
        id INTEGER PRIMARY KEY UNIQUE,
        name TEXT (ASSU 2000 etc)
        client_type TEXT
        FOREIGN KEY (client_type) REFERENCES client_types(type)
    )

- interlocuteurs : 
    (
        id INTEGER PRIMARY KEY UNIQUE,
        lastname TEXT,
        firstname TEXT,
        phone TEXT,
        mail TEXT,
        client_id INTEGER 
        FOREIGN KEY (client_id) REFERENCES client(id)
    )

- case_witnesses :
    (
        case_id INTEGER 
        witness_id INTEGER 
    )

- witnesses : 
    (
        id INTEGER PRIMARY KEY,
        location TEXT
        occupation TEXT
        phone INT
        mail TEXT
        case_id UUID FOREIGN KEY UNIQUE
        help_description TEXT (un petit texte pour se rappeler en quoi la personne a aide)
    )

- helpers
    (
        id INTERGER PRIMARY KEY     
        firstname TEXT
        lastname TEXT
        company_name TEXT
        facture_id
        devis_id
        mandat_id
        mission_description TEXT
    )

- cases_helpers :
    (
        case_id UUID
        helper_id INTEGER
        UNIQUE(case_id, helper_id)
    )

- factures:
    (
        id UUID PRIMARY KEY
        case_id TEXT
        tva FLOAT
        workhours INT
        workhours_paid INT
        solved BOOLEAN
    )
- mandat:
    (
        id UUID PRIMARY KEY
        case_id TEXT
    )
- devis:
    (
        id UUID PRIMARY KEY
        case_id TEXT
    )
