(function () {
  let template = document.createElement('template')
  template.innerHTML = `
		<style>
      :host {
        contain: content;
        all: initial;
        font-family: Avenir, arial;
      }
      .app {
        display: grid;
        grid-template-columns: 320px 1fr;
        grid-column-gap: 10px;
      }
      .chat {
        width: 320px;
        height: 480px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, .15);
        border-radius: 7px;
        overflow: hidden;

        display: grid;
        grid-template-rows: 30px 1fr 40px;
        grid-template-columns: 1fr;
        justify-content: space-around;
      }
      .header {
        background: #4488FF;
        min-height: 30px;
        line-height: 30px;
        color: white;
        padding: 0 10px;
      }
      .footer {
        min-height: 40px;
        line-height: 40px;
      }
      .search {
        width: 100%;
        border: 1px solid #DDDDDD;
        border-radius: 0 0 7px 7px;
        -webkit-appearance: none;
        height: 40px;
        padding: 0 10px;
      }
      .dialog.is-active {
        background: #4488FF;
        color: white;
      }
		
		</style>
    <div class='app'>
      <div class='chat'>
        <div class='header'>
          chat 
        </div>
        <div class='rooms'></div>
        <div class='footer'>
          <input class='search' type='text' placeholder='Search user'/>
        </div>
      </div>
      <div>
        <div class='dialogs'>
          <div class='placeholder'>No messages yet</div> 
        </div>
        <input class='input-message' type='text' placeholder='Enter message'/>
        <button class='send'>Send</button>
      </div>
    </div>
	`

  const ChatApp = (function() {
    const map = new WeakMap()

    const internal = function(key) {
      if (!map.has(key)) {
        map.set(key, {})
      }
      return map.get(key)
    }

    class ChatApp extends HTMLElement {
      constructor () {
        super()
        this.attachShadow({ mode: 'open' })
          .appendChild(template.content.cloneNode(true))

        internal(this).state = {
          connected: false,
          // The current user.
          user: '',

          // The rooms the user is in.
          rooms: new WeakMap(),

          // The View of the rooms.
          $rooms: new WeakMap(),

          // The websocket connection.
          socket: null,

          // State of isTyping.
          isTyping: false,
          
          // Throttle for typing.
          isTypingInterval: null,

          // Current chat conversation.
          chattingWith: '',

          // The current room.
          room: null, 

          // Hold all the current conversation.
          conversations: new Map(),
          
          // Timeouts for each room.
          roomTimeouts: {}
        }
      }

      async connectedCallback () {
        if (!this.hasAttribute('socket_uri')) {
          throw new Error('attribute "socket_uri" required')
          return
        }

        // Get application state.
        const state = internal(this).state

        const socketUri = this.getAttribute('socket_uri')

        // Ask for username for identification.
        // TODO: Remove this after replacing with login.
        const user = window.prompt('enter username')
        if (!user.trim().length) {
          throw new Error('username is required')
        }

        // Handshake.
        const token = await authenticate(user)

        // Connect WebSocket. 
        const socket = new window.WebSocket(`${socketUri}?token=${token}`)
        state.socket = socket

        // Utility method for sending JSON through websocket.
        const send = (msg) => socket.send(JSON.stringify(msg))

        socket.onopen = async () => {
          // Request for the user id.
          send({ type: 'auth' })

          // Fetch rooms for user.
          const rooms = await fetchRooms(user, token)
          this.rooms = rooms 

          // For each room, fetch the last 10 conversations.
          const promises = rooms 
            .map(({ room_id }) => room_id)
            .map(roomId => fetchConversations(roomId, token))

          const conversations = await Promise.all(promises)
          this.conversations = conversations
          console.log('got conversations', conversations)
          console.log(internal(this).state)
        }

        socket.onmessage = (evt) => {
          const state = internal(this).state
          const {$rooms, rooms} = state
          try {
            const msg = JSON.parse(evt.data)
            console.log('got message', msg)
            switch (msg.type) {
              case 'typing':
              {
                if (!rooms.has(msg.room)) {
                  return
                }
                const roomId = msg.room
                const room = rooms.get(roomId)
                const $room = $rooms.get(room)
                if (!$room) {
                  return
                }

                // Update view...
                $room.message = `...${room.name} is typing`

                internal(this).state.roomTimeouts[roomId] && window.clearTimeout(internal(this).state.roomTimeouts[roomId])
                internal(this).state.roomTimeouts[roomId] = window.setTimeout(() => {
                  // This ensure that the last conversations will always be retrieved.
                  const conversations = internal(this).state.conversations.get(roomId)
                  const last = conversations[conversations.length - 1]

                  $room.message = last.text
                }, 2000)
                break
              }
              case 'auth':
              {
                state.user = msg.data
                break
              }
              // Checks the current status (online/offline) of the user. 
              case 'status':
              {
                // Bad. We need to update the state first, then only update the
                // view. 
                const $room = this.getRoomView(msg.room)
                if (!$room) return 
                $room.status = msg.data === '1'
                $room.timestamp = new Date()
                break
              }
              case 'presence':
              {
                const $room = this.getRoomView(msg.room)
                if (!$room) return 
                $room.status = msg.data === '1'
                $room.timestamp = new Date()
                break
              }
              case 'message':
                {
                  const { 
                    room, 
                    user,
                    rooms,
                    $rooms
                  } = internal(this).state
                  if (room === msg.room) {
                    const isSelf = user === msg.from
                    const $dialogs = this.shadowRoot.querySelector('.dialogs')
                    const $dialog = document.createElement('chat-dialog')
                    $dialog.isSelf = isSelf
                    $dialog.message = msg.data
                    $dialogs.appendChild($dialog)
                  }

                  // Add new conversation.
                  const conversations = internal(this).state.conversations.get(msg.room)
                  const isNew = false
                  if (!conversations) {
                    isNew = true
                    conversations = []
                    internal(this).state.conversations.set(msg.room, conversations)
                    return
                  }

                  const newMessage = {
                    user_id: msg.from,
                    text: msg.data,
                    created_at: new Date()
                  }
                  conversations.push(newMessage)
                  if (isNew) {
                    this.renderDialogs(conversations)
                  }
                  // Update last message for the room.
                  const $room = this.getRoomView(msg.room)
                  $room.message = msg.data
                  $room.timestamp = msg.created_at
                }
                break
              default:
            }
          } catch (error) {
            console.error(error)
          }
        }

        this.shadowRoot.querySelector('.send').addEventListener('click', (evt) => {
          let $input = this.shadowRoot.querySelector('.input-message')
          if (!$input.value.trim().length) {
            return
          }
          send({
            room: `${internal(this).state.room}`,
            data: $input.value,
            type: 'message'
          })
          $input.value = ''
        })

        this.shadowRoot.querySelector('.input-message').addEventListener('keyup', (evt) => {
          const state = internal(this)
          const {isTyping, room, chattingWith} = state
          if (isTyping) {
            return
          }
          send({
            type: 'typing',
            room,
            data: chattingWith
          })
          state.isTyping = true
          window.clearTimeout(state.isTyping)
          state.isTypingInterval = window.setTimeout(() => {
            state.isTyping = false
          }, 1000)
        })
      }


      set conversations (items) {
        const state = internal(this).state
        const { rooms, $rooms, conversations, room } = state
        for (let { data, room: roomId } of items) {
          // Set the conversations.
          conversations.set(roomId, data)

          // Display the last message for each room.
          const $room = this.getRoomView(roomId)
          if (!data || !$room) {
            continue
          }
          const [head] = data
          $room.message = head.text
          $room.timestamp = head.created_at
        }

        // Render the first conversation.
        if (conversations.size) {
          const data = conversations.get(room)
          this.renderDialogs(data)
        }
      }
      renderDialogs (conversations = []) {
        const $dialogs = this.shadowRoot.querySelector('.dialogs')
        // Reset the view.
        $dialogs.innerHTML = ''

        // Sort in ascending order. The newest message will be last.
        conversations.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())

        const {user} = internal(this).state
        conversations.forEach((conversation) => {
          const isSelf = user === conversation.user_id
          const $dialog = document.createElement('chat-dialog')
          $dialog.isSelf = isSelf
          $dialog.message = conversation.text
          $dialogs.appendChild($dialog)
        })
      }

      // Clears the room data, and the view associated with it.
      deleteRoom(roomId) {
        const state = internal(this).state
        const room = state.rooms.get(roomId)
        // state.rooms will aut
        state.$rooms.remove(room)
      }
      getRoom(roomId) {
        return internal(this).state.rooms.get(roomId)
      }
      getRoomView(roomId) {
        const state = internal(this).state
        const room = state.rooms.get(roomId)
        return state.$rooms.get(room)
      }
      newRoomView(room) {
            const $room = document.createElement('chat-room')
            $room.user = room.name
            $room.userId = room.user_id
            $room.room = room.room_id
            $room.timestamp = new Date().toISOString()
            return $room
      }
      
      set rooms (newRooms) {
        const state = internal(this).state
        const {rooms, conversations, socket, $rooms: $roomsView} = state
        const $rooms = this.shadowRoot.querySelector('.rooms')

        // Perform diffing on the rooms.
        const prevState = rooms
        const nextState = new Set(newRooms.map(room => room.room_id))

        for (let prev of prevState.keys()) {
          if (!nextState.has(prev)) {
            this.deleteRoom(prev)
          }
        }

        for (let room of newRooms) {
          if (!prevState.has(room.room_id)) {
            // Set into collection.
            rooms.set(room.room_id, room)

            // Create a new room view.
            const $room = this.newRoomView(room)
            $room.addEventListener('select-group', (evt) => {
              const prevRoom = state.room
              console.log('select group', evt.detail.room(), evt.detail.user())
              state.room = evt.detail.room()
              state.chattingWith = evt.detail.user()

              // Render the conversations.
              this.renderDialogs(conversations.get(state.room))
              // Update room view.
              const $prevRoom = this.getRoomView(prevRoom)
              $prevRoom.selected = false
              const $currRoom = this.getRoomView(state.room)
              $currRoom.selected = true
            })

            $roomsView.set(room, $room)

            // For each user in the room, request the current status (online/offline).
            socket.send(JSON.stringify({
              type: 'status',
              data: `${room.user_id}`,
              room: `${room.room_id}`
            }))
            $rooms.appendChild($room)
          }
        }

        // Select the first room as the main room for chatting.
        if (rooms.size) {
          const [room] = [...rooms.values()]
          const $room = $roomsView.get(room)
          $room.selected = true
          state.room = room.room_id
          state.chattingWith = room.user_id
        }
      }
    }
    return ChatApp
  })()

  window.customElements.define('chat-app', ChatApp)

  async function authenticate (user) {
    const response = await window.fetch('/auth', {
      method: 'POST',
      body: JSON.stringify({ user_id: user })
    })
    if (!response.ok) {
      console.error(await response.text())
      return
    }
    const { token } = await response.json()
    return token
  }

  async function fetchRooms (user, token) {
    const response = await window.fetch('/rooms', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) {
      const msg = await response.text()
      console.error(msg)
      return []
    }
    const { data } = await response.json()
    return data || []
  }

  async function fetchConversations (room, token) {
    const response = await window.fetch(`/conversations/${room}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) {
      const msg = await response.text()
      console.error(msg)
      return []
    }
    return response.json()
  }
})()
