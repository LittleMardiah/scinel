#!/bin/bash

echo "=========================================="
echo "  FIX FRONTEND ERRORS                    "
echo "=========================================="
echo ""

echo "--- 1. FIX InventoryTable.tsx (use client) ---"
FILE="frontend/components/inventory/InventoryTable.tsx"
if [ -f "$FILE" ]; then
  # Hapus 'use client' yang ada
  sed -i '/"use client"/d' "$FILE"
  # Tambahkan di baris pertama
  echo "'use client';" | cat - "$FILE" > /tmp/InventoryTable.tsx && mv /tmp/InventoryTable.tsx "$FILE"
  echo "✅ $FILE fixed"
else
  echo "❌ $FILE not found"
fi

echo ""
echo "--- 2. FIX lib/api.ts (default export) ---"
API_FILE="frontend/lib/api.ts"
if [ -f "$API_FILE" ]; then
  # Cek apakah ada export default
  if ! grep -q "export default" "$API_FILE"; then
    echo "export default apiClient;" >> "$API_FILE"
    echo "✅ Added default export to $API_FILE"
  else
    echo "✅ $API_FILE already has default export"
  fi
else
  echo "❌ $API_FILE not found"
  echo "   Buat file dengan isi:"
  echo "   import axios from 'axios';"
  echo "   const apiClient = axios.create({ ... });"
  echo "   export default apiClient;"
fi

echo ""
echo "=========================================="
echo "  SELESAI - Restart frontend dev server  "
echo "=========================================="
