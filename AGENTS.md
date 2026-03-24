# Autochitect Agentic Era: Distributed Sandbox Infrastructure (PoC)

## 1. Vision
To create a decentralized execution layer for AI agents (Open-SWE). By utilizing community PCs as "Landlords," we can run heavy agentic coding tasks in a distributed network rather than centralized, expensive cloud clusters.

---

## 2. Architecture Overview (C1/C2)

### System Context (C1)
The system connects **Developers** (who need code written) to **Community Landlords** (who provide compute) via the **Autochitect Control Plane**.

### Container Diagram (C2)
* **Autochitect Server (Node.js):** The central switchboard. Manages the job queue, tracks landlord health, and streams logs to the frontend.
* **Landlord Daemon (Go):** The "Glue" installed on local machines. It talks to the server via WebSockets and manages local Docker/OpenShell sandboxes.
* **Execution Sandbox:** An isolated Docker container running **Open-SWE**.

---

## 3. The Landlord Daemon (Go Client)

The Go daemon is a lightweight, single-binary "middleman."

### Core Modules
| Module | Responsibility |
| :--- | :--- |
| **The Dialer** | Maintains a persistent WebSocket connection to `autochitect.com`. Handles auto-reconnect. |
| **Task Processor** | Receives `EXEC_TASK` JSON, pulls the required Docker image, and starts execution. |
| **Streamer** | Pipes `stdout/stderr` from the local container back to the server in real-time. |
| **Garbage Collector (GC)** | Enforces a Time-To-Live (TTL). Kills and removes containers after task completion or timeout. |

### Simplified State Machine
1.  **IDLE:** Authenticated and waiting for a task.
2.  **BUSY:** Executing an agent task; streaming logs.
3.  **CLEANUP:** Task finished or failed; purging local resources.
4.  **OFFLINE:** Socket disconnected; Server marks as inactive.

---

## 4. The Control Plane (Node.js Server)

Node.js is used for its high-speed I/O and vast ecosystem of WebSocket libraries (Socket.io).

### Responsibilities
* **Landlord Registry:** A memory-resident map of all connected Landlords and their hardware specs.
* **Agent Orchestration:** Receives coding requirements from the UI and assigns them to an `IDLE` Landlord.
* **Log Proxy:** Receives log chunks from the Go Daemon and broadcasts them to the user's browser.
* **Heartbeat Monitor:** If a Landlord misses 3 heartbeats, the server triggers a "Fail" state for the task and re-assigns it.

---

## 5. Technical Stack

| Component | Technology | Why? |
| :--- | :--- | :--- |
| **Server** | Node.js (TypeScript) | Fast development, excellent WebSocket support. |
| **Client (Daemon)** | Go (Golang) | Single binary distribution; native concurrency for handling logs + GC. |
| **Communication** | WebSockets + JWT | Real-time, bi-directional, and secure. |
| **Sandbox** | Docker / OpenShell | Hardware abstraction (CPU/GPU) and process isolation. |
| **Security** | mTLS / PSK | Ensures only verified Landlords join the network. |

---

## 6. Communication Protocol (PoC JSON)

### Task Request (Server -> Client)
```json
{
  "type": "TASK_ASSIGN",
  "payload": {
    "task_id": "99b-001",
    "image": "autochitect/open-swe:latest",
    "env_vars": { "REPO_URL": "[github.com/user/repo](https://github.com/user/repo)" },
    "timeout_sec": 600
  }
}
```
