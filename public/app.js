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

  class ChatApp extends HTMLElement {
    constructor () {
      super()
      this.attachShadow({ mode: 'open' })
        .appendChild(template.content.cloneNode(true))

      this.state = {
        connected: false,
        user: '',
        // rooms: [],
        // roomsCache: new WeakMap(),
        rooms: new Map(),
        $rooms: new WeakMap(),
        socket: null,
        room: null, // The selected room
        isTyping: false,
        isTypingInterval: null,
        chattingWith: '',
        conversations: new Map(),
        roomTimeouts: {}
      }
    }

    async connectedCallback () {
      if (!this.hasAttribute('socket_uri')) {
        console.error('attribute "socket_uri" required')
        return
      }

      let socketUri = this.getAttribute('socket_uri')

      // Ask for username for identification.
      let user = window.prompt('enter username')
      let token = await authenticate(user)

      // Initialize websocket connection.
      let socket = new window.WebSocket(`${socketUri}?token=${token}`)

      // Utility method for sending JSON through websocket.
      let send = (msg) => {
        socket.send(JSON.stringify(msg))
      }

      this.socket = socket

      socket.onopen = async () => {
        // Request for the user id.
        send({ type: 'auth' })

        // Fetch rooms for user.
        let rooms = await fetchRooms(user, token)
        this.rooms = rooms

        // For each room, fetch the last 10 conversations.
        let promises = rooms
          .map(({ room_id }) => room_id)
          .map(roomId => fetchConversations(roomId, token))
        let conversations = await Promise.all(promises)
        // let results = conversations.reduce((acc, { room, data }) => {
        //   acc[room] = data
        //   return acc
        // }, {})
        this.conversations = conversations
      }

      socket.onmessage = (evt) => {
        try {
          let msg = JSON.parse(evt.data)
          switch (msg.type) {
            case 'typing':
            {
              if (!this.state.rooms.has(msg.room)) {
                return
              }
              let roomId = msg.room
              let room = this.state.rooms.get(roomId)
              let $room = this.state.$rooms.get(room)
              if (!$room) {
                return
              }

              // This won't retrieve the last conversation.
              // let prevMessage = $room.message

              $room.message = `...${room.name} is typing`
              this.state.roomTimeouts[roomId] && window.clearTimeout(this.state.roomTimeouts[roomId])
              this.state.roomTimeouts[roomId] = window.setTimeout(() => {
                // This ensure that the last conversations will always be retrieved.
                let conversations = this.state.conversations.get(roomId)
                let last = conversations[conversations.length - 1]
                $room.message = last.text
              }, 2000)
              break
            }
            case 'auth':
            {
              this.user = msg.data
              break
            }
            case 'status':
            {
              let roomId = msg.room
              let room = this.state.rooms.get(roomId)
              let $room = this.state.$rooms.get(room)
              if (!$room) {
                return
              }
              $room.status = msg.data === '1'
              $room.timestamp = new Date()
              break
            }
            case 'presence':
            {
              let roomId = msg.room
              let room = this.state.rooms.get(roomId)
              let $room = this.state.$rooms.get(room)
              if (!$room) {
                return
              }
              $room.status = msg.data === '1'
              $room.timestamp = new Date()
              break
            }
            case 'message':
              {
                if (this.room === msg.room) {
                  const isSelf = this.state.user === msg.from
                  let $dialogs = this.shadowRoot.querySelector('.dialogs')
                  let $dialog = document.createElement('chat-dialog')
                  $dialog.isSelf = isSelf
                  $dialog.message = msg.data
                  $dialogs.appendChild($dialog)
                }

                // Add new conversation.
                let conversations = this.state.conversations.get(msg.room)
                let isNew = false
                if (!conversations) {
                  isNew = true
                  conversations = []
                  this.state.conversations.set(msg.room, conversations)

                  return
                }

                let newMessage = {
                  user_id: msg.from,
                  text: msg.data,
                  created_at: new Date()
                }
                conversations.push(newMessage)
                if (isNew) {
                  this.renderDialogs(conversations)
                }
                // Update last message for the room.
                let room = this.state.rooms.get(msg.room)
                let $room = this.state.$rooms.get(room)
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
          room: `${this.state.room}`,
          data: $input.value,
          type: 'message'
        })
        $input.value = ''
      })

      this.shadowRoot.querySelector('.input-message').addEventListener('keyup', (evt) => {
        if (this.state.isTyping) {
          return
        }
        send({
          type: 'typing',
          room: this.state.room,
          data: this.state.chattingWith
        })
        this.state.isTyping = true
        window.clearTimeout(this.state.isTyping)
        this.state.isTypingInterval = window.setTimeout(() => {
          this.state.isTyping = false
        }, 1000)
      })
    }

    set user (value) {
      this.state.user = value
    }

    set socket (value) {
      this.state.socket = value
    }

    set room (value) {
      this.state.room = value
    }

    get room () {
      return this.state.room
    }

    set conversations (items) {
      for (let { data, room: roomId } of items) {
        // Set the conversations.
        this.state.conversations.set(roomId, data)

        // Display the last message for each room.
        let room = this.state.rooms.get(roomId)
        let $room = this.state.$rooms.get(room)
        if (!data) {
          continue
        }
        let [head] = data
        $room.message = head.text
        $room.timestamp = head.created_at
      }

      // Render the first conversation.
      if (this.state.conversations.size) {
        let data = this.state.conversations.get(this.state.room)
        this.renderDialogs(data)
      }
    }
    renderDialogs (conversations = []) {
      let $dialogs = this.shadowRoot.querySelector('.dialogs')
      // Reset the view.
      $dialogs.innerHTML = ''

      // Sort in ascending order. The newest message will be last.
      conversations.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
      conversations.forEach((conversation) => {
        let isSelf = this.state.user === conversation.user_id
        let $dialog = document.createElement('chat-dialog')
        $dialog.isSelf = isSelf
        $dialog.message = conversation.text
        $dialogs.appendChild($dialog)
      })
    }
    set rooms (rooms) {
      let $rooms = this.shadowRoot.querySelector('.rooms')

      let prevState = this.state.rooms
      let nextState = new Set(rooms.map(room => room.room_id))

      for (let prev of prevState.keys()) {
        if (!nextState.has(prev)) {
          this.state.rooms.remove(prev)
          this.state.$rooms.remove(prevState.get(prev))
        }
      }

      for (let room of rooms) {
        if (!prevState.has(room.room_id)) {
          const $room = document.createElement('chat-room')
          $room.user = room.name
          $room.userId = room.user_id
          $room.room = room.room_id
          $room.timestamp = new Date().toISOString()

          $room.addEventListener('select-group', (evt) => {
            this.room = evt.detail.room()
            this.state.chattingWith = evt.detail.user()

            this.renderDialogs(this.state.conversations.get(this.room))
          })

          this.state.rooms.set(room.room_id, room)
          this.state.$rooms.set(room, $room)
          this.state.socket.send(JSON.stringify({
            type: 'status',
            data: `${room.user_id}`,
            room: `${room.room_id}`
          }))
          $rooms.appendChild($room)
        }
      }

      // Select the first room as the main room for chatting.
      if (this.state.rooms.size) {
        let [room] = [...this.state.rooms.values()]
        let $room = this.state.$rooms.get(room)
        $room.selected = true
        this.room = room.room_id
        this.state.chattingWith = room.user_id
      }
    }
  }

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
    // const { data, room } = await response.json()
    // return data || []
    return response.json()
  }
})()
