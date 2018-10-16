# go-chat

## Requirements

**Functional requirements:**

- users can have private chat with another user 
- users can have group chat

## Questions

1. How to join room from the client side?

  You can't join the the room from the client-side, it has to happen from the server side. The default room you should be in is the one that has contains for example, your user id. You will always be subscribed to that particular room, so that you can receive messages from others when they join your room.

  You can join multiple room. Let's say you are User A with id `1`, so you will be in room `1`. You want to talk with User `B` with id `2`, so you join his room. Room 2 now has both User A and B, and any messages broadcasted can only be seen by the two users.

2. How to prevent other users from joining your room?

  Probably add a scope-like mechanism. This can be implemented using JWT. But first, you need to identify the relationship between users, and whether you would like to have the requirements to be able to message others that are not in your friend list.

3. How to scale?

  Use Redis pub-sub. This would also scale when we have multiple servers, as different servers might be holding different websocket connections. The last history message can also be stored in redis.

4. How to handle scenario where user's are offline?

  One possibility is still to allow the message to be sent, and notify the other user through push notification.

## Design 

### Handling relationships

In chat, we have concept of friends. Basically, a user will have a list of friends he/she can chat with. When the chat is initialized, the list of friends would be added, so that the users can choose who to chat with. Communications needs to be established between the two users in order for them to send messages to one another.

A naive approach is to simply load all the user's friends when initializing the chat and attempting to establish connections with them. This might not be ideal when the users have a lot of friends, as they would only chat with a select few at a time. Hence, there are some limits that needs to be establish here. One way is to rank the friends by the last chat, and adding only the last 15 friends into the chat list (offline and online). If the user wishes to chat with someone else, he/she can search for the name and start a new chat.

### Optimizing Delivery

Whenever a new message is created, it should not be send directly to the other party. We first need to check if the other party is available before sending them a message. Hence, the message should first be stored in a persistent database. Then, if the other party is not online within n minutes, they should receive a notification message.


### Chat Rooms

When users are connected to the chat server, they have no knowledge on any of their peers. They cannot send messages to one another.

```
User A joined
User B joined
User C joined
```

In order for them to send messages to one another, we can implement a few rules:
- user can only send messages to their known contacts
- user can send message to a group they own
- user cannot send message if they are blocked by the receiving party
- user can only send send message to someone not in the contact list, but they may be blocked by the receiving party

Also, to avoid spam, we can set a few other rules:
- only 1 message per second (this rate is sufficient enough, in we take into consideration the average typing speed of users)
- cannot contain uncensored word
- spam links to external sources

But back to the topic of rooms, how do we create an actual room that contains only the user to chat to?

```
createRoom(owner)
addUsersToRoom(listOfUsers)
getRooms(owner)
getUsersFromRoom(room)
broadcast(room)
```

The API may look like this:

```
/v1/groups/:groupname/users/:id
/v1/groups/:groupname/users
/v1/messages
/v1/messages/:id
```

If a room has only 1 user, then it is a private chat, else it would be a group chat. When the user join the chat server, we can do the following operations to enable them to chat:

```
// Pseudo-code in any language.
func join():
  // Perform handshake to check origin
  handshake() 
  
  // Perform authorization of the user. This could be done by passing a valid token that can validated by the chat server.
  // A successful authorization will return the id of the user.
  id = authz() 
  
  // Register the id to the pubsub, could be a external redis storage to make the server distributed.
  pubsub.subscribe(id)
  
  // Set the user's connection to the in-memory storage. This will tie the client to a specific chat-server, but messages could still be distributed across servers through redis pub-sub.
  localconn.set(id, conn)
  
  // Read event loop. 
  for {
    let msg = conn.read()
    
    // Store the message in persistent storage.
    snapshot(msg)
    
    // Send to external pub/sub.
    pubsub.publish(msg.to, msg)
  }
  
  // Write event loop. Theoretically running in a separate thread to avoid blocking the main thread.
  for {
    msg = pubsub.receive()
    conn = localconn.get(msg.to)
    conn.write(msg)
  }

  // Cleanup.
  // We have to remove the user at the end of the session. Note that this step might be skipped if an error happened, and the user might not be successfully unsubscribed. To ensure that the user is always removed from the session, ping em periodically.
  pubsub.unsubscribe(id)
  localconn.unset(id, conn)
```

## References
- https://www.thepolyglotdeveloper.com/2016/12/create-real-time-chat-app-golang-angular-2-websockets/
- https://devcenter.heroku.com/articles/go-websockets
- https://www.jonathan-petitcolas.com/2015/01/27/playing-with-websockets-in-go.html
