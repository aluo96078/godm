package odm

// observer.go - 定義 Observer 架構與註冊機制
// Defines the Observer architecture and registration mechanisms.

// ModelObserver - 定義模型觀察者介面
// Defines the model observer interface
type ModelObserver interface {
	Creating(model interface{}) error
	Created(model interface{}) error
	Updating(model interface{}) error
	Updated(model interface{}) error
	Deleting(model interface{}) error
	Deleted(model interface{}) error
}

// EventFilter - 定義事件過濾器介面
// Defines the event filter interface
type EventFilter interface {
	InterestedIn(stage string) bool
}

// PrioritizedObserver - 定義優先級觀察者介面
// Defines the prioritized observer interface
type PrioritizedObserver interface {
	Priority() int
}

// TypedObserver - 定義類型觀察者介面
// Defines the typed observer interface
type TypedObserver interface {
	Accepts(model interface{}) bool
}

// ObservedModel - 定義被觀察模型介面
// Defines the observed model interface
type ObservedModel interface {
	Observers() []ModelObserver
}

var globalObservers []ModelObserver // globalObservers - 儲存全局觀察者
// globalObservers - Stores global observers
var observerErrorHandler func(err error, stage string, model interface{}) // observerErrorHandler - 處理觀察者錯誤的函數
// observerErrorHandler - Function to handle observer errors

// RegisterGlobalObserver - 註冊全局觀察者
// Registers a global observer
func RegisterGlobalObserver(o ModelObserver) {
	globalObservers = append(globalObservers, o)
}

// RegisterObserverErrorHandler - 註冊觀察者錯誤處理函數
// Registers the observer error handler function
func RegisterObserverErrorHandler(handler func(err error, stage string, model interface{})) {
	observerErrorHandler = handler
}

// getObserverPriority - 獲取觀察者優先級
// Retrieves the observer's priority
func getObserverPriority(o ModelObserver) int {
	if p, ok := o.(PrioritizedObserver); ok {
		return p.Priority()
	}
	return 0
}
