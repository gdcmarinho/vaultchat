const socket = new WebSocket('ws://localhost:8080/session');

let publicKeyJWK;

socket.onopen = async function(event) {
    console.log('Conexão WebSocket estabelecida.');

    const { pemPublicKey, publicKey } = await generateRSAKey();
    publicKeyJWK = publicKey;

    const publicKeyMessage = {
        type: 'publicKey',
        publicKey: pemPublicKey
    };
    socket.send(JSON.stringify(publicKeyMessage));
};

socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('Mensagem recebida:', message);
};

socket.onerror = function(error) {
    console.error('Erro na conexão WebSocket:', error);
};

socket.onclose = function(event) {
    console.log('Conexão WebSocket fechada.');
};

document.getElementById('sendMessageButton').addEventListener('click', async function() {
    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value;

    if (!message) return;

    const encryptedMessage = await encryptMessage(message, publicKeyJWK);

    if (!encryptedMessage) {
        console.error('FATAL ERROR! Message not encrypted.');
        return;
    }

    if (socket.readyState === WebSocket.OPEN) {
        const messageData = {
            type: 'message',
            content: encryptedMessage
        };
        socket.send(JSON.stringify(messageData));
    } else {
        console.error('WebSocket is not open.');
    }

    messageInput.value = '';
});
