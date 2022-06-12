# go-chat

See the improvised version [here](https://github.com/alextanhongpin/go-chat.v2).



## Requirements

**Functional requirements:**

- system should allow users to have a private conversation
- system should allow users to have a group conversation
- system should not expose the user identifier
- system should notify users in group when member is online
- system should notify users in group when member is offline
- system should notify users when message is read
- system should display count of unread messages
- system should mark messages as read
- system should allow only authenticated users to chat
- system should allow users to login with multiple sessions
- system should handle stickiness of the user's session

**Non-functional requirements**

- system should be reliable (messages should be stored)
- system should be available (users should not be disconnected)
- system should be secure (authentication)

**Extended requirements**
- system should handle validation on the messages sent to avoid spam
- system should allow users to block unwanted chat requests
- system should send notification when user is not online

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

The API may look like this, and all endpoints are authorized:

```
/v1/groups/:groupname/users/:id
/v1/groups/:groupname/users
/v1/groups/:groupname/messages
or
/v1/messages/groups
/v1/messages
/v1/messages/:id
/v1/users/:id/contacts # To get the contact lists to display in the chat
```

If a room has only 1 user, then it is a private chat, else it would be a group chat. When the user join the chat server, we can do the following operations to enable them to chat:

```python
# Pseudo-code in any language.
def join():
  # Perform handshake to check origin or valid ws endpoint.
  handshake() 
  
  # Perform authorization of the user. This could be done by passing a valid token that can validated by the chat server.
  # A successful authorization will return the id of the user.
  id = authz() 
  
  # Register the id to the pubsub, could be a external redis storage to make the server distributed.
  pubsub.subscribe(id)
  
  # Set the user's connection to the in-memory storage. This will tie the client to a specific chat-server, but messages could still be distributed across servers through redis pub-sub.
  localconn.set(id, conn)
  
  # Additionally, we might want to store the user's online status in an external storage. This would allow us to check if the user is online or not. The cache will expire after 1 minute, meaning the user has to consistently ping the redis storage to tell they are online.
  usercache.set(id, exp=1min)
  
  # Get the user's group from the database. Note that this can be done on the client-side, but we might need to validate if the client's request is valid, and hence need to be loaded on the server side. This can be cached to improve performance.
  groups = groupDB.get(id)
  
  # Read event loop. 
  while True:
    let msg = conn.read()
    
    # Store the message in persistent storage.
    snapshot(msg)
    
    match msg.type {
      case group:
        # Get all users in the group chat.
        ids = groups.get(msg.to)
        
        # The delivery can be concurrent.
        for id in ids:
           # We can perform additional checking to see if the user is online before sending the message.
           if usercache.get(id):
             pubsub.publish(id, msg)
      case single:
        # Send to external pub/sub.
        pubsub.publish(msg.to, msg)
    }
  
  # Another background task that runs periodically to set the user's status to online.
  now = time()
  while True:
    nextTimer(now + 55sec):
      now = now
      usercache.set(id, exp=now+1min)
  
  # Write event loop. Theoretically running in a separate thread to avoid blocking the main thread.
  while True:
    msg = pubsub.receive()
    # Note that the other user might not be online to receive the messages. This will was reads from the pub/sub.
    # Another alternative is to check if the user is online before triggering the pub/sub send.
    conn = localconn.get(msg.to)
    
    # User might not be online, or the connection might not be on this server.
    if conn:
      conn.write(msg)

  # Cleanup.
  # We have to remove the user at the end of the session. Note that this step might be skipped if an error happened, and the user might not be successfully unsubscribed. To ensure that the user is always removed from the session, ping em periodically.
  pubsub.unsubscribe(id)
  localconn.unset(id, conn)
  
  # We don't need to cleanup user cache. This would expire automatically. Also, this might provide additional buffer if the user just connect/disconnect frequently.
  # usercache.remove(id)
  delete group
```

The message request may look like this:
```js
{
  "from": "", // Excluded. Detected by chat server.
  "to": "group_id", // Send from client side.
  "msg": "hello world",
  "ts": 111, // Timestamp of sending. We probably need another timestamp on the server when chat server receive the message.
  "type": "group" // Or single
}
```

The message response may look like this:

```js
{
  "from": "", // Can be excluded, only the chat server should know who the `from` and `to` originate.
  "to": "", // Excluded. Only used in the chat server.
  "group": "Group A", // The name of the displayed group. Or the sender name.
  "participants": ["User A", "User B", "User C"],
  "msg": "hello world",
  "ts": 1111111, // Timestamp, shortform to save bytes.  
}
```

With this architecture, the user's may be distributed to different servers when they are online:

```
                                |-> Server 1: User A, User B
msg -> pub/sub -> RedisCluster -|
                                |-> Server 2: User C
```

Can the user exist on two different chat server, because their history is not cleaned properly? Possible. The outcome is the user might receive two same notifications.

### Client-side Architecture

The client-side architecture is slightly simpler. 

```
async function connect() {
  let token = await authzApi()
  let socket = new SocketAuthzConn(token)
  socket.on('connected', () => {
    // Fetch contact lists + online status
    let groups = await groupsApi()
    renderUI(groups)
  })
  
  // Send message whenever the user submits a new chat. Note that the client is responsible for detecting which group to send to, and the server is responsible to validate that the group is a valid group.
  input.on('submit', socket.send)
}
```

## Set Theory

User has a many-to-many relationship with rooms. A room can have multiple user, and a user can have multiple room. We can use a _joint function_ in mysql to describe this relationship. We can store the user-room pair in Redis too rather than storing it in the memory. We do so by simply creating two Redis sets, one for users, another for rooms. Consider the following user A:

```
user A: {room 1, room 2}
room 1: {user A, user B}
room 2: {user A, user C}
```

User A belongs to room 1 and room2. When user A join the chat server, user A will be automatically added into all the rooms, and the other user will be notified of the existence of user A (presence indicator). When user A exits/disconnect from the chat server, we will first remove user A from each room, notify the other party, and then remove user A. 

## Creating a new room

Users can create a new `room` or `group`, but only when they initiate a new conversation (must have at least one message). Else, it would be redundant to create a `group` as it will only consume unnecessary storage.

## Presence Indicator

To detect if the user is online, we can use redis key and set the status to be online every t duration, and set the key to expire at t + n duration. We have to utilise both `pull` and `push` method to detect the presence. `Push` model will notify other users when a user join or leave a group. However, they cannot notify users that are not online. Hence, a `pull` method is required to enquire the status of a particular user. We can have an `active` or `passive` `pull` model, by either pinging the server continuosly or just asking the server once when interacting with a particular user.

The status of the user should be stored in a distributed cache like redis, so that scaling is easier and the calls can be made by other means (API calls).

# Miscelleanous

- Each tab will create a new WebSocket connection. If we need to create only a single shared connection across different tabs, it is possible to do so with SharedWorker or BroadcastChannel API. 
- There are several ways to perform authentication, one is to use good ol' cookie, another is to use a ticker server.
- Checkout Server Side Events for read-only events 

## database design

How do we go about creating the schema to store the chat groups, as well as chat messages?

- create private chat
- create chat with groups
- store chat messages
- block chat request


Groups
- user id: creator of the group
- name: name of the group, default to names of the participant

Group participants
- group id (unique)
- user id (unique)
- blocked: null bool, whether they block the conversation or not


Group participants private chat view
- group by group id, only two person


Group messages
- user id
- text
- type: text, media, photo, video, audio

To view chats,

1) list the available chat groups if exists
2) otherwise, list your friends

To chat with someone (private)

1) select the person to chat with
2) check if there exists a group with the individuals
3) create group if not exists, add participants
4) send message to the group, but exclude yourself

To chat with a group
1) select the list of person to add to the group
2) the rest is the same as above


## Redis

We use Redis to store the nodes (backend server ) where the user is logged in to.

When a user opens a chat, we want to track
- which device they are on (each device is a different id)
- which server they are on (the pub sub should best send to only the particular server the user is on)

The naive approach is to just broadcast to all server listening to one namespace

## References
- https://www.thepolyglotdeveloper.com/2016/12/create-real-time-chat-app-golang-angular-2-websockets/
- https://devcenter.heroku.com/articles/go-websockets
- https://www.jonathan-petitcolas.com/2015/01/27/playing-with-websockets-in-go.html
- https://blog.arnellebalane.com/sending-data-across-different-browser-tabs-6225daac93ec
