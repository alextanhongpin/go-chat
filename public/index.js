
(function () {
  const socket = new window.WebSocket('ws://localhost:8080/ws')

  socket.onopen = onOpen(socket)
  socket.onmessage = onMessage(socket)
})()

function onOpen (socket) {
  return function (event) {
    // const msg = {
    //   type: 'authenticate',
    //   payload: { token: 'xxx' }
    // }
    const msg = {
      handle: 'hello socket',
      text: 'this is a new text message',
      room: '123'
    }
    socket.send(JSON.stringify(msg))
  }
}
function onMessage (socket) {
  return function (event) {
    console.log('got message:', JSON.parse(event.data))
  }
}
