async function main () {

  const user = window.prompt('What is your username?')
  const body = await window.fetch('/auth', {
    method: 'POST',
    body: JSON.stringify({
      user_id: user
    })
  })

  const response = await body.json()
  if (response) {
    console.log(`user ${user} is authenticated`)
    const { token } = response

    const socket = new window.WebSocket(`ws://localhost:4000/ws?token=${token}`)
    socket.onopen = async () => {

      const userId = user === 'john' ? 1 : 2
      let body = await window.fetch(`/rooms?user_id=${userId}`)
      let response = await body.json()
      response.data.forEach(({ room_id: room, user_id }) => {

        // Ask for the room status.
        send({
          type: 'status',
          data: user_id,
          room: room
        })

        let chatRow = document.createElement('chat-row')
        chatRow.value = room
        chatRow.addEventListener('onmessage', function (evt) {
          const payload = {
            type: 'message',
            user,
            data: evt.detail.message(), 
            room: evt.currentTarget.value 
          }
          console.log('sending:', payload)
          this.conversations = this.conversations.concat([payload])
          console.log(this.conversations)
          send(payload)
        })
        $('chat').appendChild(chatRow)
        // Check for presence in room.
      })
      // socket.send(JSON.stringify({
      //   type: 'handshake'
      // }))
    }
    socket.onmessage = (event) => {
      try {
        let message = JSON.parse(event.data)
        console.log('got message', message)
        let { data, room, type } = message
        let el = document.querySelector(`chat-row[value="${room}"]`)
        el && (el.status = data)
        // el.status = data
      } catch (error) {
        console.error(error)
      }
    }

    function send(msg) {
      console.log('sending', msg)
      socket.send(JSON.stringify(msg))
    }

    $('submit').addEventListener('click', (evt) => {
      send({
        type: 'message',
        user,
        data: $('message').value,
        room: $('room').value
      })
      $('message').value = '' 
      $('room').value = '' 
    })
  }
}

function $ (el) {
  return document.getElementById(el)
}

main().catch(console.error)
