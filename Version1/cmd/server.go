package main

import (
	"ascii-art-web/internal/asciiart"
	"ascii-art-web/internal/hash"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const originalHash = "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf"

var templates *template.Template

func main() {
	var err error
	// Загружаем шаблоны при запуске сервера
	templates, err = template.ParseFiles("web/templates/index.html", "web/templates/404.html", "web/templates/500.html", "web/templates/400.html")
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов при запуске: %v", err)
	}

	// Проверяем хеш файла баннера
	fileName := "internal/asciiart/banners/standard.txt"
	currentHash := hash.ComputeFileHash(fileName)
	if currentHash != originalHash {
		log.Printf("Предупреждение: Файл %s был изменен. Ожидаемый хеш: %s, получен: %s", fileName, originalHash, currentHash)
	}

	// Настраиваем маршруты
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/generate-ascii", generateAsciiHandler)
	mux.HandleFunc("/400", badRequestPageHandler)
	mux.HandleFunc("/test-500", test500Handler) // Тестовый маршрут для ошибки 500

	// Обслуживание статических файлов
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обработчик для перехвата всех запросов и проверки допустимых путей
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			mux.ServeHTTP(w, r)
			return
		}

		allowedPaths := map[string]bool{
			"/":               true,
			"/generate-ascii": true,
			"/400":            true,
			"/test-500":       true, // Разрешаем тестовый маршрут
		}

		if !allowedPaths[r.URL.Path] {
			notFoundHandler(w, r)
			return
		}

		mux.ServeHTTP(w, r)
	})

	log.Println("Сервер запущен по адресу http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	fileName := "internal/asciiart/banners/standard.txt"
	currentHash := hash.ComputeFileHash(fileName)
	if currentHash != originalHash {
		log.Printf("Предупреждение: Несоответствие хеша файла в homeHandler. Ожидаемый: %s, получен: %s", originalHash, currentHash)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Не удалось отрендерить шаблон index: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		internalErrorHandler(w, r, "Ошибка рендеринга шаблона")
	}
}

func generateAsciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Не удалось распарсить форму: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}

	text := r.FormValue("text")
	bannerChoice := r.FormValue("banner")

	log.Printf("Получен запрос - Текст: %q, Баннер: %q, URL: %s", text, bannerChoice, r.URL.Path)

	if text == "" {
		badRequestHandler(w, r, "Текст не может быть пустым")
		return
	}

	if !isValidInput(text) {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}

	validBanners := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}

	if !validBanners[bannerChoice] {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}

	fileName := "internal/asciiart/banners/standard.txt"
	currentHash := hash.ComputeFileHash(fileName)
	if currentHash != originalHash {
		log.Printf("Предупреждение: Несоответствие хеша файла в generateAsciiHandler. Ожидаемый: %s, получен: %s", originalHash, currentHash)
	}

	output, err := asciiart.RenderText(text, bannerChoice)
	if err != nil {
		log.Printf("Не удалось сгенерировать ASCII-арт: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		if strings.Contains(err.Error(), "Load failed") {
			http.Redirect(w, r, "/400", http.StatusSeeOther)
		} else {
			internalErrorHandler(w, r, "Не удалось загрузить файл баннера")
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	if _, err := w.Write([]byte(output)); err != nil {
		log.Printf("Не удалось записать ответ: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		internalErrorHandler(w, r, "Ошибка записи ответа")
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	log.Printf("404 Не найдено: %s, Метод: %s", r.URL.Path, r.Method)
	if err := templates.ExecuteTemplate(w, "404.html", nil); err != nil {
		log.Printf("Не удалось отрендерить шаблон 404: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		internalErrorHandler(w, r, "Ошибка рендеринга шаблона 404")
	}
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request, errorMessage string) {
	// Проверяем, не обрабатываем ли мы ошибку повторно
	if w.Header().Get("X-Error-Handled") == "true" {
		log.Printf("Рекурсивная ошибка в internalErrorHandler: %s, URL: %s, Метод: %s", errorMessage, r.URL.Path, r.Method)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки и статус
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Error-Handled", "true")
	w.WriteHeader(http.StatusInternalServerError)

	// Логируем ошибку с деталями
	log.Printf("Внутренняя ошибка сервера: %s, URL: %s, Метод: %s, Параметры: %v", errorMessage, r.URL.Path, r.Method, r.Form)

	// Рендерим шаблон
	if err := templates.ExecuteTemplate(w, "500.html", nil); err != nil {
		log.Printf("Не удалось отрендерить шаблон 500: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func badRequestHandler(w http.ResponseWriter, r *http.Request, errorMessage string) {
	log.Printf("Неверный запрос: %s, URL: %s, Метод: %s", errorMessage, r.URL.Path, r.Method)
	http.Redirect(w, r, "/400", http.StatusSeeOther)
}

func badRequestPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	log.Printf("Отображена страница неверного запроса, URL: %s, Метод: %s", r.URL.Path, r.Method)
	if err := templates.ExecuteTemplate(w, "400.html", nil); err != nil {
		log.Printf("Не удалось отрендерить шаблон 400: %v, URL: %s, Метод: %s", err, r.URL.Path, r.Method)
		internalErrorHandler(w, r, "Ошибка рендеринга шаблона 400")
	}
}

func test500Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Тестовый маршрут 500 активирован, URL: %s, Метод: %s", r.URL.Path, r.Method)
	internalErrorHandler(w, r, "Тестовая ошибка 500")
}

func isValidInput(text string) bool {
	for _, char := range text {
		if (char < 32 || char > 126) && char != 10 && char != 13 {
			return false
		}
	}
	return true
}
