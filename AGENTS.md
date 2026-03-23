# AGENTS.md

## 🎯 Project Goal: The Open-SWE Backend
We are building a high-performance **Backend Layer** for Open-SWE. Our goal is to serve as a decentralized alternative to Daytona, managing the orchestration between AI agents and **OpenShell Landlords**.

---

## 🛠 Tech Stack
* **Runtime:** [Bun.js](https://bun.sh) (TypeScript).
* **Execution:** **NVIDIA OpenShell** (Hardware-agnostic, kernel-isolated sandboxes).
* **Architecture:** Modular Backend. All external systems (OpenShell, GitHub) must be decoupled via Interfaces.

---

## 🧪 Verification (TDD)
We use a **Test-Driven Development** approach. No feature is complete without verification.

### 1. Requirements
* **Unit Tests:** Mock all Sandbox and Landlord interactions.
* **Integration Tests:** Verify the full cycle from Task Trigger → OpenShell Execution → Result.
* **Security:** Every input to the shell must be validated to protect the Landlord's host.

### 2. Commands
* **Test:** `bun test`
* **Lint:** `bun x biome check .`
* **Audit:** `openshell policy-check ./configs/policy.yaml`

---

## 🏗 Modular Rules
1. **Abstraction over Implementation:** Business logic interacts only with the `SandboxProvider` interface, never directly with the shell.
2. **Stateless Logic:** Treat every OpenShell instance as ephemeral. 
3. **Async First:** Use Bun’s native HTTP and WebSocket capabilities to handle long-running SWE tasks.

---

## 🤖 Instructions for the Agent
1. **TDD First:** Write the test before the implementation.
2. **Respect the Policy:** If a task requires new permissions, update the OpenShell `policy.yaml`.
3. **Environment Aware:** Detect if the Landlord provides NVIDIA hardware to optimize for local inference.
