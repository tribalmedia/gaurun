package gcm

import (
	"fmt"
	"strings"
)

// Message is used by the application server to send a message to
// the GCM server. See the documentation for GCM Architectural
// Overview for more information:
// http://developer.android.com/google/gcm/gcm.html#send-msg
type Message struct {
	// To GCM
	RegistrationIDs []string               `json:"registration_ids,omitempty"`
	CollapseKey     string                 `json:"collapse_key,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	// To FCM Topic
	To           string                 `json:"to,omitempty"`
	Priority     string                 `json:"priority,omitempty"`
	Notification map[string]interface{} `json:"notification,omitempty"`
	Badge            int          `json:"badge,omitempty"`
	Category         string       `json:"category,omitempty"`
	ContentAvailable bool         `json:"content_available,omitempty"`

	DelayWhileIdle        bool   `json:"delay_while_idle,omitempty"`
	TimeToLive            int    `json:"time_to_live,omitempty"`
	RestrictedPackageName string `json:"restricted_package_name,omitempty"`
	DryRun                bool   `json:"dry_run,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration IDs.
func NewMessage(data map[string]interface{}, notification map[string]interface{}, regIDs ...string) *Message {
	if len(regIDs) == 1 && strings.Contains(regIDs[0], "/topics/") {
		// To FCM Topic
		return &Message{To: regIDs[0], Priority: "high", ContentAvailable: true, Notification: notification, Data: data}
	}
	// To Other
	return &Message{RegistrationIDs: regIDs, Data: data, Notification: notification}

}

// validate validates message format. If not well-formated returns error.
func (m *Message) validate() error {
	if m == nil {
		return fmt.Errorf("the message must not be nil")
	}
	// GCMとFCM Topic向けのvalidationの分岐を行う
	if m.To != "" {
		// To FCM Topic
		return nil
	}
	if m.RegistrationIDs == nil {
		return fmt.Errorf("the message's RegistrationIDs field must not be nil")
	}

	if len(m.RegistrationIDs) == 0 {
		return fmt.Errorf("the message must specify at least one registration ID")
	}

	if len(m.RegistrationIDs) > maxRegistrationIDs {
		return fmt.Errorf("the message may specify at most %d registration IDs",
			maxRegistrationIDs)
	}

	if m.TimeToLive < 0 || maxTimeToLive < m.TimeToLive {
		return fmt.Errorf(
			"the message's TimeToLive field must be an integer between 0 and %d (4 weeks)",
			maxTimeToLive,
		)
	}

	return nil
}
