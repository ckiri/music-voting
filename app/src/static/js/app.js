let ws = new WebSocket("ws://" + location.host + "/ws");
let wsLiveReload = new WebSocket("ws://" + location.host + "/livereload");

wsLiveReload.onmessage = () => location.reload();

function vote(box) {
    ws.send(JSON.stringify({vote: box}));
}

ws.onmessage = function(event) {
    let data = JSON.parse(event.data);
    for (let box in data.counts) {
        document.getElementById(box).innerText = data.counts[box];
    }
};
