
function onOpen (socket) {
  return function (event) {
    // const msg = {
    //   type: 'authenticate',
    //   payload: { token: 'xxx' }
    // }
    const msg = {
      handle: 'hello socket',
      text: 'this is a new text message',
      room: 'car'
    }
    socket.send(JSON.stringify(msg))
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
      handle: model.get('username'),
      text: model.get('message'),
      room: model.get('room')
    })
    evt.currentTarget.value = ''
  }
}

(function () {
  const room = window.prompt('What room do you want to join?')
  const username = window.prompt('What is your username?')
  console.log(`hello, ${username}!`)

  const socket = new window.WebSocket(`ws://localhost:3000/ws?room=${room}`)

  socket.onopen = onOpen(socket)
  socket.onmessage = onMessage(socket)

  const view = new View()
  const model = new Model()
  model.set('room', room)
  model.set('username', username)

  const controller = new Controller({ model, view, publish: publish(socket) })
  controller.bindEvents()

  // const bus = new EventBus()
  // bus.on('hello', ({ name }) => {
  //   console.log(`greetings, ${name}!`)
  // })

  // bus.trigger('hello', { name: 'john' })
})()
