(function () {
  let template = document.createElement('template')
  template.innerHTML = `
    <style>
      :host {
        all: inherit;
        contain: content;
      }
      .dialog {
        display: grid;
        text-align: left;
        justify-content: flex-start;
      }

      .dialog.is-self {
        justify-content: flex-end;
      }

      .dialog .message {
        background: #4488FF;
        color: white;
        padding: 0 15px;
        border-radius: 5px 15px 15px 5px;
        line-height: 30px;
        margin: 2.5px 0;
      }

      .dialog.is-self .message {
        border-radius: 15px 5px 5px 15px;
      }

      .dialog .message:first-child {
        border-radius: 15px 15px 15px 5px;
      }

      .dialog.is-self .message:first-child {
        border-radius: 15px 15px 5px 15px;
      }

      .dialog .message:last-child {
        border-radius: 5px 15px 15px 15px;
      }

      .dialog.is-self .message:last-child {
        border-radius: 15px 5px 15px 15px;
      }

      .dialog.is-self .message {
        background: #EEEEEE;
        color: #222222;
      }
    </style>
    <div class='dialog'>
      <div class='message'>hello</div>
    </div>
  `

  class ChatDialog extends HTMLElement {
    constructor () {
      super()
      this.attachShadow({ mode: 'open' })
        .appendChild(template.content.cloneNode(true))

      this.state = {
        message: '',
        isSelf: false
      }
    }

    set message (value) {
      this.state.message = value

      let $message = this.shadowRoot.querySelector('.message')
      $message.textContent = value
    }

    set isSelf (value) {
      this.state.isSelf = value
      let $dialog = this.shadowRoot.querySelector('.dialog')
      value
        ? $dialog.classList.add('is-self')
        : $dialog.classList.remove('is-self')
    }
  }

  window.customElements.define('chat-dialog', ChatDialog)
})()
