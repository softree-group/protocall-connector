package entity

import "github.com/google/btree"

type Conference struct {
	ID           string  `json:"id"`
	Participants []*User `json:"participants"`
	HostUserID   string  `json:"host_user_id"`
	BridgeID     string  `json:"-"`
<<<<<<< HEAD
<<<<<<< HEAD
	IsRecording  bool    `json:"is_recording"`
=======
>>>>>>> 977da2b (rebase inbloud)
=======
>>>>>>> 9dff4d40660aa86a5ea66aa72ef8bfb947260101
}

func (c Conference) Less(then btree.Item) bool {
	return c.ID < then.(*Conference).ID
}
