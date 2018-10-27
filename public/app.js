(function() {
  let mapUserToId = {
    john: 1,
    jane: 2
  }
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
		constructor() {
			super()
			this.attachShadow({ mode: 'open' })
				.appendChild(template.content.cloneNode(true))

			this.state = {
				connected: false,
				user: '',
				rooms: [],
				roomsCache: new WeakMap(),
				socket: null,
        room: null, // The selected room
        isTyping: false,
        isTypingInterval: null,
        chattingWith: ''
			}
		}

		async connectedCallback() {
      if (!this.hasAttribute('socket_uri')) {
        console.error('attribute "socket_uri" required')
        return
      }

			// Perform authentication to obtain a token first
			// before connecting to the websocket server.
			let socketUri = this.getAttribute('socket_uri')
			let user = window.prompt('enter username')
			let token = await authenticate(user)
			let socket = new window.WebSocket(`${socketUri}?token=${token}`)
			let send = (msg) => {
				socket.send(JSON.stringify(msg))
			}

			this.socket = socket

			socket.onopen = async () => {
        send({ type: 'auth' })

				let rooms = await fetchRooms(user, token)
				this.rooms = rooms
			}

			socket.onmessage = (evt) => {
				try {
					let msg = JSON.parse(evt.data)
					switch (msg.type) {
            case 'is_typing':
						{
							let [room] = this.state.rooms.filter(room => room.room_id === msg.room)
							if (room && this.state.roomsCache.has(room)) {
								let $room = this.state.roomsCache.get(room)
                $room.message = 'x is typing' 
                window.setTimeout(() => {
                  $room.message = ''
                }, 2000)
							}
              break
						}
            case 'auth':
              {
                this.user = msg.data
                break
              }
						case 'status':
						{
							let [room] = this.state.rooms.filter(room => room.room_id === msg.room)
							if (room && this.state.roomsCache.has(room)) {
								let $room = this.state.roomsCache.get(room)
								$room.status = msg.data === '1'
								$room.timestamp = new Date()
							}
              break
						}
						case 'presence':
						{
							let [room] = this.state.rooms.filter(room => room.room_id === msg.room)
							if (room && this.state.roomsCache.has(room)) {
								let $room = this.state.roomsCache.get(room)
								$room.status = msg.data === '1'
								$room.timestamp = new Date()
							}
              break
						}
						case 'message':
              if (this.state.room === msg.room) {
                const isSelf = this.state.user === msg.from 
                let $dialogs = this.shadowRoot.querySelector('.dialogs')
                let $dialog = document.createElement('chat-dialog')
                $dialog.isSelf = isSelf
                $dialog.message = msg.data
                $dialogs.appendChild($dialog)
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
          type: 'message',
        })
        $input.value = ''
      })

      this.shadowRoot.querySelector('.input-message').addEventListener('keyup', (evt) => {
        if (this.state.isTyping) {
          return
        }
        send({
          type: 'is_typing',
          room: this.state.room,
          to: this.state.chattingWith
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

		set socket(value) {
			this.state.socket = value
		}

    set room(value) {
      this.state.room = value
    }

		set rooms(rooms) {
			// Diff!
			let $rooms = this.shadowRoot.querySelector('.rooms')

			let prevState = this.state.rooms
      let prevStateSet = new Set(prevState.filter(item => item.room_id))

      let currState = rooms
      let currStateSet = new Set(currState.filter(item => item.room_id))

      let nextState = []
      for (let room of prevState) {
        if (!currStateSet.has(room)) {
          this.state.roomsCache.remove(room)
        } else {
          nextState.push(room)
        }
      }
      for (let room of currState) {
        if (!prevStateSet.has(room)) {
          nextState.push(room)
          
          // Create a new element.
          const $room = document.createElement('chat-room')
          // $room.user = room.user_id
          $room.user = room.name
          $room.room = room.room_id
          $room.selected = nextState.length === 1
          if (nextState.length === 1) {
            this.room = room.room_id
            this.state.chattingWith = room.user_id
          }
          $room.timestamp = new Date().toISOString()
          $rooms.appendChild($room)
          this.state.roomsCache.set(room, $room)
        }
      }

      nextState.forEach(({ user_id: user, room_id: room })=> {
        return this.state.socket.send(JSON.stringify({
          type: 'status',
          data: `${user}`,
          room: `${room}`
        }))
      })

      this.state.rooms = nextState
		}

		attributeChangedCallback(attrName, oldValue, newValue) {
			switch (attrName) {
				case 'key':
					break
			}
		}

		render () {
		}
	}

	window.customElements.define('chat-app', ChatApp)

	async function authenticate(user_id) {
		const response = await window.fetch('/auth', {
			method: 'POST',
			body: JSON.stringify({ user_id })
		})
		const { token } = await response.json()
		return token
	}

	async function fetchRooms(user, token) {
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
})()
