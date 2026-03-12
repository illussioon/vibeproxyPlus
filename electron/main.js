const { app, BrowserWindow, shell, ipcMain, dialog } = require('electron');
const path = require('path');
const fs = require('fs');
const net = require('net');
const { spawn } = require('child_process');

let backendProcess;
let backendManagedByApp = false;

function createWindow() {
  const win = new BrowserWindow({
    width: 860,
    height: 1050,
    resizable: false,
    title: 'VibeProxy',
    backgroundColor: '#0b1119',
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false
    }
  });

  win.removeMenu();
  win.loadFile(path.join(__dirname, '..', 'renderer', 'index.html'));
}

function isBackendListening(host = '127.0.0.1', port = 7890) {
  return new Promise((resolve) => {
    const socket = net.createConnection({ host, port }, () => {
      socket.end();
      resolve(true);
    });

    socket.setTimeout(500);
    socket.on('timeout', () => {
      socket.destroy();
      resolve(false);
    });
    socket.on('error', () => resolve(false));
  });
}

function hookBackendProcess(proc) {
  backendProcess = proc;
  backendManagedByApp = true;

  backendProcess.stdout?.on('data', (chunk) => {
    process.stdout.write(`[backend] ${chunk}`);
  });

  backendProcess.stderr?.on('data', (chunk) => {
    process.stderr.write(`[backend] ${chunk}`);
  });

  backendProcess.on('error', (err) => {
    console.error(`Backend process error: ${err.message}`);
    dialog.showErrorBox(
      'Backend startup error',
      `Failed to start backend process.\n\n${err.message}`
    );
  });

  backendProcess.on('exit', (code) => {
    console.log(`Backend exited with code ${code}`);
    backendProcess = undefined;
    backendManagedByApp = false;
  });
}

function startBackendFromBinary() {
  const candidatePaths = [
    path.join(__dirname, '..', 'backend', 'bin', 'vibeproxyplus-backend.exe'),
    path.join(process.resourcesPath || '', 'backend', 'bin', 'vibeproxyplus-backend.exe')
  ];

  const exePath = candidatePaths.find((p) => p && fs.existsSync(p));
  if (!exePath) return false;

  hookBackendProcess(
    spawn(exePath, [], {
      windowsHide: true,
      stdio: 'pipe'
    })
  );

  return true;
}

function startBackendFromGoRun() {
  const backendDir = path.join(__dirname, '..', 'backend');
  const useShell = process.platform === 'win32';

  hookBackendProcess(
    spawn('go', ['run', './cmd/server'], {
      cwd: backendDir,
      windowsHide: true,
      stdio: 'pipe',
      shell: useShell
    })
  );
}

async function ensureBackend() {
  if (await isBackendListening()) {
    console.log('Backend already running on 127.0.0.1:7890, skip auto-start.');
    return;
  }

  const startedFromBinary = startBackendFromBinary();
  if (startedFromBinary) {
    console.log('Backend started from prebuilt binary.');
    return;
  }

  console.log('Backend binary not found, fallback to `go run ./cmd/server`.');
  startBackendFromGoRun();
}

app.whenReady().then(async () => {
  await ensureBackend();
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit();
});

app.on('before-quit', () => {
  if (backendManagedByApp && backendProcess && !backendProcess.killed) {
    backendProcess.kill();
  }
});

ipcMain.handle('open-external', async (_, url) => {
  await shell.openExternal(url);
});
