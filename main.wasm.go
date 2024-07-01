package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("wasmSaveToIndexedDB", js.FuncOf(saveToIndexedDB))
	select {}
}

func saveToIndexedDB(this js.Value, p []js.Value) interface{} {
	content := p[0].String()
	// JavaScript code to save to IndexedDB
	js.Global().Call("eval", `
		let request = indexedDB.open("EditorDB", 1);
		request.onupgradeneeded = function(event) {
			let db = event.target.result;
			db.createObjectStore("content", { keyPath: "id", autoIncrement: true });
		};
		request.onsuccess = function(event) {
			let db = event.target.result;
			let tx = db.transaction("content", "readwrite");
			let store = tx.objectStore("content");
			store.put({id: 1, text: '`+content+`'});
		};
	`)
	return nil
}
