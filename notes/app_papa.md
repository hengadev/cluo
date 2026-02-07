---
id: app_papa
aliases:
  - TODO
tags: []
---

## Todo: organiser tout cela dans le folder en plusieurs folders :
1. organisation = todo
2. feature to implement
3. Ressources to watch
9. Refacto the current part using typescript before getting into the part to fetch data.
10. Refacto the css so that we can use utilities, css variables and I have everything more clear.
11. Les templates pour les differents elements constituants un dossier (important c'est la majeure partie du bail en fait).
    a. Gerer la pagination des elements du main, ceux avec editeur et ceux sans.
    b. Faire le contenu stylise pour le reste.
12. Gerer la conversion ce qu'il y a dans le texte vers un apercu en pdf. (I think I need pagination for that)
13. Faire le basic backend de l'application en golang et connecte les deux.
14. Faire le machine learning pour l'ajout du watermark.
15. Faire le machine learning pour gerer la refacto du texte.
16. Faire l'application mobile pour la gestion des enregistrements.
17. Les videos de Litanhui sur  svelte, je veux faire le truc le plus idiomatique possible pour reduire la complexite du code a l'ecran. On va aussi utiliser un LLM pour apprendre en simultanne. C'est une refacto en fait.
30. Have some github action ready to server and the client with the necessary code

## En cours 

7. On commence le design de l'API et de la base de donnee afin de faire une quelconque refacto en typescript.
    - design de l'API
    - design de la BDD
8. Comment encrypter le dossier a l'envoi pour que je ne me fasse pas avoir, juste trouver la solution c'est bien.

<!-- FIX: Done -->
1. Les icones de la search bar
2. Le cote actif de la sidebar
3. Le texte editor ie comment le configurer plus les icones dans la main bar (Manque toute la partie avec les icones)
    a. Refaire fonctionner les icones avec le text editeur. 
    b. Ajouter 
    les autres icones.
    c. Faire un styling coherent de la mainbar. 
4. Rendre les elements du searchModal cliquable tout en gardant le onblur.
5. Faire passer le texte editeur dans un slot pour rendre le contenu du main modulable.
- [ ] 6. Fix the different component that got messed up due to me using the stash commits and not handle well the merge
conflicts :
    - The header component. 
    - The sidebar component. 
    - The mainbar component. 
7. Realiser les autres modals dont le drag and drop avec les photos pour les organiser avec l'api svelte.
    a. Il me faut des images pour mock le comportement.
    <!-- NOTE: Pour la realisation de mon drag and drop component je viens de download "svelte-dnd-action" -->
    b. Voir comment je peux gerer tout ce qui est en lien avec enregistrer la postion (des images) parce que je dois 
    envoyer cela au backend en fait. -> Il faut un store et des que la valeur est enregistre comme changeant, on envoie
    au backend !

### Des idees de nom pour l'application.
The name curio already exists so I have to make it more complex
- Curio or change the way you spell it like Kurio or Kuryo.
- [ ] Seek
- Quest
- Lens
- UnveilX
- Inquisio (inquiry in latin)
- Vero (truth in latin)
- Pry
- Lens
- [!] SeekForge (curiosity with seek and forge for discovery and creation)

### Notes sur comment trouver un nom de company (win the name game).
- Pick a category of name (7 categories)
    - epononyme
    - descriptive (the american airlines, the home depot, le repere des ados)
    - acronymic
    - suggestive 
        - real : slack, uber
        - composite: facebook , rayban
        - invented: kleenex, pinterest
    - associative (amazon, redbull, siriusxm)
    - non english (lego, samsung, hulu etc...)
    - abstract (rolex, kodak)

### Mes tests : 
- BLAD : Best logiciels aux detectives
- Un truc avec l'inspecteur Columbo qu'ils appellent Cocol / Kokol / Kkol / KO2l (se prononce kokol).  
- CloakCode (the cloak used by PIs and the code for the software stuff) -> Better logo with one.
- CaseAce (a case you know what that is and ace means that you destroy that shit. Expertise in handling cases)
On peut penser au trench coat de colombo qui est aussi tres proche de l'image des detectives.
google is mispelled, it should be googol

Kloak

### Probleme sur le projet. 
>Parfois j'ai des composants qui se chargent pas et cela me cree des problemes e layout, je suis oblige de sauvegarder le fichier que je pense qui est oublie et ensuite raffraichir la page.

### Pour le backend de l'application mobile, j'ai vu que pas mal de personne parle de backendless et je veux m'y interesser
On le voit dans cette [video](https://www.youtube.com/watch?v=2P0q1EdH_oQ) de t3dotgg et cette [video](https://www.youtube.com/watch?v=hfGtyd5nmwQ&t=340s) de Ben Awad

### Juste pour le plaisir, je peux me faire une library UI a partir des composants de l'application en m'inspirant de  shadcn ui.

### Un projet qui traite de streamer de larges files avec golang 
video: https://www.youtube.com/watch?v=3mO5MUbCzKQ

### Changer le display du component pour deplacer les images.  
>De telle sorte qu'il soit au dessus de l'editeur d'image parce que cela permettrait de profiter du fait que l'ecran
est plus large que long. En fait je pense que je devrais rendre cette partie configurable selon les besoins de l'utilisateur.

### Faire la documentation de l'API. 
video : https://www.youtube.com/watch?v=0CSyIBHQy9g

### Une video qui fait des review de projets, je peux m'en inspirer pour display l'app de papa. 
video : https://www.youtube.com/watch?v=CYdsh1FhpI4

### J'aime bien l'animation des tooltips de la sidebar d'Obsidian.

### Video sur l'architecture et le deploiement 
https://www.youtube.com/watch?v=10UU6umqqv8
https://www.youtube.com/watch?v=y2ICZYOU09Q

### La solution que je peux utiliser pour les mails ou dont je peux m'inspirer. 
video : https://www.youtube.com/watch?v=HyDwVN1AFwY

### Rendre la sidebar responsive. 
>Je veux que quelque soit la taille de l'ecran toutes les icones tiennent sur la page sans scroll.
Une autre solution serait de fusionner certains elements, je veux par exemple que les elements comme facture, devis, mandat etc... soit tous sous le meme dossier. Au click un truc  s'expand pour nous montrer les icones que je veux cacher.

### Pas sur de celui, j'avais une idee de use case mais je l'ai oublie 
>Je peux faire une icone requete/notification de l'application (en gros, si j'infere des choses mais je veux que  l'utilisateur les confirme). Par exemple : 
- Je sais que je dois delete un dossier 3ans apres sa creation. Je prefere notfifier l'utilisateur qu'il a une requete de suppression afin qu'il puisse confirmer la suppression du dossier.
- Imaginons que quand je fais des photos, je puisse un algo de ML pour pouvoir reconnaitre une face. Je peux ensuite faire une requete si papa n'a pas fait le travail d'identifier la personne sur toutes les photos de proposer de le faire moi meme via l'algo et qu'il ne lui reste plus qu'a confirmer.

>Les notificatons auxquelles je pensais
1. https://icon-sets.iconify.design/?query=todo
2. https://icon-sets.iconify.design/?query=notification

### Dans le choix de la base de donnee 
>J'ai une partie recherche a implementer. Je sais que je vais surement tout stocker avec du SQL mais je me demandais si  elasticseachDB ne serait pas une solution viable. Je peux l'implementer et voir si il y a un interet a le garder.
Je veux garder une historique des derniers dossiers ouverts tels que lorsque je focus la barre de recherche et   que je n'ai rien encore taper alors j'ai l'historique des derniers dossiers qui s'affichent, le plus recent en premier. 

### Ajouter un state pour gerer si on peut on cliquer ou non les boutons de la sidebar 
>Tout ceci pour deux cas de figure:
1. Quand un bouton n'est pas accessible parce que je n'ai pas implemente la feature ?
2. Quand je n'ai pas de dossiers actifs ouverts (aucun bouton ne devrait etre cliquable).

### How to acces the file system to access data ? 
- https://tauri.app/v1/api/js/fs/

### Je veux faire une presentation d'un workflow avec un especte de truc schema dans le style de prezi. 

### Faire un component popup pour les messages de confirmation. 
Si on va sur cette page et que l'on appuie sur f comme un bourrin on a un exemple de ce que je veux :
https://www.lemondeinformatique.fr/actualites/lire-10-extensions-chatgpt-dans-chrome-a-essayer-91848.html

### How to make the documentaion 
Use mdbook to make this thing so that we get a and extended documentation

### Sur la [video](https://www.youtube.com/watch?v=bpFZL4blkaA) backend du projet de Benjamin Code 
>Une mission ne peut bien se derouler que si elle a ete bien defini en amont par son cahier des chawrges qui contient :
- Objectif
- Fonctionnalites requises (uttilise par le developpeur backend)
- Contraintes techniques comme le langage de programmation ou le framework utilise.
- Type de base de donne utilise 
- Gestion de donne,
- Securite 
- Maniere dont on s'authentifie a l'application
- Veut on des tests ? A quoi ils doivent ressembler ?
- Deadline
- Preciser ce a quoi doivent ressembler les livrables finaux : la premiere version sera une sorte de microsoft word  sans les IA. Je veux qu'ils puissent utiliser word et avoir la gestion des images parfaites.

A quoi ressemble la description des fonctionnalites: On peut le retrouver [ici](https://youtu.be/bpFZL4blkaA?t=309)

### Developper un outil pour optimiser mon backend  
En gros je veux analyser les actions qu'il fait sur le front ainsi que les requetes realises par consequent pour 
optimiser le design front et back pour avoir des reponses plus rapides.
En gros je veux de l'analyse de l'acitivite, il faut penser Graphane et Prometheus

### Tauri content for the backend of my app 
https://www.youtube.com/watch?v=y5zl7cK8Ls4
https://www.youtube.com/watch?v=RunHr-uhUjA


### Une [video](https://www.youtube.com/watch?v=v-9AZKp-Ljo) sur de l'infra. 

### Penser optimisation au niveau du dev avec cette [video](https://www.youtube.com/watch?v=xbylNxx_hi8)
Voir ce que mon framework met en place de base.
frontend:
- format image
- minification (reduire le plus possible le fichier, renommer variables, regrouper les scripts.)
- ce que svelte et tauri envoie pour ne garder que l'essentiel.
- lazyloading.
- mise en cache.
- cdn ?
backend: 
- optimiser les requetes de BDD.
- mise en cache 
- compression 
- cloud pour optimiser les requetes ?
- regrouper vos infrastructures.
Je peux en faire un article

<!-- NOTE: Les couleurs que je veux utiliser que je tiens du site de login de leonardo.ai.  -->
main background : #171717;
button: #27272d;
button hover : #3d3c3d;
placeholder : #11111;
text placeholer : rgb(189, 189, 189);
Toutes les autres coleurs de texte sont blanc.

Certains boutons comme les signs in etc utilisent un gradient comme celui que j'ai deja utilise.
On peut modifier les differentes degres etc... En fait c'est une couleur accent. 
background-image: 
    linear-gradient(
        45deg, #12d2e9, #c471ed, #f64f58
);
Ce truc permet aussi d'identifier dans quel onglet sur le site on est. Les icones sont blanches par default et ensuite
devienne comme cela ie en gradient si on est sur le dit lien.

Les anciennes couleurs que j'avais: 
#393939 | #525252 | 


<!-- NOTE: basic system design -->
je veux pouvoir upload des photos et tous les fichiers textes.
fetch des photos et tous les fichiers tests.
chercher les dossiers en utilisant differents criteres comme le nom, la date, la ville
proposer des suggestions a partir des audios et a partir des entrees de texte realise par l'utilisateur.
envoyer les dossiers finis par mail.
mettre un filigranne sur les photos.
encoder tous les fichiers texte
creer des utilisateurs avec des droits d'acces differents notamment sur les photos.
faire des recherches sur les reseaux sociaux de quelqu'un

<!-- TODO: Regrouper les icones factures, devis, mandat etc.. -->
Etant que ce sont des templates dont je n'ai besoin que de certaines informations, je peux les regroupere car je ne
vais pas y acceder regulierement, je n'ai besoin que de certaines informations qui son renseignees dana l'espace
dedie. En gros ces infos viennent de query de la BDD.


<!-- NOTE: Toute sorte de lien qui ont a voir  avec les richs text editor. -->
https://www.youtube.com/watch?v=EEF2DlOUkag
https://www.youtube.com/watch?v=tfuU0Ra5yHE
https://softwareengineering.stackexchange.com/questions/187229/text-editor-document-model
https://prosemirror.net/
https://www.youtube.com/watch?v=hJrnIrsZcEs&t=91s
https://www.youtube.com/watch?v=EwoS0dIx_OI
https://www.youtube.com/watch?v=bBCVI2e18dE

<!-- NOTE: Pour me lancer dans le machine learning parce que cela fait longtemps. -->
https://www.youtube.com/watch?v=tHL5STNJKag

<!-- NOTE: How to do an auth system, I want to get inspired by Oauth2. -->
https://www.youtube.com/watch?v=uj_4vxm9u90

<!-- TODO: Creer un bouton contact ou tu rentres toutes les infos que papa a rencontre dans son affaire. -->
Il faut que ce qui est dans contact puisse etre cherchable a partir de la barre de recherche.
En gros j'ai un contact avec un gars qui travaille a la meteo, faire en sorte que je peux chercher les contacts par
leur metier par exemple. "Ce gars la etait macon mais je ne me rappelle pas son nom."
code postal, les contacts

<!-- NOTE: Pour la partie ML avec les editeurs de texte j'ai tout interet a regarder ce qui se fait au niveau des librairies -->
<!-- ou framework du style react/svelte/vue.  -->
C'est pour voir comment tout cela est organise parce que dans l'idee je veux taper
du texte, que ce texte soit process puis qu' il y ait un processus qui update l'UI en se basant sur ce que le modele 
de machine learning renvoie en fait. Mais j'avais une video youtube qui traitait exacement de cela en fait ie comment
ajouter un LLM a un editeur de texte.


<!-- NOTE: Nom de l'applcation  -->
Un jeu de mot avec Sherlock Holmes et Serge par exemple SergeLock Holmes Assistant.

<!-- NOTE: Un editeur de texte qui peut etre interessant, ne serait ce que de voir le code -->
https://www.onlyoffice.com/fr/developer-edition.aspx?utm_source=google&utm_medium=cpc&utm_campaign=pmax_dev_fr&utm_content=aid_&utm_term=&gclid=Cj0KCQjwsp6pBhCfARIsAD3GZuZOHk0xoRQo4Ef538Q4g2Z-iwLPyCKF2-QND2_hLlnZ-UWCkfhEYZAaAkDQEALw_wcB

<!-- NOTE: Pour gerer la fermeture des modales je veux une croix comme ce qui permet de fermer l'outil de gestion des -->
<!-- playlists sous les videos sur Youtube. -->

<!-- NOTE: Pour la partie template ce ne sont que des page en html que je vais display dans la partie main avec des input -->
<!-- dont je vais recuperer la valeur pour ensuite en deduire les donnes a entree tout simplement. -->


<!-- NOTE: Les boutons qu'il me faut pour mon editeur de texte sur la base de ce que j'ai vu dans les differents dossiers. -->
1. bold
2. italic
3. underline
4. De quoi  gerer les sizes (Combien de size ?)
5. Des emojis un peu comme le tick ou encore le circle devant les ol.
6. Y a pas mal d'article et tout 
7. Color : blue, 
8. Un footer avec le logo et tout le bordel.

<!-- TODO: Faire des recherches sur comment rendre les dossiers et rendre cela un peu propre comme ce que l'on pourrait -->
<!-- trouver dans des cabinets d'avocat. -->


NOTE: State
1. folder open : bool
2. if folder open (y a une dependance dedans dans ce que je peux maintenant afficher) :
mandat
rapport
photo
etc...


<!-- TODO: Rust est tres bon pour tout ce qui est manipulation de liste. -->
Je peux tres bien utiiliser cela pour recoder l'algo me permettant de faire la recherche flou de document.
D'ailleurs je peux tres bien m'inspirer de ce qu'il y a sur telescope qui utilise grep en vrai.
Je vois un probleme dans tout cela dans la mesure ou les algos doivent surement dans leur scoring prendre en 
compte un clavier qwerty alors que papa travaille en azerty.

<!-- NOTE: Pour les animations  -->
Y en a de natives avec svelte et si je veux plus complique je peux utiliser gsap ou encore motion one.
Je ne sais plus d'ou je tiens ces informtions.

<!-- NOTE: When to use stores -->
We use stores for global states or info required by mutiple unrelated components sush as the logged in user or theme. 

<!-- NOTE: La taille d'une feuille de papier est donnee sur cette video -->
backgroound-color: du main (ce qui n'est pas le papier) : #f3f3f3;
width: 8.5in
min-height: 11in
<!-- FIX: J'ai prefere donne un aspect ratio de telle sorte que je puisse modifier la width sans que j'ai a calculer pour  -->
<!-- modifier la height derriere. -->
padding: 1in
box shadow: 0 0 5px 0 rgba(0 ,0 ,0, 0.5)
background-color: white;

<!-- TODO: Probleme sur le menu image ? -->
Quand je previsualise l'application sur brave j'ai ce message d'erreur des que j'essaie de scroll le menu d'image :
Uncaught TypeError: Cannot read properties of undefined (reading 'scrollTop') at HTMLDivElement.scrollMenu

<!-- NOTE: Pour les legendes des images -->
Je peux utiliser un label pour les mettre et pouvoir les changer en double click, le t rick c'est d'avoir un
input que l'on display en double click et on cache le precedent label.  Le truc est explique ici :
https://www.youtube.com/watch?v=c0fzQRDyPXw&list=PLA9WiRZ-IS_xz1M4H0fjZAii4vzRffTno&index=2

<!-- NOTE: Quel modele pour mon LLM ? -->
Ici il parle de Mistral qui serait mieux que Llamav2. Il parle aussi de commentin staller ce genre de truc.
https://www.youtube.com/watch?v=tK1Pivdcl3U

<!-- NOTE: I need to add typesafety to my frontend. -->
First choice typescript but can we use something else like JsDoc or even vanillajavscript if I want to.
Ici il donne des pistes a oberserver pour ma refacto : https://www.youtube.com/watch?v=OQIsQDFtEnI 


NOTE: Je viens de voir comment on set les variables paths, voila ce que l'on ajoute.
1. Dans le fichier vite.config on ajoute dans le defineConfig un champ : 
resolve : {
    alias: {
        $root : path.resolve('./src'),
    },
},
Avant cela on prendra soin au debut du fichier de rajouter le import path from "path";

2. Dans le fichier tsconfig on ajoute dans le co mpilerOptions :
"paths": {
    "$root": ["./src/*"]
}
On fera attention a bien respecter tout ce qui est , et tout pour que tout passe bien.

<!-- NOTE: Je  rajoute un autofocus sur le input de la barre de recherche ? -->
Typescript me lance entendre et c'est tout a fait logique que cela n'est pas bon pour l'accessbilite mais osef
dans le cas de l'application.

<!-- NOTE: Quelle est la difference semantique entre un input checkbox et un button ? -->

<!-- NOTE: creer un fichier type dans lequel on met toutes les interfaces dont on a bbesoin et ne pas bloater le component. -->
On cree le folder types et dedans on peut chercher le fichier todo.ts qui continet l'interface que l'on exporte 
export interface ITodo {
    id: string,
    text: string,
    completed: string,
}
C'est une facon un peu plus interessante de gerer tout cela, et surtout cela permet de bien gerer tout ce qui est en 
lien avec les types de facon propre et cela evite le bordel qu'il y a dans mon projet.

<!-- NOTE: Parfois ils utilisent cette notation que je ne comprends pas. -->
let todos: ITodo[] = JSON.parse(localStorage.getItem('todos')) ?? []
Je parle du "??", je vois ce a quoi cela peut servir mais je comprends pas le double ?
-> C'est le coalescing operator, en gros c'est comme un ternary operator, on evalue ce qu'il y avant les ?? et si cette valeur est
differente de null alors on la garde sinon on prend ce qu'il y a apres le ??
Je pense que je peux utiliser a chaque fois que je vais fetch des elements avec des querySelector ou par leur id.

<!-- NOTE: Je lis un peu plus sur comment gerer les datas. -->
Je viens de cet article : https://coderpad.io/blog/development/a-guide-to-svelte-stores/
On aborde le context api qui est une facon de gerer les donnees mais on presenteles incovenients. Il n'y a pas de
builtin fonction pour gerer les changements et il est diffile de gerer le state de gerer le state au fur et a mesure
que l'app scale (du au fait de gerer des keys tout cela).
Pour ces raison on peut penser a utiliser des stores.

<!-- TODO: Understand the different way to understand exporting. -->

<!-- NOTE: Ici je veux faire de la pagination donc je veux pouvoir utilise des projets deja existants. -->
<!-- TODO: Si je fais de la pagination, pour la partie rapport, me faut il un editor par page ? -->
On me parle de pagedjs, le fait est que je n'arrive pas a le faire fonctionner avec svelte, a moi d'insister peut etre.
3 modules : 
1. Le chunker : fragmente le contenu en des pages discretes. 
2. the polisher: transform the css into one that the browser can understand.
3. the previewer : create the content that you see on the browser. 


Je vais le faire en rust :
Un faire une application bidon en rust, la compiler en wasm et voir si je peux la run dans svelte.
ensuite je passe le html au calme. et le css eventuellement.

Des sites qui permet de voir comment coder : 
https://stackoverflow.com/questions/74070328/paging-in-react
https://stackoverflow.com/questions/50965683/how-would-you-implement-a-rich-text-editor-with-pagination
https://pagedjs.org/documentation/4-how-paged.js-works/#the-chunker%3A-fragmenting-the-content


<!-- TODO: Faire une IA pour savoir dans quel cadre juridique intervient la mission. -->
Affaire : suspicion abus de confiance
Cadre juridique 
article 314-1 code penal etc... -> Trouver avec l'IA.

NOTE: En ce qui concerne le champ avec un label et que l'on peut edit, la structure est la suivante :
<div class="view">
    <input type="checkbox" />
    <label>Le texte par default</label>
</div>
Apres un double click on insere juste apres la div un input avec la class "edit".

### Add some path components to make things easy. 
Watch this video if you need help : https://www.youtube.com/watch?v=2CQm1rC7IDk&list=PLA9WiRZ-IS_xz1M4H0fjZAii4vzRffTno&index=4

### Comment gerer le fait que je dois gerer l'IA qui prend en compte les modifs dans le texte et le programme   qui gere la pagination comment integrer tout cela ? 
>Je ne sais pas, ne serait ce pas un truc lie au hot module reload ?

### Petite modification a faire sur la barre de recherche. 
>Il faut que je mette tout cela dans un form comme cela je pourrais submit directement en fonction de ce qu'il y a dans l'input, cela sera plus simple pour gerer le backend.
On peut voir une illustration de cela dans ce projet avec la barre dans laquelle on insere le contenu de l todo que
l'on veut entrer : https://joyofcode.xyz/svelte-todo-app
On peut aussi penser aussi penser a un select : 
<select>
    <option>a</option>
    <option>b</option>
    <option>c</option>
    <option>d</option>
</select>
On en parle ici : https://svelte.dev/docs/element-directives

### Commment gerer des resultats de recherche rapide dans la search bar
>Un gars dans un [tweet](https://twitter.com/thdxr/status/1736835759355634065) a montrer la difference entre la recherche dans notion  et ce que lui a pu [implementer](https://twitter.com/thdxr/status/1736780708759175487)
Dans le second tweet j'y ai vu les notions suivantes :
- SST : serverless stack 
- replicache  
- local-first software, video a chercher sur youtube.
Je veux les chercher ou tout du moins trovuer quelque chose qui me permette de faire rapidement des recherches rapides

<!-- NOTE: Dans le store que je veux creer, il me faut un store pour gerer l'input de la brre de recherche. -->
De plus si je mets cela dans un store et que je veux le submit je peux utiliser (comme je suis dans svelte), la
fonction suivante: 
function handleSubmit() {
    some_function_to_fetch_the_database(searchValue);
    // Cela on le ferait avec un reset de la todo si je mettais tout cela dans un store en fait.
    searchValue = "";
}
J'ai deja bind la value a l'input donc c'est ok de ce cote.

### Je peux changer la semantique de ce qui se apsse dans le search bar modal 
En effet, on peut voir cela comme un form avec chacun des choix comme un input de type checkbox, le form n'ayant qu'une seule valeur possible / realisable.
 ### Une video qui parle de comment gerer les changements de schema dans une BDD. 
Lien : https://www.youtube.com/watch?v=y2J00F19OaE
Il faudra penser a faire des mirgations pour ma base de donnee
### Une video interessante pour tout ce qui a trait aux notificatins.  
>En gros des que certaines actions sont  effectues on envoie un message a une queue qui se charge de lancer le 
service, ici les notifictions. En plus il dit que c'est bon pour le machine lear ning par exemple pour entrainer un
modele de machine learning en se basant  sur les donnes venant de l'application. (Je ne sais pas si entrainer est 
le bon mot). #todo: Ou est la video ?

### Je peux chanegr le cote  actif de la sidebar avec quelque chose de plus idiomatique
>Je peux utiliser une class directive:
class:active={sidebarState === "rapport"}
ici on va juste ajouter la class active a l'element que je constate.
Dans ma tete je crois que cela ne va pas marcher alors que si vu que la classe depend d'un state qui est  reactif
donc a chque fois que le state change tous les componenents vont etre re-render.
En fait j'aurais une liste d'objet de ce type 
- [ ] let button = {
    id: 1,
    name: "information",
    icon: Info,
    canBeHighlight: false,
    modal: InfoModal,
};

Du coup je vais faire une for loop sur les buttons comme cela :
{#each buttons as button (button.id)}
    {#if button.canBeHighlight}
        <Modal>
            {button.Modal}
        </Modal>
        <button
            on:click={(e) => {
                openModal(e);
            }}
            class="sidebar__button highlight"
        >
            <img class="sidebar__img" src={button.icon} alt="" />
            <div class="tooltip">{button.name}</div>
        </button>
    {/if}
    {#if !button.canBeHighlight}

        <button
            class:actif={sidebarState === button.name}
            <!-- id="facture" -->
            class="sidebar__button highlight"
            on:click={setActif}
        >
            <img class="sidebar__img" src={button.icon} alt="" />
            <div class="tooltip">{button.name}</div>
        </button>
    {/if}
{/each}
La fonction setActif mais en gros elle prend juste l'element en question et passe la valeur de cet element au state
defini dans le sidebarState

### Utiliser le nullish coalescing operator 
let foo = null ?? "right-hand";
https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Nullish_coalescing

### Je dois faire des tests. 
> Je cherche un test coverage de plus de 80%, car je suis seul sur le projet et il me faut de la pratique.
Maintenant, qu'est ce que je vais tester.
1. Le systeme de notification car beaucuop de services.
En vrai je vais tester le front apres en apprenant de cette serie : https://www.youtube.com/watch?v=dO6Xk30jyoc&list=PLA9WiRZ-IS_z7KpqhPELfEMbhAGRwZrzn&index=5

### Refaire les composants graphiques pour avoir une certaine coherence. 
>En gros, je veux utiliser shadcn pour refaire tout cela et avoir un truc coherent.

### Je peux rajouter une partie pour gerer les elements de sous traitance. 
>En gros, il veut travailler avec des gens, ces gens vont lui donner un rapport ainsi que des photos ainsi que des
factures etc... donc il me faut un espace pour pouvoir gerer le/les sous traitants.
### Pour les notifications via mon backend il me faut exactement cela. 
https://svelte.dev/repl/2254c3b9b9ba4eeda05d81d2816f6276?version=3.32.2

### Pour la gestion des input 
Je prefere on focus changer la couleur du border et mettre le outline a none.

### Mettre plus de coherence dans mes templates. 
>Je me suis fait plein de classes utilitaires mais je veux un style coherent a travers tous mes templates et donc je veux le meme style pour un article quelque soit le template. Pour cela je vais me faire un readable store pour lequel je vais definir tout mes styles :
template {
    article: "bold uppercase color-red";  
    title: "bold underline tacenter";  
}

En gros dans chacun de mes template je vais importer le store et je vais mettre une grosse classe qui ne sera qu'une collection d'utilitaire.
-> Autre idee je peux aussi definir un style global mais par ou je le mets du coup parce que je n'ai pas de composants qui gere tous mes templates donc pour l'acces et la modification c'est plus complique.
### Un projet ou ii y a de l'IA et tiptap : 
Github : https://github.com/steven-tey/novel
videos Youtube : 
1. https://www.youtube.com/watch?v=bBCVI2e18dE&t=3s
2. https://www.youtube.com/watch?v=OHbiCqPSyRM

### Pour le template du devis 
Je le refais parce que celui de papa n'est pas fou, je me base sur cette photo: 
1. https://www.google.com/url?sa=i&url=https%3A%2F%2Fboby.net%2Fdocuments%2Fdevis%2F&psig=AOvVaw2drJi2_P-uHUCzOvtt5VsM&ust=1699977607980000&source=images&cd=vfe&opi=89978449&ved=0CBEQjRxqFwoTCKDTvp2swYIDFQAAAAAdAAAAABAl
2. https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.obat.fr%2Fblog%2Fmodele-devis%2F&psig=AOvVaw2drJi2_P-uHUCzOvtt5VsM&ust=1699977607980000&source=images&cd=vfe&opi=89978449&ved=0CBEQjRxqFwoTCKDTvp2swYIDFQAAAAAdAAAAABBT

### Comment optimiser le projet pour qu'il soit plus rapide 
video: https://www.youtube.com/watch?v=5Crrc9X8K5A


### Pagination avec tiptap
Implementing pagination in a TipTap editor within a Svelte app involves creating a mechanism to handle page breaks and maintain access to previous pages. Here's a step-by-step approach:

    1. Create a Page Data Structure: Define a data structure to store the content of each page. This could be an array of objects, where each object represents a page with its corresponding content.

    2. Track Page Breaks: As the user types or inserts content, monitor for specific triggers that indicate a page break should be inserted. This could be based on character count, line breaks, or specific commands or actions within the editor.

    3. Create New Pages: When a page break is detected, create a new page object and add it to the page data structure. The editor's content should then be updated to reflect the new page.

    4. Render Pages: Implement a rendering mechanism to display the current page's content within the editor. This could involve creating individual components for each page or using a dynamic rendering approach.

    <!-- NOTE: Done through the display of all the pages -->
    5. Provide Page Navigation: Provide a way for users to navigate between pages. This could be through buttons, keyboard shortcuts, or a page navigation menu.

    6. Maintain Access to Previous Pages: Ensure that users can access and edit previous pages' content. This could involve storing the page data structure in a reactive variable or using a state management library like Svelte Store.

    7. Synchronize Content: When editing content on a previous page, make sure the changes are reflected in the corresponding page object and the editor's content is updated accordingly.

    8. Handle Content Persistence: Implement a mechanism to save the page data structure, including the content of each page. This could be done using local storage, a database, or a file system.

    - [ ] 9. Restore Content on Load: Upon loading the editor, retrieve the saved page data structure and reconstruct the pages within the editor.

    10. Handle Undo/Redo: Implement undo/redo functionality to allow users to revert or redo changes made to the content across all pages.

