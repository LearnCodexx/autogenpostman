# 📦 Example: Working AutoGenPostman Demo

Contoh lengkap penggunaan autogenpostman package.

## 🚀 Quick Test

```bash
# 1. Masuk ke folder example
cd example/

# 2. Install dependencies  
go mod tidy
go install github.com/swaggo/swag/cmd/swag@latest

# 3. Setup postman generator
go run cmd/setup-postman/main.go

# 4. Generate postman collection
go run cmd/postman-generator/main.go

# 5. Check result
ls docs/
# Output: docs.go  postman_collection.json  swagger.json  swagger.yaml
```

## 📄 Files Generated

- `docs/postman_collection.json` - Postman collection yang bisa diimport 
- `docs/swagger.json` - Swagger spec yang di-generate
- `docs/docs.go` - Go docs untuk swagger UI

## 🧪 Test API

```bash
# 1. Start server
go run main.go

# 2. Test endpoints
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/users

# 3. View swagger UI
open http://localhost:8080/swagger/index.html
```

## 📝 What This Example Shows

✅ **Complete API** dengan CRUD operations  
✅ **Swagger annotations** yang proper  
✅ **Auto-generation** swagger → postman  
✅ **Real structure** yang bisa dipakai di production  
✅ **Error handling** dan response types  

## 🔄 Modify Example

Edit `main.go` untuk:
- Tambah endpoints baru
- Ubah response structure  
- Add authentication
- Modify swagger info

Then re-run:
```bash
go run cmd/postman-generator/main.go
```

Postman collection akan terupdate otomatis! 🎉