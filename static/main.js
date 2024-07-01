document.addEventListener('DOMContentLoaded', (event) => {
    let quill = new Quill('#editor-container', {
        theme: 'snow',
        modules: {
            toolbar: '#toolbar'
        }
    });

    document.getElementById('save-button').addEventListener('click', () => {
        const content = quill.root.innerHTML;
        saveContent(content);
    });

    document.getElementById('open-button').addEventListener('click', () => {
        htmx.ajax('GET', '/list-notes', {
            target: '#editor-container',
            swap: 'innerHTML',
            onBeforeSwap: (detail) => {
                // Assuming the response is a list of notes with titles and IDs
                quill.root.innerHTML = detail.xhr.responseText;
                return false;
            }
        });
    });

    document.getElementById('new-button').addEventListener('click', () => {
        quill.root.innerHTML = '';
    });

    // Load content from the server on page load
    /*htmx.ajax('GET', '/load', {
        target: '#editor-container',
        swap: 'innerHTML',
        onBeforeSwap: (detail) => {
            quill.root.innerHTML = detail.xhr.responseText;
            return false;
        }
    });*/

    // Save content to the server on editor change
    /*quill.on('text-change', () => {
        const content = quill.root.innerHTML;
        saveContent(content);
    });*/

    function saveContent(content) {
        htmx.ajax('POST', '/save', {
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            values: { text: content }
        });

        // Save content to IndexedDB (WebAssembly)
        saveToIndexedDB(content);
    }

    // Initialize WebAssembly
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject).then((result) => {
        go.run(result.instance);
    });
});

// Function to save content to IndexedDB via WebAssembly
function saveToIndexedDB(content) {
    wasmSaveToIndexedDB(content);
}
