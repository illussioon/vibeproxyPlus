const API = 'http://127.0.0.1:7890';

const services = ['Antigravity', 'Claude Code', 'Codex', 'Gemini', 'GitHub Copilot', 'Qwen'];

const serviceList = document.getElementById('service-list');
const statusEl = document.getElementById('server-status');
const launchToggle = document.getElementById('launch-toggle');

for (const name of services) {
  const row = document.createElement('div');
  row.className = 'service';
  row.innerHTML = `
    <div class="service-left">
      <div class="service-title">${name}</div>
      <div class="service-sub">No connected accounts</div>
    </div>
    <button class="btn">Add Account</button>
  `;
  serviceList.appendChild(row);
}

async function request(path, init = {}) {
  const res = await fetch(`${API}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...init
  });

  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

function renderStatus(state) {
  statusEl.textContent = state.running ? '● Running' : '● Stopped';
  statusEl.classList.toggle('running', state.running);
  statusEl.classList.toggle('stopped', !state.running);
  launchToggle.checked = Boolean(state.launchAtLogin);
}

async function refresh() {
  try {
    const state = await request('/api/state');
    renderStatus(state);
  } catch (err) {
    console.error(err);
  }
}

document.getElementById('start-btn').addEventListener('click', async () => {
  await request('/api/server/start', { method: 'POST' });
  await refresh();
});

document.getElementById('stop-btn').addEventListener('click', async () => {
  await request('/api/server/stop', { method: 'POST' });
  await refresh();
});

launchToggle.addEventListener('change', async (event) => {
  await request('/api/settings/launch-at-login', {
    method: 'POST',
    body: JSON.stringify({ enabled: event.target.checked })
  });
  await refresh();
});

document.getElementById('open-auth').addEventListener('click', async () => {
  await request('/api/auth/open', { method: 'POST' });
});

document.getElementById('cli-link').addEventListener('click', (e) => {
  e.preventDefault();
  window.desktop.openExternal('https://github.com/router-for-me/CLIProxyAPIPlus');
});

document.getElementById('issue-link').addEventListener('click', (e) => {
  e.preventDefault();
  window.desktop.openExternal('https://github.com/automazeio/vibeproxy/issues');
});

setInterval(refresh, 3000);
refresh();
