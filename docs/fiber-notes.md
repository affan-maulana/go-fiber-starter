# Catatan Fiber

## Beberapa method penting di *fiber.Ctx:

Request/Response:

* c.BodyParser(&obj) — Parse JSON body ke struct
* c.JSON(data) — Kirim response JSON
* c.SendString("text") — Kirim response string
* c.Status(code) — Set HTTP status code
* c.SendStatus(code) — Kirim status tanpa body

Header & Query:

* c.Get("Header-Name") — Ambil header
* c.Set("Header-Name", "value") — Set header
* c.Query("param") — Ambil query string (?param=...)

Path & Params:

* c.Path() — Ambil path URL
* c.Params("id") — Ambil path param (misal /users/:id)

Cookies:

* c.Cookies("name") — Ambil cookie
* c.Cookie(&fiber.Cookie{...}) — Set cookie

Locals (context per-request):

* c.Locals("key") — Ambil data dari context
* c.Locals("key", value) — Set data ke context

Auth:

* c.IP() — Ambil IP client
* c.Protocol() — http/https

Lainnya:

* c.Method() — GET, POST, dll
* c.OriginalURL() — URL lengkap
* c.FormValue("field") — Ambil form value

## Go routine 
