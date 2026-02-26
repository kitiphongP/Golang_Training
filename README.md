# Golang_Training

## API Endpoints

Documentation:

- GitHub language flow note: `docs/github-language-flow.md`

- `GET /health`
- `GET /users`
- `POST /users`
- `GET /users?search=<keyword>`
- `GET /github/skills?username=<github-username>`
- `GET /github/skills?email=<github-public-email>`
- `GET /github/languages?username=<github-username>`
- `GET /github/languages?email=<github-public-email>`

`/github/languages` จะคืนข้อมูลโปรไฟล์และรีโพทั้งหมด เช่น `username`, `name`, `email`, `repos[]` (ชื่อโปรเจกต์/ภาษา) และ `languages` (สรุปจำนวนภาษา)
