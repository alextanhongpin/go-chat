
(function() {
	const socket = new WebSocket('ws://localhost:8080/ws')

	socket.onopen = function (evt) {
		const msg = {
			type: 'authenticate',
			payload: { token: 'xxx' }
		}
		socket.send(JSON.stringify(msg))
	}
		// socket.send(JSON.stringify({
		// 	handle: 'hello socket',
		// 	text: 'this is a new text message'
		// }))

})()