let currentChapterIndex = 0;
let readingOrder = [];

function initializeReader(filename, chapters) {
    readingOrder = chapters;
    
    // Load the first chapter
    if (readingOrder.length > 0) {
        loadChapter(0);
    }
    
    // Set up event listeners
    document.getElementById('prev-btn').addEventListener('click', () => {
        if (currentChapterIndex > 0) {
            loadChapter(currentChapterIndex - 1);
        }
    });
    
    document.getElementById('next-btn').addEventListener('click', () => {
        if (currentChapterIndex < readingOrder.length - 1) {
            loadChapter(currentChapterIndex + 1);
        }
    });
    
    document.getElementById('chapter-select').addEventListener('change', (e) => {
        const href = e.target.value;
        const index = readingOrder.findIndex(chapter => chapter.href === href);
        if (index !== -1) {
            loadChapter(index);
        }
    });
}

function loadChapter(index) {
    currentChapterIndex = index;
    const chapter = readingOrder[index];
    
    // Update the select dropdown
    const select = document.getElementById('chapter-select');
    select.value = chapter.href;
    
    // Load the content
    const contentFrame = document.getElementById('content-frame');
    const contentUrl = `/epub-content/${filename}${chapter.href}`;
    
    if (chapter.href.endsWith('.html') || chapter.href.endsWith('.xhtml')) {
        // For HTML content, use an iframe
        contentFrame.innerHTML = `<iframe src="${contentUrl}"></iframe>`;
    } else {
        // For other content types (images, etc.), display directly
        fetch(contentUrl)
            .then(response => {
                if (response.headers.get('Content-Type').includes('image')) {
                    return response.blob().then(blob => {
                        const url = URL.createObjectURL(blob);
                        contentFrame.innerHTML = `<img src="${url}" alt="${chapter.title}">`;
                    });
                } else {
                    return response.text().then(text => {
                        contentFrame.innerHTML = text;
                    });
                }
            })
            .catch(error => {
                contentFrame.innerHTML = `<p>Error loading content: ${error.message}</p>`;
            });
    }
    
    // Update browser history
    history.pushState({ chapterIndex: index }, '', `/read/${filename}?chapter=${index}`);
}

// Handle browser back/forward
window.addEventListener('popstate', (event) => {
    if (event.state && event.state.chapterIndex !== undefined) {
        loadChapter(event.state.chapterIndex);
    }
});