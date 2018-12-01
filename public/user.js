(function () {
  let template = document.createElement('template')
  template.innerHTML = `
		<style>
			:host {
        all: inherit;
        contain: content;
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
        this.shadowRoot.getElementById('submit').addEventListener('click', async() => {
          if (!this.status) {
            const result = await postFriendship({ 
              action: this.status,
              accessToken: window.localStorage.access_token,
              friendId: this.id
            })
            console.log(result)
            return
          } else {
            const result = await postFriendship({ 
              accessToken: window.localStorage.access_token,
              friendId: this.id
            })
            console.log(result)
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
        this.renderName(value)
      }
      get name () {
        return internal(this).state.name
      }
      set status(value) {
        internal(this).state.status = value
        this.renderStatus(value)
      }
      get status() {
        return internal(this).state.status
      }
      renderName(name) {
        this.shadowRoot.getElementById('name').innerHTML = name 
      }
      renderStatus(status) {
        this.shadowRoot.getElementById('status').innerHTML = status
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
    const { data } = await response.json()
    return data
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
    const {data} = await response.json()
    return data
  }
})()
