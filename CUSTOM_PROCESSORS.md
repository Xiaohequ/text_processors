# Processeurs Personnalisés

Cette fonctionnalité permet de créer vos propres processeurs de texte en utilisant JavaScript.

## Comment créer un processeur personnalisé

### 1. Cliquer sur "Add custom processor"
Sur la page principale de l'application, cliquez sur le bouton **"Add custom processor"**.

### 2. Configurer le processeur
Dans la boîte de dialogue qui s'ouvre :

- **Nom du processeur** : Donnez un nom descriptif à votre processeur
- **Script JavaScript** : Écrivez votre code de traitement

### 3. Écrire le script JavaScript

Votre script doit contenir une fonction `process(input)` qui :
- Prend en paramètre `input` (le texte à traiter)
- Retourne le texte transformé

#### Exemples de scripts :

**Convertir en majuscules :**
```javascript
function process(input) {
    return input.toUpperCase();
}
```

**Ou plus simplement :**
```javascript
return input.toUpperCase();
```

**Compter les mots :**
```javascript
function process(input) {
    var words = input.split(/\s+/).filter(word => word.length > 0);
    return 'Nombre de mots: ' + words.length;
}
```

**Inverser le texte :**
```javascript
return input.split('').reverse().join('');
```

**Traitement JSON :**
```javascript
function process(input) {
    try {
        var obj = JSON.parse(input);
        return "Nom: " + obj.name + ", Age: " + obj.age;
    } catch (e) {
        return "Erreur: " + e.message;
    }
}
```

**Ajouter des numéros de ligne :**
```javascript
function process(input) {
    var lines = input.split('\n');
    return lines.map(function(line, index) {
        return (index + 1) + '. ' + line;
    }).join('\n');
}
```

### 4. Tester le processeur
Utilisez la zone de test dans la boîte de dialogue pour vérifier que votre script fonctionne correctement.

### 5. Ajouter le processeur
Cliquez sur **"Ajouter"** pour enregistrer votre processeur personnalisé.

## Utiliser les processeurs personnalisés

### Dans le Pipeline Builder
1. Ouvrez le **Pipeline Builder**
2. Dans la liste des outils disponibles, vous verrez vos processeurs personnalisés avec le préfixe **"Custom: "**
3. Sélectionnez votre processeur et ajoutez-le au pipeline comme n'importe quel autre outil

### Export/Import
Les processeurs personnalisés sont automatiquement inclus dans l'export/import des pipelines :
- Lors de l'export, la configuration complète (nom + script) est sauvegardée
- Lors de l'import, les processeurs sont automatiquement recréés

## Fonctionnalités JavaScript supportées

### Variables et fonctions disponibles
- `input` : Le texte d'entrée à traiter
- Toutes les fonctions JavaScript standard (String, Array, JSON, etc.)

### Exemples avancés

**Formatage de liste :**
```javascript
function process(input) {
    var items = input.split('\n').filter(item => item.trim().length > 0);
    return items.map(function(item, index) {
        return '• ' + item.trim();
    }).join('\n');
}
```

**Extraction d'emails :**
```javascript
function process(input) {
    var emailRegex = /\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b/g;
    var emails = input.match(emailRegex) || [];
    return emails.join('\n');
}
```

**Statistiques de texte :**
```javascript
function process(input) {
    var chars = input.length;
    var words = input.split(/\s+/).filter(w => w.length > 0).length;
    var lines = input.split('\n').length;
    
    return 'Caractères: ' + chars + '\n' +
           'Mots: ' + words + '\n' +
           'Lignes: ' + lines;
}
```

## Conseils et bonnes pratiques

### 1. Gestion d'erreurs
Toujours inclure une gestion d'erreurs pour les opérations qui peuvent échouer :
```javascript
function process(input) {
    try {
        // Votre code ici
        return result;
    } catch (e) {
        return "Erreur: " + e.message;
    }
}
```

### 2. Validation d'entrée
Vérifiez que l'entrée est valide avant de la traiter :
```javascript
function process(input) {
    if (!input || input.trim().length === 0) {
        return "Entrée vide";
    }
    // Traitement...
}
```

### 3. Performance
Pour de gros volumes de texte, optimisez vos scripts :
- Évitez les boucles imbriquées complexes
- Utilisez les méthodes natives JavaScript quand possible

### 4. Réutilisabilité
Créez des processeurs génériques qui peuvent être réutilisés dans différents contextes.

## Limitations

- Pas d'accès aux APIs externes (fetch, XMLHttpRequest)
- Pas d'accès au système de fichiers
- Pas d'accès aux modules Node.js
- Exécution dans un environnement JavaScript isolé

## Exemples de cas d'usage

1. **Formatage de données** : Convertir des formats de données (CSV vers JSON, etc.)
2. **Nettoyage de texte** : Supprimer caractères indésirables, normaliser espaces
3. **Extraction d'informations** : Extraire emails, URLs, numéros de téléphone
4. **Transformation de contenu** : Convertir Markdown vers HTML, etc.
5. **Calculs sur texte** : Statistiques, comptages, analyses

Les processeurs personnalisés offrent une flexibilité infinie pour adapter l'application à vos besoins spécifiques !
