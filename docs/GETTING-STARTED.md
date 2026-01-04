# Getting Started with Minty Clean Architecture

This guide will walk you through building your first application using the Minty System with clean architecture principles.

## 🎯 What You'll Build

A simple task management application that demonstrates:
- ✅ Clean architecture with separated concerns  
- ✅ Pure business logic with zero UI dependencies
- ✅ Theme-based UI components
- ✅ Testable domain services

## 📋 Prerequisites

```bash
go version # Go 1.21 or higher required
```

## 🚀 Step 1: Create Your Project

```bash
mkdir task-manager
cd task-manager
go mod init github.com/yourusername/task-manager
```

## 📦 Step 2: Install Dependencies

```bash
# Core packages
go get github.com/ha1tch/minty
go get github.com/ha1tch/mintyex
go get github.com/ha1tch/mintyui

# Theme
go get github.com/ha1tch/minty-bootstrap-theme
```

## 💼 Step 3: Create Your Domain (Pure Business Logic)

Create `domain/tasks.go`:

```go
package domain

import (
    "errors"
    "time"
    "github.com/ha1tch/mintyex"
)

// Pure business types - no UI dependencies
type Task struct {
    ID          string
    Title       string
    Description string
    Status      string // pending, in_progress, completed
    Priority    string // low, medium, high
    DueDate     time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// TaskService contains pure business logic
type TaskService struct {
    tasks []Task
}

func NewTaskService() *TaskService {
    return &TaskService{tasks: make([]Task, 0)}
}

// Pure business operations
func (ts *TaskService) CreateTask(title, description string, priority string, dueDate time.Time) (*Task, error) {
    // Business validation
    if title == "" {
        return nil, errors.New("title is required")
    }
    
    task := Task{
        ID:          generateID(),
        Title:       title,
        Description: description,
        Status:      "pending",
        Priority:    priority,
        DueDate:     dueDate,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    ts.tasks = append(ts.tasks, task)
    return &task, nil
}

func (ts *TaskService) GetTasks() []Task {
    return ts.tasks
}

func (ts *TaskService) UpdateTaskStatus(taskID, status string) error {
    for i, task := range ts.tasks {
        if task.ID == taskID {
            ts.tasks[i].Status = status
            ts.tasks[i].UpdatedAt = time.Now()
            return nil
        }
    }
    return errors.New("task not found")
}

// Pure business logic helper
func generateID() string {
    return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

// Data preparation for presentation (still pure)
type TaskDisplayData struct {
    Task             Task
    FormattedDueDate string
    StatusClass      string
    PriorityClass    string
    IsOverdue        bool
}

func PrepareTaskForDisplay(task Task) TaskDisplayData {
    return TaskDisplayData{
        Task:             task,
        FormattedDueDate: task.DueDate.Format("Jan 2, 2006"),
        StatusClass:      getStatusClass(task.Status),
        PriorityClass:    getPriorityClass(task.Priority),
        IsOverdue:        time.Now().After(task.DueDate) && task.Status != "completed",
    }
}

func getStatusClass(status string) string {
    switch status {
    case "completed": return "success"
    case "in_progress": return "warning"
    default: return "secondary"
    }
}

func getPriorityClass(priority string) string {
    switch priority {
    case "high": return "danger"
    case "medium": return "warning"  
    default: return "info"
    }
}
```

## 🌐 Step 4: Create Presentation Adapters (UI Components)

Create `presentation/taskui.go`:

```go
package presentation

import (
    "fmt"
    mi "github.com/ha1tch/minty"
    mui "github.com/ha1tch/mintyui"
    "github.com/ha1tch/mintyex"
    "github.com/yourusername/task-manager/domain"
)

// Presentation adapters - convert domain data to UI

func TaskCard(theme mui.Theme, task domain.Task) mi.H {
    displayData := domain.PrepareTaskForDisplay(task)
    
    return theme.Card(task.Title, func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("task-card"),
            b.P(mi.Class("task-description"), task.Description),
            b.Div(mi.Class("task-meta"),
                theme.Badge(task.Status, displayData.StatusClass)(b),
                theme.Badge(task.Priority, displayData.PriorityClass)(b),
                b.Small(mi.Class("due-date"), "Due: ", displayData.FormattedDueDate),
                mintyex.If(displayData.IsOverdue,
                    theme.Badge("OVERDUE", "danger"),
                )(b),
            ),
            b.Div(mi.Class("task-actions"),
                mintyex.Unless(task.Status == "completed",
                    theme.Button("Complete", "success", 
                        mi.DataAttr("task-id", task.ID))(b),
                )(b),
            ),
        )
    })
}

func TaskList(theme mui.Theme, tasks []domain.Task) mi.H {
    return func(b *mi.Builder) mi.Node {
        return b.Div(mi.Class("task-list"),
            mi.NewFragment(mintyex.Each(tasks, func(task domain.Task) mi.H {
                return TaskCard(theme, task)
            })...),
        )
    }
}

func TaskDashboard(theme mui.Theme, service *domain.TaskService) mi.H {
    tasks := service.GetTasks()
    
    pendingCount := 0
    completedCount := 0
    overdueCount := 0
    
    for _, task := range tasks {
        switch task.Status {
        case "completed":
            completedCount++
        case "pending":
            pendingCount++
        }
        
        displayData := domain.PrepareTaskForDisplay(task)
        if displayData.IsOverdue {
            overdueCount++
        }
    }
    
    return mui.Dashboard(theme, "Task Manager",
        // Sidebar
        func(b *mi.Builder) mi.Node {
            return theme.Sidebar(func(b *mi.Builder) mi.Node {
                return theme.Nav([]mui.NavItem{
                    {Text: "📋 All Tasks", URL: "/", Active: true},
                    {Text: "➕ New Task", URL: "/new"},
                    {Text: "✅ Completed", URL: "/completed"},
                    {Text: "⚠️ Overdue", URL: "/overdue"},
                }))(b)
            })(b)
        },
        
        // Main content
        func(b *mi.Builder) mi.Node {
            return b.Div(mi.Class("dashboard-main"),
                // Metrics
                b.Section(mi.Class("metrics"),
                    b.H2("Task Overview"),
                    mintyex.GridLayout(3, "1rem")(
                        mui.StatsCard(theme, "Pending Tasks", 
                            fmt.Sprintf("%d", pendingCount), "Tasks to do"),
                        mui.StatsCard(theme, "Completed Tasks", 
                            fmt.Sprintf("%d", completedCount), "Tasks done"),
                        mui.StatsCard(theme, "Overdue Tasks", 
                            fmt.Sprintf("%d", overdueCount), "Needs attention"),
                    )(b),
                ),
                
                // Task list
                b.Section(mi.Class("tasks"),
                    b.H2("All Tasks"),
                    TaskList(theme, tasks)(b),
                ),
            )
        },
    )
}
```

## 🏢 Step 5: Create Application Service (Orchestration)

Create `app/application.go`:

```go
package app

import (
    "github.com/yourusername/task-manager/domain"
)

// Application layer - orchestrates domain services
type Application struct {
    TaskService *domain.TaskService
}

func NewApplication() *Application {
    return &Application{
        TaskService: domain.NewTaskService(),
    }
}

// Initialize with sample data
func (app *Application) InitSampleData() {
    app.TaskService.CreateTask(
        "Learn Clean Architecture", 
        "Study the principles and implement them in Go",
        "high",
        time.Now().AddDate(0, 0, 7), // Due in 1 week
    )
    
    app.TaskService.CreateTask(
        "Build Task Manager",
        "Create a task management app using Minty System", 
        "medium",
        time.Now().AddDate(0, 0, 3), // Due in 3 days
    )
    
    app.TaskService.CreateTask(
        "Write Tests",
        "Add comprehensive tests for business logic",
        "medium", 
        time.Now().AddDate(0, 0, 5), // Due in 5 days
    )
}
```

## 🌐 Step 6: Create Web Application (HTTP Interface)

Create `main.go`:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    mi "github.com/ha1tch/minty"
    "github.com/ha1tch/minty-bootstrap-theme"
    
    "github.com/yourusername/task-manager/app"
    "github.com/yourusername/task-manager/presentation"
)

func main() {
    // 1. Initialize application (domain services)
    application := app.NewApplication()
    application.InitSampleData()
    
    // 2. Choose theme (infrastructure)
    theme := bootstrap.NewBootstrapTheme()
    
    // 3. Create dashboard
    dashboard := presentation.TaskDashboard(theme, application.TaskService)
    
    // 4. Render HTML
    html := mi.Render(func(b *mi.Builder) mi.Node {
        return bootstrap.BootstrapDocument("Task Manager", dashboard)(b)
    })
    
    // 5. Serve web application
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, html)
    })
    
    fmt.Println("🚀 Task Manager running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 🧪 Step 7: Add Tests (Business Logic)

Create `domain/tasks_test.go`:

```go
package domain

import (
    "testing"
    "time"
)

func TestTaskCreation(t *testing.T) {
    service := NewTaskService()
    
    task, err := service.CreateTask(
        "Test Task",
        "This is a test task", 
        "high",
        time.Now().AddDate(0, 0, 1),
    )
    
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if task.Title != "Test Task" {
        t.Errorf("Expected title 'Test Task', got '%s'", task.Title)
    }
    
    if task.Status != "pending" {
        t.Errorf("Expected status 'pending', got '%s'", task.Status)
    }
}

func TestTaskValidation(t *testing.T) {
    service := NewTaskService()
    
    // Test empty title validation
    _, err := service.CreateTask("", "Description", "low", time.Now())
    
    if err == nil {
        t.Error("Expected error for empty title, got nil")
    }
}

func TestStatusUpdate(t *testing.T) {
    service := NewTaskService()
    
    task, _ := service.CreateTask("Test", "Description", "low", time.Now())
    
    err := service.UpdateTaskStatus(task.ID, "completed")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    tasks := service.GetTasks()
    if tasks[0].Status != "completed" {
        t.Errorf("Expected status 'completed', got '%s'", tasks[0].Status)
    }
}
```

## 🏃 Step 8: Run Your Application

```bash
# Run tests (business logic only - no UI dependencies!)
go test ./domain

# Run the application
go run main.go

# Visit http://localhost:8080
```

## 🎉 What You've Accomplished

✅ **Clean Architecture**: Business logic is completely separate from UI  
✅ **Testable**: Domain logic can be tested without any UI dependencies  
✅ **Theme Support**: Can easily switch from Bootstrap to Tailwind  
✅ **Maintainable**: Clear separation makes changes safe and predictable  
✅ **Scalable**: Can easily add new domains following the same pattern  

## 🚀 Next Steps

1. **Add More Features**: Try adding task editing, filtering, or categories
2. **Switch Themes**: Replace Bootstrap with Tailwind theme
3. **Add Persistence**: Create a repository interface and database implementation  
4. **Add More Domains**: Create a user management domain
5. **Add HTMX**: Make the UI interactive with HTMX attributes

## 🎓 Key Takeaways

- **Business Logic First**: Always start with pure domain models
- **Zero UI Dependencies**: Domain packages should never import UI code
- **Presentation Adapters**: Convert domain data to UI components
- **Theme System**: Abstract UI behind interfaces for flexibility
- **Test Business Logic**: Pure functions are easy to test

You now have a solid foundation for building maintainable HTML applications with clean architecture principles!
