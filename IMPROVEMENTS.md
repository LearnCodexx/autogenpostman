# Peningkatan Fleksibilitas Package

## ✅ Perbaikan yang Telah Dibuat

### 1. Auto-Detection Main File
Package sekarang bisa mencari `main.go` di berbagai lokasi umum:
- `main.go` (root)
- `cmd/main.go`  
- `cmd/api/main.go`
- `cmd/server/main.go`
- `cmd/app/main.go`
- `app/main.go`
- `api/main.go` 
- `server/main.go`

### 2. Fleksibilitas Output Path
- Default output: `docs/postman_collection.json` (lebih standard)
- Swagger output: `docs/` (lebih standard dari cmd/postman)
- Bisa dikustomisasi via flag: `--output=custom/path.json`

### 3. Auto-Detection Swagger Files  
Package mencari file OpenAPI/Swagger di berbagai lokasi:
- `docs/swagger.json`, `docs/openapi.yaml`
- `api/docs/swagger.json`, `api/docs/openapi.yaml`  
- `swagger/swagger.json`, `openapi/openapi.yaml`
- Dan berbagai kombinasi lainnya

### 4. Better Error Messages
Error sekarang memberikan saran path yang bisa digunakan:
```
auto swagger generation failed (...) and no OpenAPI file found; 
provide SwaggerInputPath or create one of: docs/swagger.json, docs/openapi.yaml, docs/openapi.yml
```

## 🚀 Cara Penggunaan untuk Berbagai Struktur Project

### Struktur Standard (Laravel-style):
```
project/
├── main.go              # ← Auto-detected
├── docs/
│   └── postman_collection.json  # ← Generated
```

### Struktur Go Standard:
```
project/
├── cmd/
│   └── api/
│       └── main.go      # ← Auto-detected
├── docs/
│   └── postman_collection.json  # ← Generated  
```

### Microservice Structure:
```
project/
├── cmd/
│   └── server/
│       └── main.go      # ← Auto-detected
├── api/
│   └── docs/
│       └── openapi.yaml # ← Auto-detected
```

### Custom Structure:
```bash
# Specify custom paths jika auto-detection gagal:
go run cmd/postman-generator/main.go \
  --swagger-input=custom/api.yaml \
  --output=build/postman.json
```

## 📝 Test Results

✅ Auto-detection main.go: **WORKING**  
✅ Flexible output paths: **WORKING**  
✅ Multiple swagger locations: **WORKING**  
✅ Better error messages: **WORKING**  
✅ Backward compatibility: **MAINTAINED**

Package sekarang cocok untuk digunakan di berbagai struktur project yang berbeda!