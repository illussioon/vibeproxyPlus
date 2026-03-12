# VibeProxyPlus (Windows-first fork)

Windows-only desktop fork inspired by VibeProxy UI/UX.

## Stack
- **Desktop shell:** Electron + vanilla JS renderer
- **Local API:** Go (gin-free stdlib HTTP server)
- **Proxy engine adapter:** prepared wrapper for CLIProxyAPIPlus-compatible process

## Features (MVP)
- UI matching original VibeProxy layout (dark glassy panel + services list)
- Start/stop local proxy service
- Launch at login toggle (stored locally)
- Open auth folder action
- Service account list placeholders for Antigravity, Claude Code, Codex, Gemini, GitHub Copilot, Qwen

## Quick start

### 1) Install dependencies
```bash
npm install
```

### 2) Run Go backend
```bash
npm run backend:dev
```

### 3) Run Electron app
```bash
npm run dev
```

## Windows build
```bash
npm run build:win
```

## Notes
- This repo is prepared as **Windows-first**. Linux/macOS support is not a priority.
- Adapter for external proxy binary can be connected in `backend/internal/proxy/manager.go`.
