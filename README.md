https://www.thepolyglotdeveloper.com/2016/12/create-real-time-chat-app-golang-angular-2-websockets/

https://devcenter.heroku.com/articles/go-websockets

https://www.jonathan-petitcolas.com/2015/01/27/playing-with-websockets-in-go.html

# go-chat

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

## Handling relationships

In chat, we have concept of friends. Basically, a user will have a list of friends he/she can chat with. When the chat is initialized, the list of friends would be added, so that the users can choose who to chat with. Communications needs to be established between the two users in order for them to send messages to one another.

A naive approach is to simply load all the user's friends when initializing the chat and attempting to establish connections with them. This might not be ideal when the users have a lot of friends, as they would only chat with a select few at a time. Hence, there are some limits that needs to be establish here. One way is to rank the friends by the last chat, and adding only the last 15 friends into the chat list (offline and online). If the user wishes to chat with someone else, he/she can search for the name and start a new chat.
