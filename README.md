# XDR Cleaner

Outil de gestion d'alertes XDR avec support de téléchargement parallèle, filtrage avancé et clôture automatique.

## Fonctionnalités

- ✅ **Téléchargement parallèle** : Récupère jusqu'à 50 pages simultanément
- ✅ **Filtrage avancé** : Recherche dans les champs imbriqués (Observable, Rule, BaseEvent, Alert)
- ✅ **Clôture automatique** : Ferme les alertes filtrées via l'API
- ✅ **Gestion des erreurs** : Retry automatique avec backoff
- ✅ **Mode debug** : Logs détaillés pour le dépannage

## Installation

```bash
go build -o xdr-cleaner
```

## Configuration

Le fichier `config.json` est créé automatiquement au premier lancement. Éditez-le avec vos paramètres :

```json
{
  "pageNumber": 1,
  "tenantID": "3d6c4203-7328-4fb2-98fa-f37d385ffbde",
  "token": "VOTRE_TOKEN_BEARER",
  "baseURL": "https://api.louni.priv/xdr/api/v1",
  "maxConcurrentPages": 50,
  "filterMode": true,
  "filteredOutfile": "/etc/xdr-cleaner/filtered.json",
  "filters": [
    {
      "field": "Observable|Value",
      "value": "192.168.1.100"
    }
  ],
  "closeAlerts": true,
  "closeReason": "falsePositive"
}
```

### Paramètres principaux

| Paramètre | Type | Description |
|-----------|------|-------------|
| `tenantID` | string | ID du tenant (requis) |
| `token` | string | Bearer token d'authentification (requis) |
| `baseURL` | string | URL de base de l'API XDR |
| `maxConcurrentPages` | int | Nombre de pages à télécharger en parallèle (défaut: 50) |
| `outfile` | string | Fichier de sortie pour toutes les alertes |
| `filterMode` | bool | Active le filtrage des alertes |
| `filteredOutfile` | string | Fichier de sortie pour les alertes filtrées |
| `closeAlerts` | bool | Active la clôture automatique des alertes filtrées |
| `closeReason` | string | Raison de clôture (falsePositive, resolved, duplicate, etc.) |
| `flushEvery` | int | Nombre d'alertes avant flush sur disque (défaut: 1000) - limite l'utilisation mémoire |
| `debug` | bool | Active les logs détaillés |

### Filtres disponibles

Les filtres utilisent le format `"Section|Champ"` :

#### Section "Observable"
- `Observable|Value` - Valeur de l'observable (IP, hash, URL, etc.)
- `Observable|Type` - Type d'observable
- `Observable|Details` - Détails de l'observable

#### Section "Rule"
- `Rule|Name` - Nom de la règle de détection
- `Rule|ID` - ID de la règle
- `Rule|Type` - Type de règle
- `Rule|Severity` - Sévérité de la règle
- `Rule|Confidence` - Niveau de confiance

#### Section "BaseEvent"
- `BaseEvent|DestinationAddress` - IP de destination
- `BaseEvent|SourceAddress` - IP source
- `BaseEvent|DeviceAddress` - IP du device
- `BaseEvent|DeviceHostName` - Nom d'hôte du device
- `BaseEvent|DeviceAction` - Action du device
- `BaseEvent|DeviceVendor` - Vendeur du device
- `BaseEvent|DeviceProduct` - Produit du device
- `BaseEvent|TransportProtocol` - Protocole de transport
- `BaseEvent|ApplicationProtocol` - Protocole applicatif
- `BaseEvent|Message` - Message de l'événement
- `BaseEvent|DestinationPort` - Port de destination
- `BaseEvent|SourcePort` - Port source

#### Section "Alert"
- `Alert|Name` - Nom de l'alerte
- `Alert|Severity` - Sévérité de l'alerte
- `Alert|Status` - Statut de l'alerte
- `Alert|InternalID` - ID interne
- `Alert|IncidentID` - ID de l'incident

### Filtres multiples

Les filtres sont combinés avec un **ET logique**. Une alerte doit correspondre à **tous** les filtres pour être sélectionnée.

Exemple : filtrer les alertes avec IP source 10.0.0.5 **ET** règle contenant "Malware"

```json
"filters": [
  {
    "field": "BaseEvent|SourceAddress",
    "value": "10.0.0.5"
  },
  {
    "field": "Rule|Name",
    "value": "Malware"
  }
]
```

## Utilisation

### 1. Télécharger toutes les alertes

```bash
./xdr-cleaner
```

Résultat : `out.log` contient toutes les alertes récupérées.

### 2. Télécharger et filtrer

Activez `filterMode: true` et configurez vos filtres dans `config.json` :

```bash
./xdr-cleaner
```

Résultat :
- `out.log` - Toutes les alertes
- `filtered.json` - Alertes correspondant aux filtres

### 3. Télécharger, filtrer ET clôturer

Activez `filterMode: true` et `closeAlerts: true` :

```bash
./xdr-cleaner
```

Résultat :
- `out.log` - Toutes les alertes
- `filtered.json` - Alertes filtrées
- Les alertes filtrées sont automatiquement clôturées via l'API

## Exemples de scénarios

### Scénario 1 : Clôturer les faux positifs pour une IP interne

```json
{
  "filterMode": true,
  "filters": [
    {
      "field": "BaseEvent|DestinationAddress",
      "value": "192.168.1.50"
    }
  ],
  "closeAlerts": true,
  "closeReason": "falsePositive"
}
```

### Scénario 2 : Trouver toutes les alertes d'une règle spécifique

```json
{
  "filterMode": true,
  "filters": [
    {
      "field": "Rule|Name",
      "value": "Windows Defender Alert"
    }
  ],
  "closeAlerts": false
}
```

### Scénario 3 : Clôturer les alertes low severity d'un device particulier

```json
{
  "filterMode": true,
  "filters": [
    {
      "field": "Alert|Severity",
      "value": "low"
    },
    {
      "field": "BaseEvent|DeviceHostName",
      "value": "srv-backup-01"
    }
  ],
  "closeAlerts": true,
  "closeReason": "resolved"
}
```

## Options de dates

Pour filtrer par période :

```json
{
  "fromDate": "2024-01-01T00:00:00Z",
  "toDate": "2024-12-31T23:59:59Z"
}
```

Format requis : RFC3339 (ISO 8601)

## Performance

### Téléchargement parallèle

Le paramètre `maxConcurrentPages` contrôle la concurrence :

- **10-20** : Serveurs avec limitations de rate
- **50** (défaut) : Bon équilibre performance/charge serveur
- **100+** : Serveurs haute performance

### Gestion de la mémoire

Le paramètre `flushEvery` limite l'utilisation mémoire en écrivant périodiquement sur disque :

- **500** : Faible utilisation mémoire (~50-100 MB selon la taille des alertes)
- **1000** (défaut) : Équilibre mémoire/performance (~100-200 MB)
- **5000** : Haute performance, plus de mémoire (~500 MB - 1 GB)

**Comment ça fonctionne :**
1. Les alertes sont téléchargées en parallèle
2. Toutes les `flushEvery` alertes, elles sont écrites sur disque
3. Le buffer mémoire est vidé
4. Le téléchargement continue sans surcharge mémoire

**Exemple :** Pour 100,000 alertes avec `flushEvery: 1000` :
- 100 flush opérations
- Mémoire maximale : ~1000 alertes en RAM
- Fichier écrit progressivement

### Clôture parallèle

Les alertes sont clôturées avec :
- **10 requêtes simultanées** maximum
- **3 tentatives** en cas d'erreur
- **Retry avec backoff** pour les erreurs 5xx

## Mode Debug

Activez `debug: true` pour obtenir :

```
Fetching page 1: https://api.louni.priv/xdr/api/v1?page=1&tenantID=...
Page 1 completed: 100 alerts
Matched Observable.Value: 192.168.1.100 contains 192.168.1
✓ Closed: c445d5bb-d426-46d2-8c91-9ff4a8cb044c (Suspicious Activity)
```

## API de clôture

### Endpoint

```
POST /xdr/api/v1/alerts/close?tenantID={tenantID}
```

### Payload

```json
{
  "ID": "c445d5bb-d426-46d2-8c91-9ff4a8cb044c",
  "TenantID": "3d6c4203-7328-4fb2-98fa-f37d385ffbde",
  "Reason": "falsePositive"
}
```

### Raisons de clôture supportées

- `falsePositive` - Faux positif
- `resolved` - Résolu
- `duplicate` - Doublon
- `testing` - Test
- `accepted_risk` - Risque accepté

## Dépannage

### Erreur "TenantID is required"

Vérifiez que `tenantID` est renseigné dans `config.json`.

### Erreur "Token is required"

Ajoutez votre Bearer token dans le champ `token`.

### Erreur HTTP 401

Le token est invalide ou expiré. Générez un nouveau token.

### Erreur HTTP 403

Le token n'a pas les permissions pour clôturer les alertes.

### Aucune alerte filtrée

- Vérifiez le format des filtres : `"Section|Champ"`
- Activez `debug: true` pour voir les matchs
- Les filtres sont sensibles à la casse (par défaut en mode insensible)

## Structure des fichiers

```
xdr-cleaner/
├── xdr-cleaner          # Binaire exécutable
├── config.json          # Configuration (généré automatiquement)
├── config.example.json  # Exemple de configuration
├── out.log              # Toutes les alertes (JSON)
├── filtered.json        # Alertes filtrées (JSON)
├── main.go              # Point d'entrée
├── Config.go            # Gestion de la configuration
├── Filter.go            # Logique de filtrage
├── Close.go             # API de clôture
├── Tools.go             # Utilitaires HTTP
└── structs.go           # Structures de données
```

## Sécurité

⚠️ **Attention** : Le client HTTP désactive la vérification TLS (`InsecureSkipVerify: true`) dans `Tools.go:61`.

Cela est acceptable pour :
- Réseaux internes avec certificats auto-signés
- Environnements de test

Pour la production avec certificats valides, modifiez `BuilClient()` :

```go
func BuilClient() *http.Client {
    // Supprimez TLSClientConfig pour une vérification complète
    return &http.Client{}
}
```

## Licence

Ce projet est fourni tel quel sans garantie.
