package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)
const (
	imageDir  = "/usr/src/app/files"
	imagePath = "/usr/src/app/files/image.jpg"
	imageURL  = "https://picsum.photos/1200"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func getBackendURL() string {
	if url := os.Getenv("BACKEND_URL"); url != "" {
		return url
	}
	return "http://todo-backend-svc:5000"
}

func fetchAndCacheImage() error {
	_ = os.MkdirAll(imageDir, 0755)

	info, err := os.Stat(imagePath)

	// ၁။ ဖိုင်ရှိပြီးသားဖြစ်ပြီး ၁၀ မိနစ် မပြည့်သေးပါက API ထပ်မခေါ်ဘဲ ရှိပြီးသားကို သုံးမည်
	if err == nil {
		if time.Since(info.ModTime()) < 10*time.Minute {
			return nil
		}
	}

	// ၂။ ဖိုင်မရှိသေးပါက (သို့) ၁၀ မိနစ်ကျော်သွားပါက Picsum မှ ပုံအသစ် ဒေါင်းလုဒ်ဆွဲမည်
	log.Println("Fetching a new image from Lorem Picsum...")
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	// ၃။ ပုံကို PV Mount ထားသည့် File Path ထဲသို့ သွားသိမ်းမည်
	out, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	err := fetchAndCacheImage()
	if err != nil {
		http.Error(w, "Could not load image", http.StatusInternalServerError)
		return
	}
	// Image File ကို Client ဆီ ပြန်ပို့ပေးခြင်း
	http.ServeFile(w, r, imagePath)
}

const indexTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Todo App</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; max-width: 600px; }
        img { width: 100%; height: auto; border-radius: 8px; margin-bottom: 20px; }
        input[type="text"] { width: 70%; padding: 8px; font-size: 14px; }
        button { padding: 8px 12px; font-size: 14px; cursor: pointer; }
        ul { margin-top: 20px; line-height: 1.6; }
    </style>
</head>
<body>
    <h1>Todo App</h1>
    <img src="/image.jpg" alt="Hourly Random Image" style="max-width: 600px; height: auto;" />
    
    <form action="/" method="POST">
        <input type="text" name="todo" maxlength="140" placeholder="Enter a new todo (max 140 characters)..." required />
        <button type="submit">Send</button>
    </form>

    <ul>
        {{range .}}
            <li>{{.Text}}</li>
        {{else}}
            <li>No todos available</li>
        {{end}}
    </ul>
</body>
</html>
`;

// Main HTML Page Handler (`/`)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	backendURL := getBackendURL()

	// ၁။ Form Submit (POST Method) ပြုလုပ်လိုက်ပါက Backend ဆီ Todo အသစ် လှမ်းပို့မည်
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		todoText := r.FormValue("todo")
		if todoText != "" {
			newTodo := Todo{Text: todoText}
			jsonData, _ := json.Marshal(newTodo)

			// http://todo-backend-svc:5000/todos ဆီ POST Request ပို့ခြင်း
			resp, err := http.Post(backendURL+"/todos", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Error posting todo to backend: %v", err)
			} else {
				resp.Body.Close()
			}
		}

		// Form Resubmission မဖြစ်စေရန် GET "/" သို့ Redirect ပြန်လုပ်ခြင်း
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// ၂။ GET Request ဖြစ်ပါက Backend မှ Todos စာရင်းကို Fetch လုပ်ယူမည်
	resp, err := http.Get(backendURL + "/todos")
	var todos []Todo

	if err == nil {
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&todos)
	} else {
		log.Printf("Error fetching todos from backend: %v", err)
	}

	// ၃။ Template ထဲ Todos ထည့်သွင်း၍ HTML Render လုပ်ခြင်း
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, todos)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Todo App v1.5")
	// })

	http.HandleFunc("/image.jpg", imageHandler)
	http.HandleFunc("/", indexHandler)

	// Exercise 1.2 လိုအပ်ချက်: "Server started in port NNNN"
	log.Printf("Server started in port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}