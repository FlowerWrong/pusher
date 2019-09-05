// https://github.com/pusher/pusher-js#nodejs
const Pusher = require('pusher-js');

APP_KEY = "app-key"

Pusher.log = (msg) => {
	console.log(msg);
};

for (var i = 0; i < 100; i++) {
	const socket = new Pusher(APP_KEY, {
		authEndpoint: 'http://127.0.0.1:8080/pusher/auth',
		wsHost: 'localhost',
		wsPort: 8100,
		encrypted: false,
		enabledTransports: ["ws"]
	});

	socket.connection.bind('error', function (err) {
		console.log(err);
	});

	const channel = socket.subscribe('my-channel');
	channel.bind('new-message', function (data) {
		console.log(data.message);
	});

	const privateChannel = socket.subscribe('private-my-channel');
	privateChannel.bind('new-message', function (data) {
		console.log(data.message);
	});

	const presenceChannel = socket.subscribe('presence-my-channel');
	presenceChannel.bind('new-message', function (data, metadata) {
		console.log('received data from', metadata.user_id, ':', data);
	});
}

console.log("bye")

// socket.allChannels().forEach(channel => console.log(channel.name));
