
(function() {
	const socket = new WebSocket('ws://localhost:8080/ws')

	socket.onopen = function (evt) {
		socket.send(JSON.stringify({
			handle: 'hello socket',
			text: 'this is a new text message'
		}))
	}

})()