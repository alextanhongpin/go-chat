
function onOpen (socket) {
  return function (event) {
    // const msg = {
    //   type: 'authenticate',
    //   payload: { token: 'xxx' }
    // }
    // const msg = {
    //   type: 'authenticate',
    //   data: 'some token',
    //   room: ''
    // }
    // socket.send(JSON.stringify(msg))
  }
}

function onMessage (socket) {
  return function (event) {
    console.log('got message:', JSON.parse(event.data))
  }
}

function publish (socket) {
  return function (payload) {
    socket.send(JSON.stringify(payload))
  }
}

class EventBus {
  constructor () {
    this.events = {}
  }
  on (event, fn) {
    if (!this.events[event]) {
      this.events[event] = []
    }
    this.events[event].push(fn)
  }
  trigger (event, params) {
    if (!this.events[event]) {
      return
    }
    this.events[event].forEach(fn => {
      fn(params)
    })
  }
}

class View {
  constructor () {
    this.message = document.getElementById('message')
    this.submitMessage = document.getElementById('submit_message')
    this.username = document.getElementById('username')
    this.room = document.getElementById('room')
  }
}

class Model {
  constructor () {
    this.data = {}
  }
  set (key, value) {
    this.data[key] = value
  }
  get (key) {
    return this.data[key]
  }
}

class Controller {
  constructor ({ model, view, publish }) {
    this.model = model
    this.view = view
    this.publish = publish
  }

  bindEvents () {
    const view = this.view
    view.message.addEventListener('keyup', this.onEnterMessage.bind(this))
    view.submitMessage.addEventListener('click', this.onSubmitMessage.bind(this))
  }

  onEnterMessage (evt) {
    this.model.set('message', evt.currentTarget.value)
  }

  onSubmitMessage (evt) {
    const model = this.model
    this.publish({
      type: model.get('username'),
      data: model.get('message'),
      room: model.get('room'),
      token: model.get('token')
    })
    evt.currentTarget.value = ''
  }
}

(async function () {
  try {
    const body = await window.fetch('/auth', {
      method: 'POST'
    })
    const response = await body.json()
    if (response) {
      const { ticket } = response
      console.log('ticket:', ticket)

      const username = window.prompt('What is your username?')
      console.log(`hello, ${username}!`)

      const socket = new window.WebSocket(`ws://localhost:4000/ws?ticket=${ticket}`)

      socket.onopen = onOpen(socket)
      socket.onmessage = onMessage(socket)

      const view = new View()
      const model = new Model()
      model.set('username', username)
      model.set('room', 'abc123')
      model.set('token', ticket)

      const controller = new Controller({ model, view, publish: publish(socket) })
      controller.bindEvents()
    }
  } catch (error) {
    console.error(error.message)
  }
})()
