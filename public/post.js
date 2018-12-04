(function () {
  let template = document.createElement('template')
  template.innerHTML = `
		<style>
			:host {
        all: inherit;
        contain: content;
			}	

      #status {
      }

		</style>
    <div id='post-form'>
      <form>
        <textarea id='post'></textarea>
        <button id='submit'>Submit</button>

        <div id='posts'></div>
      </form>
		</div>
	`

  const PostForm = (function () {
    const cache = new WeakMap()
    const internal = function(key) {
      if (!cache.has(key)) {
        cache.set(key, {})
      }
      return cache.get(key)
    }

    class PostForm extends HTMLElement {
      static get observedAttributes() {
        return []
      }
      constructor() {
        super()
        this.attachShadow({mode: 'open'})
        .appendChild(template.content.cloneNode(true))

        internal(this).state = {}
      }
      async connectedCallback() {
        const state = internal(this)
        state.$post = this.shadowRoot.getElementById('post')
        state.$posts = this.shadowRoot.getElementById('posts')
        state.$submit = this.shadowRoot.getElementById('submit')

        state.$submit.addEventListener('click', async(evt) => {
          evt.preventDefault()
          const text = state.$post.value
          if (!text.trim().length) {
            return
          }
          const result = await postPost(window.localStorage.access_token, text)
          console.log(result)
        })

        const posts = await fetchPosts(window.localStorage.access_token)
        for (let post of posts) {
          const div = document.createElement('div')
          div.innerHTML = `${post.id} ${post.text} ${post.created_at}`
          state.$posts.appendChild(div)
        }
      }
    }

    return PostForm
  })()

  window.customElements.define('post-form', PostForm)

  async function postPost(accessToken, text) {
    const response = await window.fetch(`/posts`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      },
      body: JSON.stringify({
        text
      })
    })
    if (!response.ok) {
      console.error(await response.text())
      return
    }
    const {status} = await response.json()
    return status 
  }


  async function fetchPosts(accessToken) {
    const response = await window.fetch('/posts', {
      method: 'GET',
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
