# Text Processors

Une application de bureau développée en Go avec Fyne pour le traitement de texte et de données JSON.

## Description de l'application

Text Processors est une suite d'outils de traitement de texte qui offre une interface graphique intuitive pour effectuer des opérations courantes sur du texte et des données JSON. L'application propose plusieurs outils intégrés et la possibilité de créer des processeurs personnalisés :

## Outils intégrés

1. **JSON Formatter** : Formate et valide du JSON avec différentes options d'indentation
2. **Text Splitter** : Divise du texte selon un délimiteur spécifié
3. **Text Joiner** : Joint des lignes de texte avec un délimiteur personnalisé
4. **Pipeline Builder** : Enchaîne plusieurs outils pour créer des workflows complexes

## Processeurs personnalisés

5. **Custom Processors** : Créez vos propres processeurs de texte en JavaScript
   - Interface intuitive pour créer des processeurs personnalisés
   - Support complet de JavaScript pour la transformation de texte
   - Intégration transparente dans les pipelines
   - Export/Import automatique avec les pipelines
   - Exemples prédéfinis pour démarrer rapidement

Voir [CUSTOM_PROCESSORS.md](CUSTOM_PROCESSORS.md) pour un guide détaillé.

## Architecture du projet

### Structure des fichiers

```
text_processors/
├── main.go                 # Point d'entrée de l'application
├── go.mod                  # Gestion des dépendances Go
├── go.sum                  # Checksums des dépendances
├── build/                  # Dossier des artifacts de compilation
│   └── text_processors.exe
├── build.ps1               # Script PowerShell pour compiler dans ./build
└── ui/                     # Package contenant l'interface utilisateur
    ├── app.go              # Interface principale et navigation
    ├── pipeline.go         # Gestion des pipelines et configurations
    ├── pipeline_builder.go # Interface de construction de pipelines
    ├── tools_grid.go       # Grille de sélection des outils
    └── processors/         # Package contenant les processeurs de texte
        ├── processor.go            # Interfaces communes des processeurs
        ├── json_formatter_ui.go    # Processeur de formatage JSON
        ├── text_splitter.go        # Processeur de division de texte
        ├── text_joiner.go          # Processeur de jointure de texte
        ├── formatter.go            # Logique de formatage JSON
        └── validator.go            # Validation et gestion d'erreurs JSON
```

### Architecture logicielle

#### 1. Point d'entrée (`main.go`)
- Initialise l'application Fyne
- Crée la fenêtre principale (800x600 pixels)
- Configure le titre "JSON Formatter"
- Lance l'interface utilisateur

#### 2. Package UI (`ui/`)

**Composants principaux :**

- **`app.go`** : Gère la navigation entre les outils et l'interface principale
  - `MakeUI()` : Crée l'interface principale avec navigation
  - `MakeJSONFormatterUI()` : Interface spécifique au formatage JSON

- **`tools_grid.go`** : Système de sélection des outils
  - Structure `Tool` pour définir chaque outil
  - Grille de cartes pour la sélection
  - Système de callbacks pour la navigation

- **`formatter.go`** : Logique de formatage JSON
  - Structure `Formatter` avec options d'indentation
  - Support pour 2 espaces, 4 espaces, ou tabulations
  - Validation et formatage JSON

- **`validator.go`** : Validation et gestion d'erreurs
  - `ValidateJSON()` : Validation de syntaxe JSON
  - `PrettyValidationError()` : Formatage des erreurs pour l'affichage

- **`text_splitter.go`** : Outil de division de texte
  - Interface pour diviser du texte selon un délimiteur
  - Affichage ligne par ligne du résultat

- **`text_joiner.go`** : Outil de jointure de texte
  - Interface pour joindre des lignes avec un délimiteur
  - Filtrage automatique des lignes vides

### Technologies utilisées

- **Go 1.22** : Langage de programmation principal
- **Fyne v2.6.1** : Framework d'interface graphique multiplateforme
- **encoding/json** : Package standard Go pour le traitement JSON
- **strings** : Package standard Go pour la manipulation de chaînes

## Fonctionnalités

### JSON Formatter
- **Validation** : Vérifie la syntaxe JSON et affiche des erreurs détaillées
- **Formatage** : Indentation configurable (2 espaces, 4 espaces, tabulations)
- **Copie** : Bouton pour copier le résultat formaté dans le presse-papiers
- **Mise à jour automatique** : Reformate automatiquement lors du changement d'indentation

### Text Splitter
- **Division flexible** : Utilise n'importe quel délimiteur
- **Affichage structuré** : Une ligne par partie dans le résultat
- **Délimiteur par défaut** : Virgule si aucun délimiteur n'est spécifié

### Text Joiner
- **Jointure personnalisée** : Délimiteur configurable
- **Nettoyage automatique** : Supprime les lignes vides
- **Délimiteur par défaut** : ", " (virgule + espace)

## Règles de développement

### 1. Structure du code

- **Séparation des responsabilités** : Chaque fichier a une responsabilité claire
- **Package UI isolé** : Toute la logique d'interface dans le package `ui/`
- **Fonctions publiques** : Commencent par une majuscule (convention Go)
- **Nommage descriptif** : Noms de fonctions et variables explicites

### 2. Gestion des erreurs

- **Validation systématique** : Toujours valider les entrées utilisateur
- **Messages d'erreur clairs** : Utiliser `PrettyValidationError()` pour le JSON
- **Gestion gracieuse** : Ne jamais faire planter l'application sur une erreur utilisateur

### 3. Interface utilisateur

- **Cohérence visuelle** : Utiliser les mêmes patterns pour tous les outils
- **Responsive design** : Utiliser `container.NewBorder()` pour la mise en page
- **Feedback utilisateur** : Boutons de copie et messages de statut
- **Navigation intuitive** : Bouton retour toujours visible

### 4. Conventions de code

```go
// Structure type
type Tool struct {
    Name        string
    Description string
    Icon        fyne.Resource
    MakeUI      func() fyne.CanvasObject
}

// Fonction constructeur
func NewFormatter(indentType string) *Formatter

// Fonction d'interface
func MakeToolNameUI() fyne.CanvasObject
```

### 5. Ajout de nouveaux outils

Pour ajouter un nouvel outil :

1. **Créer le fichier UI** : `ui/nom_outil.go`
2. **Implémenter la fonction** : `MakeNomOutilUI() fyne.CanvasObject`
3. **Ajouter à la grille** : Modifier `tools_grid.go`
4. **Respecter le pattern** : Entrée → Traitement → Sortie → Copie

### 6. Tests et qualité

- **Validation des entrées** : Toujours tester les cas limites
- **Gestion mémoire** : Éviter les fuites avec les callbacks Fyne
- **Performance** : Optimiser pour les gros volumes de texte
- **Accessibilité** : Labels clairs et navigation au clavier

### 7. Compilation et distribution

```bash
# Via script PowerShell (Windows)
powershell -NoProfile -ExecutionPolicy Bypass -File .\build.ps1
# Mode release (binaire plus petit)
powershell -NoProfile -ExecutionPolicy Bypass -File .\build.ps1 -Release

# Ou directement avec go build
go build -o build/text_processors.exe .

# Cross-compilation (exemple Linux)
GOOS=linux GOARCH=amd64 go build -o build/text_processors .
```

## Installation et utilisation

### Prérequis
- Go 1.22 ou supérieur
- Dépendances automatiquement gérées par `go mod`

### Compilation
```bash
# Option 1: Script PowerShell
powershell -NoProfile -ExecutionPolicy Bypass -File .\build.ps1

# Option 2: Commande go
go mod tidy
go build -o build/text_processors.exe .
```

### Exécution
```bash
./build/text_processors.exe  # Windows
./build/text_processors      # Linux/macOS
```

## Nouveautés

### Processeurs Personnalisés (Custom Processors)
La dernière version introduit la possibilité de créer des processeurs personnalisés en JavaScript :

- **Bouton "Add custom processor"** sur la page principale
- **Interface de création intuitive** avec exemples prédéfinis
- **Support JavaScript complet** pour la transformation de texte
- **Intégration dans les pipelines** comme n'importe quel autre outil
- **Export/Import automatique** avec les configurations de pipeline
- **Zone de test intégrée** pour valider vos scripts

#### Exemples de processeurs personnalisés :
- Convertisseur majuscules/minuscules
- Compteur de mots/caractères
- Inverseur de texte
- Extracteur d'emails/URLs
- Formateur de listes
- Traitement JSON personnalisé

Consultez [CUSTOM_PROCESSORS.md](CUSTOM_PROCESSORS.md) pour un guide complet.

## Contribution

Lors de l'ajout de nouvelles fonctionnalités :
1. Respecter l'architecture existante
2. Maintenir la cohérence de l'interface
3. Ajouter la validation appropriée
4. Tester avec différents types d'entrées
5. Documenter les nouvelles fonctions
6. Tester l'intégration avec les processeurs personnalisés
