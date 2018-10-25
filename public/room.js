(function () {
	let template = document.createElement('template')
	template.innerHTML = `
		<style>
			:host {
        all: inherit;
				contain: content;
			}	

      .group {
        display: grid;
        grid-template-columns: 30px 1fr 60px;
        grid-column-gap: 5px;
        justify-content: center;
        align-items: center;

        min-height: 60px;
      }
      .group:hover,
      .group.is-selected {
        background: #EEEEEE;
        cursor: pointer;
      }

      .group:not(:last-child) {
        border-bottom: 1px solid grey;
      }

      .status {
        background: #999999;
        height: 10px;
        width: 10px;
        display: inline-block;
        border-radius: 50%;
        justify-self: center;
      }
      .status.is-online {
        background: #4caf50;
      }

      .user {
        font-weight: bold;
        display: block;
      }
      .message {
        color: #444444;
        font-size: 14px;
        display: block;
      }

      .timestamp {
        color: #444444;
        font-size: 14px;
        text-align: center;
      }
		</style>
		<div class='group'>
			<div class='status'></div>
			<div class='info'>
        <div class='user'></div>
        <div class='message'>No message</div>
      </div>
			<div class='timestamp'></div>
		</div>
	`

	class ChatRoom extends HTMLElement {
		static get observedAttributes () {
			return []
		}
		constructor() {
			super()

			this.attachShadow({mode: 'open'})
			.appendChild(template.content.cloneNode(true))
			this.state = {
				status: false,
				room: null,
				timestamp: null,
				user: null,
        selected: false
			}
		}

		connectedCallback() {

		}

		set status (value) {
			this.state.status = value
			let $status = this.shadowRoot.querySelector('.status')
			value ? $status.classList.add('is-online') : $status.classList.remove('is-online')
		}

		get status () {
			return this.state.status
		}

		set timestamp(value) {
			this.state.timestamp = value
			let $timestamp = this.shadowRoot.querySelector('.timestamp')
			$timestamp.textContent = timeDifference(Date.now(), new Date(value))
		}
		
		get timestamp() {
			return this.state.timestamp
		}

		set user (value) {
			this.state.user = value
			let $user = this.shadowRoot.querySelector('.user')
			$user.textContent = value
		}

		set room (value) {
			this.state.room = value 
		}

    set selected (value) {
      this.state.selected = value
      let $group = this.shadowRoot.querySelector('.group')
      value 
        ? $group.classList.add('is-selected')
        : $group.classList.remove('is-selected')
    }

		attributeChangedCallback(attrName, oldValue, newValue) {

		}
	}

	window.customElements.define('chat-room', ChatRoom)

  function timeDifference(current, previous) {

      var msPerMinute = 60 * 1000;
      var msPerHour = msPerMinute * 60;
      var msPerDay = msPerHour * 24;
      var msPerMonth = msPerDay * 30;
      var msPerYear = msPerDay * 365;

      var elapsed = current - previous;

      if (elapsed < msPerMinute) {
           return Math.round(elapsed/1000) + 's ago';   
      }

      else if (elapsed < msPerHour) {
           return Math.round(elapsed/msPerMinute) + 'm ago';   
      }

      else if (elapsed < msPerDay ) {
           return Math.round(elapsed/msPerHour ) + 'h ago';   
      }

      else if (elapsed < msPerMonth) {
          return Math.round(elapsed/msPerDay) + 'days ago';   
      }

  }
})()
