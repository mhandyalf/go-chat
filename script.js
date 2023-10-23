// script.js
const chatBox = document.getElementById('chat');
const messageInput = document.getElementById('message');
const sendButton = document.getElementById('send');
const messageDisplay = document.getElementById('messages');

const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = (event) => {
    console.log("WebSocket connection opened");
};

ws.onmessage = (event) => {
    const message = event.data;
    messageDisplay.innerHTML += `<p>${message}</p>`;
};

ws.onclose = (event) => {
    console.log("WebSocket connection closed");
};

sendButton.addEventListener('click', () => {
    const message = messageInput.value;
    ws.send(message);
    messageInput.value = '';
});

messageInput.addEventListener('keyup', (event) => {
    if (event.key === "Enter") {
        sendButton.click();
    }
});
