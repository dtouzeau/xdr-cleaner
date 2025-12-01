# Flush Feature - Memory Optimization

## Vue d'ensemble

La fonctionnalité de flush périodique limite l'utilisation mémoire lors du téléchargement de grandes quantités d'alertes.

## Principe

Au lieu de stocker toutes les alertes en mémoire avant d'écrire sur disque, le système écrit périodiquement par lots (flush) et libère la mémoire.

## Configuration

### Paramètre `flushEvery`

```json
{
  "flushEvery": 1000
}
```

- **Type**: int
- **Défaut**: 1000
- **Description**: Nombre d'alertes à accumuler avant d'écrire sur disque

## Recommandations

| Scénario | flushEvery | maxConcurrentPages | RAM estimée |
|----------|------------|-------------------|-------------|
| Faible mémoire | 500 | 20 | ~50-100 MB |
| Équilibre (défaut) | 1000 | 50 | ~100-200 MB |
| Haute performance | 5000 | 100 | ~500 MB - 1 GB |

## Architecture

### FlushManager (Flush.go)

```go
type FlushManager struct {
    mu            sync.Mutex
    outfile       string
    flushEvery    int
    currentAlerts []Alert
    totalFlushed  int
    debug         bool
    firstWrite    bool
}
```

### Méthodes principales

1. **NewFlushManager(outfile, flushEvery, debug)** - Initialisation
2. **AddAlerts(alerts)** - Ajoute des alertes, flush si nécessaire
3. **flush()** - Écrit le buffer sur disque et le vide
4. **Finalize()** - Finalise le fichier JSON

## Workflow

```
1. Téléchargement page 1-50 (parallèle)
   ↓
2. AddAlerts() → buffer += alertes
   ↓
3. Si buffer >= flushEvery:
   - Écriture sur disque (append JSON)
   - Vidage du buffer
   - Mémoire libérée
   ↓
4. Répétition jusqu'à fin
   ↓
5. Finalize() → fermeture JSON valide
```

## Format du fichier

Le fichier est écrit progressivement avec un JSON valide :

```json
{
  "Alerts": [
    { "ID": 1, ... },    ← Flush 1 (alertes 1-1000)
    { "ID": 2, ... },
    ...
    { "ID": 1000, ... },
    { "ID": 1001, ... }, ← Flush 2 (alertes 1001-2000)
    ...
  ]
}
```

## Thread Safety

- Utilise `sync.Mutex` pour la synchronisation
- Safe avec téléchargement parallèle
- Pas de race conditions

## Exemple de sortie (debug: true)

```
Flushing 1000 alerts to disk (total flushed: 1000)...
Page 12 completed: 100 alerts (total in memory: 1200)
Flushing 1000 alerts to disk (total flushed: 2000)...
Page 24 completed: 100 alerts (total in memory: 2400)
...
Finalized: total 50000 alerts written to /path/to/out.log
```

## Calcul de l'utilisation mémoire

### Facteurs

- Taille moyenne d'une alerte avec événements : ~100 KB
- Nombre d'alertes en buffer : `flushEvery`
- Overhead Go : ~20-30%

### Formule approximative

```
RAM_max = (flushEvery × 100 KB) × 1.3
```

### Exemples

| flushEvery | RAM max estimée |
|------------|-----------------|
| 500 | ~65 MB |
| 1000 | ~130 MB |
| 2000 | ~260 MB |
| 5000 | ~650 MB |
| 10000 | ~1.3 GB |

## Cas d'usage

### Scénario 1 : Serveur avec 2 GB RAM

```json
{
  "flushEvery": 500,
  "maxConcurrentPages": 20
}
```

Permet de télécharger des millions d'alertes sans problème.

### Scénario 2 : Workstation avec 16 GB RAM

```json
{
  "flushEvery": 5000,
  "maxConcurrentPages": 100
}
```

Performance maximale.

### Scénario 3 : Container avec limite 512 MB

```json
{
  "flushEvery": 200,
  "maxConcurrentPages": 10
}
```

Utilisation mémoire minimale.

## Monitoring

En mode debug, surveillez :

```
Flush every: 1000 alerts
Flushing 1000 alerts to disk (total flushed: 1000)...
Page 12 completed: 100 alerts (total in memory: 1200)
```

Si `total in memory` augmente constamment sans flush, il y a un problème.

## Limitations

1. Le fichier final contient toutes les alertes (filtrage s'applique après)
2. Le filtrage nécessite toujours les alertes en mémoire
3. Pour filtrer avec faible mémoire, utilisez une approche streaming (non implémentée)

## Performance

### Impact sur la vitesse

- Flush périodique : ~10-50ms par opération
- Négligeable par rapport au temps réseau
- Pas de ralentissement notable

### Benchmarks estimés

| Alertes | Sans flush | Avec flush (1000) |
|---------|-----------|-------------------|
| 1,000 | ~2s | ~2s |
| 10,000 | ~15s | ~15s |
| 100,000 | ~2m30s | ~2m35s |
| 1,000,000 | OOM crash | ~25min |

## Troubleshooting

### "out of memory" malgré flush

Augmentez `flushEvery` à 500 ou moins :

```json
{
  "flushEvery": 500
}
```

### Fichier JSON invalide

Le fichier peut être invalide si le programme crash avant `Finalize()`.

Solution : Le programme gère les erreurs et finalise toujours.

### Performance lente

Si trop de flush, augmentez `flushEvery` :

```json
{
  "flushEvery": 2000
}
```

## Code exemple

### Utilisation du FlushManager

```go
// Création
flushMgr := NewFlushManager("/path/to/out.log", 1000, true)

// Ajout d'alertes
alerts := fetchSomeAlerts()
if err := flushMgr.AddAlerts(alerts); err != nil {
    log.Fatal(err)
}

// Finalisation
if err := flushMgr.Finalize(); err != nil {
    log.Fatal(err)
}
```

## Références

- Fichier: `Flush.go`
- Configuration: `Config.go` (ligne 30, 68-70)
- Intégration: `main.go` (ligne 39-40, 177-237)
- Documentation: `README.md` (section "Gestion de la mémoire")
