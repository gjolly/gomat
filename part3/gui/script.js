let identifier, messageTextarea, messagesList;

function sendMessage() {
	const msg = {
		"Origin": identifier.value,
		"PeerMessage": {
			"Text": messageTextarea.value,
		}
	};
	const headers = new Headers();
	fetch("http://localhost:8080/newMessage", {
		method: 'POST',
		mode: "no-cors",
		body: JSON.stringify(msg)
	});
}

function displayMessage(messages) {
	messagesList.innerHTML = "";

	for (msg of messages) {
		const originIdDOM = document.createElement("div");
		originIdDOM.classList.add("msgOriginId");
		originIdDOM.textContent = msg.origin + "(" + msg.id + ")";

		const textDOM = document.createElement("div");
		textDOM.classList.add("msgText");
		textDOM.textContent = msg.text;

		const msgDOM = document.createElement("div");
		msgDOM.classList.add("msg");
		msgDOM.appendChild(originIdDOM);
		msgDOM.appendChild(textDOM);

		messagesList.insertBefore(msgDOM, messagesList.firstChild);
	}
}

async function fetchNewMessages() {
	try {
		const messages = await (await fetch("http://localhost:8080/getMessages", {
			method: 'GET'
		})).json();

		console.log(messages);

		displayMessage(messages.map(msg => ({
			origin: msg.Origin,
			id: msg.PeerMessage.ID,
			text: msg.PeerMessage.Text
		})));
	} catch (err) { }
}


window.onload = (() => {
	messagesList = document.querySelector(".messages-list");
	identifier = document.querySelector("#identifier");
	messageTextarea = document.querySelector("#messageTextarea");
	messageTextarea.addEventListener('keydown', function (e) {
		var key = e.which || e.keyCode;
		if (key === 13) {
			e.preventDefault();
			if (messageTextarea.value != "") {
				sendMessage();
			}
		}
	});
	messageTextarea.addEventListener('keyup', function (e) {
		var key = e.which || e.keyCode;
		if (key === 13 && messageTextarea.value != "") {
			messageTextarea.value = "";
		}
	});

	setInterval(fetchNewMessages, 500);
});
