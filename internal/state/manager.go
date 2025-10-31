package state

import (
	"sync"
	"time"
)

// UserState represents user's current state
type UserState string

const (
	// User states
	StateNone                 UserState = "NONE"
	StateWaitingPhone         UserState = "WAITING_PHONE"
	StateWaitingInvoicePayload UserState = "WAITING_INVOICE_PAYLOAD"
	StateWaitingInvoiceContract UserState = "WAITING_INVOICE_CONTRACT"
	StateInvalidPayload       UserState = "INVALID_PAYLOAD"
	StateWaitingForSupport    UserState = "WAITING_FOR_SUPPORT"
	StateChatting             UserState = "CHATTING"
	
	// Admin states
	StateSearchContract       UserState = "SEARCH_CONTRACT"
	StateSearchPhone          UserState = "SEARCH_PHONE"
	StateSearchName           UserState = "SEARCH_NAME"
	StateSearchAddress        UserState = "SEARCH_ADDRESS"
	StateAccountMenuList      UserState = "ACCOUNT_MENU_LIST"
	StateAdminChangeBalance   UserState = "ADMIN_CHANGE_BALANCE"
	StateSendMessagePhone     UserState = "SEND_MESSAGE_PHONE"
	StateSendMessageText      UserState = "SEND_MESSAGE_TEXT"
	StateAnswer               UserState = "ANSWER"
	StateMessageHistory       UserState = "MESSAGE_HISTORY"
)

// StateData holds additional data for state
type StateData map[string]interface{}

// StateEntry represents a state entry with expiration
type StateEntry struct {
	State     UserState
	Data      StateData
	ExpiresAt time.Time
}

// StateManager manages user states
type StateManager struct {
	states map[int64]*StateEntry
	mu     sync.RWMutex
	ttl    time.Duration
}

// NewStateManager creates a new state manager
func NewStateManager(ttl time.Duration) *StateManager {
	sm := &StateManager{
		states: make(map[int64]*StateEntry),
		ttl:    ttl,
	}
	
	// Start cleanup goroutine
	go sm.cleanup()
	
	return sm
}

// SetState sets user state with optional data
func (sm *StateManager) SetState(userID int64, state UserState, data StateData) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.states[userID] = &StateEntry{
		State:     state,
		Data:      data,
		ExpiresAt: time.Now().Add(sm.ttl),
	}
}

// GetState gets user state
func (sm *StateManager) GetState(userID int64) (UserState, StateData, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	entry, exists := sm.states[userID]
	if !exists {
		return StateNone, nil, false
	}
	
	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return StateNone, nil, false
	}
	
	return entry.State, entry.Data, true
}

// UpdateData updates state data without changing state
func (sm *StateManager) UpdateData(userID int64, data StateData) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if entry, exists := sm.states[userID]; exists {
		if entry.Data == nil {
			entry.Data = make(StateData)
		}
		for k, v := range data {
			entry.Data[k] = v
		}
		entry.ExpiresAt = time.Now().Add(sm.ttl)
	}
}

// ClearState clears user state
func (sm *StateManager) ClearState(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, userID)
}

// GetData gets specific data from state
func (sm *StateManager) GetData(userID int64, key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	entry, exists := sm.states[userID]
	if !exists || entry.Data == nil {
		return nil, false
	}
	
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	
	value, exists := entry.Data[key]
	return value, exists
}

// cleanup removes expired states periodically
func (sm *StateManager) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for userID, entry := range sm.states {
			if now.After(entry.ExpiresAt) {
				delete(sm.states, userID)
			}
		}
		sm.mu.Unlock()
	}
}

// GetStateManager returns a state manager instance (singleton pattern)
var globalStateManager *StateManager
var once sync.Once

// GetStateManagerInstance returns singleton state manager
func GetStateManagerInstance() *StateManager {
	once.Do(func() {
		globalStateManager = NewStateManager(30 * time.Minute) // 30 minutes TTL
	})
	return globalStateManager
}

