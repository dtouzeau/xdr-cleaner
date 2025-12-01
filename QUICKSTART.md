# Quick Start Guide - XDR Cleaner

## 🚀 Démarrage rapide en 3 étapes

### Étape 1 : Configuration initiale

```bash
# Compiler le binaire
go build -o xdr-cleaner

# Générer config.json
./xdr-cleaner
```

### Étape 2 : Éditer config.json

Remplissez les champs requis :

```json
{
  "tenantID": "0a0a0000-0000-0aa0-00aa-a00a000aaaa",
  "token": "VOTRE_TOKEN_ICI",
  "baseURL": "https://your.xdr.addr/xdr/api/v1"
}
```

### Étape 3 : Exécuter

```bash
./xdr-cleaner
```

---

## 📋 Cas d'usage courants

### 1️⃣ Télécharger toutes les alertes

**config.json minimal :**
```json
{
  "tenantID": "votre-tenant-id",
  "token": "votre-token",
  "baseURL": "https://your.xdr.addr/xdr/api/v1",
  "filterMode": false,
  "closeAlerts": false
}
```

**Commande :**
```bash
./xdr-cleaner
```

**Résultat :** Fichier `out.log` avec toutes les alertes

---

### 2️⃣ Filtrer les alertes par IP de destination

**config.json :**
```json
{
  "tenantID": "votre-tenant-id",
  "token": "votre-token",
  "baseURL": "https://your.xdr.addr/xdr/api/v1",
  "filterMode": true,
  "filters": [
    {
      "field": "BaseEvent|DestinationAddress",
      "value": "10.0.0.5"
    }
  ],
  "closeAlerts": false
}
```

**Commande :**
```bash
./xdr-cleaner
```

**Résultat :**
- `out.log` - Toutes les alertes
- `filtered.json` - Alertes avec destination 10.0.0.5

---

### 3️⃣ Clôturer automatiquement les faux positifs

**config.json :**
```json
{
  "tenantID": "0a0a0000-0000-0aa0-00aa-a00a000aaaa",
  "token": "votre-token",
  "baseURL": "https://your.xdr.addr/xdr/api/v1",
  "filterMode": true,
  "filters": [
    {
      "field": "Rule|Name",
      "value": "Windows Defender False Positive"
    }
  ],
  "closeAlerts": true,
  "closeReason": "falsePositive"
}
```

**Commande :**
```bash
./xdr-cleaner
```

**Résultat :**
- `out.log` - Toutes les alertes
- `filtered.json` - Alertes filtrées
- **Alertes filtrées sont clôturées automatiquement**

---

### 4️⃣ Filtres multiples (ET logique)

**Exemple :** Alertes low severity sur un serveur spécifique

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
      "value": "srv-backup"
    }
  ],
  "closeAlerts": true,
  "closeReason": "accepted_risk"
}
```

Une alerte doit correspondre aux **deux** filtres pour être sélectionnée.

---

## 🔍 Champs de filtrage populaires

### Par adresse IP
```json
{"field": "BaseEvent|SourceAddress", "value": "192.168.1.100"}
{"field": "BaseEvent|DestinationAddress", "value": "10.0.0.5"}
{"field": "Observable|Value", "value": "8.8.8.8"}
```

### Par règle
```json
{"field": "Rule|Name", "value": "Malware"}
{"field": "Rule|Severity", "value": "low"}
```

### Par alerte
```json
{"field": "Alert|Name", "value": "Suspicious"}
{"field": "Alert|Severity", "value": "medium"}
{"field": "Alert|Status", "value": "open"}
```

### Par device
```json
{"field": "BaseEvent|DeviceHostName", "value": "firewall-01"}
{"field": "BaseEvent|DeviceVendor", "value": "Palo Alto"}
```

---

## 🐛 Debug

Activez les logs détaillés :

```json
{
  "debug": true
}
```

Vous verrez :
```
Fetching page 1: https://your.xdr.addr/xdr/api/v1?page=1
Page 1 completed: 100 alerts
Matched BaseEvent.DestinationAddress: 10.0.0.5 contains 10.0.0.5
✓ Closed: c445d5bb-d426-46d2-8c91-9ff4a8cb044c (Alert Name)
```

---

## ⚙️ Performance

### Ajuster la concurrence

```json
{
  "maxConcurrentPages": 20
}
```

- **10-20** : Serveurs lents ou avec rate limiting
- **50** (défaut) : Serveurs normaux
- **100+** : Serveurs haute performance

### Limiter l'utilisation mémoire

```json
{
  "flushEvery": 500
}
```

- **500** : Faible mémoire (~50-100 MB)
- **1000** (défaut) : Équilibre (~100-200 MB)
- **5000** : Haute performance (~500 MB - 1 GB)

Les alertes sont écrites sur disque tous les `flushEvery` téléchargements, libérant la mémoire.

---

## 📊 Workflow complet

```
1. Téléchargement parallèle
   ↓
2. Sauvegarde de toutes les alertes (out.log)
   ↓
3. Application des filtres (si activés)
   ↓
4. Sauvegarde des alertes filtrées (filtered.json)
   ↓
5. Clôture automatique (si activée)
   ↓
6. Rapport de synthèse
```

---

## 📝 Exemple de sortie

```
Fetching alerts...
Total alerts fetched: 1250

Saved 1250 alerts to /etc/xdr-cleaner/out.log

=== Filtering Alerts ===
Filtered 45 alerts from 1250 total
Saved 45 filtered alerts to /etc/xdr-cleaner/filtered.json

=== Closing Filtered Alerts ===
Starting to close 45 alerts...

Close Summary:
  Success: 43
  Failed:  2
  Total:   45
```

---

## ❓ FAQ

**Q: Les filtres sont-ils sensibles à la casse ?**
R: Non, la recherche est insensible à la casse.

**Q: Puis-je utiliser des regex dans les filtres ?**
R: Non, les filtres utilisent `contains` (recherche de sous-chaîne).

**Q: Que se passe-t-il si une clôture échoue ?**
R: L'outil fait 3 tentatives avec backoff. Les échecs sont loggés.

**Q: Puis-je clôturer sans filtrer ?**
R: Non, vous devez activer `filterMode` pour utiliser `closeAlerts`.

**Q: Comment récupérer seulement les alertes ouvertes ?**
R: Utilisez le paramètre dans config.json :
```json
{"status": "open"}
```

---

## 🔗 Pour aller plus loin

Consultez le [README.md](README.md) complet pour :
- Liste exhaustive des champs de filtrage
- API de clôture détaillée
- Options avancées de configuration
- Dépannage approfondi
