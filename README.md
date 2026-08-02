# 🔐 EnvGuard Native Go CLI (`envg`)

The official standalone Go-native command line tool for **[EnvGuard](https://getenvguard.com)** — Zero dependencies, standalone binary execution for modern DevOps & CI/CD pipelines.

---

## 🌟 Key Features

* **Zero Runtime Dependencies**: No Node.js, Python, or npm required. Standalone binary executable.
* **Direct 1-Line Kubernetes Sync (`envg sync`)**: Fetch & apply secrets directly into Kubernetes clusters without temporary disk files or ANSI control codes.
* **In-Memory Process Injection (`envg run`)**: Zero-disk footprint secret injection directly into process environment variables.
* **Multi-Platform Support**: Binaries for Linux (amd64/arm64), macOS (Apple Silicon & Intel), and Windows.

---

## ⚡ Installation

Download the binary for your operating system:

```bash
# macOS (Apple Silicon M1/M2/M3/M4)
curl -fsSL https://github.com/getenvguard/envguard-cli/releases/latest/download/envg-darwin-arm64 -o envg
chmod +x envg && sudo mv envg /usr/local/bin/

# Linux (amd64 / Jenkins / CI Runners)
curl -fsSL https://github.com/getenvguard/envguard-cli/releases/latest/download/envg-linux-amd64 -o envg
chmod +x envg && sudo mv envg /usr/local/bin/
```

---

## 🚀 Commands & Usage

### 1. Login

```bash
envg login --token "eg_pat_YOUR_PAT_TOKEN"
```

### 2. Direct 1-Line Kubernetes Sync

```bash
envg sync -p "payment-api" -e "production" -n "production"
```

### 3. In-Memory Process Execution

```bash
envg run -p "payment-api" -e "production" -- ./my-server
```

### 4. Pull to Local File (.env, k8s, json)

```bash
envg pull -p "payment-api" -e "production" --format k8s -o secret.yaml
```

---

## 📄 License

[MIT License](LICENSE) &copy; 2026 EnvGuard Team.
