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
make vet
make build
make resume
```

Les PDF sont créés dans `output/` :

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
./bin/cv extract-text --lang en
./bin/cv extract-text --pdf output/example.pdf
./bin/cv tailor --job job.txt --lang en
```

`build` valide toujours le YAML avant le rendu. `extract-text` affiche exactement le texte vu par `pdftotext`.

## Structure des données

`data/resume.yaml` contient `profile`, `contact`, `summaries`, `skills`, `experience`, `education`, `languages` et `links`. Les champs `text`, `roles`, `titles`, `descriptions` et `degrees` sont des objets localisés :

```yaml
text:
  fr: Développement d'une API REST.
  en: Developed a REST API.
```

Les technologies d'une expérience référencent les `id` des skills. La validation rejette toute référence inconnue, les identifiants dupliqués, les traductions obligatoires manquantes, un email incorrect et une URL invalide.

### Ajouter une expérience

Ajouter un objet à `experience`, avec un `id` stable, les faits, les textes FR/EN et des bullets identifiés et tagués. Ne mettre dans `technologies` que des IDs définis dans `skills`, puis lancer `./bin/cv validate`.

### Ajouter une compétence

Ajouter à `skills` un `id` stable, son libellé, sa catégorie et ses tags :

```yaml
- id: java
  name: Java
  category: Langages
  tags: [java, backend]
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

Les priorités réordonnent les compétences et bullets par correspondance avec leurs IDs/tags. `max_bullets` peut limiter les bullets par expérience. Un target ne crée aucun contenu.

## Tailoring

`cv tailor` normalise le texte de l'offre, recherche les compétences connues, produit un score indicatif et affiche les correspondances ainsi qu'un résumé existant sûr. La V1 est volontairement déterministe et locale. Les mots absents ne sont jamais injectés dans le CV. Le score mesure la part du catalogue de compétences mentionnée par l'offre ; ce n'est pas une probabilité de recrutement.

## Tests ATS

Le test `tests/ats_integration_test.go` rend un `.tex`, compile un PDF, lance `pdftotext`, vérifie la présence de contenu et l'ordre `NAME → SUMMARY → TECHNICAL SKILLS → EXPERIENCE → EDUCATION`. Il est ignoré avec un motif explicite si `lualatex` ou `pdftotext` manque.
