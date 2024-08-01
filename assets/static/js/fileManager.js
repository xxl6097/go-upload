// fileManager.js

document.addEventListener('DOMContentLoaded', function() {
    fetch('/files')
        .then(response => response.json())
        .then(files => {
            const fileList = document.getElementById('file-list');
            files.forEach(file => {
                const fileLink = document.createElement('a');
                fileLink.href = `#`;
                fileLink.textContent = file.name;
                fileLink.onclick = function() {
                    loadFile(file.name);
                    return false;
                };
                fileList.appendChild(fileLink);
                fileList.appendChild(document.createElement('br'));
            });
        });
});

function loadFile(fileName) {
    fetch(`/files/${fileName}`)
        .then(response => response.text())
        .then(content => {
            document.getElementById('file-content').value = content;
        });
}

function saveFile() {
    const currentFile = document.getElementById('file-content').value;
    const fileName = document.getElementById('file-list').innerText.trim();
    fetch(`/files/${fileName}`, {
        method: 'PUT',
        body: currentFile,
        headers: {
            'Content-Type': 'text/plain'
        }
    }).then(response => {
        if (response.ok) {
            alert('File saved successfully!');
        } else {
            alert('Failed to save file.');
        }
    });
}