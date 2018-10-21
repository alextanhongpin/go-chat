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

    // socket.onopen = () => {
    //   socket.send(JSON.stringify({
    //     type: 'handshake'
    //   }))
    // }
    socket.onmessage = (event) => {
      console.log(JSON.parse(event.data))
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
