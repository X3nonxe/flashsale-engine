#!/bin/bash

echo "🔥 MULAI PERANG FLASH SALE (50 Request Paralel)..."

# Kita akan looping 50 kali
for i in {1..50}
do
   # Generate User ID acak (1000 - 9000) agar tidak kena error Duplicate Key DB
   USER_ID=$((1000 + i))
   
   # Kirim request di background (&) supaya jalan berbarengan/paralel
   curl -s -X POST http://localhost:8080/api/v1/flash-sale/purchase \
   -H "Content-Type: application/json" \
   -d "{\"user_id\": $USER_ID, \"product_id\": 1, \"quantity\": 1}" &
done

wait
echo "✅ PERANG SELESAI."