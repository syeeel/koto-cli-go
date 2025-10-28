# koto - ToDo管理CLIツール 詳細設計書

## 1. ディレクトリ構成

```
koto/
├── cmd/
│   └── koto/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── model/
│   │   └── todo.go              # Todoデータモデル
│   ├── repository/
│   │   ├── repository.go        # Repositoryインターフェース
│   │   └── sqlite.go            # SQLite実装
│   ├── service/
│   │   └── todo_service.go      # ビジネスロジック
│   ├── tui/
│   │   ├── model.go             # Bubbletea Model
│   │   ├── commands.go          # コマンドパーサー
│   │   ├── views.go             # ビューレンダリング
│   │   ├── styles.go            # Lipglossスタイル
│   │   └── update.go            # Update関数
│   └── config/
│       └── config.go            # 設定管理
├── pkg/
│   └── utils/
│       ├── time.go              # 時刻関連ユーティリティ
│       └── validation.go        # バリデーション
├── migrations/
│   └── 001_init.sql             # DBマイグレーション
├── docs/
│   ├── basic_design.md          # 基本設計書
│   └── detailed_design.md       # 詳細設計書（本書）
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── Makefile
```

## 2. パッケージ構成

### 2.1 cmd/koto
- アプリケーションのエントリーポイント
- 設定の初期化
- TUIアプリケーションの起動

### 2.2 internal/model
- データモデルの定義
- ドメインロジックを持たないPlain Old Go Object (POGO)

### 2.3 internal/repository
- データ永続化の抽象化レイヤー
- SQLite実装
- 将来的な他のストレージへの対応を容易にする

### 2.4 internal/service
- ビジネスロジック層
- Repositoryを使用してデータ操作
- バリデーション実行

### 2.5 internal/tui
- Bubbletea によるTUI実装
- Model-View-Update パターン
- コマンドパーサー
- スタイリング

### 2.6 internal/config
- アプリケーション設定
- データベースパスの管理
- 設定ファイルの読み込み（オプション）

### 2.7 pkg/utils
- 汎用ユーティリティ関数
- プロジェクト外からも利用可能な関数

## 3. データモデル詳細

### 3.1 Todo構造体

```go
package model

import "time"

type TodoStatus int

const (
    StatusPending TodoStatus = iota
    StatusCompleted
)

type Priority int

const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
)

type Todo struct {
    ID           int64      `db:"id"`
    Title        string     `db:"title"`
    Description  string     `db:"description"`
    Status       TodoStatus `db:"status"`
    Priority     Priority   `db:"priority"`
    DueDate      *time.Time `db:"due_date"`
    WorkDuration int        `db:"work_duration"` // 累積作業時間（分）
    CreatedAt    time.Time  `db:"created_at"`
    UpdatedAt    time.Time  `db:"updated_at"`
}

func (t Todo) IsCompleted() bool {
    return t.Status == StatusCompleted
}

func (t Todo) IsPending() bool {
    return t.Status == StatusPending
}

func (t Todo) IsOverdue() bool {
    if t.DueDate == nil {
        return false
    }
    return time.Now().After(*t.DueDate) && t.IsPending()
}
```

### 3.2 データベーススキーマ

```sql
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 0,
    due_date DATETIME,
    work_duration INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_due_date ON todos(due_date);
CREATE INDEX idx_todos_created_at ON todos(created_at);
```

## 4. Repository層詳細

### 4.1 Repositoryインターフェース

```go
package repository

import (
    "context"
    "github.com/yourusername/koto/internal/model"
)

type TodoRepository interface {
    // Create creates a new todo item
    Create(ctx context.Context, todo *model.Todo) error

    // GetByID retrieves a todo by ID
    GetByID(ctx context.Context, id int64) (*model.Todo, error)

    // GetAll retrieves all todos
    GetAll(ctx context.Context) ([]*model.Todo, error)

    // GetByStatus retrieves todos by status
    GetByStatus(ctx context.Context, status model.TodoStatus) ([]*model.Todo, error)

    // Update updates a todo
    Update(ctx context.Context, todo *model.Todo) error

    // Delete deletes a todo by ID
    Delete(ctx context.Context, id int64) error

    // MarkAsCompleted marks a todo as completed
    MarkAsCompleted(ctx context.Context, id int64) error

    // AddWorkDuration adds work duration (in minutes) to a todo
    AddWorkDuration(ctx context.Context, id int64, minutes int) error

    // Close closes the repository connection
    Close() error
}
```

### 4.2 SQLite実装

modernc.org/sqlite を使用することで、CGO不要の Pure Go 実装を実現します。

```go
package repository

import (
    "context"
    "database/sql"
    "time"

    _ "modernc.org/sqlite"  // Pure Go SQLite driver
    "github.com/yourusername/koto/internal/model"
)

type SQLiteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
    // modernc.org/sqliteを使用（Pure Go、CGO不要）
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    // ファイルパーミッションを設定（セキュリティ）
    if err := os.Chmod(dbPath, 0600); err != nil {
        return nil, err
    }

    // Initialize schema
    if err := initSchema(db); err != nil {
        return nil, err
    }

    return &SQLiteRepository{db: db}, nil
}

// 各メソッドの実装...
```

**注意点:**
- `CGO_ENABLED=0` でビルド可能
- クロスコンパイルが容易
- 依存関係がシンプル（Pure Go）

## 5. Service層詳細

### 5.1 TodoService

```go
package service

import (
    "context"
    "errors"
    "time"

    "github.com/yourusername/koto/internal/model"
    "github.com/yourusername/koto/internal/repository"
)

var (
    ErrTodoNotFound    = errors.New("todo not found")
    ErrInvalidTitle    = errors.New("title cannot be empty")
    ErrInvalidPriority = errors.New("invalid priority")
)

type TodoService struct {
    repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}

func (s *TodoService) AddTodo(ctx context.Context, title, description string, priority model.Priority, dueDate *time.Time) (*model.Todo, error) {
    if err := s.validateTitle(title); err != nil {
        return nil, err
    }

    todo := &model.Todo{
        Title:       title,
        Description: description,
        Status:      model.StatusPending,
        Priority:    priority,
        DueDate:     dueDate,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := s.repo.Create(ctx, todo); err != nil {
        return nil, err
    }

    return todo, nil
}

func (s *TodoService) EditTodo(ctx context.Context, id int64, title, description string, priority model.Priority, dueDate *time.Time) error {
    todo, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if todo == nil {
        return ErrTodoNotFound
    }

    if err := s.validateTitle(title); err != nil {
        return err
    }

    todo.Title = title
    todo.Description = description
    todo.Priority = priority
    todo.DueDate = dueDate
    todo.UpdatedAt = time.Now()

    return s.repo.Update(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int64) error {
    return s.repo.Delete(ctx, id)
}

func (s *TodoService) CompleteTodo(ctx context.Context, id int64) error {
    return s.repo.MarkAsCompleted(ctx, id)
}

func (s *TodoService) ListTodos(ctx context.Context) ([]*model.Todo, error) {
    return s.repo.GetAll(ctx)
}

func (s *TodoService) ListPendingTodos(ctx context.Context) ([]*model.Todo, error) {
    return s.repo.GetByStatus(ctx, model.StatusPending)
}

func (s *TodoService) ListCompletedTodos(ctx context.Context) ([]*model.Todo, error) {
    return s.repo.GetByStatus(ctx, model.StatusCompleted)
}

func (s *TodoService) validateTitle(title string) error {
    if title == "" {
        return ErrInvalidTitle
    }
    return nil
}

func (s *TodoService) ExportToJSON(ctx context.Context, filepath string) error {
    todos, err := s.repo.GetAll(ctx)
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(todos, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filepath, data, 0600)
}

func (s *TodoService) ImportFromJSON(ctx context.Context, filepath string) error {
    data, err := os.ReadFile(filepath)
    if err != nil {
        return err
    }

    var todos []*model.Todo
    if err := json.Unmarshal(data, &todos); err != nil {
        return err
    }

    // インポート前に既存データを削除するかユーザーに確認する実装も検討
    for _, todo := range todos {
        if err := s.repo.Create(ctx, todo); err != nil {
            return err
        }
    }

    return nil
}

func (s *TodoService) AddWorkDuration(ctx context.Context, id int64, minutes int) error {
    // ToDoが存在するか確認
    todo, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if todo == nil {
        return ErrTodoNotFound
    }

    // 作業時間を追加
    return s.repo.AddWorkDuration(ctx, id, minutes)
}
```

## 6. TUI層詳細設計

### 6.1 Bubbletea Model

```go
package tui

import (
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourusername/koto/internal/model"
    "github.com/yourusername/koto/internal/service"
)

type ViewMode int

const (
    ViewModeList ViewMode = iota
    ViewModeAdd
    ViewModeEdit
    ViewModeDelete
    ViewModeHelp
    ViewModePomodoro
)

type Model struct {
    service         *service.TodoService
    todos           []*model.Todo
    cursor          int
    viewMode        ViewMode
    input           textinput.Model
    message         string
    err             error
    width           int
    height          int
    selectedID      int64
    pomodoroTodoID  *int64        // ポモドーロタイマーに紐づくToDo ID (nilの場合は紐づけなし)
    pomodoroStarted time.Time     // ポモドーロ開始時刻
    pomodoroDuration time.Duration // ポモドーロの長さ（デフォルト25分）
}

func NewModel(service *service.TodoService) Model {
    ti := textinput.New()
    ti.Placeholder = "コマンドを入力してください (/help でヘルプ)"
    ti.Focus()
    ti.CharLimit = 500
    ti.Width = 80

    return Model{
        service:  service,
        todos:    []*model.Todo{},
        cursor:   0,
        viewMode: ViewModeList,
        input:    ti,
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(textinput.Blink, loadTodos(m.service))
}
```

### 6.2 Update関数

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "esc":
            if m.viewMode != ViewModeList {
                m.viewMode = ViewModeList
                m.input.SetValue("")
                return m, nil
            }
            return m, tea.Quit

        case "enter":
            return m.handleEnter()

        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
            return m, nil

        case "down", "j":
            if m.cursor < len(m.todos)-1 {
                m.cursor++
            }
            return m, nil
        }

    case todosLoadedMsg:
        m.todos = msg.todos
        m.err = msg.err
        return m, nil
    }

    m.input, cmd = m.input.Update(msg)
    return m, cmd
}

func (m *Model) handleEnter() (tea.Model, tea.Cmd) {
    value := m.input.Value()
    m.input.SetValue("")

    return m, parseAndExecuteCommand(m.service, value)
}
```

### 6.3 View関数

```go
func (m Model) View() string {
    switch m.viewMode {
    case ViewModeList:
        return m.renderListView()
    case ViewModeAdd:
        return m.renderAddView()
    case ViewModeEdit:
        return m.renderEditView()
    case ViewModeDelete:
        return m.renderDeleteView()
    case ViewModeHelp:
        return m.renderHelpView()
    case ViewModePomodoro:
        return m.renderPomodoroView()
    default:
        return m.renderListView()
    }
}

func (m Model) renderListView() string {
    var s string

    s += titleStyle.Render("📝 koto - ToDo Manager") + "\n\n"

    if len(m.todos) == 0 {
        s += emptyStyle.Render("ToDoがありません。/add で新しいToDoを追加してください。") + "\n"
    } else {
        // ヘッダー行
        s += headerStyle.Render("No.  Title                Description          Created At") + "\n"
        s += separatorStyle.Render(strings.Repeat("─", 80)) + "\n"

        for i, todo := range m.todos {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }

            // No. (ID)
            no := fmt.Sprintf("%-4d", todo.ID)

            // Title (最大20文字、切り詰め)
            title := truncateString(todo.Title, 20)

            // Description (最大20文字、切り詰め)
            desc := truncateString(todo.Description, 20)

            // Created At (日付のみ)
            created := todo.CreatedAt.Format("2006-01-02")

            // 行の組み立て
            line := fmt.Sprintf("%s %s %s %s %s", cursor, no, title, desc, created)

            if m.cursor == i {
                s += selectedStyle.Render(line) + "\n"
            } else {
                s += line + "\n"
            }

            // 区切り線
            s += separatorStyle.Render(strings.Repeat("─", 80)) + "\n"
        }
    }

    s += "\n" + m.input.View() + "\n"

    if m.message != "" {
        s += "\n" + messageStyle.Render(m.message) + "\n"
    }

    if m.err != nil {
        s += "\n" + errorStyle.Render("エラー: "+m.err.Error()) + "\n"
    }

    s += "\n" + helpStyle.Render("使い方: /help | 終了: Ctrl+C")

    return s
}

func (m Model) renderPomodoroView() string {
    var s string

    // タイマーの残り時間を計算
    elapsed := time.Since(m.pomodoroStarted)
    remaining := m.pomodoroDuration - elapsed

    if remaining < 0 {
        remaining = 0
    }

    minutes := int(remaining.Minutes())
    seconds := int(remaining.Seconds()) % 60

    s += pomodoroTitleStyle.Render("🍅 ポモドーロタイマー") + "\n\n"

    // タイマー表示（大きく表示）
    timerText := fmt.Sprintf("%02d:%02d", minutes, seconds)
    s += pomodoroTimerStyle.Render(timerText) + "\n\n"

    // 紐づくToDoがある場合、タイトルを表示
    if m.pomodoroTodoID != nil {
        for _, todo := range m.todos {
            if todo.ID == *m.pomodoroTodoID {
                s += pomodoroTaskStyle.Render(fmt.Sprintf("作業中: %s", todo.Title)) + "\n\n"
                break
            }
        }
    } else {
        s += pomodoroTaskStyle.Render("フリータイマーモード") + "\n\n"
    }

    // タイマー終了時
    if remaining == 0 {
        s += pomodoroCompleteStyle.Render("🎉 ポモドーロ完了！") + "\n"
        s += "Enterキーを押してメイン画面に戻る\n"
    } else {
        s += helpStyle.Render("Escキーでキャンセル")
    }

    return s
}
```

### 6.4 コマンドパーサー

```go
package tui

import (
    "context"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourusername/koto/internal/service"
)

type commandExecutedMsg struct {
    message string
    err     error
}

type todosLoadedMsg struct {
    todos []*model.Todo
    err   error
}

func parseAndExecuteCommand(svc *service.TodoService, input string) tea.Cmd {
    return func() tea.Msg {
        input = strings.TrimSpace(input)

        if !strings.HasPrefix(input, "/") {
            return commandExecutedMsg{
                err: errors.New("コマンドは / で始める必要があります"),
            }
        }

        parts := strings.Fields(input)
        command := parts[0]
        args := parts[1:]

        ctx := context.Background()

        switch command {
        case "/add":
            return handleAddCommand(ctx, svc, args)
        case "/edit":
            return handleEditCommand(ctx, svc, args)
        case "/delete":
            return handleDeleteCommand(ctx, svc, args)
        case "/done":
            return handleDoneCommand(ctx, svc, args)
        case "/list":
            return handleListCommand(ctx, svc)
        case "/export":
            return handleExportCommand(ctx, svc, args)
        case "/import":
            return handleImportCommand(ctx, svc, args)
        case "/pomo":
            return handlePomodoroCommand(ctx, svc, args)
        case "/help":
            return commandExecutedMsg{message: getHelpText()}
        case "/quit":
            return tea.Quit()
        default:
            return commandExecutedMsg{
                err: errors.New("不明なコマンド: " + command),
            }
        }
    }
}

func loadTodos(svc *service.TodoService) tea.Cmd {
    return func() tea.Msg {
        todos, err := svc.ListTodos(context.Background())
        return todosLoadedMsg{todos: todos, err: err}
    }
}

// 各コマンドハンドラーの実装...
```

### 6.5 スタイル定義

```go
package tui

import "github.com/charmbracelet/lipgloss"

var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        MarginBottom(1)

    emptyStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Italic(true)

    messageStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Bold(true)

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)

    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        MarginTop(1)

    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("212")).
        Bold(true)

    pomodoroTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("196")).
        MarginBottom(1)

    pomodoroTimerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("226")).
        FontSize(48).  // 大きなフォント（TUIでは擬似的に表現）
        Align(lipgloss.Center)

    pomodoroTaskStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Italic(true)

    pomodoroCompleteStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("46"))
)
```

## 7. コマンド仕様詳細

### 7.1 /add - ToDo追加

**構文:**
```
/add <タイトル> [--desc=<説明>] [--priority=<low|medium|high>] [--due=<YYYY-MM-DD>]
```

**例:**
```
/add レポートを書く
/add レポートを書く --desc=月次報告書 --priority=high --due=2024-12-31
```

**処理フロー:**
1. コマンドパース
2. バリデーション（タイトル必須）
3. Service層のAddTodo呼び出し
4. 成功メッセージ表示
5. 一覧を再読み込み

### 7.2 /edit - ToDo編集

**構文:**
```
/edit <ID> [--title=<新しいタイトル>] [--desc=<新しい説明>] [--priority=<low|medium|high>] [--due=<YYYY-MM-DD>]
```

**例:**
```
/edit 1 --title=新しいタイトル
/edit 2 --priority=high --due=2024-12-25
```

**処理フロー:**
1. ID指定でTodoを検索
2. 存在チェック
3. 指定されたフィールドを更新
4. Service層のEditTodo呼び出し
5. 成功メッセージ表示
6. 一覧を再読み込み

### 7.3 /delete - ToDo削除

**構文:**
```
/delete <ID>
```

**例:**
```
/delete 1
```

**処理フロー:**
1. ID指定でTodoを検索
2. 存在チェック
3. 削除確認プロンプト表示（オプション）
4. Service層のDeleteTodo呼び出し
5. 成功メッセージ表示
6. 一覧を再読み込み

### 7.4 /done - ToDo完了

**構文:**
```
/done <ID>
```

**例:**
```
/done 1
```

**処理フロー:**
1. ID指定でTodoを検索
2. 存在チェック
3. Service層のCompleteTodo呼び出し
4. 成功メッセージ表示
5. 一覧を再読み込み

### 7.5 /list - 一覧表示

**構文:**
```
/list [--status=<pending|completed|all>]
```

**例:**
```
/list
/list --status=pending
/list --status=completed
```

**処理フロー:**
1. ステータスフィルターの解析
2. Service層の対応するメソッド呼び出し
3. 一覧を表示

### 7.6 /help - ヘルプ表示

**構文:**
```
/help
```

**処理フロー:**
1. ヘルプテキストを表示

### 7.7 /export - JSONエクスポート

**構文:**
```
/export [ファイルパス]
```

**例:**
```
/export
/export ~/backups/todos-2024-12-01.json
```

**処理フロー:**
1. エクスポート先ファイルパスの解析（デフォルト: `~/.koto/export.json`）
2. Service層のExportToJSON呼び出し
3. 全ToDoデータをJSON形式でファイルに書き込み
4. 成功メッセージ表示

**エクスポート形式:**
```json
[
  {
    "id": 1,
    "title": "レポートを書く",
    "description": "月次報告書",
    "status": 0,
    "priority": 2,
    "due_date": "2024-12-31T23:59:59Z",
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  }
]
```

### 7.8 /import - JSONインポート

**構文:**
```
/import <ファイルパス>
```

**例:**
```
/import ~/backups/todos-2024-12-01.json
```

**処理フロー:**
1. インポート元ファイルパスの解析
2. ファイルの存在確認
3. JSON形式のバリデーション
4. 確認プロンプト表示（既存データとの重複について）
5. Service層のImportFromJSON呼び出し
6. 成功メッセージ表示（インポート件数を含む）
7. 一覧を再読み込み

**注意事項:**
- IDが重複する場合の処理（上書き or スキップ）をユーザーに確認
- 大量データのインポート時の進捗表示

### 7.9 /pomo - ポモドーロタイマー

**構文:**
```
/pomo [ToDo ID]
```

**例:**
```
/pomo              # フリータイマーモード（25分）
/pomo 1            # ID 1のToDoに紐づけて25分タイマー開始
```

**処理フロー:**
1. 引数チェック（ToDo IDが指定されている場合、存在確認）
2. ポモドーロモードに遷移
3. タイマー開始（25分 = 1500秒）
4. 1秒ごとに画面を更新（tea.Tick使用）
5. タイマー終了時
   - アラーム表示
   - ToDo IDが指定されている場合、作業時間を記録（25分追加）
   - Enterキー待ち
6. Enterキーでメイン画面に戻る

**キャンセル:**
- Escキーでタイマーをキャンセルし、メイン画面に戻る
- キャンセル時は作業時間を記録しない

**注意事項:**
- タイマー実行中は専用画面が表示される
- バックグラウンドでの実行はサポートしない
- タイマー完了時のみ作業時間を記録

### 7.10 /quit - アプリ終了

**構文:**
```
/quit
```

**処理フロー:**
1. データベース接続をクローズ
2. アプリケーション終了

## 8. エラーハンドリング

### 8.1 エラー種別

| エラー種別 | 説明 | 対応方法 |
|-----------|------|---------|
| バリデーションエラー | 入力値が不正 | エラーメッセージを表示し、再入力を促す |
| データベースエラー | DB操作失敗 | エラーログを出力し、ユーザーに通知 |
| 存在しないID | 指定IDのTodoが見つからない | エラーメッセージを表示 |
| コマンド解析エラー | 不正なコマンド | ヘルプを表示し、正しい構文を案内 |
| ファイルI/Oエラー | エクスポート/インポート失敗 | ファイルパスやパーミッションを確認するよう案内 |
| JSON解析エラー | 不正なJSON形式 | 正しいJSON形式でエクスポートされたファイルを指定するよう案内 |

### 8.2 エラーメッセージ例

```go
var (
    ErrMessages = map[error]string{
        service.ErrTodoNotFound:    "指定されたToDoが見つかりません",
        service.ErrInvalidTitle:    "タイトルは必須です",
        service.ErrInvalidPriority: "優先度は low, medium, high のいずれかを指定してください",
        service.ErrFileNotFound:    "ファイルが見つかりません",
        service.ErrInvalidJSON:     "JSONの形式が不正です",
        service.ErrExportFailed:    "エクスポートに失敗しました",
        service.ErrImportFailed:    "インポートに失敗しました",
    }
)
```

## 9. データ永続化

### 9.1 データベースファイル配置

- **パス**: `~/.koto/koto.db`
- **パーミッション**: 0600 (所有者のみ読み書き可能)

### 9.2 初期化処理

```go
func initDatabase() (*sql.DB, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    kotoDir := filepath.Join(homeDir, ".koto")
    if err := os.MkdirAll(kotoDir, 0700); err != nil {
        return nil, err
    }

    dbPath := filepath.Join(kotoDir, "koto.db")
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    return db, nil
}
```

## 10. テスト戦略

### 10.1 ユニットテスト

- **対象**: 各パッケージの関数・メソッド
- **ツール**: Go標準のtestingパッケージ
- **カバレッジ目標**: 80%以上

**例:**
```go
func TestTodoService_AddTodo(t *testing.T) {
    // Setup
    repo := repository.NewMockRepository()
    svc := service.NewTodoService(repo)

    // Execute
    todo, err := svc.AddTodo(context.Background(), "Test Todo", "", model.PriorityMedium, nil)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, todo)
    assert.Equal(t, "Test Todo", todo.Title)
}
```

### 10.2 統合テスト

- **対象**: Repository層とデータベースの統合
- **ツール**: Go testing + インメモリSQLite

### 10.3 E2Eテスト

- **対象**: アプリケーション全体のフロー
- **ツール**: Go testing + Bubbletea test utilities

## 11. ビルドとリリース

### 11.1 Makefile

```makefile
.PHONY: build test clean install release-local

# Pure Goビルド（CGO不要）
build:
	CGO_ENABLED=0 go build -o bin/koto ./cmd/koto

# クロスコンパイル用ビルド
build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/koto-linux-amd64 ./cmd/koto
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/koto-darwin-amd64 ./cmd/koto
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/koto-darwin-arm64 ./cmd/koto
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/koto-windows-amd64.exe ./cmd/koto

test:
	go test -v -cover ./...

clean:
	rm -rf bin/

install: build
	cp bin/koto $(GOPATH)/bin/koto

lint:
	golangci-lint run

release:
	goreleaser release --clean
```

### 11.2 GoReleaserの設定

```yaml
# .goreleaser.yml
project_name: koto

builds:
  - id: koto
    main: ./cmd/koto
    binary: koto
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
```

## 12. データ可搬性の実現

### 12.1 エクスポート機能
SQLiteはバイナリ形式ですが、`/export`コマンドにより以下を実現：
- 人間が読めるJSON形式への変換
- 他の環境への移行が容易
- バックアップの作成
- 他のツールとの連携

### 12.2 インポート機能
`/import`コマンドにより以下を実現：
- バックアップからの復元
- 他の環境からのデータ移行
- 外部ツールで作成したToDoデータの取り込み

### 12.3 データ連携の例
```bash
# 環境A: データをエクスポート
/export ~/todos-backup.json

# 環境B: データをインポート
/import ~/todos-backup.json

# スクリプトでの加工も可能
cat ~/todos-backup.json | jq '.[] | select(.priority == 2)' > high-priority.json
/import high-priority.json
```

## 13. 将来の拡張

### 13.1 フェーズ2（v2.0）
- タグ機能（複数タグの付与）
- 高度なフィルタリング・検索（全文検索、複数条件）
- カテゴリ分類

### 13.2 フェーズ3（v3.0）
- 定期タスク（リカーリング）
- CSVエクスポート
- 設定ファイル対応（YAML/TOML）
- サブタスク機能

### 13.3 フェーズ4（v4.0以降）
- クラウド同期
- チーム共有
- Web UI
- モバイルアプリ連携
