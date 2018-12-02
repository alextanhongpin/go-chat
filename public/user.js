(function () {
  let template = document.createElement('template')
  template.innerHTML = `
		<style>
			:host {
        all: inherit;
        contain: content;
			}	

      #user-list {
        border: 1px solid #EEEEEE;
        padding: 10px;
        min-width: 160px;
        display: inline-block;
      }

      #name {
        font-weight: bold;
      }
      #status {
      }

		</style>
    <div id='user-list'>
      <div id='name'></div>
      <div id='status'>None</div>
      <button id='submit'>Add</button>
		</div>
	`

  const UserRow = (function () {
    const cache = new WeakMap()
    const internal = function(key) {
      if (!cache.has(key)) {
        cache.set(key, {})
      }
      return cache.get(key)
    }

    class UserRow extends HTMLElement {
      static get observedAttributes() {
        return ['id', 'name']
      }
      constructor() {
        super()
        this.attachShadow({mode: 'open'})
        .appendChild(template.content.cloneNode(true))

        internal(this).state = {}
      }
      connectedCallback() {
        this.shadowRoot.getElementById('submit').addEventListener('click', async(evt) => {

          switch (evt.target.dataset.action) {
            case 'add': {
              const result = await postFriendship({ 
                accessToken: window.localStorage.access_token,
                friendId: this.id
              })
              this.isRequested = true
              this.status = 'request'
              console.log('add friend', evt.target.dataset.action)
              return
            }
            case 'unfriend':
            case 'cancel': {
              const result = await patchFriendship({ 
                action: 'reject',
                accessToken: window.localStorage.access_token,
                friendId: this.id
              })
              this.isRequested = false
              this.status = ''
              console.log('reject friend', evt.target.dataset.action)
              return
            }
            case 'accept': {
              const result = await patchFriendship({ 
                action: 'accept',
                accessToken: window.localStorage.access_token,
                friendId: this.id
              })
              this.status = 'friend'
              console.log('accept friend', evt.target.dataset.action)
              return
            }
          }
        })
      }
      set id (value) {
        internal(this).state.id = value
      }
      get id () {
        return internal(this).state.id
      }
      set name (value) {
        internal(this).state.name = value
        this.renderName()
      }
      get name () {
        return internal(this).state.name
      }
      set status(value) {
        internal(this).state.status = value
        this.renderStatus()
      }
      get status() {
        return internal(this).state.status
      }
      set isRequested (value) {
        internal(this).state.isRequested = value
        this.renderStatus()
      }
      get isRequested () {
        return internal(this).state.isRequested
      }
      renderName() {
        this.shadowRoot.getElementById('name').innerHTML = this.name
      }
      renderStatus() {
        const status = parseStatus(this.status, this.isRequested)
        this.shadowRoot.getElementById('status').innerHTML = status 

        const action = parseAction(this.status, this.isRequested)
        const $submit = this.shadowRoot.getElementById('submit')
        $submit.innerHTML = action 
        $submit.dataset.action = action 
      }
      // attributeChangedCallback(attrName, oldValue, newValue) {
      //   switch (attrName) {
      //     case 'name':
      //       this.shadowRoot.getElementById('name').innerHTML = newValue
      //       break
      //     default:
      //   }
      // }
    }

    return UserRow
  })()

  window.customElements.define('user-row', UserRow)

  async function patchFriendship({ accessToken, action, friendId }) {
    const response = await window.fetch(`/friends/${friendId}`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      },
      body: JSON.stringify({
        action 
      })
    })
    if (!response.ok) {
      console.error(await response.text())
      return null
    }
    const { status } = await response.json()
    return status 
  }

  async function postFriendship({accessToken, friendId}) {
    const response = await window.fetch(`/friends/${friendId}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      }
    })

    if (!response.ok) {
      console.error(await response.text())
      return
    }
    const {status} = await response.json()
    return status 
  }

  function parseStatus(status = '', isRequested = false) {
    switch (status) {
      case 'request':
        return isRequested ? 'requested' : 'pending'
      case 'friend': 
      case 'block':
        return status
      default:
        return 'not contact'
    }
  }

  function parseAction(status = '', isRequested = false) {
    switch (status) {
      case 'request':
        return isRequested ? 'cancel' : 'accept'
      case 'friend':
        return 'unfriend'
      case 'block':
        return 'unblock'
      default:
        return 'add'
    }
  }

})()
