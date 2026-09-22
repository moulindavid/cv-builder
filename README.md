# CV Builder

CLI Go générant les CV français et anglais de David Moulin depuis une source YAML unique, avec un PDF LaTeX sobre et vérifiable par les ATS.

## Principes

- `data/resume.yaml` est l'unique source de vérité.
- Les faits (entreprises, dates, technologies) ne sont pas dupliqués par langue ; seuls les textes sont localisés.
- Le tailoring classe uniquement des éléments déjà présents. Il n'ajoute jamais une compétence, une expérience ou un résultat.
- Le PDF utilise une colonne, du texte Unicode réel, des titres standards et aucune table, image, sidebar ou information essentielle en en-tête/pied de page.

## Prérequis Linux

- Go 1.24 ou ultérieur : voir <https://go.dev/doc/install>
- LuaLaTeX et les paquets utilisés par le template
- `pdftotext` fourni par Poppler

Debian/Ubuntu :

```sh
sudo apt update
sudo apt install texlive-luatex texlive-latex-recommended texlive-fonts-recommended poppler-utils
```

Fedora :

```sh
sudo dnf install golang texlive-luatex texlive-collection-latexrecommended poppler-utils
```

## Démarrage

```sh
make setup
make test
make lint
make build
make resume
```

Les noms de fichiers sont dérivés du nom du profil, du target et de la langue. Les PDF fournis par `make resume` sont créés dans `output/` :

- `david-moulin-senior-backend-fr.pdf`
- `david-moulin-senior-backend-en.pdf`

Les PDF et fichiers LaTeX générés ne sont pas versionnés : ils sont reproductibles, augmentent inutilement la taille du dépôt et créeraient des diffs binaires. `output/.gitkeep` conserve le répertoire.

## CLI

```sh
./bin/cv validate
./bin/cv build                         # français par défaut
./bin/cv build --lang en
./bin/cv build --lang en --target senior-backend
./bin/cv build --target targets/senior-backend.yaml
./bin/cv extract-text --lang en --target senior-backend
./bin/cv extract-text --pdf output/example.pdf
./bin/cv tailor --job job.txt --lang en
./bin/cv tailor --job job.txt --lang en --build
./bin/cv tailor --job job.txt --lang en --build --max-bullets 4
```

`build` valide toujours le YAML avant le rendu. Sans target, il produit un fichier `*-resume-<lang>.pdf`. `extract-text` accepte le même `--target` que `build`; `--pdf` permet toujours de viser un fichier explicite.

## Structure des données

`data/resume.yaml` contient `profile`, `contact`, `summaries`, `skills`, `experience`, `education`, `languages`, `links` et la configuration `tailoring`. Les champs `text`, `roles`, `titles`, `descriptions` et `degrees` sont des objets localisés :

```yaml
text:
  fr: Développement d'une API REST.
  en: Developed a REST API.
```

Les technologies d'une expérience référencent les `id` des skills. La validation rejette notamment les références inconnues, les identifiants dupliqués, les traductions obligatoires manquantes, les catégories inconnues, les dates invalides ou incohérentes, un email incorrect et une URL invalide.

### Ajouter une expérience

Ajouter un objet à `experience`, avec un `id` stable, les faits, les textes FR/EN et des bullets identifiés et tagués. Ne mettre dans `technologies` que des IDs définis dans `skills`, puis lancer `./bin/cv validate`.

### Ajouter une compétence

Ajouter à `skills` un `id` stable, son libellé, sa catégorie et ses tags. Le champ optionnel `aliases` contient uniquement des noms explicites équivalents utilisés pour le matching d'offres :

```yaml
- id: java
  name: Java
  category: Langages
  tags: [java, backend]
  aliases: [JVM]
```

### Créer un target

Créer `targets/mon-target.yaml` :

```yaml
name: mon-target
titles:
  fr: Senior Backend Engineer
  en: Senior Backend Engineer
priorities: [java, postgresql, messaging]
max_bullets: 5
```

Les priorités réordonnent les compétences et bullets par correspondance avec leurs IDs/tags. `max_bullets` peut limiter les bullets par expérience. Un target ne crée aucun contenu. Lors d’un build, son nom, ses langues, ses priorités et `max_bullets` sont validés contre le CV.

## Tailoring

`cv tailor` normalise le texte de l'offre, recherche de manière déterministe les noms, IDs et aliases explicites des compétences connues (y compris les expressions multi-mots), produit un score indicatif et affiche les correspondances ainsi qu'un résumé existant sûr. Les tags présents dans l'offre servent uniquement à classer les skills et bullets, pas à déclarer artificiellement une compétence comme correspondante.

Le catalogue `tailoring.external_keywords` du YAML définit les technologies externes que le rapport peut signaler comme absentes. `--build` génère `output/<nom>-tailored-<lang>.pdf`; `--max-bullets` peut limiter les bullets après classement. Le tailoring reste local et déterministe : une technologie absente du CV n'est jamais injectée. Le score mesure la part des compétences détectées parmi les compétences détectées et les mots-clés externes absents ; ce n'est pas une probabilité de recrutement.

## Tests ATS

Le test `tests/ats_integration_test.go` rend un `.tex`, compile un PDF, lance `pdftotext`, vérifie l’ordre `NAME → SUMMARY → TECHNICAL SKILLS → EXPERIENCE → EDUCATION`, les informations essentielles et leur absence de duplication. Il est ignoré avec un motif explicite si `lualatex` ou `pdftotext` manque.


## Vérification et qualité

```sh
make test       # tests unitaires et intégration ATS si les outils sont présents
make lint       # gofmt puis go vet
make resume     # validation, rendu et compilation FR/EN
./bin/cv extract-text --lang en --target senior-backend
```

La sortie de `extract-text` doit conserver un ordre linéaire, des coordonnées et URLs extractibles, puis les sections résumé, compétences, expérience et formation. Le test ATS est skipped proprement si LuaLaTeX ou `pdftotext` n’est pas installé.
