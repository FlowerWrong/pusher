// https://github.com/pusher/pusher-js#nodejs
const Pusher = require('pusher-js');

APP_KEY = "5c107909cb0b804d6e21"

Pusher.log = (msg) => {
    console.log(msg);
};

const socket = new Pusher(APP_KEY, {
    authEndpoint: 'http://127.0.0.1:8080/pusher/auth',
    cluster: 'us3',
    forceTLS: true
});

socket.connection.bind('error', function (err) {
    console.log(err);
});

// const channel = socket.subscribe('my-channel');
// channel.bind('new-message', function (data) {
//     console.log(data.message);
// });

// const privateChannel = socket.subscribe('private-my-channel');
// privateChannel.bind('new-message', function (data) {
//     console.log(data.message);
// });

const presenceChannel = socket.subscribe('presence-my-channel');
presenceChannel.bind('new-message', function (data, metadata) {
    console.log('received data from', metadata.user_id, ':', data);
});
