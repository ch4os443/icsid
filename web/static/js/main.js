// Configuração do WebSocket
const ws = new WebSocket(`wss://${window.location.host}/ws`);

// Elementos da interface
const sessionsTable = document.getElementById('sessions-table');
const commandInput = document.getElementById('command-input');
const commandOutput = document.getElementById('command-output');
const systemInfo = document.getElementById('system-info');
const alertsContainer = document.getElementById('alerts-container');

// Estado da aplicação
let currentSession = null;
let sessions = new Map();

// Funções auxiliares
function showAlert(message, type = 'info') {
    const alert = document.createElement('div');
    alert.className = `alert alert-${type}`;
    alert.textContent = message;
    alertsContainer.appendChild(alert);
    setTimeout(() => alert.remove(), 5000);
}

function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function formatDate(date) {
    return new Date(date).toLocaleString();
}

// Manipuladores de eventos WebSocket
ws.onopen = () => {
    showAlert('Conectado ao servidor', 'success');
    ws.send(JSON.stringify({ type: 'get_sessions' }));
};

ws.onclose = () => {
    showAlert('Desconectado do servidor', 'danger');
};

ws.onerror = (error) => {
    showAlert('Erro na conexão WebSocket', 'danger');
    console.error('WebSocket error:', error);
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    
    switch (data.type) {
        case 'sessions_update':
            updateSessionsTable(data.sessions);
            break;
        case 'command_output':
            appendCommandOutput(data.output);
            break;
        case 'system_info':
            updateSystemInfo(data.info);
            break;
        case 'error':
            showAlert(data.message, 'danger');
            break;
    }
};

// Atualização da interface
function updateSessionsTable(sessionsData) {
    sessions.clear();
    sessionsTable.innerHTML = `
        <thead>
            <tr>
                <th>ID</th>
                <th>Usuário</th>
                <th>IP</th>
                <th>Conectado</th>
                <th>Última Atividade</th>
                <th>Ações</th>
            </tr>
        </thead>
        <tbody>
        </tbody>
    `;
    
    sessionsData.forEach(session => {
        sessions.set(session.id, session);
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${session.id}</td>
            <td>${session.username}</td>
            <td>${session.ip}</td>
            <td>${formatDate(session.connected_at)}</td>
            <td>${formatDate(session.last_activity)}</td>
            <td>
                <button class="btn btn-primary btn-sm" onclick="selectSession('${session.id}')">
                    Selecionar
                </button>
                <button class="btn btn-danger btn-sm" onclick="terminateSession('${session.id}')">
                    Terminar
                </button>
            </td>
        `;
        sessionsTable.querySelector('tbody').appendChild(row);
    });
}

function appendCommandOutput(output) {
    const pre = document.createElement('pre');
    pre.textContent = output;
    commandOutput.appendChild(pre);
    commandOutput.scrollTop = commandOutput.scrollHeight;
}

function updateSystemInfo(info) {
    systemInfo.innerHTML = `
        <div class="row">
            <div class="col-md-6">
                <h5>CPU</h5>
                <p>Uso: ${info.cpu.usage}%</p>
                <p>Núcleos: ${info.cpu.cores}</p>
            </div>
            <div class="col-md-6">
                <h5>Memória</h5>
                <p>Total: ${formatBytes(info.memory.total)}</p>
                <p>Usada: ${formatBytes(info.memory.used)}</p>
                <p>Livre: ${formatBytes(info.memory.free)}</p>
            </div>
        </div>
        <div class="row mt-3">
            <div class="col-md-6">
                <h5>Disco</h5>
                <p>Total: ${formatBytes(info.disk.total)}</p>
                <p>Usado: ${formatBytes(info.disk.used)}</p>
                <p>Livre: ${formatBytes(info.disk.free)}</p>
            </div>
            <div class="col-md-6">
                <h5>Rede</h5>
                <p>Enviado: ${formatBytes(info.network.sent)}</p>
                <p>Recebido: ${formatBytes(info.network.received)}</p>
            </div>
        </div>
    `;
}

// Manipuladores de eventos da interface
function selectSession(sessionId) {
    currentSession = sessionId;
    showAlert(`Sessão ${sessionId} selecionada`, 'info');
    commandOutput.innerHTML = '';
}

function terminateSession(sessionId) {
    if (confirm('Tem certeza que deseja terminar esta sessão?')) {
        ws.send(JSON.stringify({
            type: 'terminate_session',
            session_id: sessionId
        }));
    }
}

document.getElementById('command-form').addEventListener('submit', (e) => {
    e.preventDefault();
    
    if (!currentSession) {
        showAlert('Selecione uma sessão primeiro', 'warning');
        return;
    }
    
    const command = commandInput.value.trim();
    if (!command) {
        showAlert('Digite um comando', 'warning');
        return;
    }
    
    ws.send(JSON.stringify({
        type: 'execute_command',
        session_id: currentSession,
        command: command
    }));
    
    commandInput.value = '';
});

// Atualização periódica do sistema
setInterval(() => {
    ws.send(JSON.stringify({ type: 'get_system_info' }));
}, 5000);

// Inicialização
document.addEventListener('DOMContentLoaded', () => {
    showAlert('Inicializando interface...', 'info');
}); 