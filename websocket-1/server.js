const WebSocket = require("ws");

const server = new WebSocket.Server({
    port: 8080
});


const clients = new Set();


server.on("connection", (socket) => {

    console.log("client connected");


    clients.add(socket);


    socket.on("message", (message) => {

        console.log(
            "receive:",
            message.toString()
        );


        // broadcast
        for (const client of clients) {

            if (client.readyState === WebSocket.OPEN) {

                client.send(
                    message.toString()
                );

            }
        }
    });


    socket.on("close", () => {

        console.log("client disconnected");

        clients.delete(socket);

    });


});


console.log(
    "WebSocket server running :8080"
);