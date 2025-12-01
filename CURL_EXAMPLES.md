# Exemples de requêtes cURL pour l'API XDR

## Variables d'environnement

```bash
export API_URL="https://api.louni.priv/xdr/api/v1"
export TENANT_ID="3d6c4203-7328-4fb2-98fa-f37d385ffbde"
export TOKEN="AebikvWr5Qx-QmPipsqjHEs_6RUh0ovx5DVm7zXX4iY.d82UaJHL0asAnyYluVE46BZmazV9ILYN3pUQBgVhpf5Vh3hgisOI5mGKxBlyITgtKB-zP9G_3hsbI41czPj9LQ"
```

## 1. Récupérer les alertes (GET)

### Première page
```bash
curl -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json"
```

### Avec filtres de date
```bash
curl -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}&from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json"
```

### Avec événements inclus
```bash
curl -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}&withEvents=true" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json"
```

### Filtrer par statut
```bash
curl -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}&status=open" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json"
```

## 2. Clôturer une alerte (POST)

### Clôture simple
```bash
curl -v -k -X POST \
  "${API_URL}/alerts/close?tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "ID": "c445d5bb-d426-46d2-8c91-9ff4a8cb044c",
    "TenantID": "'"${TENANT_ID}"'",
    "Reason": "falsePositive"
  }'
```

### Clôture avec raison "resolved"
```bash
curl -v -k -X POST \
  "${API_URL}/alerts/close?tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "ID": "alert-id-here",
    "TenantID": "'"${TENANT_ID}"'",
    "Reason": "resolved"
  }'
```

### Clôture avec raison "duplicate"
```bash
curl -v -k -X POST \
  "${API_URL}/alerts/close?tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "ID": "alert-id-here",
    "TenantID": "'"${TENANT_ID}"'",
    "Reason": "duplicate"
  }'
```

## 3. Raisons de clôture supportées

| Raison | Description |
|--------|-------------|
| `falsePositive` | Faux positif - alerte incorrecte |
| `resolved` | Résolu - problème corrigé |
| `duplicate` | Doublon - alerte en double |
| `testing` | Test - alerte de test |
| `accepted_risk` | Risque accepté - confirmé mais accepté |

## 4. Test de connectivité

### Vérifier l'authentification
```bash
curl -v -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n"
```

### Tester avec un mauvais token
```bash
curl -v -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer INVALID_TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n"
```

## 5. Scripts de test

### Script pour clôturer plusieurs alertes
```bash
#!/bin/bash
# close_alerts.sh

API_URL="https://api.louni.priv/xdr/api/v1"
TENANT_ID="3d6c4203-7328-4fb2-98fa-f37d385ffbde"
TOKEN="votre-token"

ALERT_IDS=(
  "c445d5bb-d426-46d2-8c91-9ff4a8cb044c"
  "d556e6cc-e537-57e3-9da2-a0f5b9dc155d"
  "e667f7dd-f648-68f4-ae03-b1a6c8ed266e"
)

for ALERT_ID in "${ALERT_IDS[@]}"; do
  echo "Closing alert: $ALERT_ID"

  curl -k -X POST \
    "${API_URL}/alerts/close?tenantID=${TENANT_ID}" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d '{
      "ID": "'"${ALERT_ID}"'",
      "TenantID": "'"${TENANT_ID}"'",
      "Reason": "falsePositive"
    }'

  echo -e "\n---"
  sleep 1
done
```

### Script pour récupérer toutes les pages
```bash
#!/bin/bash
# fetch_all_pages.sh

API_URL="https://api.louni.priv/xdr/api/v1"
TENANT_ID="3d6c4203-7328-4fb2-98fa-f37d385ffbde"
TOKEN="votre-token"
OUTPUT_DIR="./alerts"

mkdir -p "$OUTPUT_DIR"

PAGE=1
while true; do
  echo "Fetching page $PAGE..."

  RESPONSE=$(curl -k -s -X GET \
    "${API_URL}?page=${PAGE}&tenantID=${TENANT_ID}" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json")

  echo "$RESPONSE" > "${OUTPUT_DIR}/page_${PAGE}.json"

  ALERT_COUNT=$(echo "$RESPONSE" | jq '.Alerts | length')
  echo "  Received $ALERT_COUNT alerts"

  if [ "$ALERT_COUNT" -lt 100 ]; then
    echo "Last page reached"
    break
  fi

  PAGE=$((PAGE + 1))
done

echo "All pages saved to $OUTPUT_DIR"
```

## 6. Extraction de données avec jq

### Compter les alertes
```bash
curl -k -s -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | jq '.Alerts | length'
```

### Extraire les IDs des alertes
```bash
curl -k -s -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | jq -r '.Alerts[].InternalID'
```

### Filtrer par sévérité
```bash
curl -k -s -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | jq '.Alerts[] | select(.Severity == "high")'
```

### Extraire les Observable values
```bash
curl -k -s -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | jq -r '.Alerts[].Observables[].Value'
```

## 7. Débogage

### Voir les headers de réponse
```bash
curl -v -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" 2>&1 | grep -E '^(<|>)'
```

### Mesurer le temps de réponse
```bash
curl -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nTime total: %{time_total}s\n" \
  -o /dev/null -s
```

### Enregistrer la réponse complète
```bash
curl -v -k -X GET \
  "${API_URL}?page=1&tenantID=${TENANT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  > response.json 2> debug.log
```

## 8. Codes de retour HTTP

| Code | Signification | Action |
|------|---------------|--------|
| 200 | OK | Succès |
| 401 | Unauthorized | Vérifier le token |
| 403 | Forbidden | Vérifier les permissions |
| 404 | Not Found | Vérifier l'URL et l'ID |
| 429 | Too Many Requests | Ralentir les requêtes |
| 500 | Internal Server Error | Réessayer plus tard |

## Notes

- L'option `-k` désactive la vérification SSL (pour certificats auto-signés)
- L'option `-v` active le mode verbose
- L'option `-s` active le mode silencieux
- Remplacez les variables `${API_URL}`, `${TENANT_ID}`, `${TOKEN}` par vos valeurs
