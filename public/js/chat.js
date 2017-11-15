const Chat = (function () {
    class Chat {
        constructor(placeholder, host) {
            this.$ = placeholder.querySelector.bind(placeholder);
            this.socket = new WebSocket(`ws://${host}/room`);
        }

        init () {
            this.getDomElements();
            this.setEvents()
        }

        getDomElements() {
            this.dom = {
                textArea: this.$('textarea'),
                listMessage: this.$('#messages'),
                form: this.$('#chatbox')
            };
        }

        setEvents() {
            this.dom.form.addEventListener('submit', this.onSubmitForm.bind(this));
            this.socket.addEventListener('close', () => {
                console.log('Connection has been closed');
            });
            this.socket.addEventListener('message', this.onReceiveMessage.bind(this))
        }

        onSubmitForm (event) {
            event.preventDefault();
            const textAreaValue = this.dom.textArea.value;
            if (textAreaValue === "" || !this.socket) {
                return
            }
            this.socket.send(JSON.stringify({"Message": textAreaValue}));
            this.dom.textArea.value = ""
        }

        onReceiveMessage (event) {
            const item = document.createElement('li');
            const strong = document.createElement('strong');
            const span = document.createElement('span');
            const small = document.createElement('small');
            const img = document.createElement('img');
            const msg = JSON.parse(event.data);
            img.src = msg.AvatarUrl;
            span.textContent = msg.Message;
            strong.textContent = msg.Name + ": ";
            small.textContent = new Date(msg.When).toLocaleTimeString();
            item.appendChild(img);
            item.appendChild(strong);
            item.appendChild(span);
            item.appendChild(small);
            this.dom.listMessage.appendChild(item)
        }


    }
    return Chat
})();

window.addEventListener('load', () => {
    const placeholder = document.querySelector('#main-section');
    const chat = new Chat(placeholder, params.host);
    chat.init();
});