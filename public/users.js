(function () {
  let template = document.createElement('template')
  template.innerHTML = `
		<style>
			:host {
        all: inherit;
        contain: content;
			}	

		</style>
    <div id='user'>
      <div id='name'></div>
		</div>
	`

  const UserList = (function () {
    const cache = new WeakMap()
    const internal = function(key) {
      if (!cache.has(key)) {
        cache.set(key, {})
      }
      return cache.get(key)
    }

    class UserList extends HTMLElement {
      static get observedAttributes() {
        return ['users']
      }
      constructor() {
        super()
        this.attachShadow({mode: 'open'})
        .appendChild(template.content.cloneNode(true))

        internal(this).state = {}
      }
      async connectedCallback() {
        this.users = []
        const users = await fetchUsers(window.localStorage.access_token)
        this.users = users
      }
      set users(value) {
        internal(this).state.users = value
        this.renderUsers(value)
      }
      get users () {
        return internal(this).state.users
      }
      renderUsers(users) {
            const $users = this.shadowRoot.getElementById('user')
            $users.innerHTML = ''
            for (let user of users) {
              const $user = document.createElement('user-row')
              $user.id = user.id
              $user.name = user.name
              $user.status = user.status
              $user.isRequested = user.is_requested
              $users.appendChild($user)
            }
      }
      // attributeChangedCallback(attrName, oldValue, newValue) {
      //   switch (attrName) {
      //     case 'users':
      //       break
      //     default:
      //   }
      // }
    }

    return UserList
  })()

  window.customElements.define('user-list', UserList)

  async function fetchUsers(accessToken) {
    const response = await window.fetch('/users', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      }
    })
    if (!response.ok) {
      console.error(await response.text())
      return []
    }
    const {data} = await response.json()
    return data
  }
})()
