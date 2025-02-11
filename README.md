curl http://localhost:8000/blockchain

curl -X POST http://localhost:8000/transaction -H "Content-Type: application/json" -d "{\"sender\":\"Alice\",\"recipient\":\"Bob\",\"amount\":10}"

curl http://localhost:8000/mine
