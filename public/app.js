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


  const ChatApp = (function () {
    const map = new WeakMap()

    const internal = function (key) {
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
          rooms: new Map(),

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
        state.user = user
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
            .map(({ roomId }) => roomId)
            .map(roomId => fetchConversations(roomId, token))

          const conversations = await Promise.all(promises)
          this.conversations = conversations
        }

        socket.onmessage = (evt) => {
          const state = internal(this).state
          const { $rooms, rooms } = state
          try {
            const msg = JSON.parse(evt.data)
            // console.log('got message', msg)
            switch (msg.type) {
              case 'typing':
              {
                const { roomTimeouts, conversations } = internal(this).state
                const roomId = msg.room
                this.updateRoom(roomId, {
                  message: `...${msg.data} is typing`
                })
                roomTimeouts[roomId] && window.clearTimeout(roomTimeouts[roomId])
                roomTimeouts[roomId] = window.setTimeout(() => {
                  // This ensure that the last conversations will always be retrieved.
                  const convs = conversations.get(roomId)
                  const last = convs[convs.length - 1]
                  this.updateRoom(msg.room, {
                    message: last.text
                  })
                }, 1000)
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
                this.updateRoom(msg.room, {
                  status: msg.data === '1',
                  timestamp: new Date()
                })
                break
              }
              case 'presence':
              {
                // updateRoom.
                // Controller: Update the data first, then update the view.
                this.updateRoom(msg.room, {
                  status: msg.data === '1',
                  timestamp: new Date()
                })
                break
              }
              case 'message':
                {
                  const { room, user } = internal(this).state
                  if (room === msg.room) {
                    const isSelf = user === msg.from
                    const $dialogs = this.shadowRoot.querySelector('.dialogs')
                    const $dialog = document.createElement('chat-dialog')
                    $dialog.isSelf = isSelf
                    $dialog.message = msg.data
                    $dialogs.appendChild($dialog)
                  }

                  // Add new conversation.
                  let conversations = internal(this).state.conversations.get(msg.room)
                  let isNew = false
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
                  this.updateRoom(msg.room, {
                    message: msg.data,
                    timestamp: msg.created_at
                  })
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
          const state = internal(this).state
          const { isTyping, room, chattingWith } = state
          if (isTyping) {
            return
          }
          send({
            type: 'typing',
            room,
            data: state.user 
          })
          state.isTyping = true
          window.clearTimeout(state.isTypingInterval)
          state.isTypingInterval = window.setTimeout(() => {
            state.isTyping = false
          }, 1000)
        })
      }

      updateRoom(roomId = '', nextState = {}) {
        const state = internal(this).state
        const {rooms, $rooms} = state
        const room = rooms.get(roomId)
        Object.assign(room, nextState)

        // Update the view.
        const $room = $rooms.get(room)
        $room.timestamp = room.timestamp
        $room.status = room.status
        $room.message = room.message
        $room.selected = room.selected
      }

      set conversations (items) {
        const state = internal(this).state
        const { conversations, room } = state
        for (let { data, room } of items) {
          // Set the conversations.
          conversations.set(room, data)

          // Display the last message for each room.
          if (!data) {
            continue
          }
          const [head] = data
          this.updateRoom(room, {
            message: head.text,
            timestamp: head.created_at
          })
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

        const { user } = internal(this).state
        conversations.forEach((conversation) => {
          const isSelf = user === conversation.user_id
          const $dialog = document.createElement('chat-dialog')
          $dialog.isSelf = isSelf
          $dialog.message = conversation.text
          $dialogs.appendChild($dialog)
        })
      }

      // Clears the room data, and the view associated with it.
      deleteRoom (roomId) {
        const { rooms, $rooms } = internal(this).state
        const room = rooms.get(roomId)
        $rooms.remove(room)
        rooms.delete(roomId)
      }

      newRoomView (room) {
        const $room = document.createElement('chat-room')
        $room.user = room.name
        $room.userId = room.userId
        $room.room = room.roomId
        $room.timestamp = new Date().toISOString()
        return $room
      }

      onChangeRoomsState (newRooms, onChangeFn) {
        const state = internal(this).state
        const { rooms, $rooms } = state

        // Perform diffing on the rooms.
        const prevState = rooms
        const nextState = new Set(newRooms.map(room => room.roomId))

        for (let prev of prevState.keys()) {
          if (!nextState.has(prev)) {
            const room = rooms.get(prev)
            const $room = $rooms.get(room) 
            // rooms.delete(prev)
            this.deleteRoom(prev)
            onChangeFn && onChangeFn({type:'delete', room, $room})
          }
        }

        for (let room of newRooms) {
          if (!prevState.has(room.roomId)) {
            // Set into collection.
            rooms.set(room.roomId, room)
            onChangeFn && onChangeFn({type:'add', room })
          }
        }
      }

      set rooms (newRooms) {
        const state = internal(this).state
        const { conversations, rooms, socket, $rooms: $roomsView } = state
        const $rooms = this.shadowRoot.querySelector('.rooms')

        this.onChangeRoomsState(newRooms, (({ type, room, $room }) => {
          switch (type) {
            case 'delete':
              $rooms.remove($room)
              break
            case 'add':
              const $room = this.newRoomView(room)
              $room.addEventListener('select-group', (evt) => {
                const prevRoom = state.room
                state.room = evt.detail.room()
                state.chattingWith = evt.detail.user()

                // Render the conversations.
                this.renderDialogs(conversations.get(state.room))

                // Update room view.
                this.updateRoom(prevRoom, {selected: false})
                this.updateRoom(state.room, {selected: true})
              })

              $roomsView.set(room, $room)

              // For each user in the room, request the current status (online/offline).
              socket.send(JSON.stringify({
                type: 'status',
                data: `${room.userId}`,
                room: `${room.roomId}`
              }))
              $rooms.appendChild($room)
              break
          }
        }))
        // Select the first room as the main room for chatting.
        if (rooms.size) {
          const [room] = [...rooms.values()]
          this.updateRoom(room.roomId, { selected: true })

          state.room = room.roomId
          state.chattingWith = room.userId
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

  class Room {
    constructor(roomId, userId, name) {
      // From API.
      this.roomId = roomId
      this.userId = userId
      this.name = name

      // true=online, false=offline
      this.status = false 
      this.timestamp = new Date()
      // message to display for the room.
      this.message = ''
      // state of the selected room for conversations.
      this.selected = false
    }
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
    return data.map(({ user_id, room_id, name }) => new Room(room_id, user_id, name))
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
