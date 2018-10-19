# Pseudo-code with golang

```go
package main

import (
	"fmt"
)

// Conn is a websocket connection.
type Conn struct {}

func main() {
	// Sending messages in chat can be simple.
	// Say we have subject A and subject B, our websocket connections may look like this:
	conns := map[string]*Conn {
		"A": connA,
		"B": connB,
	}
	
	// For B to send message to A, B performs:
	conns["A"].write(msg)
	
	// The opposite applies:
	conns["B"].write(msg)
	
	
	// How about messaging in group? Let's add C to the picture.
	conns := map[string]*Conn {
		"A": connA,
		"B": connB,
		"C": connC,
	}
	// A, B, and C are now in the same room. When A sends a message, A needs to perform:
	conns["A"].write(msg) // Write to ownself to display on UI.
	conns["B"].write(msg)
	conns["C"].write(msg)
	
	// The system must have the knowledge that A, B and C is in a room.
	rooms := map[string][]string {
		"room1": []string{"A", "B", "C"},
	}
	
	// And now we can do this. Alternative way is to use map if we want to perform lookup in the room.
	participants := rooms["room1"] 
	for _, p := range participants {
		conns[p].write(msg)
	}
	
	// But there's another issue - how does A, B and C knows the room they are in?
	// We can create a database to store the rooms. The primary id can be the room id. How to query the rooms for A, B, and C?
	| RoomID | Participants | CreatedBy |
	| room1  | A, B, C      | A         |
	| room2  | A, B         | B         |
	
	// This schema makes adding/removing participants in a room easier.
	| Room  | Participants |
	| room1 | A            |
	| room1 | B            |
	| room1 | C            |
	| room2 | A            |
	| room2 | B            |
	| room3 | C            | <- There should not exist a room with only one participants.
	
	// But that leads to another question, how to find room with only 2 participants, say A and B?
	
	
	// What do we need?
	// Get contacts list (with or without room)
	// A: {B, C, D}
	// B: {A, C}
	
	// Get rooms id?
	// A:B room id?
	
	// Let's say we want to get the rooms for user A.
	
	// ConversationTable
	| $RoomID | CreatedBy | Participants | Enabled |
	| room1   | A         | B, C, D      | 1       |
	| room2   | B         | A            | 0       |
	// Can find room created by, but if the room is created by B and A is in the room, it would not be inserted.
	// Hence we need additional table below.
	// This table is still required for creation/deletion of rooms. 
	// Participants cannot be added/created/removed after the room is created.
	
	// ConversationParticipantTable
	| $ParticipantID | UserID | RoomID |
	| 1              | A      | room1  |
	| 2              | B      | room1  |
	| 3              | C      | room1  |
	| 4              | D      | room1  |
	| 5              | A      | room2  |
	| 6              | B      | room2  |
	// To get the rooms for user A, `SELECT * FROM table where UserID = A`
	// Returns room1, room2.
	
	// To check if user A is in room1 `SELECT * FROM table where UserID = A and RoomID = room1`.
	
	// ConversationReplyTable
	| $ConversationID | RoomID | ReplyBy | Msg         | Type 
	| 1               | room1  | A       | hello world | msg, attachment, image, file
	| 2               | room1  | B       | hi          |
	// Get conversation history
	// Get conversation counts
	// Get conversation for user A
	// Get conversation for user B
}
```
