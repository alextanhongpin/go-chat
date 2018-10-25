(function() {
	let template = document.createElement('template')
	template.innerHTML = `
		<style>
		
		</style>
		<div class='rooms'>
			...loading	
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
				socket: null
			}
		}

		async connectedCallback() {
			// if (!this.hasAttribute('socket_uri'))
      //   console.error('attribute "socket_uri" required')
			//   return

			// Perform authentication to obtain a token first
			// before connecting to the websocket server.
			let socketUri = this.getAttribute('socket_uri')
			let user = window.prompt('enter username')
			let token = await authenticate(user)
      console.log('authenticated', user)
			let socket = new window.WebSocket(`${socketUri}?token=${token}`)
			let send = (msg) => {
				socket.send(JSON.stringify(msg))
			}

			socket.onopen = async () => {
				let rooms = await fetchRooms(user)
				this.rooms = rooms
				console.log('socket opened')
			}

			socket.onmessage = (evt) => {
				try {
					let msg= JSON.parse(evt.data)
					console.log('receive message:', msg)
					switch (msg.type) {
						case 'status':
						{
							let [room] = this.state.rooms.filter(room => room.room_id === parseInt(msg.room, 10))
							if (room && this.state.roomsCache.has(room)) {
								let $room = this.state.roomsCache.get(room)
								$room.status = msg.data === '1'
								$room.timestamp = new Date()
							}
						}
						case 'presence':
						{
							let [room] = this.state.rooms.filter(room => room.room_id === parseInt(msg.room, 10))
							if (room && this.state.roomsCache.has(room)) {
								let $room = this.state.roomsCache.get(room)
								$room.status = msg.data === '1'
								$room.timestamp = new Date()
							}
						}
						case 'message':
						default:
							console.log(msg)
					}
				} catch (error) {
					console.error(error)
				}
			}

			this.socket = socket
			this.user = user
		}

		set user (value) {
			this.state.user = value
		}

		set socket(value) {
			this.state.socket = value
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
          console.log('removed')
          this.state.roomsCache.remove(room)
        } else {
          console.log('existed')
          nextState.push(room)
        }
      }
      for (let room of currState) {
        if (!prevStateSet.has(room)) {
          console.log('added')
          nextState.push(room)
          
          // Create a new element.
          const $room = document.createElement('chat-room')
          // $room.user = room.user_id
          $room.user = room.name
          $room.room = room.room_id
          $room.timestamp = new Date().toISOString()
          $rooms.appendChild($room)
          this.state.roomsCache.set(room, $room)
        }
      }

      nextState.forEach(({ user_id: user, room_id: room })=> {
        console.log('enquiring status', user)
        return this.state.socket.send(JSON.stringify({
          type: 'status',
          data: `${user}`,
          room: `${room}`
        }))
      })

      this.state.rooms = nextState

			// console.log('found rooms', $rooms, rooms)
      //
			// let set = new WeakSet()
			// // Add all the new room.
			// rooms.forEach(room => {
			//   let exist = this.state.roomsCache.has(room)
			//   let $room = exist
			//     ? this.state.roomsCache.get(room)
			//     : document.createElement('chat-room')
			//   $room.user = room.user_id
			//   $room.room = room.room_id
			//
			//   !exist && $rooms.appendChild($room)
			//   !exist && this.state.roomsCache.set(room, $room)
			//   set.add(room)
			// })
			// // For each old room, remove them from the view.
			// prevState.forEach(room => {
			//   if (!set.has(room) && this.state.roomsCache.has(room)) {
			//     const $room = this.state.roomsCache.get(room)
			//     $room.remove()
			//   }
			// })
			// set = null
      //
		}

		attributeChangedCallback(attrName, oldValue, newValue) {
			switch (attrName) {
				case 'key':
					break
			}
		}

		render () {
			console.log('rendering components')
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

	async function fetchRooms(user) {
		let mapUserToId = {
			john: 1,
			jane: 2
		}
		let user_id = mapUserToId[user]
		const response = await window.fetch(`/rooms?user_id=${user_id}`)
		const { data } = await response.json()
		return data || []
	}
})()
