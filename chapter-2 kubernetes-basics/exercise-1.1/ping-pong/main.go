package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// ၁။ Package-level မှာ Variable တွေ ကြေညာခြင်း
var (
	pongs int
	mu    sync.Mutex
)
// ၂။ init() function ထဲမှာ ဖိုလ်ဒါဆောက်ခြင်းနဲ့ Path သတ်မှတ်ခြင်း
func main() {
	// မူလ Endpoint - Browser မှ ဝင်ကြည့်လျှင် Pong အရေအတွက် တိုးမည်
	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		pongs++
		currentPongs := pongs
		mu.Unlock()
		fmt.Fprintf(w, "pong %d", currentPongs)
	})

	// ✨ အသစ်ထည့်ရမည့် Endpoint - Log-output မှ အရေအတွက်ကိုသာ လှမ်းတောင်းမည်
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		currentPongs := pongs
		mu.Unlock()
		fmt.Fprintf(w, "%d", currentPongs)
	})

	log.Println("Ping-pong app listening on port 8080...")
	http.ListenAndServe(":8080", nil) // Port ကို ကိုယ်လိုသလို ပြင်နိုင်သည်
}