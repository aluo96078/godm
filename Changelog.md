## [v0.0.1] - 2025-03-27 ~ 2025-03-28

### 新增
- ✨ 引入完整 Observer 系統（靈感來自 Laravel）：
  - 定義 `ModelObserver` 介面，包含 `creating`、`created`、`updating`、`updated`、`deleting`、`deleted` 等模型生命週期事件。
  - 支援每個模型多個 observer，以及全域 observer 註冊。
  - 支援 Observer 過濾功能：
    - `EventFilter`：可過濾感興趣的事件類型。
    - `TypedObserver`：可過濾特定模型類型。
    - `PrioritizedObserver`：可設定觸發順序。
  - 提供全域錯誤攔截器 `RegisterObserverErrorHandler`。
- 🧪 範例程式已更新，展示實際 observer 註冊與觸發流程。
- 📄 README 與程式碼已加入中英雙語註解，提升可讀性。

### Added
- ✨ Introduced full Observer system inspired by Laravel:
  - `ModelObserver` interface with `creating`, `created`, `updating`, `updated`, `deleting`, `deleted` lifecycle hooks.
  - Supports multiple observers per model and global observers.
  - Added support for observer filtering:
    - `EventFilter`: filter by event type.
    - `TypedObserver`: filter by model type.
    - `PrioritizedObserver`: control execution order.
  - Built-in observer error interception via `RegisterObserverErrorHandler`.
- 🧪 Example updated with real Observer usage and global registration.
- 📄 README and code now include full bilingual (中文/English) annotations.